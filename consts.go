package dbi

/********************************************************************
created:    2020-05-31
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

const (
	None = iota
	BeginTx
	QueryContext
	ExecContext
	TxQueryContext
	TxExecContext
	TxCommit
	TxRollback
)
