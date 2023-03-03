package wasm

import (
	"github.com/ionos-cloud/v8go-polyfills/utils"

	v8 "rogchap.com/v8go"
)

// Module ...
type Module struct {
	utils.Injector
}

// New ...
func New() *Module {
	return &Module{}
}

// Inject ...
func (m *Module) Inject(*v8.Isolate, *v8.ObjectTemplate) error {
	return nil
}
