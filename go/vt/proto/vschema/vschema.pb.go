//
//Copyright 2019 The Vitess Authors.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

// This file contains the types needed to define a vschema.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v3.21.3
// source: vschema.proto

package vschema

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	query "vitess.io/vitess/go/vt/proto/query"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// RoutingRules specify the high level routing rules for the VSchema.
type RoutingRules struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// rules should ideally be a map. However protos dont't allow
	// repeated fields as elements of a map. So, we use a list
	// instead.
	Rules []*RoutingRule `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
}

func (x *RoutingRules) Reset() {
	*x = RoutingRules{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vschema_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoutingRules) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoutingRules) ProtoMessage() {}

func (x *RoutingRules) ProtoReflect() protoreflect.Message {
	mi := &file_vschema_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoutingRules.ProtoReflect.Descriptor instead.
func (*RoutingRules) Descriptor() ([]byte, []int) {
	return file_vschema_proto_rawDescGZIP(), []int{0}
}

func (x *RoutingRules) GetRules() []*RoutingRule {
	if x != nil {
		return x.Rules
	}
	return nil
}

// RoutingRule specifies a routing rule.
type RoutingRule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromTable string   `protobuf:"bytes,1,opt,name=from_table,json=fromTable,proto3" json:"from_table,omitempty"`
	ToTables  []string `protobuf:"bytes,2,rep,name=to_tables,json=toTables,proto3" json:"to_tables,omitempty"`
}

func (x *RoutingRule) Reset() {
	*x = RoutingRule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vschema_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RoutingRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RoutingRule) ProtoMessage() {}

func (x *RoutingRule) ProtoReflect() protoreflect.Message {
	mi := &file_vschema_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RoutingRule.ProtoReflect.Descriptor instead.
func (*RoutingRule) Descriptor() ([]byte, []int) {
	return file_vschema_proto_rawDescGZIP(), []int{1}
}

func (x *RoutingRule) GetFromTable() string {
	if x != nil {
		return x.FromTable
	}
	return ""
}

func (x *RoutingRule) GetToTables() []string {
	if x != nil {
		return x.ToTables
	}
	return nil
}

// Keyspace is the vschema for a keyspace.
type Keyspace struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// If sharded is false, vindexes and tables are ignored.
	Sharded  bool               `protobuf:"varint,1,opt,name=sharded,proto3" json:"sharded,omitempty"`
	Vindexes map[string]*Vindex `protobuf:"bytes,2,rep,name=vindexes,proto3" json:"vindexes,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Tables   map[string]*Table  `protobuf:"bytes,3,rep,name=tables,proto3" json:"tables,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// If require_explicit_routing is true, vindexes and tables are not added to global routing
	RequireExplicitRouting bool `protobuf:"varint,4,opt,name=require_explicit_routing,json=requireExplicitRouting,proto3" json:"require_explicit_routing,omitempty"`
}

func (x *Keyspace) Reset() {
	*x = Keyspace{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vschema_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Keyspace) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Keyspace) ProtoMessage() {}

