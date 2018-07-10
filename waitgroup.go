package goutils

import (
	"reflect"
	"sync"
	"errors"
	"sync/atomic"
)

var (
	ERROR_PARAM_1_IS_NOT_FUNC = errors.New("first parameter is not function")
	ERROR_INPUT_PARAM_NUM_NOT_MATCH = errors.New("input parameter number is not match")
	ERROR_OUTPUT_RESULT_NUM_GREATER_THAN_ONE = errors.New("output result number is greater than 1")
)

type WaitGroupWrapper struct {
	sync.WaitGroup
	sync.Mutex
	currentId uint64
	promise map[uint64]interface{}
}

func NewWaitGroupWrapper() *WaitGroupWrapper{
	return &WaitGroupWrapper{
		promise : make(map[uint64]interface{}),
	}
}

func (wg *WaitGroupWrapper) Run(fn interface{},arg... interface{}) (id uint64,done <- chan struct{},err error) {
	t := reflect.TypeOf(fn)
	if t.Kind() != reflect.Func {
		return 0,nil,ERROR_PARAM_1_IS_NOT_FUNC
	}
	if len(arg) != t.NumIn() {
		return 0,nil,ERROR_INPUT_PARAM_NUM_NOT_MATCH
	}

	if t.NumOut() > 1 {
		return 0,nil,ERROR_OUTPUT_RESULT_NUM_GREATER_THAN_ONE
	}

	var ch chan struct{}
	if t.NumOut() > 0 {
		id = atomic.AddUint64(&wg.currentId,1)
		ch = make(chan struct{})
		done = ch
	}

	v := reflect.ValueOf(fn)
	wg.Add(1)
	go func(args... interface{}) {
		defer wg.Done()
		va := make([]reflect.Value,len(args))
		for i := range args {
			va[i] = reflect.ValueOf(args[i])
		}
		r := v.Call(va)
		if len(r) != 0 {
			wg.Lock()
			wg.promise[id] = r[0].Interface()
			wg.Unlock()
			close(ch)
		}
	}(arg...)

	return
}
