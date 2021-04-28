package call

import "gb_go_best_practics/homework/sqlcsv/parser"

type Result struct {
	View *View
}

func Execute(input string) ([]Result, error) {
	results := []Result{}

	res, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	for _, stmt := range res {
		switch stmt.(type) {
		case parser.SelectReq:
			view, err := ExecuteSelect(stmt.(parser.SelectReq), nil)
			if err != nil {
				return nil, err
			}
			results = append(results, Result{
				View: view,
			})
		}
	}

	return results, nil
}

func ExecuteSelect(query parser.SelectReq, parentFilter Filter) (*View, error) {
	var view *View

	v, err := NewView(query.FromExpr.(parser.FromExpr), parentFilter)
	if err != nil {
		return nil, err
	}
	view = v

	if query.WhereExpr != nil {
		if err := view.Where(query.WhereExpr.(parser.WhereExpr)); err != nil {
			return nil, err
		}
		view.Extract()
	}

	if err := view.Select(query.SelectExpr.(parser.SelectExpr)); err != nil {
		return nil, err
	}

	if query.LimitExpr != nil {
		view.Limit(query.LimitExpr.(parser.LimitExpr))
	}

	view.FillView()

	return view, nil
}
