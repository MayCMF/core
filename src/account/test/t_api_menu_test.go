package test

import (
	"net/http/httptest"
	"testing"

	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/util"
	"github.com/stretchr/testify/assert"
)

func TestAPIPermission(t *testing.T) {
	const router = apiPrefix + "v1/permissions"
	var err error

	w := httptest.NewRecorder()

	// post /permissions
	addItem := &schema.Permission{
		Name:     util.MustUUID(),
		Sequence: 9999999,
		Router:   "/system/permission",
		Actions: []*schema.PermissionAction{
			{Code: "query", Name: "query"},
		},
		Resources: []*schema.PermissionResource{
			{Code: "query", Name: "query", Method: "GET", Path: "/test/v1/permissions"},
		},
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)

	var addNewItem schema.Permission
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.Name, addNewItem.Name)
	assert.Equal(t, addItem.Router, addNewItem.Router)
	assert.Equal(t, addItem.Sequence, addNewItem.Sequence)
	assert.Equal(t, len(addItem.Actions), len(addNewItem.Actions))
	assert.Equal(t, len(addItem.Resources), len(addNewItem.Resources))
	assert.NotEmpty(t, addNewItem.UUID)

	// query /permissions
	engine.ServeHTTP(w, newGetRequest(router, newPageParam()))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.Permission
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	if len(pageItems) > 0 {
		assert.Equal(t, addNewItem.UUID, pageItems[0].UUID)
		assert.Equal(t, addNewItem.Name, pageItems[0].Name)
	}

	// put /permissions/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.UUID))
	assert.Equal(t, 200, w.Code)
	var putItem schema.Permission
	err = parseReader(w.Body, &putItem)
	assert.Nil(t, err)
	putItem.Name = util.MustUUID()
	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, addNewItem.UUID))
	assert.Equal(t, 200, w.Code)
	var putNewItem schema.Permission
	err = parseReader(w.Body, &putNewItem)
	assert.Nil(t, err)
	assert.Equal(t, putItem.Name, putNewItem.Name)
	assert.Equal(t, len(putItem.Actions), len(putNewItem.Actions))
	assert.Equal(t, len(putItem.Resources), len(putNewItem.Resources))

	// delete /permissions/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.UUID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}
