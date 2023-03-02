package utils

import (
	v8 "rogchap.com/v8go"
)

// Injector ...
type Injector func(*v8.Isolate, *v8.ObjectTemplate) error
