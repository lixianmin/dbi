package dbi

import "context"

/********************************************************************
created:    2020-05-31
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

type Context struct {
	context.Context
	Kind int
	Text string
	//StartTime time.Time
	err error
}

func newContext(ctx context.Context, kind int, text string) *Context {
	return &Context{
		Context: ctx,
		Kind:    kind,
		Text:    text,
		//StartTime: time.Now(),
	}
}

func (ctx *Context) Err() error {
	if ctx.err != nil {
		return ctx.err
	}

	return ctx.Context.Err()
}