func (x *Keyspace) ProtoReflect() protoreflect.Message {
	mi := &file_vschema_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Keyspace.ProtoReflect.Descriptor instead.
func (*Keyspace) Descriptor() ([]byte, []int) {
	return file_vschema_proto_rawDescGZIP(), []int{2}
}

func (x *Keyspace) GetSharded() bool {
	if x != nil {
		return x.Sharded
	}
	return false
}

func (x *Keyspace) GetVindexes() map[string]*Vindex {
	if x != nil {
		return x.Vindexes
	}
	return nil
}

func (x *Keyspace) GetTables() map[string]*Table {
	if x != nil {
		return x.Tables
	}
	return nil
}

func (x *Keyspace) GetRequireExplicitRouting() bool {
	if x != nil {
		return x.RequireExplicitRouting
	}
	return false
}

// Vindex is the vindex info for a Keyspace.
type Vindex struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The type must match one of the predefined
	// (or plugged in) vindex names.
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// params is a map of attribute value pairs
	// that must be defined as required by the
	// vindex constructors. The values can only
	// be strings.
	Params map[string]string `protobuf:"bytes,2,rep,name=params,proto3" json:"params,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// A lookup vindex can have an owner table defined.
	// If so, rows in the lookup table are created or
	// deleted in sync with corresponding rows in the
	// owner table.
	Owner string `protobuf:"bytes,3,opt,name=owner,proto3" json:"owner,omitempty"`
}

func (x *Vindex) Reset() {
	*x = Vindex{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vschema_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vindex) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vindex) ProtoMessage() {}

func (x *Vindex) ProtoReflect() protoreflect.Message {
	mi := &file_vschema_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vindex.ProtoReflect.Descriptor instead.
func (*Vindex) Descriptor() ([]byte, []int) {
	return file_vschema_proto_rawDescGZIP(), []int{3}
}

func (x *Vindex) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Vindex) GetParams() map[string]string {
	if x != nil {
		return x.Params
	}
	return nil
}

func (x *Vindex) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

// Table is the table info for a Keyspace.
type Table struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// If the table is a sequence, type must be
	// "sequence".
	//
	// If the table is a reference, type must be
	// "reference".
	// See https://vitess.io/docs/reference/features/vschema/#reference-tables.
	//
	// Otherwise, it should be empty.
	Type string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	// column_vindexes associates columns to vindexes.
	ColumnVindexes []*ColumnVindex `protobuf:"bytes,2,rep,name=column_vindexes,json=columnVindexes,proto3" json:"column_vindexes,omitempty"`
	// auto_increment is specified if a column needs
	// to be associated with a sequence.
	AutoIncrement *AutoIncrement `protobuf:"bytes,3,opt,name=auto_increment,json=autoIncrement,proto3" json:"auto_increment,omitempty"`
	// columns lists the columns for the table.
	Columns []*Column `protobuf:"bytes,4,rep,name=columns,proto3" json:"columns,omitempty"`
	// pinned pins an unsharded table to a specific
	// shard, as dictated by the keyspace id.
	// The keyspace id is represented in hex form
	// like in keyranges.
	Pinned string `protobuf:"bytes,5,opt,name=pinned,proto3" json:"pinned,omitempty"`
	// column_list_authoritative is set to true if columns is
	// an authoritative list for the table. This allows
	// us to expand 'select *' expressions.
	ColumnListAuthoritative bool `protobuf:"varint,6,opt,name=column_list_authoritative,json=columnListAuthoritative,proto3" json:"column_list_authoritative,omitempty"`
	// reference tables may optionally indicate their source table.
	Source string `protobuf:"bytes,7,opt,name=source,proto3" json:"source,omitempty"`
}

func (x *Table) Reset() {
	*x = Table{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vschema_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Table) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Table) ProtoMessage() {}

func (x *Table) ProtoReflect() protoreflect.Message {
	mi := &file_vschema_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Table.ProtoReflect.Descriptor instead.
func (*Table) Descriptor() ([]byte, []int) {
	return file_vschema_proto_rawDescGZIP(), []int{4}
}

func (x *Table) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Table) GetColumnVindexes() []*ColumnVindex {
	if x != nil {
		return x.ColumnVindexes
	}
	return nil
}

func (x *Table) GetAutoIncrement() *AutoIncrement {
	if x != nil {
		return x.AutoIncrement
	}
	return nil
}

func (x *Table) GetColumns() []*Column {
	if x != nil {
		return x.Columns
	}
	return nil
}

func (x *Table) GetPinned() string {
	if x != nil {
		return x.Pinned
	}
	return ""
}

func (x *Table) GetColumnListAuthoritative() bool {
	if x != nil {
		return x.ColumnListAuthoritative
	}
	return false
}

func (x *Table) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

// ColumnVindex is used to associate a column to a vindex.
type ColumnVindex struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Legacy implementation, moving forward all vindexes should define a list of columns.
	Column string `protobuf:"bytes,1,opt,name=column,proto3" json:"column,omitempty"`
	// The name must match a vindex defined in Keyspace.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// List of columns that define this Vindex
	Columns []string `protobuf:"bytes,3,rep,name=columns,proto3" json:"columns,omitempty"`
}

func (x *ColumnVindex) Reset() {
	*x = ColumnVindex{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vschema_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ColumnVindex) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ColumnVindex) ProtoMessage() {}

func (x *ColumnVindex) ProtoReflect() protoreflect.Message {
	mi := &file_vschema_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ColumnVindex.ProtoReflect.Descriptor instead.
func (*ColumnVindex) Descriptor() ([]byte, []int) {
	return file_vschema_proto_rawDescGZIP(), []int{5}
}

func (x *ColumnVindex) GetColumn() string {
	if x != nil {
		return x.Column
	}
	return ""
}

func (x *ColumnVindex) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ColumnVindex) GetColumns() []string {
	if x != nil {
		return x.Columns
	}
	return nil
}

// Autoincrement is used to designate a column as auto-inc.
type AutoIncrement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Column string `protobuf:"bytes,1,opt,name=column,proto3" json:"column,omitempty"`
	// The sequence must match a table of type SEQUENCE.
	Sequence string `protobuf:"bytes,2,opt,name=sequence,proto3" json:"sequence,omitempty"`
}

func (x *AutoIncrement) Reset() {
	*x = AutoIncrement{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vschema_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AutoIncrement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AutoIncrement) ProtoMessage() {}

func (x *AutoIncrement) ProtoReflect() protoreflect.Message {
	mi := &file_vschema_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AutoIncrement.ProtoReflect.Descriptor instead.
func (*AutoIncrement) Descriptor() ([]byte, []int) {
	return file_vschema_proto_rawDescGZIP(), []int{6}
}

func (x *AutoIncrement) GetColumn() string {
	if x != nil {
		return x.Column
	}
	return ""
}

func (x *AutoIncrement) GetSequence() string {
	if x != nil {
		return x.Sequence
	}
	return ""
}

// Column describes a column.
type Column struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Type query.Type `protobuf:"varint,2,opt,name=type,proto3,enum=query.Type" json:"type,omitempty"`
}

func (x *Column) Reset() {
	*x = Column{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vschema_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Column) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Column) ProtoMessage() {}

func (x *Column) ProtoReflect() protoreflect.Message {
	mi := &file_vschema_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Column.ProtoReflect.Descriptor instead.
func (*Column) Descriptor() ([]byte, []int) {
	return file_vschema_proto_rawDescGZIP(), []int{7}
}

func (x *Column) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Column) GetType() query.Type {
	if x != nil {
		return x.Type
	}
	return query.Type(0)
}

// SrvVSchema is the roll-up of all the Keyspace schema for a cell.
type SrvVSchema struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// keyspaces is a map of keyspace name -> Keyspace object.
	Keyspaces         map[string]*Keyspace `protobuf:"bytes,1,rep,name=keyspaces,proto3" json:"keyspaces,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	RoutingRules      *RoutingRules        `protobuf:"bytes,2,opt,name=routing_rules,json=routingRules,proto3" json:"routing_rules,omitempty"` // table routing rules
	ShardRoutingRules *ShardRoutingRules   `protobuf:"bytes,3,opt,name=shard_routing_rules,json=shardRoutingRules,proto3" json:"shard_routing_rules,omitempty"`
}

