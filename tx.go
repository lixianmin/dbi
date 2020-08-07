package dbi

import (
	"database/sql"
	"github.com/lixianmin/dbi/reflectx"
)

/********************************************************************
created:    2020-05-31
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

type Tx struct {
	TX     *sql.Tx
	mapper *reflectx.Mapper
	ctx    *Context
}

func (tx *Tx) QueryContext(ctx *Context, query string, args ...interface{}) (*sql.Rows, error) {
	var ctx1 = ensureContext(ctx)
	var rows, err = tx.TX.QueryContext(ctx1, query, args...)
	err = ctx1.ErrorFilter(err)
	return rows, err
}

func (tx *Tx) ExecContext(ctx *Context, query string, args ...interface{}) (sql.Result, error) {
	var ctx1 = ensureContext(ctx)
	var result, err = tx.TX.ExecContext(ctx1, query, args...)
	err = ctx1.ErrorFilter(err)
	return result, err
}

func (tx *Tx) GetContext(ctx *Context, dest interface{}, query string, args ...interface{}) error {
	var ctx1 = ensureContext(ctx)
	var err = getContextInner(ctx1, tx.TX, tx.mapper, dest, query, args...)
	err = ctx1.ErrorFilter(err)
	return err
}

// Any placeholder parameters are replaced with supplied args.
func (tx *Tx) SelectContext(ctx *Context, dest interface{}, query string, args ...interface{}) error {
	var ctx1 = ensureContext(ctx)
	var err = selectContextInner(ctx1, tx.TX, tx.mapper, dest, query, args...)
	err = ctx1.ErrorFilter(err)
	return err
}

func (tx *Tx) Commit() error {
	var err = tx.TX.Commit()
	err = tx.ctx.ErrorFilter(err)
	return err
}

func (tx *Tx) Rollback() error {
	var err = tx.TX.Rollback()
	err = tx.ctx.ErrorFilter(err)
	return err
}
