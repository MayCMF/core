package filemanager

import (
	"github.com/MayCMF/core/src/filemanager/controllers"
	"github.com/MayCMF/core/src/filemanager/controllers/implement"
	"github.com/MayCMF/core/src/filemanager/model"
	imodel "github.com/MayCMF/core/src/filemanager/model/impl/gorm/model"
	"go.uber.org/dig"
)

// Inject - injection controllers implementation
func InjectControllers(container *dig.Container) error {
	_ = container.Provide(implement.NewFile)
	_ = container.Provide(func(b *implement.File) controllers.IFile { return b })
	return nil
}

// Inject - Injection of gorm
func InjectStarage(container *dig.Container) error {
	_ = container.Provide(imodel.NewFile)
	_ = container.Provide(func(m *imodel.File) model.IFile { return m })
	return nil
}
