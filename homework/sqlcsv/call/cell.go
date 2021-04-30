package call

import "gb_go_best_practics/homework/sqlcsv/parser"

type Cell []parser.Prime

type Record []Cell

func NewCell(value parser.Prime) Cell {
	return []parser.Prime{value}
}

func NewRecord(values []parser.Prime) Record {
	record := make(Record, len(values))

	for i, v := range values {
		record[i] = NewCell(v)
	}
	return record
}

func (cell Cell) Zero() parser.Prime {
	return cell[0]
}

func (cell Cell) Idx(index int) parser.Prime {
	return cell[index]
}

func (cell Cell) Len() int {
	return len(cell)
}
