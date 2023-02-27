package console

import (
	"errors"
	"fmt"
	"io"
	"os"

	v8 "rogchap.com/v8go"
)

// Option ...
type Option func(*console)

// WithOutput ...
func WithOutput(output io.Writer) Option {
	return func(c *console) {
		c.out = output
	}
}

type console struct {
	out        io.Writer
	methodName string
}

// AddTo ...
func AddTo(ctx *v8.Context, opt ...Option) error {
	if ctx == nil {
		return errors.New("v8-polyfills/console: ctx is required")
	}

	c := New(opt...)

	iso := ctx.Isolate()
	con := v8.NewObjectTemplate(iso)

	logFn := v8.NewFunctionTemplate(iso, c.GetFunctionCallback())

	if err := con.Set(c.methodName, logFn, v8.ReadOnly); err != nil {
		return fmt.Errorf("v8-polyfills/console: %w", err)
	}

	conObj, err := con.NewInstance(ctx)
	if err != nil {
		return fmt.Errorf("v8-polyfills/console: %w", err)
	}

	global := ctx.Global()

	if err := global.Set("console", conObj); err != nil {
		return fmt.Errorf("v8-polyfills/console: %w", err)
	}

	return nil
}

// New ...
func New(opt ...Option) *console {
	c := new(console)

	c.out = os.Stdout
	c.methodName = "log"

	for _, o := range opt {
		o(c)
	}

	return c
}

// Console is a JavaScript console object.
type Console interface{}

// GetFunctionCallback ...
func (c *console) GetFunctionCallback() v8.FunctionCallback {
	return func(info *v8.FunctionCallbackInfo) *v8.Value {
		if args := info.Args(); len(args) > 0 {
			inputs := make([]interface{}, len(args))
			for i, input := range args {
				inputs[i] = input
			}

			fmt.Fprintln(c.out, inputs...)
		}

		return nil
	}
}