func (x *SrvVSchema) Reset() {
	*x = SrvVSchema{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vschema_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SrvVSchema) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SrvVSchema) ProtoMessage() {}

func (x *SrvVSchema) ProtoReflect() protoreflect.Message {
	mi := &file_vschema_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SrvVSchema.ProtoReflect.Descriptor instead.
func (*SrvVSchema) Descriptor() ([]byte, []int) {
	return file_vschema_proto_rawDescGZIP(), []int{8}
}

func (x *SrvVSchema) GetKeyspaces() map[string]*Keyspace {
	if x != nil {
		return x.Keyspaces
	}
	return nil
}

func (x *SrvVSchema) GetRoutingRules() *RoutingRules {
	if x != nil {
		return x.RoutingRules
	}
	return nil
}

func (x *SrvVSchema) GetShardRoutingRules() *ShardRoutingRules {
	if x != nil {
		return x.ShardRoutingRules
	}
	return nil
}

// ShardRoutingRules specify the shard routing rules for the VSchema.
type ShardRoutingRules struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Rules []*ShardRoutingRule `protobuf:"bytes,1,rep,name=rules,proto3" json:"rules,omitempty"`
}

func (x *ShardRoutingRules) Reset() {
	*x = ShardRoutingRules{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vschema_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShardRoutingRules) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShardRoutingRules) ProtoMessage() {}

