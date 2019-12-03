package test

import (
	"net/http/httptest"
	"testing"

	"github.com/MayCMF/core/src/common/util"
	"github.com/MayCMF/core/src/i18n/schema"
	"github.com/stretchr/testify/assert"
)

func TestAPILnguage(t *testing.T) {
	const router = apiPrefix + "v1/language"
	var err error

	w := httptest.NewRecorder()

	// post /language
	addItem := &schema.Lnguage{
		Code:   util.MustUUID(),
		Name:   util.MustUUID(),
		Status: 1,
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)

	var addNewItem schema.Lnguage
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.Code, addNewItem.Code)
	assert.Equal(t, addItem.Code, addNewItem.Code)
	assert.Equal(t, addItem.Status, addNewItem.Status)
	assert.NotEmpty(t, addNewItem.UUID)

	// query /language
	engine.ServeHTTP(w, newGetRequest(router, newPageParam()))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.Lnguage
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	if len(pageItems) > 0 {
		assert.Equal(t, addNewItem.UUID, pageItems[0].UUID)
		assert.Equal(t, addNewItem.Name, pageItems[0].Name)
	}

	// put /language/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.UUID))
	assert.Equal(t, 200, w.Code)
	var putItem schema.Lnguage
	err = parseReader(w.Body, &putItem)
	assert.Nil(t, err)

	putItem.Name = util.MustUUID()
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, addNewItem.UUID))
	assert.Equal(t, 200, w.Code)

	var putNewItem schema.Lnguage
	err = parseReader(w.Body, &putNewItem)
	assert.Nil(t, err)
	assert.Equal(t, putItem.Name, putNewItem.Name)

	// delete /language/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.UUID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}
