package locale

import "github.com/nicksnyder/go-i18n/v2/i18n"

func WithTplData(data interface{}) MessageOption {
	return func(c *i18n.LocalizeConfig) {
		c.TemplateData = data
	}
}
