package i18n

import (
	"github.com/MayCMF/core/src/i18n/controllers"
	"github.com/MayCMF/core/src/i18n/controllers/implement"
	"github.com/MayCMF/core/src/i18n/model"
	imodel "github.com/MayCMF/core/src/i18n/model/impl/gorm/model"
	"go.uber.org/dig"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// Inject - injection controllers implementation
func InjectControllers(container *dig.Container) error {
	_ = container.Provide(implement.NewLanguage)
	_ = container.Provide(func(b *implement.Language) controllers.ILanguage { return b })
	_ = container.Provide(implement.NewCountry)
	_ = container.Provide(func(b *implement.Country) controllers.ICountry { return b })
	return nil
}

// Inject - Injection of gorm
func InjectStarage(container *dig.Container) error {
	_ = container.Provide(imodel.NewLanguage)
	_ = container.Provide(func(m *imodel.Language) model.ILanguage { return m })
	_ = container.Provide(imodel.NewCountry)
	_ = container.Provide(func(m *imodel.Country) model.ICountry { return m })
	return nil
}
