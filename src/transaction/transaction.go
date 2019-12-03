package transaction

import (
	"github.com/MayCMF/core/src/transaction/controllers"
	"github.com/MayCMF/core/src/transaction/model"
	"go.uber.org/dig"
)

// Inject - injection controllers implementation
func InjectControllers(container *dig.Container) error {
	_ = container.Provide(controllers.NewTrans)
	_ = container.Provide(func(b *controllers.Trans) controllers.ITrans { return b })
	return nil
}

// Inject - Injection of gorm
func InjectStarage(container *dig.Container) error {
	_ = container.Provide(model.NewTrans)
	_ = container.Provide(func(m *model.Trans) model.ITrans { return m })
	return nil
}
