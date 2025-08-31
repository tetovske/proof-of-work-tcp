package quote

import "github.com/tetovske/proof-of-work-tcp/internal/model"

func (r *Repository) GetRandom() *model.Quote {
	return r.c.GetRandom()
}
