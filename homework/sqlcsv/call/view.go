package call

import (
	"bufio"
	"fmt"
	"gb_go_best_practics/homework/sqlcsv/parser"
	"gb_go_best_practics/homework/sqlcsv/reader"
	"os"
	"strconv"
)

func CheckCsvFile(filename string) error {
	var info os.FileInfo
	var err error
	if info, err = os.Stat(filename); err != nil {
		if info, err = os.Stat(filename + ".csv"); err == nil {
			filename = filename + ".csv"
		} else {
			// should replace `errors.New(fmt.Sprintf(...))` with `fmt.Errorf(...)` (golint)
			return fmt.Errorf("file %s absent", filename)
		}
	}
	// should replace `errors.New(fmt.Sprintf(...))` with `fmt.Errorf(...)` (golint)
	if info.IsDir() {
		return fmt.Errorf("%s is directory", filename)
	}
	return nil
}

type View struct {
	Header          Header
	Records         []Record
	Filename        string
	selectFields    []int
	parentFilter    Filter
	filteredIndices []int
}

func NewView(clause parser.FromExpr, parentFilter Filter) (*View, error) {
	views := make([]*View, len(clause.Tables))
	for i, v := range clause.Tables {
		view, err := loadView(v.(parser.Table))
		if err != nil {
			return nil, err
		}
		views[i] = view
	}

	joinedView := views[0]

	if parentFilter != nil {
		joinedView.parentFilter = parentFilter
	}
	return joinedView, nil
}

func loadView(table parser.Table) (*View, error) {
	var view *View
	var err error

	switch table.Object.(type) {
	case parser.Identifier:
		file := table.Object.(parser.Identifier)
		view, err = loadViewFromFile(file.String(), table.Name())
	}

	if err != nil {
		return nil, err
	}
	return view, nil
}

func loadViewFromFile(filename string, reference string) (*View, error) {
	err := CheckCsvFile(filename)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := bufio.NewReader(f)

	view := new(View)

	reader := reader.NewReader(r)

	var header []string
	header, err = reader.ReadHeader()
	if err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	view.Records = make([]Record, len(records))
	for i, v := range records {
		view.Records[i] = NewRecord(v)
	}

	if header == nil {
		header = make([]string, reader.FieldsPerRec)
		for i := 0; i < reader.FieldsPerRec; i++ {
			header[i] = "c" + strconv.Itoa(i+1)
		}
	}
	view.Header = NewHeader(reference, header)

	view.Filename = filename

	return view, nil
}

func (view *View) Where(clause parser.WhereExpr) error {
	indices, err := view.filter(clause.Filter)
	if err != nil {
		return err
	}

	view.filteredIndices = indices
	return nil
}

func (view *View) filter(condition parser.Expression) ([]int, error) {
	indices := []int{}
	for i := range view.Records {
		var filter Filter = append([]FilterRecord{{View: view, RecordIndex: i}}, view.parentFilter...)
		primary, err := filter.Evaluate(condition)
		if err != nil {
			return nil, err
		}
		if primary.Logic() == true {
			indices = append(indices, i)
		}
	}
	return indices, nil
}

func (view *View) Extract() {
	records := make([]Record, len(view.filteredIndices))
	for i, idx := range view.filteredIndices {
		records[i] = view.Records[idx]
	}
	view.Records = records
	view.filteredIndices = nil
}

func (view *View) Select(clause parser.SelectExpr) error {
	var parseAllColumns = func(view *View, fields []parser.Expression) []parser.Expression {
		insertIdx := -1

		for i, field := range fields {
			if _, ok := field.(parser.Field).Object.(parser.Asterisk); ok {
				insertIdx = i
				break
			}
		}

		if insertIdx < 0 {
			return fields
		}

		columns := view.Header.TableColumns()
		insert := make([]parser.Expression, len(columns))
		for i, c := range columns {
			insert[i] = parser.Field{
				Object: c,
			}
		}

		return append(append(fields[:insertIdx], insert...), fields[insertIdx+1:]...)
	}

	var evalFields = func(view *View, fields []parser.Expression) ([]Record, error) {
		records := make([]Record, len(view.Records))
		for i := range view.Records {
			var record Record
			var filter Filter = append([]FilterRecord{{View: view, RecordIndex: i}}, view.parentFilter...)

			for _, f := range fields {
				field := f.(parser.Field)
				primary, err := filter.Evaluate(field.Object)
				if err != nil {
					return nil, err
				}
				if _, ok := field.Object.(parser.Identifier); !ok {
					record = append(record, NewCell(primary))
				}
			}
			records[i] = record
		}
		return records, nil
	}

	fields := parseAllColumns(view, clause.Fields)
	records, err := evalFields(view, fields)
	if err != nil {
		//if _, ok := err.(*NotGroupedError); ok {
		//	view.group(nil)
		//	records, err = evalFields(view, fields)
		//	if err != nil {
		//		return err
		//	}
		//} else {
		//	return err
		//}
		return err
	}

	view.selectFields = make([]int, len(fields))
	addIndex := view.Header.Len()
	for i, f := range fields {
		field := f.(parser.Field)
		if ident, ok := field.Object.(parser.Identifier); ok {
			ref, field, _ := ident.Field()
			idx, _ := view.Header.Contains(ref, field)
			view.selectFields[i] = idx
		} else {
			view.Header = AddHeaderField(view.Header, field.Name())
			view.selectFields[i] = addIndex
			addIndex++
		}
	}

	for i := range view.Records {
		if 0 < len(records[i]) {
			view.Records[i] = append(view.Records[i], records[i]...)
		}
	}

	return nil
}

func (view *View) Limit(clause parser.LimitExpr) {
	if clause.Number < int64(len(view.Records)) {
		view.Records = view.Records[:clause.Number]
	}
}

func (view *View) FillView() {
	hfields := NewEmptyHeader(len(view.selectFields))
	records := make([]Record, len(view.Records))

	for i, v := range view.Records {
		record := make(Record, len(view.selectFields))
		for j, idx := range view.selectFields {
			record[j] = v[idx]
		}

		records[i] = record
	}

	for i, idx := range view.selectFields {
		hfields[i] = view.Header[idx]
	}

	view.Header = hfields
	view.Records = records
	view.selectFields = []int(nil)
	view.parentFilter = Filter(nil)
}
