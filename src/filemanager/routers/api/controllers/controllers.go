package controllers

import (
	"go.uber.org/dig"
)

// Inject - injection controllers
func Inject(container *dig.Container) error {
	_ = container.Provide(NewFile)
	return nil
}
