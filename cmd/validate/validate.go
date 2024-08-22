package validate

import (
	"fmt"
	pkg_validate "github.com/Beaelf/validate-lint/pkg/validate"
	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"
)

func init() {
	register.Plugin("validate-lint", New)
}

type MySettings struct {
	One   string    `json:"one"`
	Two   []Element `json:"two"`
	Three Element   `json:"three"`
}

type Element struct {
	Name string `json:"name"`
}

type PluginExample struct {
	settings MySettings
}

func New(settings any) (register.LinterPlugin, error) {
	// The configuration type will be map[string]any or []interface, it depends on your configuration.
	// You can use https://github.com/go-viper/mapstructure to convert map to struct.

	s, err := register.DecodeSettings[MySettings](settings)
	if err != nil {
		return nil, err
	}

	fmt.Printf("==> setting: %v\n", s)

	return &PluginExample{settings: s}, nil
}

func (f *PluginExample) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	return []*analysis.Analyzer{
		pkg_validate.Validate,
	}, nil
}

func (f *PluginExample) GetLoadMode() string {
	return register.LoadModeSyntax
}
