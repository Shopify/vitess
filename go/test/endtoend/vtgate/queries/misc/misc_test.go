/*
Copyright 2022 The Vitess Authors.

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

package misc

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"vitess.io/vitess/go/test/endtoend/cluster"
	"vitess.io/vitess/go/test/endtoend/utils"
)

func start(t *testing.T) (utils.MySQLCompare, func()) {
	mcmp, err := utils.NewMySQLCompare(t, vtParams, mysqlParams)
	require.NoError(t, err)

	deleteAll := func() {
		tables := []string{"t1", "uks.unsharded"}
		for _, table := range tables {
			_, _ = mcmp.ExecAndIgnore("delete from " + table)
		}
	}

	deleteAll()

	return mcmp, func() {
		deleteAll()
		mcmp.Close()
		cluster.PanicHandler(t)
	}
}

func TestBitVals(t *testing.T) {
	mcmp, closer := start(t)
	defer closer()

	mcmp.Exec("insert into t1(id1, id2) values (0,0)")

	mcmp.AssertMatches(`select b'1001', 0x9, B'010011011010'`, `[[VARBINARY("\t") VARBINARY("\t") VARBINARY("\x04\xda")]]`)
	mcmp.AssertMatches(`select b'1001', 0x9, B'010011011010' from t1`, `[[VARBINARY("\t") VARBINARY("\t") VARBINARY("\x04\xda")]]`)
	vtgateVersion, err := cluster.GetMajorVersion("vtgate")
	require.NoError(t, err)
	if vtgateVersion >= 19 {
		mcmp.AssertMatchesNoCompare(`select 1 + b'1001', 2 + 0x9, 3 + B'010011011010'`, `[[INT64(10) UINT64(11) INT64(1245)]]`, `[[INT64(10) UINT64(11) INT64(1245)]]`)
		mcmp.AssertMatchesNoCompare(`select 1 + b'1001', 2 + 0x9, 3 + B'010011011010' from t1`, `[[INT64(10) UINT64(11) INT64(1245)]]`, `[[INT64(10) UINT64(11) INT64(1245)]]`)
	} else {
		mcmp.AssertMatchesNoCompare(`select 1 + b'1001', 2 + 0x9, 3 + B'010011011010'`, `[[INT64(10) UINT64(11) INT64(1245)]]`, `[[UINT64(10) UINT64(11) UINT64(1245)]]`)
		mcmp.AssertMatchesNoCompare(`select 1 + b'1001', 2 + 0x9, 3 + B'010011011010' from t1`, `[[INT64(10) UINT64(11) INT64(1245)]]`, `[[UINT64(10) UINT64(11) UINT64(1245)]]`)
	}
}

func TestHexVals(t *testing.T) {
	mcmp, closer := start(t)
	defer closer()

	mcmp.Exec("insert into t1(id1, id2) values (0,0)")

	mcmp.AssertMatches(`select x'09', 0x9`, `[[VARBINARY("\t") VARBINARY("\t")]]`)
	mcmp.AssertMatches(`select X'09', 0x9 from t1`, `[[VARBINARY("\t") VARBINARY("\t")]]`)
	mcmp.AssertMatches(`select 1 + x'09', 2 + 0x9`, `[[UINT64(10) UINT64(11)]]`)
	mcmp.AssertMatches(`select 1 + X'09', 2 + 0x9 from t1`, `[[UINT64(10) UINT64(11)]]`)
}

func TestDateTimeTimestampVals(t *testing.T) {
	mcmp, closer := start(t)
	defer closer()

	mcmp.AssertMatches(`select date'2022-08-03'`, `[[DATE("2022-08-03")]]`)
	mcmp.AssertMatches(`select time'12:34:56'`, `[[TIME("12:34:56")]]`)
	mcmp.AssertMatches(`select timestamp'2012-12-31 11:30:45'`, `[[DATETIME("2012-12-31 11:30:45")]]`)
}

func TestInvalidDateTimeTimestampVals(t *testing.T) {
	mcmp, closer := start(t)
	defer closer()

	_, err := mcmp.ExecAllowAndCompareError(`select date'2022'`)
	require.Error(t, err)
	_, err = mcmp.ExecAllowAndCompareError(`select time'12:34:56:78'`)
	require.Error(t, err)
	_, err = mcmp.ExecAllowAndCompareError(`select timestamp'2022'`)
	require.Error(t, err)
}

func TestJoinWithThreeTables(t *testing.T) {
	version, err := cluster.GetMajorVersion("vtgate")
	require.NoError(t, err)
	if version != 20 {
		t.Skip("cannot run upgrade/downgrade test")
	}

	mcmp, closer := start(t)
	defer closer()

	mcmp.Exec("insert into t1(id1, id2) values (0,0), (1,1), (2,2)")
	mcmp.Exec("insert into tbl(id, unq_col, nonunq_col) values (0,0,0), (1,1,1), (2,2,1)")
	mcmp.Exec("select 42 from t1 u1, t1 u2, tbl u3 where u1.id1 = u2.id1 and u1.id1 = u3.id and (u1.id2 or u2.id2 or u3.unq_col)")
}

// TestIntervalWithMathFunctions tests that the Interval keyword can be used with math functions.
func TestIntervalWithMathFunctions(t *testing.T) {
	mcmp, closer := start(t)
	defer closer()

	// Set the time zone explicitly to UTC, otherwise the output of FROM_UNIXTIME is going to be dependent
	// on the time zone of the system.
	mcmp.Exec("SET time_zone = '+00:00'")
	mcmp.AssertMatches("select '2020-01-01' + interval month(date_sub(FROM_UNIXTIME(1234), interval 1 month))-1 month", `[[CHAR("2020-12-01")]]`)
	mcmp.AssertMatches("select date_add(MIN(FROM_UNIXTIME(1673444922)),interval -DAYOFWEEK(MIN(FROM_UNIXTIME(1673444922)))+1 DAY)", `[[DATETIME("2023-01-08 13:48:42")]]`)
}

// TestCast tests the queries that contain the cast function.
func TestCast(t *testing.T) {
	mcmp, closer := start(t)
	defer closer()

	mcmp.AssertMatches("select cast('2023-01-07 12:34:56' as date) limit 1", `[[DATE("2023-01-07")]]`)
	mcmp.AssertMatches("select cast('2023-01-07 12:34:56' as date)", `[[DATE("2023-01-07")]]`)
	mcmp.AssertMatches("select cast('3.2' as float)", `[[FLOAT32(3.2)]]`)
	mcmp.AssertMatches("select cast('3.2' as double)", `[[FLOAT64(3.2)]]`)
	mcmp.AssertMatches("select cast('3.2' as unsigned)", `[[UINT64(3)]]`)
}

func TestOuterJoinWithPredicate(t *testing.T) {
	mcmp, closer := start(t)
	defer closer()

	// This test uses a predicate on the outer side.
	// These can't be pushed down to MySQL and have
	// to be evaluated on the vtgate, so we are checking
	// that evalengine handles the predicate correctly

	mcmp.Exec("insert into t1(id1, id2) values (0,0), (1,10), (2,20), (3,30), (4,40)")

	mcmp.AssertMatchesNoOrder("select A.id1, B.id2 from t1 as A left join t1 as B on A.id1*10 = B.id2 WHERE B.id2 BETWEEN 20 AND 30",
		`[[INT64(2) INT64(20)] [INT64(3) INT64(30)]]`)
	mcmp.AssertMatchesNoOrder("select A.id1, B.id2 from t1 as A left join t1 as B on A.id1*10 = B.id2 WHERE B.id2 NOT BETWEEN 20 AND 30",
		`[[INT64(0) INT64(0)] [INT64(1) INT64(10)] [INT64(4) INT64(40)]]`)
}

// This test ensures that we support PREPARE statement with 65530 parameters.
// It opens a MySQL connection using the go-mysql driver and execute a select query
// it then checks the result contains the proper rows and that it's not failing.
func TestHighNumberOfParams(t *testing.T) {
	mcmp, closer := start(t)
	defer closer()

	mcmp.Exec("insert into t1(id1) values (0), (1), (2), (3), (4)")

	paramCount := 65530

	// create the value and argument slices used to build the prepare stmt
	var vals []any
	var params []string
	for i := 0; i < paramCount; i++ {
		vals = append(vals, strconv.Itoa(i))
		params = append(params, "?")
	}

	// connect to the vitess cluster
	db, err := sql.Open("mysql", fmt.Sprintf("@tcp(%s:%v)/%s", vtParams.Host, vtParams.Port, vtParams.DbName))
	require.NoError(t, err)
	defer db.Close()

	// run the query
	r, err := db.Query(fmt.Sprintf("SELECT id1 FROM t1 WHERE id1 in (%s) ORDER BY id1 ASC", strings.Join(params, ", ")), vals...)
	require.NoError(t, err)
	defer r.Close()

	// check the results we got, we should get 5 rows with each: 0, 1, 2, 3, 4
	// count is the row number we are currently visiting, also correspond to the
	// column value we expect.
	count := 0
	for r.Next() {
		j := -1
		err := r.Scan(&j)
		require.NoError(t, err)
		require.Equal(t, j, count)
		count++
	}
	require.Equal(t, 5, count)
}

func TestPrepareStatements(t *testing.T) {
	mcmp, closer := start(t)
	defer closer()

	mcmp.Exec("insert into t1(id1, id2) values (0,0), (1,0), (2,0)")

	// prepare query with equal sharding key
	mcmp.Exec(`prepare prep_pk from 'select count(*) from t1 where id1 = ?'`)
	mcmp.AssertMatches(`execute prep_pk using @id1`, `[[INT64(0)]]`)
	mcmp.Exec(`set @id1 = 1`)
	mcmp.AssertMatches(`execute prep_pk using @id1`, `[[INT64(1)]]`)

	// prepare query with equal non sharding key
	mcmp.Exec(`prepare prep_non_pk from 'select id1, id2 from t1 where id2 = ?'`)
	mcmp.Exec(`set @id2 = 0`)
	mcmp.AssertMatches(`execute prep_non_pk using @id1`, `[]`)
	mcmp.AssertMatchesNoOrder(`execute prep_non_pk using @id2`, `[[INT64(0) INT64(0)] [INT64(1) INT64(0)] [INT64(2) INT64(0)]]`)

	// prepare query with in on sharding key
	mcmp.Exec(`prepare prep_in_pk from 'select id1, id2 from t1 where id1 in (?, ?)'`)
	mcmp.AssertMatches(`execute prep_in_pk using @id1, @id1`, `[[INT64(1) INT64(0)]]`)
	mcmp.AssertMatchesNoOrder(`execute prep_in_pk using @id1, @id2`, `[[INT64(0) INT64(0)] [INT64(1) INT64(0)]]`)

	// Fail by providing wrong number of arguments
	_, err := mcmp.ExecAllowAndCompareError(`execute prep_in_pk using @id1, @id1, @id`)
	incorrectCount := "VT03025: Incorrect arguments to EXECUTE"
	assert.ErrorContains(t, err, incorrectCount)
	_, err = mcmp.ExecAllowAndCompareError(`execute prep_in_pk using @id1`)
	assert.ErrorContains(t, err, incorrectCount)
	_, err = mcmp.ExecAllowAndCompareError(`execute prep_in_pk`)
	assert.ErrorContains(t, err, incorrectCount)

	mcmp.Exec(`prepare prep_art from 'select 1+?, 10/?'`)
	mcmp.Exec(`set @x1 = 1, @x2 = 2.0, @x3 = "v", @x4 = 9999999999999999999999999999`)

	// We are not matching types and precision with mysql at the moment, so not comparing with `mcmp`
	// This is because of the difference in how MySQL executes a raw query with literal values and
	// the PREPARE/EXEC way that is missing type info at the PREPARE stage
	utils.AssertMatches(t, mcmp.VtConn, `execute prep_art using @x1, @x1`, `[[INT64(2) DECIMAL(10.0000)]]`)
	utils.AssertMatches(t, mcmp.VtConn, `execute prep_art using @x2, @x2`, `[[DECIMAL(3.0) DECIMAL(5.0000)]]`)
	utils.AssertMatches(t, mcmp.VtConn, `execute prep_art using @x3, @x3`, `[[FLOAT64(1) NULL]]`)
	utils.AssertMatches(t, mcmp.VtConn, `execute prep_art using @x4, @x4`, `[[DECIMAL(10000000000000000000000000000) DECIMAL(0.0000)]]`)

	mcmp.Exec(`select 1+1, 10/1 from t1 limit 1`)
	mcmp.Exec(`select 1+2.0, 10/2.0 from t1 limit 1`)
	mcmp.Exec(`select 1+'v', 10/'v' from t1 limit 1`)
	mcmp.Exec(`select 1+9999999999999999999999999999, 10/9999999999999999999999999999 from t1 limit 1`)

	mcmp.Exec("deallocate prepare prep_art")
	_, err = mcmp.ExecAllowAndCompareError(`execute prep_art using @id1, @id1`)
	assert.ErrorContains(t, err, "VT09011: Unknown prepared statement handler (prep_art) given to EXECUTE")

	_, err = mcmp.ExecAllowAndCompareError("deallocate prepare prep_art")
	assert.ErrorContains(t, err, "VT09011: Unknown prepared statement handler (prep_art) given to DEALLOCATE PREPARE")
}

// TestBuggyOuterJoin validates inconsistencies around outer joins, adding these tests to stop regressions.
func TestBuggyOuterJoin(t *testing.T) {
	mcmp, closer := start(t)
	defer closer()

	mcmp.Exec("insert into t1(id1, id2) values (1,2), (42,5), (5, 42)")
	mcmp.Exec("select t1.id1, t2.id1 from t1 left join t1 as t2 on t2.id1 = t2.id2")
}

func TestLeftJoinUsingUnsharded(t *testing.T) {
	mcmp, closer := start(t)
	defer closer()

	utils.Exec(t, mcmp.VtConn, "insert into uks.unsharded(id1) values (1),(2),(3),(4),(5)")
	utils.Exec(t, mcmp.VtConn, "select * from uks.unsharded as A left join uks.unsharded as B using(id1)")
}

// TestAnalyze executes different analyze statement and validates that they run successfully.
func TestAnalyze(t *testing.T) {
	utils.SkipIfBinaryIsBelowVersion(t, 18, "vtgate")
	utils.SkipIfBinaryIsBelowVersion(t, 18, "vttablet")
	mcmp, closer := start(t)
	defer closer()

	for _, workload := range []string{"olap", "oltp"} {
		mcmp.Run(workload, func(mcmp *utils.MySQLCompare) {
			utils.Exec(t, mcmp.VtConn, fmt.Sprintf("set workload = %s", workload))
			utils.Exec(t, mcmp.VtConn, "analyze table t1")
			utils.Exec(t, mcmp.VtConn, "analyze table uks.unsharded")
			utils.Exec(t, mcmp.VtConn, "analyze table mysql.user")
		})
	}
}
