package test

import (
	"testing"

	"github.com/ognev-dev/bits/database/factory"
	"github.com/ognev-dev/bits/database/model"
	"github.com/ognev-dev/bits/server/request"
	"github.com/ognev-dev/bits/server/response"
	. "github.com/ognev-dev/bits/test/apitest"
	"github.com/ognev-dev/bits/test/assert"
)

func TestSearchEntity(t *testing.T) {
	// Create test data
	topic1, err := factory.CreateTopic(model.Topic{Name: "topic1"})
	assert.NotError(t, err)
	topic2, err := factory.CreateTopic(model.Topic{Name: "topic2"})
	assert.NotError(t, err)
	topic3, err := factory.CreateTopic(model.Topic{Name: "topic3"})
	assert.NotError(t, err)

	// entity with 1 topic
	_, err = factory.CreateEntity(model.Entity{
		Title:  "ent1",
		Topics: []model.Topic{{ID: topic1.ID}},
	})
	assert.NotError(t, err)

	// entity with 2 topics
	ent2, err := factory.CreateEntity(model.Entity{
		Title: "ent2",
		Topics: []model.Topic{
			{ID: topic1.ID},
			{ID: topic2.ID},
		},
	})
	assert.NotError(t, err)

	// entity with 3 topics
	ent3, err := factory.CreateEntity(model.Entity{
		Title: "ent3",
		Topics: []model.Topic{
			{ID: topic1.ID},
			{ID: topic2.ID},
			{ID: topic3.ID},
		},
	})
	assert.NotError(t, err)

	var req request.SearchEntity
	var resp response.SearchEntity

	// should get ent2 and ent3
	// as they both have t1 and t2 (and ent1 missing t2)
	req.Topics = []string{topic1.Name, topic2.Name}
	TestSearch(t, "entities", req, &resp)
	assert.Equals(t, 2, resp.Count)
	for _, v := range resp.Data {
		if v.ID == ent2.ID || v.ID == ent3.ID {
			continue
		}
		t.Fatal("invalid entity in response")
	}

	// should get only ent3
	// (ent2 missing t3 and ent1 have none of them)
	req.Topics = []string{topic2.Name, topic3.Name}
	TestSearch(t, "entities", req, &resp)
	assert.Equals(t, 1, resp.Count)
	assert.Equals(t, ent3.ID, resp.Data[0].ID)

	// should get ent1, ent2, ent3
	req.Topics = []string{topic1.Name}
	TestSearch(t, "entities", req, &resp)
	assert.Equals(t, 3, resp.Count)

	// should get only ent3
	req.Topics = []string{topic1.Name, topic2.Name, topic3.Name}
	TestSearch(t, "entities", req, &resp)
	assert.Equals(t, 1, resp.Count)
	assert.Equals(t, ent3.ID, resp.Data[0].ID)
}
