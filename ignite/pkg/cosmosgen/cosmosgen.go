package cosmosgen

import (
	"context"

	gomodmodule "golang.org/x/mod/module"

	"github.com/ignite-hq/cli/ignite/pkg/cosmosanalysis/module"
)

// generateOptions used to configure code generation.
type generateOptions struct {
	includeDirs []string
	gomodPath   string

	jsOut               func(module.Module) string
	jsIncludeThirdParty bool
	tsClientRootPath    string

	specOut string

	dartOut               func(module.Module) string
	dartIncludeThirdParty bool
	dartRootPath          string
}

// TODO add WithInstall.

// Option configures code generation.
type Option func(*generateOptions)

// WithTSClientGeneration adds Typescript Client code generation. tsClientRootPath is used to determine the root path of generated
// Typescript Classes. includeThirdPartyModules and out configures the underlying JS lib generation which is
// documented in WithJSGeneration.
func WithTSClientGeneration(includeThirdPartyModules bool, out func(module.Module) (path string), tsClientRootPath string) Option {
	return func(o *generateOptions) {
		o.jsOut = out
		o.jsIncludeThirdParty = includeThirdPartyModules
		o.tsClientRootPath = tsClientRootPath
	}
}

func WithDartGeneration(includeThirdPartyModules bool, out func(module.Module) (path string), rootPath string) Option {
	return func(o *generateOptions) {
		o.dartOut = out
		o.dartIncludeThirdParty = includeThirdPartyModules
		o.dartRootPath = rootPath
	}
}

// WithGoGeneration adds Go code generation.
func WithGoGeneration(gomodPath string) Option {
	return func(o *generateOptions) {
		o.gomodPath = gomodPath
	}
}

// WithOpenAPIGeneration adds OpenAPI spec generation.
func WithOpenAPIGeneration(out string) Option {
	return func(o *generateOptions) {
		o.specOut = out
	}
}

// IncludeDirs configures the third party proto dirs that used by app's proto.
// relative to the projectPath.
func IncludeDirs(dirs []string) Option {
	return func(o *generateOptions) {
		o.includeDirs = dirs
	}
}

// generator generates code for sdk and sdk apps.
type generator struct {
	ctx          context.Context
	appPath      string
	protoDir     string
	o            *generateOptions
	sdkImport    string
	deps         []gomodmodule.Version
	appModules   []module.Module
	thirdModules map[string][]module.Module // app dependency-modules pair.
}

// Generate generates code from protoDir of an SDK app residing at appPath with given options.
// protoDir must be relative to the projectPath.
func Generate(ctx context.Context, appPath, protoDir string, options ...Option) error {
	g := &generator{
		ctx:          ctx,
		appPath:      appPath,
		protoDir:     protoDir,
		o:            &generateOptions{},
		thirdModules: make(map[string][]module.Module),
	}

	for _, apply := range options {
		apply(g.o)
	}

	if err := g.setup(); err != nil {
		return err
	}

	if g.o.gomodPath != "" {
		if err := g.generateGo(); err != nil {
			return err
		}
	}

	// js generation requires Go types to be existent in the source code. because
	// sdk.Msg implementations defined on the generated Go types.
	// so it needs to run after Go code gen.
	if g.o.jsOut != nil {
		if err := g.generateTS(); err != nil {
			return err
		}
	}

	if g.o.dartOut != nil {
		if err := g.generateDart(); err != nil {
			return err
		}
	}

	if g.o.specOut != "" {
		if err := generateOpenAPISpec(g); err != nil {
			return err
		}
	}

	return nil

}