// Copyright (c) 2019, The Emergent Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dtable

import (
	"fmt"

	"github.com/emer/emergent/etensor"
	"github.com/goki/ki/ints"
)

// dtable.Table is the DataTable structure, containing columns of etensor tensors.
// All tensors MUST have RowMajor stride layout!
type Table struct {
	Cols       []etensor.Tensor `desc:"columns of data, as etensor.Tensor tensors"`
	ColNames   []string         `desc:"the names of the columns"`
	Rows       int              `desc:"number of rows, which is enforced to be the size of the outer-most dimension of the column tensors"`
	ColNameMap map[string]int   `desc:"the map of column names to column numbers"`
}

// NumRows returns the number of rows (arrow / dframe api)
func (dt *Table) NumRows() int {
	return dt.Rows
}

// NumCols returns the number of columns (arrow / dframe api)
func (dt *Table) NumCols() int {
	return len(dt.Cols)
}

// Col returns the tensor at given column index
func (dt *Table) Col(i int) etensor.Tensor {
	return dt.Cols[i]
}

// Schema returns the Schema (column properties) for this table
func (dt *Table) Schema() Schema {
	nc := dt.NumCols()
	sc := make(Schema, nc)
	for i := range dt.Cols {
		cl := &sc[i]
		tsr := dt.Cols[i]
		cl.Name = dt.ColNames[i]
		cl.Type = etensor.Type(tsr.DataType().ID())
		cl.CellShape = tsr.Shapes()[1:]
		cl.DimNames = tsr.DimNames()[1:]
	}
	return sc
}

// ColByName returns the tensor at given column name, error if not found
func (dt *Table) ColByName(name string) (etensor.Tensor, error) {
	i, ok := dt.ColNameMap[name]
	if !ok {
		return nil, fmt.Errorf("dtable.Table ColByName: column named: %v not found", name)
	}
	return dt.Cols[i], nil
}

// ColName returns the name of given column
func (dt *Table) ColName(i int) string {
	return dt.ColNames[i]
}

// UpdateColNameMap updates the column name map
func (dt *Table) UpdateColNameMap() {
	nc := dt.NumCols()
	dt.ColNameMap = make(map[string]int, nc)
	for i := range dt.ColNames {
		dt.ColNameMap[dt.ColNames[i]] = i
	}
}

// AddCol adds the given tensor as a column to the table.
// returns error if it is not a RowMajor organized tensor, and automatically
// adjusts the shape to fit the current number of rows.
func (dt *Table) AddCol(tsr etensor.Tensor, name string) error {
	if !tsr.IsRowMajor() {
		return fmt.Errorf("tensor must be RowMajor organized")
	}
	dt.Cols = append(dt.Cols, tsr)
	dt.ColNames = append(dt.ColNames, name)
	dt.UpdateColNameMap()
	tsr.SetNumRows(dt.Rows)
	return nil
}

// AddRows adds n rows to each of the columns
func (dt *Table) AddRows(n int) {
	for _, tsr := range dt.Cols {
		tsr.AddRows(n)
	}
}

// SetNumRows sets the number of rows in the table, across all columns
// if rows = 0 then effective number of rows in tensors is 1, as this dim cannot be 0
func (dt *Table) SetNumRows(rows int) {
	dt.Rows = rows // can be 0
	rows = ints.MaxInt(1, rows)
	for _, tsr := range dt.Cols {
		tsr.SetNumRows(rows)
	}
}

// SetFromSchema configures table from given Schema.
// The actual tensor number of rows is enforced to be > 0, because we
// cannot have a null dimension in tensor shape.
// does not preserve any existing columns / data.
func (dt *Table) SetFromSchema(sc Schema, rows int) {
	nc := len(sc)
	dt.Cols = make([]etensor.Tensor, nc)
	dt.ColNames = make([]string, nc)
	dt.Rows = rows // can be 0
	rows = ints.MaxInt(1, rows)
	for i := range dt.Cols {
		cl := &sc[i]
		dt.ColNames[i] = cl.Name
		sh := append([]int{rows}, cl.CellShape...)
		dn := append([]string{"row"}, cl.DimNames...)
		tsr := etensor.New(cl.Type, sh, nil, dn)
		dt.Cols[i] = tsr
	}
	dt.UpdateColNameMap()
}

// New returns a new Table constructed from given Schema.
// The actual tensor number of rows is enforced to be > 0, because we
// cannot have a null dimension in tensor shape
func New(sc Schema, rows int) *Table {
	dt := &Table{}
	dt.SetFromSchema(sc, rows)
	return dt
}