func (x *ShardRoutingRules) ProtoReflect() protoreflect.Message {
	mi := &file_vschema_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShardRoutingRules.ProtoReflect.Descriptor instead.
func (*ShardRoutingRules) Descriptor() ([]byte, []int) {
	return file_vschema_proto_rawDescGZIP(), []int{9}
}

func (x *ShardRoutingRules) GetRules() []*ShardRoutingRule {
	if x != nil {
		return x.Rules
	}
	return nil
}

// RoutingRule specifies a routing rule.
type ShardRoutingRule struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromKeyspace string `protobuf:"bytes,1,opt,name=from_keyspace,json=fromKeyspace,proto3" json:"from_keyspace,omitempty"`
	ToKeyspace   string `protobuf:"bytes,2,opt,name=to_keyspace,json=toKeyspace,proto3" json:"to_keyspace,omitempty"`
	Shard        string `protobuf:"bytes,3,opt,name=shard,proto3" json:"shard,omitempty"`
}

func (x *ShardRoutingRule) Reset() {
	*x = ShardRoutingRule{}
	if protoimpl.UnsafeEnabled {
		mi := &file_vschema_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ShardRoutingRule) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ShardRoutingRule) ProtoMessage() {}

func (x *ShardRoutingRule) ProtoReflect() protoreflect.Message {
	mi := &file_vschema_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ShardRoutingRule.ProtoReflect.Descriptor instead.
func (*ShardRoutingRule) Descriptor() ([]byte, []int) {
	return file_vschema_proto_rawDescGZIP(), []int{10}
}

func (x *ShardRoutingRule) GetFromKeyspace() string {
	if x != nil {
		return x.FromKeyspace
	}
	return ""
}

func (x *ShardRoutingRule) GetToKeyspace() string {
	if x != nil {
		return x.ToKeyspace
	}
	return ""
}

func (x *ShardRoutingRule) GetShard() string {
	if x != nil {
		return x.Shard
	}
	return ""
}

var File_vschema_proto protoreflect.FileDescriptor

var file_vschema_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x1a, 0x0b, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3a, 0x0a, 0x0c, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67,
	0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x2a, 0x0a, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x52,
	0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65,
	0x73, 0x22, 0x49, 0x0a, 0x0b, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x75, 0x6c, 0x65,
	0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x72, 0x6f, 0x6d, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x74, 0x6f, 0x5f, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x08, 0x74, 0x6f, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x22, 0xeb, 0x02, 0x0a,
	0x08, 0x4b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x68, 0x61,
	0x72, 0x64, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x68, 0x61, 0x72,
	0x64, 0x65, 0x64, 0x12, 0x3b, 0x0a, 0x08, 0x76, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e,
	0x4b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x2e, 0x56, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x08, 0x76, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73,
	0x12, 0x35, 0x0a, 0x06, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x1d, 0x2e, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x4b, 0x65, 0x79, 0x73, 0x70,
	0x61, 0x63, 0x65, 0x2e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x06, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x38, 0x0a, 0x18, 0x72, 0x65, 0x71, 0x75, 0x69,
	0x72, 0x65, 0x5f, 0x65, 0x78, 0x70, 0x6c, 0x69, 0x63, 0x69, 0x74, 0x5f, 0x72, 0x6f, 0x75, 0x74,
	0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x16, 0x72, 0x65, 0x71, 0x75, 0x69,
	0x72, 0x65, 0x45, 0x78, 0x70, 0x6c, 0x69, 0x63, 0x69, 0x74, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e,
	0x67, 0x1a, 0x4c, 0x0a, 0x0d, 0x56, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x25, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x56, 0x69,
	0x6e, 0x64, 0x65, 0x78, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x1a,
	0x49, 0x0a, 0x0b, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x24, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x0e, 0x2e, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xa2, 0x01, 0x0a, 0x06, 0x56,
	0x69, 0x6e, 0x64, 0x65, 0x78, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x33, 0x0a, 0x06, 0x70, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x76, 0x73, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x2e, 0x56, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x2e, 0x50, 0x61, 0x72, 0x61, 0x6d,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x70, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f,
	0x77, 0x6e, 0x65, 0x72, 0x1a, 0x39, 0x0a, 0x0b, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22,
	0xb1, 0x02, 0x0a, 0x05, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x3e, 0x0a,
	0x0f, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x5f, 0x76, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61,
	0x2e, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x56, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x52, 0x0e, 0x63,
	0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x56, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x65, 0x73, 0x12, 0x3d, 0x0a,
	0x0e, 0x61, 0x75, 0x74, 0x6f, 0x5f, 0x69, 0x6e, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e,
	0x41, 0x75, 0x74, 0x6f, 0x49, 0x6e, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0d, 0x61,
	0x75, 0x74, 0x6f, 0x49, 0x6e, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x29, 0x0a, 0x07,
	0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x52, 0x07,
	0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x70, 0x69, 0x6e, 0x6e, 0x65,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x70, 0x69, 0x6e, 0x6e, 0x65, 0x64, 0x12,
	0x3a, 0x0a, 0x19, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x5f, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x61,
	0x75, 0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x61, 0x74, 0x69, 0x76, 0x65, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x17, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x74, 0x61, 0x74, 0x69, 0x76, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x22, 0x54, 0x0a, 0x0c, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x56, 0x69, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x07, 0x63, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x73, 0x22, 0x43, 0x0a, 0x0d, 0x41, 0x75, 0x74,
	0x6f, 0x49, 0x6e, 0x63, 0x72, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f,
	0x6c, 0x75, 0x6d, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6c, 0x75,
	0x6d, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x22, 0x3d,
	0x0a, 0x06, 0x43, 0x6f, 0x6c, 0x75, 0x6d, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1f, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x0b, 0x2e, 0x71, 0x75, 0x65,
	0x72, 0x79, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0xa7, 0x02,
	0x0a, 0x0a, 0x53, 0x72, 0x76, 0x56, 0x53, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x12, 0x40, 0x0a, 0x09,
	0x6b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x53, 0x72, 0x76, 0x56, 0x53, 0x63,
	0x68, 0x65, 0x6d, 0x61, 0x2e, 0x4b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x09, 0x6b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x73, 0x12, 0x3a,
	0x0a, 0x0d, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e,
	0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x52, 0x0c, 0x72, 0x6f,
	0x75, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x4a, 0x0a, 0x13, 0x73, 0x68,
	0x61, 0x72, 0x64, 0x5f, 0x72, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x5f, 0x72, 0x75, 0x6c, 0x65,
	0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x76, 0x73, 0x63, 0x68, 0x65, 0x6d,
	0x61, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x75,
	0x6c, 0x65, 0x73, 0x52, 0x11, 0x73, 0x68, 0x61, 0x72, 0x64, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e,
	0x67, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x1a, 0x4f, 0x0a, 0x0e, 0x4b, 0x65, 0x79, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x27, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x76, 0x73, 0x63, 0x68,
	0x65, 0x6d, 0x61, 0x2e, 0x4b, 0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x44, 0x0a, 0x11, 0x53, 0x68, 0x61, 0x72, 0x64,
	0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x75, 0x6c, 0x65, 0x73, 0x12, 0x2f, 0x0a, 0x05,
	0x72, 0x75, 0x6c, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x76, 0x73,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x2e, 0x53, 0x68, 0x61, 0x72, 0x64, 0x52, 0x6f, 0x75, 0x74, 0x69,
	0x6e, 0x67, 0x52, 0x75, 0x6c, 0x65, 0x52, 0x05, 0x72, 0x75, 0x6c, 0x65, 0x73, 0x22, 0x6e, 0x0a,
	0x10, 0x53, 0x68, 0x61, 0x72, 0x64, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e, 0x67, 0x52, 0x75, 0x6c,
	0x65, 0x12, 0x23, 0x0a, 0x0d, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x6b, 0x65, 0x79, 0x73, 0x70, 0x61,
	0x63, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x66, 0x72, 0x6f, 0x6d, 0x4b, 0x65,
	0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x74, 0x6f, 0x5f, 0x6b, 0x65, 0x79,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x6f, 0x4b,
	0x65, 0x79, 0x73, 0x70, 0x61, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x68, 0x61, 0x72, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x68, 0x61, 0x72, 0x64, 0x42, 0x26, 0x5a,
	0x24, 0x76, 0x69, 0x74, 0x65, 0x73, 0x73, 0x2e, 0x69, 0x6f, 0x2f, 0x76, 0x69, 0x74, 0x65, 0x73,
	0x73, 0x2f, 0x67, 0x6f, 0x2f, 0x76, 0x74, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x76, 0x73,
	0x63, 0x68, 0x65, 0x6d, 0x61, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_vschema_proto_rawDescOnce sync.Once
	file_vschema_proto_rawDescData = file_vschema_proto_rawDesc
)

