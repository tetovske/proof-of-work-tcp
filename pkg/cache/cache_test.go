package cache

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tetovske/proof-of-work-tcp/internal/model"
)

func TestCache_Fill(t *testing.T) {
	data := []*model.Quote{
		{
			Text: "test_text1",
		},
		{
			Text: "test_text2",
		},
	}

	c := New[*model.Quote](10)
	c.Fill(data)

	assert.Equal(t, 2, len(c.sl))
}

func TestCache_GetRandom(t *testing.T) {
	data := []*model.Quote{
		{
			Text: "test_text1",
		},
		{
			Text: "test_text2",
		},
	}

	c := New[*model.Quote](10)
	c.Fill(data)

	assert.Equal(t, 2, len(c.sl))

	v := c.GetRandom()

	assert.NotNil(t, v)
	assert.True(t, strings.HasPrefix(v.Text, "test_"))
}
