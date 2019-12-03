package controllers

import (
	"go.uber.org/dig"
)

func Inject(container *dig.Container) error {
	_ = container.Provide(NewLanguage)
	_ = container.Provide(NewCountry)
	return nil
}