func file_vschema_proto_rawDescGZIP() []byte {
	file_vschema_proto_rawDescOnce.Do(func() {
		file_vschema_proto_rawDescData = protoimpl.X.CompressGZIP(file_vschema_proto_rawDescData)
	})
	return file_vschema_proto_rawDescData
}

var file_vschema_proto_msgTypes = make([]protoimpl.MessageInfo, 15)
var file_vschema_proto_goTypes = []interface{}{
	(*RoutingRules)(nil),      // 0: vschema.RoutingRules
	(*RoutingRule)(nil),       // 1: vschema.RoutingRule
	(*Keyspace)(nil),          // 2: vschema.Keyspace
	(*Vindex)(nil),            // 3: vschema.Vindex
	(*Table)(nil),             // 4: vschema.Table
	(*ColumnVindex)(nil),      // 5: vschema.ColumnVindex
	(*AutoIncrement)(nil),     // 6: vschema.AutoIncrement
	(*Column)(nil),            // 7: vschema.Column
	(*SrvVSchema)(nil),        // 8: vschema.SrvVSchema
	(*ShardRoutingRules)(nil), // 9: vschema.ShardRoutingRules
	(*ShardRoutingRule)(nil),  // 10: vschema.ShardRoutingRule
	nil,                       // 11: vschema.Keyspace.VindexesEntry
	nil,                       // 12: vschema.Keyspace.TablesEntry
	nil,                       // 13: vschema.Vindex.ParamsEntry
	nil,                       // 14: vschema.SrvVSchema.KeyspacesEntry
	(query.Type)(0),           // 15: query.Type
}
var file_vschema_proto_depIdxs = []int32{
	1,  // 0: vschema.RoutingRules.rules:type_name -> vschema.RoutingRule
	11, // 1: vschema.Keyspace.vindexes:type_name -> vschema.Keyspace.VindexesEntry
	12, // 2: vschema.Keyspace.tables:type_name -> vschema.Keyspace.TablesEntry
	13, // 3: vschema.Vindex.params:type_name -> vschema.Vindex.ParamsEntry
	5,  // 4: vschema.Table.column_vindexes:type_name -> vschema.ColumnVindex
	6,  // 5: vschema.Table.auto_increment:type_name -> vschema.AutoIncrement
	7,  // 6: vschema.Table.columns:type_name -> vschema.Column
	15, // 7: vschema.Column.type:type_name -> query.Type
	14, // 8: vschema.SrvVSchema.keyspaces:type_name -> vschema.SrvVSchema.KeyspacesEntry
	0,  // 9: vschema.SrvVSchema.routing_rules:type_name -> vschema.RoutingRules
	9,  // 10: vschema.SrvVSchema.shard_routing_rules:type_name -> vschema.ShardRoutingRules
	10, // 11: vschema.ShardRoutingRules.rules:type_name -> vschema.ShardRoutingRule
	3,  // 12: vschema.Keyspace.VindexesEntry.value:type_name -> vschema.Vindex
	4,  // 13: vschema.Keyspace.TablesEntry.value:type_name -> vschema.Table
	2,  // 14: vschema.SrvVSchema.KeyspacesEntry.value:type_name -> vschema.Keyspace
	15, // [15:15] is the sub-list for method output_type
	15, // [15:15] is the sub-list for method input_type
	15, // [15:15] is the sub-list for extension type_name
	15, // [15:15] is the sub-list for extension extendee
	0,  // [0:15] is the sub-list for field type_name
}

func init() { file_vschema_proto_init() }
func file_vschema_proto_init() {
	if File_vschema_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_vschema_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoutingRules); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vschema_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RoutingRule); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vschema_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Keyspace); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vschema_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vindex); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vschema_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Table); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vschema_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ColumnVindex); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vschema_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AutoIncrement); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vschema_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Column); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vschema_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SrvVSchema); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vschema_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShardRoutingRules); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_vschema_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ShardRoutingRule); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_vschema_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   15,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_vschema_proto_goTypes,
		DependencyIndexes: file_vschema_proto_depIdxs,
		MessageInfos:      file_vschema_proto_msgTypes,
	}.Build()
	File_vschema_proto = out.File
	file_vschema_proto_rawDesc = nil
	file_vschema_proto_goTypes = nil
	file_vschema_proto_depIdxs = nil
}
