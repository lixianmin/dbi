package dbi

import (
	"database/sql"
	"github.com/lixianmin/dbi/reflectx"
	"strings"
)

/********************************************************************
created:    2020-05-13
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

type DB struct {
	DB     *sql.DB
	Mapper *reflectx.Mapper
}

func NewDB(db *sql.DB) *DB {
	var my = &DB{
		DB:     db,
		Mapper: reflectx.NewMapperFunc("db", strings.ToLower),
	}

	return my
}

// Connect to a database and verify with a ping.
func Connect(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		_ = db.Close()
		return nil, err
	}

	return NewDB(db), nil
}

func (db *DB) BeginTx(ctx *Context, opts *sql.TxOptions) (*Tx, error) {
	var ctx1 = ensureContext(ctx)
	var tx, err = db.DB.BeginTx(ctx1, opts)
	err = ctx1.ErrorFilter(err)
	var tx1 = &Tx{TX: tx, mapper: db.Mapper, ctx: ctx1}
	return tx1, err
}

func (db *DB) QueryContext(ctx *Context, query string, args ...interface{}) (*sql.Rows, error) {
	var ctx1 = ensureContext(ctx)
	var rows, err = db.DB.QueryContext(ctx1, query, args...)
	err = ctx1.ErrorFilter(err)
	return rows, err
}

func (db *DB) ExecContext(ctx *Context, query string, args ...interface{}) (sql.Result, error) {
	var ctx1 = ensureContext(ctx)
	var result, err = db.DB.ExecContext(ctx1, query, args...)
	err = ctx1.ErrorFilter(err)
	return result, err
}

func (db *DB) GetContext(ctx *Context, dest interface{}, query string, args ...interface{}) error {
	var ctx1 = ensureContext(ctx)
	var err = getContextInner(ctx1, db.DB, db.Mapper, dest, query, args...)
	err = ctx1.ErrorFilter(err)
	return err
}

// Any placeholder parameters are replaced with supplied args.
func (db *DB) SelectContext(ctx *Context, dest interface{}, query string, args ...interface{}) error {
	var ctx1 = ensureContext(ctx)
	var err = selectContextInner(ctx1, db.DB, db.Mapper, dest, query, args...)
	err = ctx1.ErrorFilter(err)
	return err
}
