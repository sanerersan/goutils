package goutils

import (
	"sync/atomic"
	"time"
	"errors"
)

var (
	ERROR_SEMAPHORE_INVALID_PARAMETER = errors.New("semaphore error:invalid parameter")
)

type SemaphoreWarpper struct {
	name string
	semaCap int32
	semaLen int32
	semaCh chan struct{}
}

func (sema* SemaphoreWarpper) GetSemaName() string {
	return sema.name
}

func (sema* SemaphoreWarpper) GetSemaCapacity() int32 {
	return sema.semaCap
}

func (sema* SemaphoreWarpper) GetSemaUsed() int32 {
	return sema.semaLen
}

func (sema *SemaphoreWarpper) WaitSync() {
	sema.semaCh <- struct{}{}
	atomic.AddInt32(&sema.semaLen,1)
}

func (sema *SemaphoreWarpper) WaitTimeOut(timeout time.Duration) bool {
	var waitSuccess = false
	select {
	case sema.semaCh <- struct{}{}:
		atomic.AddInt32(&sema.semaLen,1)
		waitSuccess = true
	case <- time.After(timeout):	
	}

	return waitSuccess	
}

func (sema *SemaphoreWarpper) Release() {
	<- sema.semaCh
	atomic.AddInt32(&sema.semaLen,-1)
}

func NewSemaphore(semaName string,semaCaps,initSemaLen int32) (*SemaphoreWarpper,error) {
	if (semaCaps < initSemaLen) || (semaCaps <= 0) || (initSemaLen < 0) {
		return nil,ERROR_SEMAPHORE_INVALID_PARAMETER
	}

	sema := &SemaphoreWarpper {
		name : semaName,
		semaCap : semaCaps,
		semaLen : initSemaLen,
		semaCh : make(chan struct{},semaCaps),
	}
	var i int32 = 0
	for ; i < sema.semaLen;i++ {
		sema.semaCh <- struct{}{}
	}

	return sema,nil
}
