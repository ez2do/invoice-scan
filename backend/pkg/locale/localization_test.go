package locale

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
	"testing"
)

func initLocalize() *Localize {
	loc := NewLocalize(language.English)
	return loc
}

func TestLocalize_T(t *testing.T) {
	loc := initLocalize()
	loc.LoadTOML(bytes.NewBuffer([]byte(`
["hello world"]
other = "Hello World!"
`)), "en")

	msg := loc.T("hello world")
	assert.Equal(t, "Hello World!", msg)
}

func TestTL_NonSetupData(t *testing.T) {
	loc := initLocalize()
	Init(loc)
	m := TL("en", "hello world")
	assert.Equal(t, m, "hello world")
}

func TestTL(t *testing.T) {
	loc := initLocalize()
	loc.LoadTOML(bytes.NewBuffer([]byte(`
"hello world" = "Hello World!"
"hello someone" = "Hello {{.Name}}!"
"password cannot contain token" = "password cannot contain token {{.Token}}"
`)), "en")
	loc.LoadTOML(bytes.NewBuffer([]byte(`
"hello world" = "Hallo Welt!"
`)), "de")
	Init(loc)
	type args struct {
		lang string
		m    string
		opts []MessageOption
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "found localize",
			args: args{
				lang: "en",
				m:    "hello world",
				opts: nil,
			},
			want: "Hello World!",
		},
		{
			name: "found parent localize",
			args: args{
				lang: "en-US",
				m:    "hello world",
				opts: nil,
			},
			want: "Hello World!",
		},
		{
			name: "found localize with template",
			args: args{
				lang: "en",
				m:    "hello someone",
				opts: []MessageOption{WithTplData(map[string]interface{}{"Name": "F"})},
			},
			want: "Hello F!",
		},
		{
			name: "found localize in german",
			args: args{
				lang: "de",
				m:    "hello world",
				opts: nil,
			},
			want: "Hallo Welt!",
		},
		{
			name: "found localize use default lang",
			args: args{
				lang: "",
				m:    "hello world",
				opts: nil,
			},
			want: "Hello World!",
		},
		{
			name: "found localize use none load lang (response default lang)",
			args: args{
				lang: "vi",
				m:    "hello world",
				opts: nil,
			},
			want: "Hello World!",
		},
		{
			name: "not found localize",
			args: args{
				lang: "en",
				m:    "another/world",
				opts: nil,
			},
			want: "another/world",
		},
		{
			name: "with template data",
			args: args{
				lang: "en",
				m:    "password cannot contain token",
				opts: []MessageOption{WithTplData(struct {
					Token string
				}{Token: "GHTK"})},
			},
			want: "password cannot contain token GHTK",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TL(tt.args.lang, tt.args.m, tt.args.opts...); got != tt.want {
				t.Errorf("TL() = %v, want %v", got, tt.want)
			}
		})
	}
}
