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

package vdiff

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"google.golang.org/protobuf/encoding/prototext"

	"vitess.io/vitess/go/vt/binlog/binlogplayer"

	"vitess.io/vitess/go/vt/schema"

	"vitess.io/vitess/go/vt/key"
	"vitess.io/vitess/go/vt/log"
	binlogdatapb "vitess.io/vitess/go/vt/proto/binlogdata"
	querypb "vitess.io/vitess/go/vt/proto/query"
	tabletmanagerdatapb "vitess.io/vitess/go/vt/proto/tabletmanagerdata"
	vtrpcpb "vitess.io/vitess/go/vt/proto/vtrpc"
	"vitess.io/vitess/go/vt/sqlparser"
	"vitess.io/vitess/go/vt/vtctl/schematools"
	"vitess.io/vitess/go/vt/vterrors"
	"vitess.io/vitess/go/vt/vttablet/tabletmanager/vreplication"
)

// workflowDiffer has metadata and state for the vdiff of a single workflow on this tablet
// only one vdiff can be running for a workflow at any time.
type workflowDiffer struct {
	ct *controller

	tableDiffers map[string]*tableDiffer // key is table name
	opts         *tabletmanagerdatapb.VDiffOptions
}

func newWorkflowDiffer(ct *controller, opts *tabletmanagerdatapb.VDiffOptions) (*workflowDiffer, error) {
	wd := &workflowDiffer{
		ct:           ct,
		opts:         opts,
		tableDiffers: make(map[string]*tableDiffer, 1),
	}
	return wd, nil
}

// If the only difference is the order in which the rows were returned
// by MySQL on each side then we'll have the same number of extras on
// both sides. If that's the case, then let's see if the extra rows on
// both sides are actually different.
func (wd *workflowDiffer) reconcileExtraRows(dr *DiffReport, maxExtraRowsToCompare int64) {
	if (dr.ExtraRowsSource == dr.ExtraRowsTarget) && (dr.ExtraRowsSource <= maxExtraRowsToCompare) {
		for i := range dr.ExtraRowsSourceDiffs {
			foundMatch := false
			for j := range dr.ExtraRowsTargetDiffs {
				if reflect.DeepEqual(dr.ExtraRowsSourceDiffs[i], dr.ExtraRowsTargetDiffs[j]) {
					dr.ExtraRowsSourceDiffs = append(dr.ExtraRowsSourceDiffs[:i], dr.ExtraRowsSourceDiffs[i+1:]...)
					dr.ExtraRowsSource--
					dr.ExtraRowsTargetDiffs = append(dr.ExtraRowsTargetDiffs[:j], dr.ExtraRowsTargetDiffs[j+1:]...)
					dr.ExtraRowsTarget--
					dr.ProcessedRows--
					dr.MatchingRows++
					foundMatch = true
					break
				}
			}
			// If we didn't find a match then the tables are in fact different and we can short circuit the second pass
			if !foundMatch {
				break
			}
		}
	}
	// We can now trim the extra rows diffs on both sides to the maxVDiffReportSampleRows value
	if len(dr.ExtraRowsSourceDiffs) > maxVDiffReportSampleRows {
		dr.ExtraRowsSourceDiffs = dr.ExtraRowsSourceDiffs[:maxVDiffReportSampleRows-1]
	}
	if len(dr.ExtraRowsTargetDiffs) > maxVDiffReportSampleRows {
		dr.ExtraRowsTargetDiffs = dr.ExtraRowsTargetDiffs[:maxVDiffReportSampleRows-1]
	}
}

