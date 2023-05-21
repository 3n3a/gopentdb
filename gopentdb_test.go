package gopentdb_test

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/3n3a/gopentdb"
)

func TestNew(t *testing.T)()  {
	o := gopentdb.New(gopentdb.Config{
		BaseUrl: "https://opentdb.com",
	})
	isUp, err := o.Ping()
	want := true
	if isUp != want || err != nil {
		t.Failed()
	}
}

func TestGetCategories(t *testing.T)() {
	o := gopentdb.New(gopentdb.Config{
		BaseUrl: "https://opentdb.com",
	})
	categories, err := o.GetCategories()
	assert.NilError(t, err)
	assert.Assert(t, len(categories) > 0)
	assert.Assert(t, categories[0].Id > 0)
	assert.Assert(t, len(categories[0].Name) > 0)
}


func TestGetCategoryCount(t *testing.T)()  {
	o := gopentdb.New(gopentdb.Config{
		BaseUrl: "https://opentdb.com",
	})
	categories, err := o.GetCategories()
	assert.NilError(t, err)

	count, err := o.GetCategoryCount(categories[0].Id)
	assert.NilError(t, err)
	assert.Assert(t, count > 0)
}