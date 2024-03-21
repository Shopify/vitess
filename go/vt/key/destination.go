/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package key

import (
	"bytes"
	"encoding/hex"
	"math/rand"
	"sort"
	"strings"

	"vitess.io/vitess/go/vt/vterrors"

	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
	vtrpcpb "vitess.io/vitess/go/vt/proto/vtrpc"
)

// AnyShardPicker makes a choice on what shard to use when any shard will do. Used for testing.
var AnyShardPicker DestinationAnyShardPicker = DestinationAnyShardPickerRandomShard{}

// Destination is an interface definition for a query destination,
// within a given Keyspace / Tablet Type. It is meant to be an internal
// data structure, with multiple possible implementations.
// The srvtopo package can resolve Destinations into actual Targets.
type Destination interface {
	// Resolve calls the callback for every shard Destination
	// resolves into, given the shards list.
	// The returned error must be generated by vterrors.
	Resolve([]*topodatapb.ShardReference, func(shard string) error) error

	// String returns a printable version of the Destination.
	String() string
}

// DestinationsString returns a printed version of the destination array.
func DestinationsString(destinations []Destination) string {
	var buffer bytes.Buffer
	buffer.WriteString("Destinations:")
	for i, d := range destinations {
		if i > 0 {
			buffer.WriteByte(',')
		}
		buffer.WriteString(d.String())
	}
	return buffer.String()
}

//
// DestinationShard
//

// DestinationShard is the destination for a single Shard.
// It implements the Destination interface.
type DestinationShard string

// Resolve is part of the Destination interface.
func (d DestinationShard) Resolve(allShards []*topodatapb.ShardReference, addShard func(shard string) error) error {
	return addShard(string(d))
}

// String is part of the Destination interface.
func (d DestinationShard) String() string {
	return "DestinationShard(" + string(d) + ")"
}

//
// DestinationShards
//

// DestinationShards is the destination for multiple shards.
// It implements the Destination interface.
type DestinationShards []string

// Resolve is part of the Destination interface.
func (d DestinationShards) Resolve(allShards []*topodatapb.ShardReference, addShard func(shard string) error) error {
	for _, shard := range d {
		if err := addShard(shard); err != nil {
			return err
		}
	}
	return nil
}

// String is part of the Destination interface.
func (d DestinationShards) String() string {
	return "DestinationShards(" + strings.Join(d, ",") + ")"
}

//
// DestinationExactKeyRange
//

// DestinationExactKeyRange is the destination for a single KeyRange.
// The KeyRange must map exactly to one or more shards, and cannot
// start or end in the middle of a shard.
// It implements the Destination interface.
// (it cannot be just a type *topodatapb.KeyRange, as then the receiver
// methods don't work. And it can't be topodatapb.KeyRange either,
// as then the methods are on *DestinationExactKeyRange, and the original
// KeyRange cannot be returned).
type DestinationExactKeyRange struct {
	KeyRange *topodatapb.KeyRange
}

// Resolve is part of the Destination interface.
func (d DestinationExactKeyRange) Resolve(allShards []*topodatapb.ShardReference, addShard func(shard string) error) error {
	return processExactKeyRange(allShards, d.KeyRange, addShard)
}

// String is part of the Destination interface.
func (d DestinationExactKeyRange) String() string {
	return "DestinationExactKeyRange(" + KeyRangeString(d.KeyRange) + ")"
}

func processExactKeyRange(allShards []*topodatapb.ShardReference, kr *topodatapb.KeyRange, addShard func(shard string) error) error {
	sort.SliceStable(allShards, func(i, j int) bool {
		return KeyRangeLess(allShards[i].GetKeyRange(), allShards[j].GetKeyRange())
	})

	shardnum := 0
	for shardnum < len(allShards) {
		if KeyRangeStartEqual(kr, allShards[shardnum].KeyRange) {
			break
		}
		shardnum++
	}
	for shardnum < len(allShards) {
		if !KeyRangeIntersect(kr, allShards[shardnum].KeyRange) {
			// If we are over the requested keyrange, we
			// can stop now, we won't find more.
			break
		}
		if err := addShard(allShards[shardnum].Name); err != nil {
			return err
		}
		if KeyRangeEndEqual(kr, allShards[shardnum].KeyRange) {
			return nil
		}
		shardnum++
	}
	return vterrors.Errorf(vtrpcpb.Code_INVALID_ARGUMENT, "keyrange %v does not exactly match shards", KeyRangeString(kr))
}

