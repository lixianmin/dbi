package dbi

import "context"

/********************************************************************
created:    2020-05-31
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

var emptyErrorFilter = func(err error) error { return err }

var background = &Context{
	Context:     context.Background(),
	ErrorFilter: emptyErrorFilter,
}

type Context struct {
	context.Context
	ErrorFilter func(err error) error // 所有的public方法，都需要在返回err的时候调用errorFilter(err)
}

func NewContext(other Context) *Context {
	var ctx = &Context{
		Context:     other.Context,
		ErrorFilter: other.ErrorFilter,
	}

	if ctx.Context == nil {
		ctx.Context = background.Context
	}

	if ctx.ErrorFilter == nil {
		ctx.ErrorFilter = background.ErrorFilter
	}

	return ctx
}

func Background() *Context {
	return background
}

// 确保返回一个非空的ctx对象
func ensureContext(ctx *Context) *Context {
	if ctx != nil {
		return ctx
	}

	return background
}
