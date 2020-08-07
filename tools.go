package dbi

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lixianmin/dbi/reflectx"
	"reflect"
)

/********************************************************************
created:    2020-08-07
author:     lixianmin

Copyright (C) - All Rights Reserved
 *********************************************************************/

// Get does a QueryRow using the provided Queryer, and scans the resulting row
// to dest.  If dest is scannable, the result must only have one column.  Otherwise,
// StructScan is used.  Get will return sql.ErrNoRows like row.Scan would.
// Any placeholder parameters are replaced with supplied args.
// An error is returned if the result set is empty.
func getContextInner(ctx context.Context, rawQuery IRawQuery, mapper *reflectx.Mapper, dest interface{}, query string, args ...interface{}) error {
	rows, err := rawQuery.QueryContext(ctx, query, args...)

	if err != nil {
		return err
	}

	if rows == nil {
		return sql.ErrNoRows
	}

	defer rows.Close()

	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr {
		return ErrNotPointer
	}

	if v.IsNil() {
		return ErrNilPointer
	}

	base := reflectx.Deref(v.Type())
	scannable := isScannable(base)

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	if scannable && len(columns) > 1 {
		return fmt.Errorf("scannable dest type %s with >1 columns (%d) in result", base.Kind(), len(columns))
	}

	if scannable {
		return scanOneRow(rows, dest)
	}

	fields := mapper.TraversalsByName(v.Type(), columns)
	// if we are not unsafe and are missing fields, return an error
	if f, err := missingFields(fields); err != nil {
		return fmt.Errorf("missing destination name %s in %T", columns[f], dest)
	}

	values := make([]interface{}, len(columns))

	err = fieldsByTraversal(v, fields, values, true)
	if err != nil {
		return err
	}

	// scan into the struct field pointers and append to our results
	return scanOneRow(rows, values...)
}

// selectContextInner scans all rows into a destination, which must be a slice of any
// type.  If the destination slice type is a Struct, then StructScan will be
// used on each row.  If the destination is some other kind of base type, then
// each row must only have one column which can scan into that type.  This
// allows you to do something like:
//
//    rows, _ := db.Query("select id from people;")
//    var ids []int
//    selectContextInner(rows, &ids)
//
// and ids will be a list of the id results.  I realize that this is a desirable
// interface to expose to users, but for now it will only be exposed via changes
// to `Get` and `Select`.  The reason that this has been implemented like this is
// this is the only way to not duplicate reflect work in the new API while
// maintaining backwards compatibility.
func selectContextInner(ctx context.Context, rawQuery IRawQuery, mapper *reflectx.Mapper, dest interface{}, query string, args ...interface{}) error {
	var v, vp reflect.Value

	value := reflect.ValueOf(dest)

	// json.Unmarshal returns errors for these
	if value.Kind() != reflect.Ptr {
		return ErrNotPointer
	}
	if value.IsNil() {
		return ErrNilPointer
	}
	direct := reflect.Indirect(value)

	slice, err := baseType(value.Type(), reflect.Slice)
	if err != nil {
		return err
	}

	isPtr := slice.Elem().Kind() == reflect.Ptr
	base := reflectx.Deref(slice.Elem())
	scannable := isScannable(base)

	rows, err := rawQuery.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}

	// if something happens here, we want to make sure the rows are Closed
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	// if it's a base type make sure it only has 1 column;  if not return an error
	if scannable && len(columns) > 1 {
		return fmt.Errorf("non-struct dest type %s with >1 columns (%d)", base.Kind(), len(columns))
	}

	if !scannable {
		var values []interface{}

		fields := mapper.TraversalsByName(base, columns)
		// if we are not unsafe and are missing fields, return an error
		if f, err := missingFields(fields); err != nil {
			return fmt.Errorf("missing destination name %s in %T", columns[f], dest)
		}
		values = make([]interface{}, len(columns))

		for rows.Next() {
			// create a new struct type (which returns PtrTo) and indirect it
			vp = reflect.New(base)
			v = reflect.Indirect(vp)

			err = fieldsByTraversal(v, fields, values, true)
			if err != nil {
				return err
			}

			// scan into the struct field pointers and append to our results
			err = rows.Scan(values...)
			if err != nil {
				return err
			}

			if isPtr {
				direct.Set(reflect.Append(direct, vp))
			} else {
				direct.Set(reflect.Append(direct, v))
			}
		}
	} else {
		for rows.Next() {
			vp = reflect.New(base)
			err = rows.Scan(vp.Interface())
			if err != nil {
				return err
			}
			// append
			if isPtr {
				direct.Set(reflect.Append(direct, vp))
			} else {
				direct.Set(reflect.Append(direct, reflect.Indirect(vp)))
			}
		}
	}

	return rows.Err()
}

func baseType(t reflect.Type, expected reflect.Kind) (reflect.Type, error) {
	t = reflectx.Deref(t)
	if t.Kind() != expected {
		return nil, fmt.Errorf("expected %s but got %s", expected, t.Kind())
	}
	return t, nil
}

// isScannable takes the reflect.Type and the actual dest value and returns
// whether or not it's Scannable.  Something is scannable if:
//   * it is not a struct
//   * it implements sql.Scanner
//   * it has no exported fields

var _scannerInterface = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

func isScannable(t reflect.Type) bool {
	if reflect.PtrTo(t).Implements(_scannerInterface) {
		return true
	}
	if t.Kind() != reflect.Struct {
		return true
	}

	return false
}

func missingFields(transversals [][]int) (field int, err error) {
	for i, t := range transversals {
		if len(t) == 0 {
			return i, ErrMissingField
		}
	}

	return 0, nil
}

// fieldsByName fills a values interface with fields from the passed value based
// on the traversals in int.  If ptrs is true, return addresses instead of values.
// We write this instead of using FieldsByName to save allocations and map lookups
// when iterating over many rows.  Empty traversals will get an interface pointer.
// Because of the necessity of requesting ptrs or values, it's considered a bit too
// specialized for inclusion in reflectx itself.
func fieldsByTraversal(v reflect.Value, traversals [][]int, values []interface{}, ptrs bool) error {
	v = reflect.Indirect(v)
	if v.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i, traversal := range traversals {
		if len(traversal) == 0 {
			values[i] = new(interface{})
			continue
		}
		f := reflectx.FieldByIndexes(v, traversal)
		if ptrs {
			values[i] = f.Addr().Interface()
		} else {
			values[i] = f.Interface()
		}
	}

	return nil
}

// Scan is a fixed implementation of sql.Row.Scan, which does not discard the
// underlying error from the internal rows object if it exists.
func scanOneRow(rows *sql.Rows, dest ...interface{}) error {
	for _, dp := range dest {
		if _, ok := dp.(*sql.RawBytes); ok {
			return errors.New("sql: RawBytes isn't allowed on Row.Scan")
		}
	}

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return err
		}

		return sql.ErrNoRows
	}

	err := rows.Scan(dest...)
	if err != nil {
		return err
	}

	// Make sure the query can be processed to completion with no errors.
	if err := rows.Close(); err != nil {
		return err
	}

	return nil
}
