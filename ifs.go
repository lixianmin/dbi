package dbi

import "database/sql"

/********************************************************************
created:    2020-06-18
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

type IQueryContext interface {
	QueryContext(ctx *Context, query string, args ...interface{}) (*sql.Rows, error)
}

type IExecContext interface {
	ExecContext(ctx *Context, query string, args ...interface{}) (sql.Result, error)
}
