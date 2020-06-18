package dbi

import "errors"

/********************************************************************
created:    2020-06-18
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

var ErrNotPointer = errors.New("must be a pointer")
var ErrNilPointer = errors.New("nil pointer passed")
var ErrMissingField = errors.New("missing field")
var ErrNotStruct = errors.New("argument not a struct")
