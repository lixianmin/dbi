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

func (tx *Tx) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return tx.QueryContext(context.Background(), query, args...)
}

func (tx *Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	var ctx1 = newContext(ctx, TxQuery, query)

	tx.db.preExecuteHandler(ctx1)
	var rows, err = tx.TX.QueryContext(ctx1, query, args...)
	ctx1.err = err
	tx.db.postExecuteHandler(ctx1)
	return rows, err
}

func (tx *Tx) Exec(query string, args ...interface{}) (sql.Result, error) {
	return tx.ExecContext(context.Background(), query, args...)
}

func (tx *Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	var ctx1 = newContext(ctx, TxExec, query)

	tx.db.preExecuteHandler(ctx1)
	var result, err = tx.TX.ExecContext(ctx1, query, args...)
	ctx1.err = err
	tx.db.postExecuteHandler(ctx1)
	return result, err
}

func (tx *Tx) Commit() error {
	return tx.CommitContext(context.Background())
}

func (tx *Tx) CommitContext(ctx context.Context) error {
	var ctx1 = newContext(ctx, TxCommit, "TxCommit")

	tx.db.preExecuteHandler(ctx1)
	var err = tx.TX.Commit()
	ctx1.err = err
	tx.db.postExecuteHandler(ctx1)
	return err
}

func (tx *Tx) Rollback() error {
	return tx.RollbackContext(context.Background())
}

func (tx *Tx) RollbackContext(ctx context.Context) error {
	var ctx1 = newContext(ctx, TxRollback, "TxRollback")

	tx.db.preExecuteHandler(ctx1)
	var err = tx.TX.Rollback()
	ctx1.err = err
	tx.db.postExecuteHandler(ctx1)
	return err
}
