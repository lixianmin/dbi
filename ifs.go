package dbi

import (
	"context"
	"database/sql"
)

/********************************************************************
created:    2020-06-18
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

type IContext interface {
	ErrorFilter(err error) error // 所有的public方法，都需要在返回err的时候调用errorFilter(err)
}

type IQueryContext interface {
	QueryContext(ctx *Context, query string, args ...interface{}) (*sql.Rows, error)
}

type IExecContext interface {
	ExecContext(ctx *Context, query string, args ...interface{}) (sql.Result, error)
}

type IRawQuery interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}