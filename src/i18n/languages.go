package i18n

import (
	"context"
	"os"

	"github.com/MayCMF/core/src/common/config"
	comschema "github.com/MayCMF/core/src/common/schema"
	"github.com/MayCMF/core/src/common/util"
	"github.com/MayCMF/core/src/i18n/controllers"
	"github.com/MayCMF/core/src/i18n/schema"
	transaction "github.com/MayCMF/core/src/transaction/controllers"
	"go.uber.org/dig"
)

// InitLanguages - Initialize Languages data
func InitLanguages(ctx context.Context, container *dig.Container) error {
	if c := config.Global().I18n; c.Enable && c.Data != "" {
		return initLanguagesData(ctx, container)
	}

	return nil
}

// initLanguagesData - Initialize language data
func initLanguagesData(ctx context.Context, container *dig.Container) error {
	return container.Invoke(func(trans transaction.ITrans, language controllers.ILanguage) error {
		// Check if there is language data, initialize if it does not exist
		languageResult, err := language.Query(ctx, schema.LanguageQueryParam{}, schema.LanguageQueryOptions{
			PageParam: &comschema.PaginationParam{PageIndex: -1},
		})
		if err != nil {
			return err
		} else if languageResult.PageResult.Total > 0 {
			return nil
		}

		data, err := readLanguagesData()
		if err != nil {
			return err
		}

		return createLanguages(ctx, trans, language, data)
	})
}

func readLanguagesData() (schema.Languages, error) {
	file, err := os.Open(config.Global().I18n.Data)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data schema.Languages
	err = util.JSONNewDecoder(file).Decode(&data)
	return data, err
}

func createLanguages(ctx context.Context, trans transaction.ITrans, language controllers.ILanguage, list schema.Languages) error {
	return trans.Exec(ctx, func(ctx context.Context) error {
		for _, item := range list {
			sitem := schema.Language{
				Code:    item.Code,
				Name:    item.Name,
				Native:  item.Native,
				Rtl:     item.Rtl,
				Default: item.Default,
				Active:  item.Active,
			}
			_, err := language.Create(ctx, sitem)
			if err != nil {
				return err
			}
		}

		return nil
	})
}
