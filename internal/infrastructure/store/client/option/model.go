package option

import (
	model "github.com/bjlag/go-keeper/internal/domain/client"
)

type row struct {
	Slug  string `db:"slug"`
	Value string `db:"value"`
}

func toRow(model model.Option) row {
	return row{
		Slug:  model.Slug,
		Value: model.Value,
	}
}

func (r *row) toModel() model.Option {
	return model.Option{
		Slug:  r.Slug,
		Value: r.Value,
	}
}
