package controllers

import (
	"go.uber.org/dig"
)

// Inject - injection ctl
func Inject(container *dig.Container) error {
	_ = container.Provide(NewLogin)
	_ = container.Provide(NewPermission)
	_ = container.Provide(NewRole)
	_ = container.Provide(NewUser)
	return nil
}