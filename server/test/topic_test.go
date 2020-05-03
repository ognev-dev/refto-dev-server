package test

import (
	"testing"

	"github.com/go-pg/pg/v9"
	"github.com/ognev-dev/bits/database"
	"github.com/ognev-dev/bits/database/factory"
	"github.com/ognev-dev/bits/database/model"
	"github.com/ognev-dev/bits/server/request"
	"github.com/ognev-dev/bits/server/response"
	. "github.com/ognev-dev/bits/test/apitest"
	"github.com/ognev-dev/bits/test/assert"
)

func TestSearchTopic(t *testing.T) {
	// Create test data
	topics := []string{"test_this", "test_that", "test_else"}
	for _, name := range topics {
		_, err := factory.CreateTopic(model.Topic{Name: name})
		assert.NotError(t, err)
	}

	var req request.SearchTopic
	var resp response.SearchTopic

	topicsCount := 0
	_, err := database.ORM().Query(pg.Scan(&topicsCount), "SELECT COUNT(id) FROM topics")
	assert.NotError(t, err)

	cases := map[string]int{
		"":          topicsCount, // all
		"test_th":   2,           // partial match of: this, that
		"test_else": 1,           // exact match
	}
	for query, count := range cases {
		req.Name = query
		TestSearch(t, "topics", req, &resp)
		assert.Equals(t, count, resp.Count)
	}
}
