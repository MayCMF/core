package test

import (
	"net/http/httptest"
	"testing"

	"github.com/MayCMF/core/src/account/schema"
	"github.com/MayCMF/core/src/common/util"
	"github.com/stretchr/testify/assert"
)

func TestAPIRole(t *testing.T) {
	const router = apiPrefix + "v1/roles"
	var err error

	w := httptest.NewRecorder()

	// post /permissions
	addPermissionItem := &schema.Permission{
		Name:     util.MustUUID(),
		Sequence: 9999999,
		Actions: []*schema.PermissionAction{
			{Code: "query", Name: "query"},
		},
		Resources: []*schema.PermissionResource{
			{Code: "query", Name: "query", Method: "GET", Path: "/test/v1/permissions"},
		},
	}
	engine.ServeHTTP(w, newPostRequest(apiPrefix+"v1/permissions", addPermissionItem))
	assert.Equal(t, 200, w.Code)
	var addNewPermissionItem schema.Permission
	err = parseReader(w.Body, &addNewPermissionItem)
	assert.Nil(t, err)

	// post /roles
	addItem := &schema.Role{
		Name:     util.MustUUID(),
		Sequence: 9999999,
		Permissions: []*schema.RolePermission{
			{
				PermissionID: addNewPermissionItem.UUID,
				Actions:      []string{"query"},
				Resources:    []string{"query"},
			},
		},
	}
	engine.ServeHTTP(w, newPostRequest(router, addItem))
	assert.Equal(t, 200, w.Code)
	var addNewItem schema.Role
	err = parseReader(w.Body, &addNewItem)
	assert.Nil(t, err)
	assert.Equal(t, addItem.Name, addNewItem.Name)
	assert.Equal(t, addItem.Sequence, addNewItem.Sequence)
	assert.Equal(t, len(addItem.Permissions), len(addNewItem.Permissions))
	assert.NotEmpty(t, addNewItem.UUID)

	// query /roles
	engine.ServeHTTP(w, newGetRequest(router, newPageParam()))
	assert.Equal(t, 200, w.Code)
	var pageItems []*schema.Role
	err = parsePageReader(w.Body, &pageItems)
	assert.Nil(t, err)
	assert.Equal(t, len(pageItems), 1)
	if len(pageItems) > 0 {
		assert.Equal(t, addNewItem.UUID, pageItems[0].UUID)
		assert.Equal(t, addNewItem.Name, pageItems[0].Name)
	}

	// put /roles/:id
	engine.ServeHTTP(w, newGetRequest("%s/%s", nil, router, addNewItem.UUID))
	assert.Equal(t, 200, w.Code)

	var putItem schema.Role
	err = parseReader(w.Body, &putItem)
	assert.Nil(t, err)
	assert.Equal(t, len(putItem.Permissions), 1)

	putItem.Name = util.MustUUID()
	putItem.Permissions = []*schema.RolePermission{
		{
			PermissionID: addNewPermissionItem.UUID,
			Actions:      []string{},
			Resources:    []string{"query"},
		},
	}

	engine.ServeHTTP(w, newPutRequest("%s/%s", putItem, router, addNewItem.UUID))
	assert.Equal(t, 200, w.Code)
	var putNewItem schema.Role
	err = parseReader(w.Body, &putNewItem)
	assert.Nil(t, err)
	assert.Equal(t, putItem.Name, putNewItem.Name)
	assert.Equal(t, len(putItem.Permissions), len(putNewItem.Permissions))
	assert.Equal(t, 0, len(putNewItem.Permissions[0].Actions))
	assert.Equal(t, 1, len(putNewItem.Permissions[0].Resources))

	// delete /permissions/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", apiPrefix+"v1/permissions", addNewPermissionItem.UUID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)

	// delete /roles/:id
	engine.ServeHTTP(w, newDeleteRequest("%s/%s", router, addNewItem.UUID))
	assert.Equal(t, 200, w.Code)
	err = parseOK(w.Body)
	assert.Nil(t, err)
}

func BenchmarkAPIRoleCreateParallel(b *testing.B) {
	const router = apiPrefix + "v1/roles"

	w := httptest.NewRecorder()

	// post /permissions
	addPermissionItem := &schema.Permission{
		Name:     util.MustUUID(),
		Sequence: 9999999,
		Actions: []*schema.PermissionAction{
			{Code: "query", Name: "query"},
		},
		Resources: []*schema.PermissionResource{
			{Code: "query", Name: "query", Method: "GET", Path: "/test/v1/permissions"},
		},
	}
	engine.ServeHTTP(w, newPostRequest(apiPrefix+"v1/permissions", addPermissionItem))
	var addNewPermissionItem schema.Permission
	_ = parseReader(w.Body, &addNewPermissionItem)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// post /roles
			addItem := &schema.Role{
				Name:     util.MustUUID(),
				Sequence: 9999999,
				Permissions: []*schema.RolePermission{
					{
						PermissionID: addNewPermissionItem.UUID,
						Actions:      []string{"query"},
						Resources:    []string{"query"},
					},
				},
			}
			w := httptest.NewRecorder()
			engine.ServeHTTP(w, newPostRequest(router, addItem))
			if w.Code != 200 {
				b.Errorf("Expected value: %d, given value: %d", 200, w.Code)
			}
		}
	})
}