func (wd *workflowDiffer) diffTable(ctx context.Context, dbClient binlogplayer.DBClient, td *tableDiffer) error {
	defer func() {
		if td.shardStreamsCancel != nil {
			td.shardStreamsCancel()
		}
		// Wait for all the shard streams to finish before returning.
		td.wgShardStreamers.Wait()
	}()

	select {
	case <-ctx.Done():
		return vterrors.Errorf(vtrpcpb.Code_CANCELED, "context has expired")
	case <-wd.ct.done:
		return vterrors.Errorf(vtrpcpb.Code_CANCELED, "vdiff was stopped")
	default:
	}

	log.Infof("Starting differ on table %s for vdiff %s", td.table.Name, wd.ct.uuid)
	if err := td.updateTableState(ctx, dbClient, StartedState); err != nil {
		return err
	}
	if err := td.initialize(ctx); err != nil {
		return err
	}
	log.Infof("Table initialization done on table %s for vdiff %s", td.table.Name, wd.ct.uuid)
	dr, err := td.diff(ctx, wd.opts.CoreOptions.MaxRows, wd.opts.ReportOptions.DebugQuery, wd.opts.ReportOptions.OnlyPks, wd.opts.CoreOptions.MaxExtraRowsToCompare)
	if err != nil {
		log.Errorf("Encountered an error diffing table %s for vdiff %s: %v", td.table.Name, wd.ct.uuid, err)
		return err
	}
	log.Infof("Table diff done on table %s for vdiff %s with report: %+v", td.table.Name, wd.ct.uuid, dr)
	if dr.ExtraRowsSource > 0 || dr.ExtraRowsTarget > 0 {
		wd.reconcileExtraRows(dr, wd.opts.CoreOptions.MaxExtraRowsToCompare)
	}

	if dr.MismatchedRows > 0 || dr.ExtraRowsTarget > 0 || dr.ExtraRowsSource > 0 {
		if err := updateTableMismatch(dbClient, wd.ct.id, td.table.Name); err != nil {
			return err
		}
	}

	log.Infof("Completed reconciliation on table %s for vdiff %s with updated report: %+v", td.table.Name, wd.ct.uuid, dr)
	if err := td.updateTableStateAndReport(ctx, dbClient, CompletedState, dr); err != nil {
		return err
	}
	return nil
}

func (wd *workflowDiffer) diff(ctx context.Context) error {
	dbClient := wd.ct.dbClientFactory()
	if err := dbClient.Connect(); err != nil {
		return err
	}
	defer dbClient.Close()

	select {
	case <-ctx.Done():
		return vterrors.Errorf(vtrpcpb.Code_CANCELED, "context has expired")
	case <-wd.ct.done:
		return vterrors.Errorf(vtrpcpb.Code_CANCELED, "vdiff was stopped")
	default:
	}

	filter := wd.ct.filter
	req := &tabletmanagerdatapb.GetSchemaRequest{}
	schm, err := schematools.GetSchema(ctx, wd.ct.ts, wd.ct.tmc, wd.ct.vde.thisTablet.Alias, req)
	if err != nil {
		return vterrors.Wrap(err, "GetSchema")
	}
	if err = wd.buildPlan(dbClient, filter, schm); err != nil {
		return vterrors.Wrap(err, "buildPlan")
	}
	if err := wd.initVDiffTables(dbClient); err != nil {
		return err
	}
	for _, td := range wd.tableDiffers {
		select {
		case <-ctx.Done():
			return vterrors.Errorf(vtrpcpb.Code_CANCELED, "context has expired")
		case <-wd.ct.done:
			return vterrors.Errorf(vtrpcpb.Code_CANCELED, "vdiff was stopped")
		default:
		}
		query := fmt.Sprintf(sqlGetVDiffTable, wd.ct.id, encodeString(td.table.Name))
		qr, err := dbClient.ExecuteFetch(query, 1)
		if err != nil {
			return err
		}
		if len(qr.Rows) == 0 {
			return fmt.Errorf("no vdiff table found for %s on tablet %v",
				td.table.Name, wd.ct.vde.thisTablet.Alias)
		}

		log.Infof("Starting diff of table %s for vdiff %s", td.table.Name, wd.ct.uuid)
		if err := wd.diffTable(ctx, dbClient, td); err != nil {
			if err := td.updateTableState(ctx, dbClient, ErrorState); err != nil {
				return err
			}
			insertVDiffLog(ctx, dbClient, wd.ct.id, fmt.Sprintf("Table %s Error: %s", td.table.Name, err))
			return err
		}
		if err := td.updateTableState(ctx, dbClient, CompletedState); err != nil {
			return err
		}
		log.Infof("Completed diff of table %s for vdiff %s", td.table.Name, wd.ct.uuid)
	}
	if err := wd.markIfCompleted(ctx, dbClient); err != nil {
		return err
	}
	return nil
}

func (wd *workflowDiffer) markIfCompleted(ctx context.Context, dbClient binlogplayer.DBClient) error {
	query := fmt.Sprintf(sqlGetIncompleteTables, wd.ct.id)
	qr, err := dbClient.ExecuteFetch(query, -1)
	if err != nil {
		return err
	}

	// Double check to be sure all of the individual table diffs completed without error
	// before marking the vdiff as completed.
	if len(qr.Rows) == 0 {
		if err := wd.ct.updateState(dbClient, CompletedState, nil); err != nil {
			return err
		}
	}
	return nil
}

