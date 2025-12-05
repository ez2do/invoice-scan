package locale

import (
	"encoding/json"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
	"invoice-scan/backend/pkg/log"
	"io"
)

type Localize struct {
	bundle      *i18n.Bundle
	defaultLang language.Tag
	locales     map[language.Tag]*i18n.Localizer
}

type MessageOption func(c *i18n.LocalizeConfig)

var (
	defaultLocalize = &Localize{}
)

func Init(l *Localize) {
	defaultLocalize = l
}

func NewLocalize(defaultLang language.Tag) *Localize {
	bundle := i18n.NewBundle(defaultLang)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	return &Localize{
		bundle:      bundle,
		defaultLang: defaultLang,
		locales:     make(map[language.Tag]*i18n.Localizer),
	}
}

func (l *Localize) LoadTOML(r io.Reader, lang string) {
	langTag := language.Make(lang)
	buf, _ := io.ReadAll(r)
	l.bundle.MustParseMessageFileBytes(buf, fmt.Sprintf("%s.toml", langTag.String()))
	l.locales[langTag] = i18n.NewLocalizer(l.bundle, langTag.String())
}

func (l *Localize) T(m string, opts ...MessageOption) string {
	return l.TL(l.defaultLang, m, opts...)
}

func (l *Localize) TL(foundLang language.Tag, m string, opts ...MessageOption) string {
	if foundLang == language.Und {
		foundLang = l.defaultLang
	}
	// find matching language
	var locale *i18n.Localizer
	for lang, localize := range l.locales {
		if lang == foundLang || (foundLang.Parent() != language.Und && foundLang.Parent() == lang) {
			locale = localize
			break
		}
	}
	if locale == nil {
		locale = l.locales[l.defaultLang]
	}
	// if language data is not setup
	if locale == nil {
		return m
	}
	lc := &i18n.LocalizeConfig{
		MessageID: m,
	}
	for _, opt := range opts {
		opt(lc)
	}
	msg, err := locale.Localize(lc)
	if err == nil {
		return msg
	}
	if _, ok := err.(*i18n.MessageNotFoundErr); !ok {
		log.Errorw("Unknown error when localize message", "error", err,
			"context", fmt.Sprintf("m: %v", m))
	}
	return m
}

func T(m string, opts ...MessageOption) string {
	return defaultLocalize.T(m, opts...)
}

func TL(lang, m string, opts ...MessageOption) string {
	lt, err := language.Parse(lang)
	if err != nil {
		lt = defaultLocalize.defaultLang
	}
	return defaultLocalize.TL(lt, m, opts...)
}
