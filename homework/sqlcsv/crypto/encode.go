package crypto

import (
	"gb_go_best_practics/homework/sqlcsv/call"
	"gb_go_best_practics/homework/sqlcsv/parser"
	"strconv"
	"strings"
)

func EncodeCSV(result call.Result) string {
	view := result.View

	var header string
	h := make([]string, view.Header.Len())

	for i := range view.Header {
		h[i] = view.Header[i].Label()
	}
	header = strings.Join(h, ",")

	records := make([]string, len(view.Records))

	for i, record := range view.Records {
		cells := make([]string, view.Header.Len())
		for j, cell := range record {
			cells[j] = formatCSVCell(cell)
		}
		records[i] = strings.Join(cells, " , ")
	}

	s := strings.Join(records, "\n")
	s = header + "\n" + s
	return s
}

func formatCSVCell(c call.Cell) string {
	primary := c.Zero()
	var s string

	switch primary.(type) {
	case parser.String:
		s = primary.(parser.String).String()
	case parser.Integer:
		s = strconv.FormatInt((primary.(parser.Integer).Value()), 10)
	case parser.Float:
		s = strconv.FormatFloat(primary.(parser.Float).Value(), 'f', -1, 64)
	case parser.Boolean:
		s = strconv.FormatBool(primary.(parser.Boolean).Logic())
	case parser.LogicOp:
		s = strconv.FormatBool(primary.(parser.LogicOp).Bool())
	case parser.Datetime:
		s = primary.(parser.Datetime).Value().Format("2006-01-02 15:04:05.999999999")
	case parser.Null:
		s = ""
	}

	return s
}