func (wd *workflowDiffer) buildPlan(dbClient binlogplayer.DBClient, filter *binlogdatapb.Filter, schm *tabletmanagerdatapb.SchemaDefinition) error {
	var specifiedTables []string
	optTables := strings.TrimSpace(wd.opts.CoreOptions.Tables)
	if optTables != "" {
		specifiedTables = strings.Split(optTables, ",")
	}

	for _, table := range schm.TableDefinitions {
		// if user specified tables explicitly only use those, otherwise diff all tables in workflow
		if len(specifiedTables) != 0 && !stringListContains(specifiedTables, table.Name) {
			continue
		}
		if schema.IsInternalOperationTableName(table.Name) {
			continue
		}
		rule, err := vreplication.MatchTable(table.Name, filter)
		if err != nil {
			return err
		}
		if rule == nil || rule.Filter == "exclude" {
			continue
		}
		sourceQuery := rule.Filter
		switch {
		case rule.Filter == "":
			buf := sqlparser.NewTrackedBuffer(nil)
			buf.Myprintf("select * from %v", sqlparser.NewIdentifierCS(table.Name))
			sourceQuery = buf.String()
		case key.IsValidKeyRange(rule.Filter):
			buf := sqlparser.NewTrackedBuffer(nil)
			buf.Myprintf("select * from %v where in_keyrange(%v)", sqlparser.NewIdentifierCS(table.Name), sqlparser.NewStrLiteral(rule.Filter))
			sourceQuery = buf.String()
		}

		td := newTableDiffer(wd, table, sourceQuery)
		lastpkpb, err := wd.getTableLastPK(dbClient, table.Name)
		if err != nil {
			return err
		}
		td.lastPK = lastpkpb
		wd.tableDiffers[table.Name] = td
		if _, err := td.buildTablePlan(); err != nil {
			return err
		}
	}
	if len(wd.tableDiffers) == 0 {
		return fmt.Errorf("no tables found to diff, %s:%s, on tablet %v",
			optTables, specifiedTables, wd.ct.vde.thisTablet.Alias)
	}
	return nil
}

// getTableLastPK gets the lastPK protobuf message for a given vdiff table.
func (wd *workflowDiffer) getTableLastPK(dbClient binlogplayer.DBClient, tableName string) (*querypb.QueryResult, error) {
	query := fmt.Sprintf(sqlGetVDiffTable, wd.ct.id, encodeString(tableName))
	qr, err := dbClient.ExecuteFetch(query, 1)
	if err != nil {
		return nil, err
	}
	if len(qr.Rows) == 1 {
		var lastpk []byte
		if lastpk, err = qr.Named().Row().ToBytes("lastpk"); err != nil {
			return nil, err
		}
		if len(lastpk) != 0 {
			var lastpkpb querypb.QueryResult
			if err := prototext.Unmarshal(lastpk, &lastpkpb); err != nil {
				return nil, err
			}
			return &lastpkpb, nil
		}
	}
	return nil, nil
}

func (wd *workflowDiffer) initVDiffTables(dbClient binlogplayer.DBClient) error {
	tableIn := strings.Builder{}
	n := 0
	for tableName := range wd.tableDiffers {
		tableIn.WriteString(encodeString(tableName))
		if n++; n < len(wd.tableDiffers) {
			tableIn.WriteByte(',')
		}
	}
	query := fmt.Sprintf(sqlGetAllTableRows, encodeString(wd.ct.vde.dbName), tableIn.String())
	isqr, err := dbClient.ExecuteFetch(query, -1)
	if err != nil {
		return err
	}
	for _, row := range isqr.Named().Rows {
		tableName, _ := row.ToString("table_name")
		tableRows, _ := row.ToInt64("table_rows")

		query := fmt.Sprintf(sqlGetVDiffTable, wd.ct.id, encodeString(tableName))
		qr, err := dbClient.ExecuteFetch(query, -1)
		if err != nil {
			return err
		}
		if len(qr.Rows) == 0 {
			query = fmt.Sprintf(sqlNewVDiffTable, wd.ct.id, encodeString(tableName), tableRows)
		} else if len(qr.Rows) == 1 {
			query = fmt.Sprintf(sqlUpdateTableRows, tableRows, wd.ct.id, encodeString(tableName))
		} else {
			return fmt.Errorf("invalid state found for vdiff table %s for vdiff_id %d on tablet %s",
				tableName, wd.ct.id, wd.ct.vde.thisTablet.Alias)
		}
		if _, err := dbClient.ExecuteFetch(query, 1); err != nil {
			return err
		}
	}
	return nil
}
