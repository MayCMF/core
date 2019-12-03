package errors

import (
	"github.com/pkg/errors"
)

// Define an alias
var (
	New         = errors.New
	Wrap        = errors.Wrap
	Wrapf       = errors.Wrapf
	WithStack   = errors.WithStack
	WithMessage = errors.WithMessage
	// WithMessagef = errors.WithMessagef
)

// Definition error
var (
	ErrBadRequest              = New400Response("Request error")
	ErrInvalidParent           = New400Response("Invalid parent node")
	ErrNotAllowDeleteWithChild = New400Response("Contains children, cannot be deleted")
	ErrNotAllowDelete          = New400Response("Resources are not allowed to delete")
	ErrInvalidUserName         = New400Response("Invalid username")
	ErrInvalidPassword         = New400Response("Invalid password")
	ErrInvalidUser             = New400Response("Invalid user")
	ErrUserDisable             = New400Response("User is disabled, please contact administrator")

	ErrNoPerm          = NewResponse(401, "No access", 401)
	ErrInvalidToken    = NewResponse(9999, "Token invalidation", 401)
	ErrNotFound        = NewResponse(404, "Resource does not exist.", 404)
	ErrMethodNotAllow  = NewResponse(405, "Method is not allowed", 405)
	ErrTooManyRequests = NewResponse(429, "Request too frequently", 429)
	ErrInternalServer  = NewResponse(500, "Server error", 500)
)
