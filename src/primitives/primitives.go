package primitives

import (
	"github.com/MayCMF/core/src/primitives/controllers"
	"github.com/MayCMF/core/src/primitives/controllers/implement"
	"github.com/MayCMF/core/src/primitives/model"
	imodel "github.com/MayCMF/core/src/primitives/model/impl/gorm/model"
	"go.uber.org/dig"
)

// Inject - injection controllers implementation
func InjectControllers(container *dig.Container) error {
	_ = container.Provide(implement.NewPrimitive)
	_ = container.Provide(func(b *implement.Primitive) controllers.IPrimitive { return b })
	_ = container.Provide(implement.NewNode)
	_ = container.Provide(func(b *implement.Node) controllers.INode { return b })
	return nil
}

// Inject - Injection of gorm
func InjectStarage(container *dig.Container) error {
	_ = container.Provide(imodel.NewPrimitive)
	_ = container.Provide(func(m *imodel.Primitive) model.IPrimitive { return m })
	_ = container.Provide(imodel.NewNode)
	_ = container.Provide(func(m *imodel.Node) model.INode { return m })
	return nil
}
