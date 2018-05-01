package debora

// Row a row in a table
type Row []Column

type row []*column

func (row row) Columns() []Column {
	cx := make([]Column, len(row))
	for i, r := range row {
		cx[i] = r
	}
	return cx
}