//
// DestinationExactKeyRanges
//

// DestinationExactKeyRanges is the destination for multiple KeyRanges.
// The KeyRanges must map exactly to one or more shards, and cannot
// start or end in the middle of a shard.
// It implements the Destination interface.
type DestinationExactKeyRanges []*topodatapb.KeyRange

// Resolve is part of the Destination interface.
func (d DestinationExactKeyRanges) Resolve(allShards []*topodatapb.ShardReference, addShard func(shard string) error) error {
	for _, kr := range d {
		if err := processExactKeyRange(allShards, kr, addShard); err != nil {
			return err
		}
	}
	return nil
}

// String is part of the Destination interface.
func (d DestinationExactKeyRanges) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("DestinationExactKeyRanges(")
	for i, kr := range d {
		if i > 0 {
			buffer.WriteByte(',')
		}
		buffer.WriteString(KeyRangeString(kr))
	}
	buffer.WriteByte(')')
	return buffer.String()
}

//
// DestinationKeyRange
//

// DestinationKeyRange is the destination for a single KeyRange.
// It implements the Destination interface.
// (it cannot be just a type *topodatapb.KeyRange, as then the receiver
// methods don't work. And it can't be topodatapb.KeyRange either,
// as then the methods are on *DestinationKeyRange, and the original
// KeyRange cannot be returned).
type DestinationKeyRange struct {
	KeyRange *topodatapb.KeyRange
}

// Resolve is part of the Destination interface.
func (d DestinationKeyRange) Resolve(allShards []*topodatapb.ShardReference, addShard func(shard string) error) error {
	return processKeyRange(allShards, d.KeyRange, addShard)
}

// String is part of the Destination interface.
func (d DestinationKeyRange) String() string {
	return "DestinationKeyRange(" + KeyRangeString(d.KeyRange) + ")"
}

func processKeyRange(allShards []*topodatapb.ShardReference, kr *topodatapb.KeyRange, addShard func(shard string) error) error {
	for _, shard := range allShards {
		if !KeyRangeIntersect(kr, shard.KeyRange) {
			// We don't need that shard.
			continue
		}
		if err := addShard(shard.Name); err != nil {
			return err
		}
	}
	return nil
}

//
// DestinationKeyRanges
//

// DestinationKeyRanges is the destination for multiple KeyRanges.
// It implements the Destination interface.
type DestinationKeyRanges []*topodatapb.KeyRange

// Resolve is part of the Destination interface.
func (d DestinationKeyRanges) Resolve(allShards []*topodatapb.ShardReference, addShard func(shard string) error) error {
	for _, kr := range d {
		if err := processKeyRange(allShards, kr, addShard); err != nil {
			return err
		}
	}
	return nil
}

// String is part of the Destination interface.
func (d DestinationKeyRanges) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("DestinationKeyRanges(")
	for i, kr := range d {
		if i > 0 {
			buffer.WriteByte(',')
		}
		buffer.WriteString(KeyRangeString(kr))
	}
	buffer.WriteByte(')')
	return buffer.String()
}

//
// DestinationKeyspaceID
//

// DestinationKeyspaceID is the destination for a single KeyspaceID.
// It implements the Destination interface.
type DestinationKeyspaceID []byte

// Resolve is part of the Destination interface.
func (d DestinationKeyspaceID) Resolve(allShards []*topodatapb.ShardReference, addShard func(shard string) error) error {
	shard, err := GetShardForKeyspaceID(allShards, d)
	if err != nil {
		return err
	}
	return addShard(shard)
}

// String is part of the Destination interface.
func (d DestinationKeyspaceID) String() string {
	return "DestinationKeyspaceID(" + hex.EncodeToString(d) + ")"
}

