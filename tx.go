package dbi

import (
	"database/sql"
)

/********************************************************************
created:    2020-05-31
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

type Tx struct {
	TX  *sql.Tx
	ctx *Context
}

func (tx *Tx) QueryContext(ctx *Context, query string, args ...interface{}) (*sql.Rows, error) {
	var ctx1 = ensureContext(ctx)
	var rows, err = tx.TX.QueryContext(ctx1, query, args...)
	err = ctx1.errorFilter(err)
	return rows, err
}

func (tx *Tx) ExecContext(ctx *Context, query string, args ...interface{}) (sql.Result, error) {
	var ctx1 = ensureContext(ctx)
	var result, err = tx.TX.ExecContext(ctx1, query, args...)
	err = ctx1.errorFilter(err)
	return result, err
}

func (tx *Tx) Commit() error {
	var err = tx.TX.Commit()
	err = tx.ctx.errorFilter(err)
	return err
}

func (tx *Tx) Rollback() error {
	var err = tx.TX.Rollback()
	err = tx.ctx.errorFilter(err)
	return err
}
