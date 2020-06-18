package dbi

import (
	"time"
)

/********************************************************************
created:    2020-05-31
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

var emptyErrorFilter = func(err error) error { return err }

var background = &Context{
	parent: nil,
	args:   ContextArgs{ErrorFilter: emptyErrorFilter},
}

type ContextArgs struct {
	ErrorFilter func(err error) error // 所有的public方法，都需要在返回err的时候调用errorFilter(err)
}

type Context struct {
	parent *Context
	args   ContextArgs
}

func NewContext(parent *Context, args ContextArgs) *Context {
	parent = ensureContext(parent)
	ensureArgs(&args)

	var ctx = &Context{
		parent: parent,
		args:   args,
	}

	return ctx
}

func (ctx *Context) errorFilter(err error) error {
	var err1 = ctx.args.ErrorFilter(err)
	if ctx.parent != nil {
		err1 = ctx.parent.errorFilter(err1)
	}

	return err1
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	if ctx.parent != nil {
		return ctx.parent.Deadline()
	}

	return
}

func (ctx *Context) Done() <-chan struct{} {
	if ctx.parent != nil {
		return ctx.parent.Done()
	}

	return nil
}

func (ctx *Context) Err() error {
	if ctx.parent != nil {
		return ctx.parent.Err()
	}

	return nil
}

func (ctx *Context) Value(key interface{}) interface{} {
	if ctx.parent != nil {
		return ctx.parent.Value(key)
	}

	return nil
}

func (ctx *Context) String() string {
	if ctx.parent != nil {
		return ctx.parent.String()
	}

	return ""
}

///////////////////////////////////////////////////////////
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

func ensureArgs(args *ContextArgs) {
	if args.ErrorFilter == nil {
		args.ErrorFilter = emptyErrorFilter
	}
}