// GetShardForKeyspaceID finds the right shard for a keyspace id.
func GetShardForKeyspaceID(allShards []*topodatapb.ShardReference, keyspaceID []byte) (string, error) {
	if len(allShards) == 0 {
		return "", vterrors.Errorf(vtrpcpb.Code_UNAVAILABLE, "no shard in keyspace")
	}

	for _, shardReference := range allShards {
		if KeyRangeContains(shardReference.KeyRange, keyspaceID) {
			return shardReference.Name, nil
		}
	}
	return "", vterrors.Errorf(vtrpcpb.Code_INVALID_ARGUMENT, "KeyspaceId %v didn't match any shards %+v", hex.EncodeToString(keyspaceID), allShards)
}

//
// DestinationKeyspaceIDs
//

// DestinationKeyspaceIDs is the destination for multiple KeyspaceIDs.
// It implements the Destination interface.
type DestinationKeyspaceIDs [][]byte

// Resolve is part of the Destination interface.
func (d DestinationKeyspaceIDs) Resolve(allShards []*topodatapb.ShardReference, addShard func(shard string) error) error {
	for _, ksid := range d {
		shard, err := GetShardForKeyspaceID(allShards, ksid)
		if err != nil {
			return err
		}
		if err := addShard(shard); err != nil {
			return err
		}
	}
	return nil
}

// String is part of the Destination interface.
func (d DestinationKeyspaceIDs) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("DestinationKeyspaceIDs(")
	for i, ksid := range d {
		if i > 0 {
			buffer.WriteByte(',')
		}
		buffer.WriteString(hex.EncodeToString(ksid))
	}
	buffer.WriteByte(')')
	return buffer.String()
}

// DestinationAnyShardPicker exposes an interface that will pick an index given a number of available shards.
type DestinationAnyShardPicker interface {
	// PickShard picks a shard given a number of shards
	PickShard(shardCount int) int
}

// DestinationAnyShardPickerRandomShard picks a random shard.
type DestinationAnyShardPickerRandomShard struct{}

// PickShard is DestinationAnyShardPickerRandomShard's implementation.
func (dp DestinationAnyShardPickerRandomShard) PickShard(shardCount int) int {
	return rand.Intn(shardCount)
}

//
// DestinationAnyShard
//

// DestinationAnyShard is the destination for any one shard in the keyspace.
// It implements the Destination interface.
type DestinationAnyShard struct{}

// Resolve is part of the Destination interface.
func (d DestinationAnyShard) Resolve(allShards []*topodatapb.ShardReference, addShard func(shard string) error) error {
	if len(allShards) == 0 {
		return vterrors.Errorf(vtrpcpb.Code_UNAVAILABLE, "no shard in keyspace")
	}
	return addShard(allShards[AnyShardPicker.PickShard(len(allShards))].Name)
}

// String is part of the Destination interface.
func (d DestinationAnyShard) String() string {
	return "DestinationAnyShard()"
}

//
// DestinationAllShards
//

// DestinationAllShards is the destination for all the shards in the
// keyspace. This usually maps to the first one in the list.
// It implements the Destination interface.
type DestinationAllShards struct{}

// Resolve is part of the Destination interface.
func (d DestinationAllShards) Resolve(allShards []*topodatapb.ShardReference, addShard func(shard string) error) error {
	for _, shard := range allShards {
		if err := addShard(shard.Name); err != nil {
			return err
		}
	}

	return nil
}

// String is part of the Destination interface.
func (d DestinationAllShards) String() string {
	return "DestinationAllShards()"
}

//
// DestinationNone
//

// DestinationNone is a destination that doesn't resolve to any shard.
// It implements the Destination interface.
type DestinationNone struct{}

// Resolve is part of the Destination interface.
func (d DestinationNone) Resolve(allShards []*topodatapb.ShardReference, addShard func(shard string) error) error {
	return nil
}

// String is part of the Destination interface.
func (d DestinationNone) String() string {
	return "DestinationNone()"
}
