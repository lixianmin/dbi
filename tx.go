package dbi

import (
	"context"
	"database/sql"
)

/********************************************************************
created:    2020-05-31
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

type Tx struct {
	TX *sql.Tx
	db *DB
}

func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var ctx1 = newContext(ctx, TxQueryContext, query)

	tx.db.preExecuteHandler(ctx1)
	var rows, err = tx.TX.QueryContext(ctx1, query, args...)
	ctx1.err = err
	tx.db.postExecuteHandler(ctx1)
	return rows, err
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var ctx1 = newContext(ctx, TxExecContext, query)

	tx.db.preExecuteHandler(ctx1)
	var result, err = tx.TX.ExecContext(ctx1, query, args...)
	ctx1.err = err
	tx.db.postExecuteHandler(ctx1)
	return result, err
}

func (tx *Tx) Commit() error {
	var ctx1 = newContext(context.Background(), TxCommit, "TxCommit")

	tx.db.preExecuteHandler(ctx1)
	var err = tx.TX.Commit()
	ctx1.err = err
	tx.db.postExecuteHandler(ctx1)
	return err
}

func (tx *Tx) Rollback() error {
	var ctx1 = newContext(context.Background(), TxRollback, "TxRollback")

	tx.db.preExecuteHandler(ctx1)
	var err = tx.TX.Rollback()
	ctx1.err = err
	tx.db.postExecuteHandler(ctx1)
	return err
}
