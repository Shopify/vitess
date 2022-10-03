/*
   Copyright 2014 Outbrain Inc.

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

package inst

import (
	"github.com/rcrowley/go-metrics"

	"vitess.io/vitess/go/vt/log"

	"github.com/openark/golib/sqlutils"

	"vitess.io/vitess/go/vt/vtorc/config"
	"vitess.io/vitess/go/vt/vtorc/db"
)

var writeResolvedHostnameCounter = metrics.NewCounter()
var writeUnresolvedHostnameCounter = metrics.NewCounter()
var readResolvedHostnameCounter = metrics.NewCounter()
var readUnresolvedHostnameCounter = metrics.NewCounter()
var readAllResolvedHostnamesCounter = metrics.NewCounter()

func init() {
	_ = metrics.Register("resolve.write_resolved", writeResolvedHostnameCounter)
	_ = metrics.Register("resolve.write_unresolved", writeUnresolvedHostnameCounter)
	_ = metrics.Register("resolve.read_resolved", readResolvedHostnameCounter)
	_ = metrics.Register("resolve.read_unresolved", readUnresolvedHostnameCounter)
	_ = metrics.Register("resolve.read_resolved_all", readAllResolvedHostnamesCounter)
}

// WriteResolvedHostname stores a hostname and the resolved hostname to backend database
func WriteResolvedHostname(hostname string, resolvedHostname string) error {
	writeFunc := func() error {
		_, err := db.ExecVTOrc(`
			insert into
					hostname_resolve (hostname, resolved_hostname, resolved_timestamp)
				values
					(?, ?, NOW())
				on duplicate key update
					resolved_hostname = VALUES(resolved_hostname),
					resolved_timestamp = VALUES(resolved_timestamp)
			`,
			hostname,
			resolvedHostname)
		if err != nil {
			log.Error(err)
			return err
		}
		if hostname != resolvedHostname {
			// history is only interesting when there's actually something to resolve...
			_, _ = db.ExecVTOrc(`
			insert into
					hostname_resolve_history (hostname, resolved_hostname, resolved_timestamp)
				values
					(?, ?, NOW())
				on duplicate key update
					hostname=values(hostname),
					resolved_timestamp=values(resolved_timestamp)
			`,
				hostname,
				resolvedHostname)
		}
		writeResolvedHostnameCounter.Inc(1)
		return nil
	}
	return ExecDBWriteFunc(writeFunc)
}

// ReadResolvedHostname returns the resolved hostname given a hostname, or empty if not exists
func ReadResolvedHostname(hostname string) (string, error) {
	var resolvedHostname string

	query := `
		select
			resolved_hostname
		from
			hostname_resolve
		where
			hostname = ?
		`

	err := db.QueryVTOrc(query, sqlutils.Args(hostname), func(m sqlutils.RowMap) error {
		resolvedHostname = m.GetString("resolved_hostname")
		return nil
	})
	readResolvedHostnameCounter.Inc(1)

	if err != nil {
		log.Error(err)
	}
	return resolvedHostname, err
}

func ReadAllHostnameResolves() ([]HostnameResolve, error) {
	res := []HostnameResolve{}
	query := `
		select
			hostname,
			resolved_hostname
		from
			hostname_resolve
		`
	err := db.QueryVTOrcRowsMap(query, func(m sqlutils.RowMap) error {
		hostnameResolve := HostnameResolve{hostname: m.GetString("hostname"), resolvedHostname: m.GetString("resolved_hostname")}

		res = append(res, hostnameResolve)
		return nil
	})
	readAllResolvedHostnamesCounter.Inc(1)

	if err != nil {
		log.Error(err)
	}
	return res, err
}

// ExpireHostnameUnresolve expires hostname_unresolve entries that haven't been updated recently.
func ExpireHostnameUnresolve() error {
	writeFunc := func() error {
		_, err := db.ExecVTOrc(`
      	delete from hostname_unresolve
				where last_registered < NOW() - INTERVAL ? MINUTE
				`, config.ExpiryHostnameResolvesMinutes,
		)
		if err != nil {
			log.Error(err)
		}
		return err
	}
	return ExecDBWriteFunc(writeFunc)
}

// ForgetExpiredHostnameResolves
func ForgetExpiredHostnameResolves() error {
	_, err := db.ExecVTOrc(`
			delete
				from hostname_resolve
			where
				resolved_timestamp < NOW() - interval ? minute`,
		2*config.ExpiryHostnameResolvesMinutes,
	)
	return err
}

// DeleteInvalidHostnameResolves removes invalid resolves. At this time these are:
// - infinite loop resolves (A->B and B->A), remove earlier mapping
func DeleteInvalidHostnameResolves() error {
	var invalidHostnames []string

	query := `
		select
		    early.hostname
		  from
		    hostname_resolve as latest
		    join hostname_resolve early on (latest.resolved_hostname = early.hostname and latest.hostname = early.resolved_hostname)
		  where
		    latest.hostname != latest.resolved_hostname
		    and latest.resolved_timestamp > early.resolved_timestamp
	   	`

	err := db.QueryVTOrcRowsMap(query, func(m sqlutils.RowMap) error {
		invalidHostnames = append(invalidHostnames, m.GetString("hostname"))
		return nil
	})
	if err != nil {
		return err
	}

	for _, invalidHostname := range invalidHostnames {
		_, err = db.ExecVTOrc(`
			delete
				from hostname_resolve
			where
				hostname = ?`,
			invalidHostname,
		)
		if err != nil {
			log.Error(err)
		}
	}
	return err
}

// writeHostnameIPs stroes an ipv4 and ipv6 associated witha hostname, if available
func writeHostnameIPs(hostname string, ipv4String string, ipv6String string) error {
	writeFunc := func() error {
		_, err := db.ExecVTOrc(`
			insert into
					hostname_ips (hostname, ipv4, ipv6, last_updated)
				values
					(?, ?, ?, NOW())
				on duplicate key update
					ipv4 = VALUES(ipv4),
					ipv6 = VALUES(ipv6),
					last_updated = VALUES(last_updated)
			`,
			hostname,
			ipv4String,
			ipv6String,
		)
		if err != nil {
			log.Error(err)
		}
		return err
	}
	return ExecDBWriteFunc(writeFunc)
}
