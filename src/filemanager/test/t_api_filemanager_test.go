package test

import (
	"net/http/httptest"
	"testing"

	"github.com/MayCMF/core/src/common/util"
	"github.com/MayCMF/core/src/filemanager/schema"
	"github.com/stretchr/testify/assert"
)

func TestAPIFile(t *testing.T) {
	const router = apiPrefix + "v1/file"
	var err error

	w := httptest.NewRecorder()

	// post /file
	addItem := &schema.File{
		Filename: util.MustUUID(),
		Filemime: util.MustUUID(),
		Uri:      util.MustUUID(),
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Filename)

	var addNewItem schema.File
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.Filename, addNewItem.Filename)
	assert.Equal(t, addItem.Filename, addNewItem.Filename)
	assert.Equal(t, addItem.Uri, addNewItem.Uri)
	assert.NotEmpty(t, addNewItem.UUID)

	// query /file
	engine.ServeHTTP(w, newGetRequest(router, newPageParam()))
	assert.Equal(t, 200, w.Filename)
	var pageItems []*schema.File
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	if len(pageItems) > 0 {
		assert.Equal(t, addNewItem.UUID, pageItems[0].UUID)
		assert.Equal(t, addNewItem.Filename, pageItems[0].Filename)
	}

	// put /file/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.UUID))
	assert.Equal(t, 200, w.Filename)
	var putItem schema.File
	err = parseReader(w.Body, &putItem)
	assert.Nil(t, err)

	putItem.Filename = util.MustUUID()
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, addNewItem.UUID))
	assert.Equal(t, 200, w.Filename)

	var putNewItem schema.File
	err = parseReader(w.Body, &putNewItem)
	assert.Nil(t, err)
	assert.Equal(t, putItem.Filename, putNewItem.Filename)

	// delete /file/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.UUID))
	assert.Equal(t, 200, w.Filename)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}
