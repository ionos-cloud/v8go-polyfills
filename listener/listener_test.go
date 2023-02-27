package listener_test

import (
	"testing"

	"github.com/ionos-cloud/v8go-polyfills/console"
	"github.com/ionos-cloud/v8go-polyfills/listener"

	"github.com/stretchr/testify/assert"

	v8 "rogchap.com/v8go"
)

func BenchmarkEventListenerCall(b *testing.B) {
	iso := v8.NewIsolate()
	global := v8.NewObjectTemplate(iso)

	events := make(chan *v8.Object)

	if err := listener.AddTo(iso, global, listener.WithEvents("auth", events)); err != nil {
		panic(err)
	}

	ctx := v8.NewContext(iso, global)

	if err := console.AddTo(ctx); err != nil {
		panic(err)
	}

	_, err := ctx.RunScript("addListener('auth', event => { return event.sourceIP === '127.0.0.1' })", "listener.js")
	if err != nil {
		panic(err)
	}

	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		obj, err := newContextObject(ctx)
		assert.NoError(b, err)

		events <- obj
	}
}

func newContextObject(ctx *v8.Context) (*v8.Object, error) {
	iso := ctx.Isolate()
	obj := v8.NewObjectTemplate(iso)

	resObj, err := obj.NewInstance(ctx)
	if err != nil {
		return nil, err
	}

	for _, v := range []struct {
		Key string
		Val interface{}
	}{
		{Key: "sourceIP", Val: "127.0.0.1"},
	} {
		if err := resObj.Set(v.Key, v.Val); err != nil {
			return nil, err
		}
	}

	return resObj, nil
}
