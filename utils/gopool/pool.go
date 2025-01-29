package gopool

import (
	"log"
	"reflect"
	"time"

	"github.com/panjf2000/ants/v2"
)

var poolSize = ants.DefaultAntsPoolSize
var defaultPool, _ = ants.NewPool(poolSize, ants.WithOptions(ants.Options{
	ExpiryDuration: (1 * time.Second),
	Nonblocking:    true,
}))

func Submit(task func()) {
	if err := defaultPool.Submit(task); err != nil {
		log.Fatalln(err)
	}
}

func Submits(f interface{}, args ...interface{}) {
	Submit(func() {
		reflect.ValueOf(f).Call(fnArgs(args))
	})
}

func fnArgs(args []interface{}) []reflect.Value {
	values := make([]reflect.Value, len(args))
	for i, arg := range args {
		values[i] = reflect.ValueOf(arg)
	}
	return values
}
