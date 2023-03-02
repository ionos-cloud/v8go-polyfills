package utils

import (
	v8 "rogchap.com/v8go"
)

// Injector ...
type Injector interface {
	Inject(*v8.Isolate, *v8.ObjectTemplate) error
}
