package quote

import (
	"github.com/tetovske/proof-of-work-tcp/internal/model"
	"github.com/tetovske/proof-of-work-tcp/pkg/cache"
)

type Repository struct {
	c *cache.Cache[*model.Quote]
}

func New(c *cache.Cache[*model.Quote]) *Repository {
	return &Repository{
		c: c,
	}
}

func (r *Repository) WarmUpCache(data []string) {
	quotes := make([]*model.Quote, 0, len(data))
	for _, d := range data {
		quotes = append(quotes, &model.Quote{Text: d})
	}

	r.c.Fill(quotes)
}
