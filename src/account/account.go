package account

import (
	"github.com/MayCMF/core/src/account/controllers"
	"github.com/MayCMF/core/src/account/controllers/implement"
	"github.com/MayCMF/core/src/account/model"
	imodel "github.com/MayCMF/core/src/account/model/impl/gorm/model"
	"go.uber.org/dig"
)

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

// Inject - injection controllers implementation
func InjectControllers(container *dig.Container) error {
	_ = container.Provide(implement.NewLogin)
	_ = container.Provide(func(b *implement.Login) controllers.ILogin { return b })
	_ = container.Provide(implement.NewPermission)
	_ = container.Provide(func(b *implement.Permission) controllers.IPermission { return b })
	_ = container.Provide(implement.NewRole)
	_ = container.Provide(func(b *implement.Role) controllers.IRole { return b })
	_ = container.Provide(implement.NewUser)
	_ = container.Provide(func(b *implement.User) controllers.IUser { return b })
	return nil
}

// Inject - Injection of gorm
func InjectStarage(container *dig.Container) error {
	_ = container.Provide(imodel.NewPermission)
	_ = container.Provide(func(m *imodel.Permission) model.IPermission { return m })
	_ = container.Provide(imodel.NewRole)
	_ = container.Provide(func(m *imodel.Role) model.IRole { return m })
	_ = container.Provide(imodel.NewUser)
	_ = container.Provide(func(m *imodel.User) model.IUser { return m })
	return nil
}
