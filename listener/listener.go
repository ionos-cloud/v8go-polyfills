package listener

import (
	"fmt"

	v8 "rogchap.com/v8go"
)

// Option ...
type Option func(*listener)

type listener struct {
	events map[string]chan *v8.Object
}

// Listener ...
type Listener interface {
	GetFunctionCallback() v8.FunctionCallback
}

// New ...
func New(opt ...Option) *listener {
	c := new(listener)
	c.events = make(map[string]chan *v8.Object)

	for _, o := range opt {
		o(c)
	}

	return c
}

// WithEvents ...
func WithEvents(name string, events chan *v8.Object) Option {
	return func(l *listener) {
		l.events[name] = events
	}
}

// AddTo ...
func AddTo(iso *v8.Isolate, global *v8.ObjectTemplate, opt ...Option) error {
	l := New(opt...)

	ctxFn := v8.NewFunctionTemplate(iso, l.GetFunctionCallback())

	if err := global.Set("addListener", ctxFn, v8.ReadOnly); err != nil {
		return fmt.Errorf("v8-polyfills/listener: %w", err)
	}

	return nil
}

// GetFunctionCallback ...
func (l *listener) GetFunctionCallback() v8.FunctionCallback {
	return func(info *v8.FunctionCallbackInfo) *v8.Value {
		ctx := info.Context()
		args := info.Args()

		if len(args) <= 1 {
			err := fmt.Errorf("addListener: expected 2 arguments, got %d", len(args))

			return newErrorValue(ctx, err)
		}

		fn, err := args[1].AsFunction()
		if err != nil {
			err := fmt.Errorf("%w", err)

			return newErrorValue(ctx, err)
		}

		chn, ok := l.events[args[0].String()]
		if !ok {
			err := fmt.Errorf("addListener: event %s not found", args[0].String())

			return newErrorValue(ctx, err)
		}

		go func(chn chan *v8.Object, fn *v8.Function) {
			for e := range chn {
				v, err := fn.Call(ctx.Global(), e)
				if err != nil {
					fmt.Printf("addListener: %v", err)
				}

				v.Release()
			}
		}(chn, fn)

		return v8.Undefined(ctx.Isolate())
	}
}

func newErrorValue(ctx *v8.Context, err error) *v8.Value {
	iso := ctx.Isolate()
	e, _ := v8.NewValue(iso, fmt.Sprintf("addListener: %v", err))

	return e
}
