// Code generated by counterfeiter. DO NOT EDIT.
package capturepageslogicfakes

import (
	"context"
	"sync"

	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
)

type FakePersistor struct {
	CreateCapturePagesStub        func(context.Context, persistence.TransactionHandler, *model.CapturePages) (*model.CapturePages, error)
	createCapturePagesMutex       sync.RWMutex
	createCapturePagesArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.CapturePages
	}
	createCapturePagesReturns struct {
		result1 *model.CapturePages
		result2 error
	}
	createCapturePagesReturnsOnCall map[int]struct {
		result1 *model.CapturePages
		result2 error
	}
	GetCapturePageTypeByIdStub        func(context.Context, persistence.TransactionHandler, int) (*model.CapturePageType, error)
	getCapturePageTypeByIdMutex       sync.RWMutex
	getCapturePageTypeByIdArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 int
	}
	getCapturePageTypeByIdReturns struct {
		result1 *model.CapturePageType
		result2 error
	}
	getCapturePageTypeByIdReturnsOnCall map[int]struct {
		result1 *model.CapturePageType
		result2 error
	}
	GetCapturePagesStub        func(context.Context, persistence.TransactionHandler, *model.CapturePagesFilters) (*model.PaginatedCapturePages, error)
	getCapturePagesMutex       sync.RWMutex
	getCapturePagesArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.CapturePagesFilters
	}
	getCapturePagesReturns struct {
		result1 *model.PaginatedCapturePages
		result2 error
	}
	getCapturePagesReturnsOnCall map[int]struct {
		result1 *model.PaginatedCapturePages
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePersistor) CreateCapturePages(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 *model.CapturePages) (*model.CapturePages, error) {
	fake.createCapturePagesMutex.Lock()
	ret, specificReturn := fake.createCapturePagesReturnsOnCall[len(fake.createCapturePagesArgsForCall)]
	fake.createCapturePagesArgsForCall = append(fake.createCapturePagesArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.CapturePages
	}{arg1, arg2, arg3})
	stub := fake.CreateCapturePagesStub
	fakeReturns := fake.createCapturePagesReturns
	fake.recordInvocation("CreateCapturePages", []interface{}{arg1, arg2, arg3})
	fake.createCapturePagesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePersistor) CreateCapturePagesCallCount() int {
	fake.createCapturePagesMutex.RLock()
	defer fake.createCapturePagesMutex.RUnlock()
	return len(fake.createCapturePagesArgsForCall)
}

func (fake *FakePersistor) CreateCapturePagesCalls(stub func(context.Context, persistence.TransactionHandler, *model.CapturePages) (*model.CapturePages, error)) {
	fake.createCapturePagesMutex.Lock()
	defer fake.createCapturePagesMutex.Unlock()
	fake.CreateCapturePagesStub = stub
}

func (fake *FakePersistor) CreateCapturePagesArgsForCall(i int) (context.Context, persistence.TransactionHandler, *model.CapturePages) {
	fake.createCapturePagesMutex.RLock()
	defer fake.createCapturePagesMutex.RUnlock()
	argsForCall := fake.createCapturePagesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) CreateCapturePagesReturns(result1 *model.CapturePages, result2 error) {
	fake.createCapturePagesMutex.Lock()
	defer fake.createCapturePagesMutex.Unlock()
	fake.CreateCapturePagesStub = nil
	fake.createCapturePagesReturns = struct {
		result1 *model.CapturePages
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) CreateCapturePagesReturnsOnCall(i int, result1 *model.CapturePages, result2 error) {
	fake.createCapturePagesMutex.Lock()
	defer fake.createCapturePagesMutex.Unlock()
	fake.CreateCapturePagesStub = nil
	if fake.createCapturePagesReturnsOnCall == nil {
		fake.createCapturePagesReturnsOnCall = make(map[int]struct {
			result1 *model.CapturePages
			result2 error
		})
	}
	fake.createCapturePagesReturnsOnCall[i] = struct {
		result1 *model.CapturePages
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) GetCapturePageTypeById(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 int) (*model.CapturePageType, error) {
	fake.getCapturePageTypeByIdMutex.Lock()
	ret, specificReturn := fake.getCapturePageTypeByIdReturnsOnCall[len(fake.getCapturePageTypeByIdArgsForCall)]
	fake.getCapturePageTypeByIdArgsForCall = append(fake.getCapturePageTypeByIdArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 int
	}{arg1, arg2, arg3})
	stub := fake.GetCapturePageTypeByIdStub
	fakeReturns := fake.getCapturePageTypeByIdReturns
	fake.recordInvocation("GetCapturePageTypeById", []interface{}{arg1, arg2, arg3})
	fake.getCapturePageTypeByIdMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePersistor) GetCapturePageTypeByIdCallCount() int {
	fake.getCapturePageTypeByIdMutex.RLock()
	defer fake.getCapturePageTypeByIdMutex.RUnlock()
	return len(fake.getCapturePageTypeByIdArgsForCall)
}

func (fake *FakePersistor) GetCapturePageTypeByIdCalls(stub func(context.Context, persistence.TransactionHandler, int) (*model.CapturePageType, error)) {
	fake.getCapturePageTypeByIdMutex.Lock()
	defer fake.getCapturePageTypeByIdMutex.Unlock()
	fake.GetCapturePageTypeByIdStub = stub
}

func (fake *FakePersistor) GetCapturePageTypeByIdArgsForCall(i int) (context.Context, persistence.TransactionHandler, int) {
	fake.getCapturePageTypeByIdMutex.RLock()
	defer fake.getCapturePageTypeByIdMutex.RUnlock()
	argsForCall := fake.getCapturePageTypeByIdArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) GetCapturePageTypeByIdReturns(result1 *model.CapturePageType, result2 error) {
	fake.getCapturePageTypeByIdMutex.Lock()
	defer fake.getCapturePageTypeByIdMutex.Unlock()
	fake.GetCapturePageTypeByIdStub = nil
	fake.getCapturePageTypeByIdReturns = struct {
		result1 *model.CapturePageType
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) GetCapturePageTypeByIdReturnsOnCall(i int, result1 *model.CapturePageType, result2 error) {
	fake.getCapturePageTypeByIdMutex.Lock()
	defer fake.getCapturePageTypeByIdMutex.Unlock()
	fake.GetCapturePageTypeByIdStub = nil
	if fake.getCapturePageTypeByIdReturnsOnCall == nil {
		fake.getCapturePageTypeByIdReturnsOnCall = make(map[int]struct {
			result1 *model.CapturePageType
			result2 error
		})
	}
	fake.getCapturePageTypeByIdReturnsOnCall[i] = struct {
		result1 *model.CapturePageType
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) GetCapturePages(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 *model.CapturePagesFilters) (*model.PaginatedCapturePages, error) {
	fake.getCapturePagesMutex.Lock()
	ret, specificReturn := fake.getCapturePagesReturnsOnCall[len(fake.getCapturePagesArgsForCall)]
	fake.getCapturePagesArgsForCall = append(fake.getCapturePagesArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.CapturePagesFilters
	}{arg1, arg2, arg3})
	stub := fake.GetCapturePagesStub
	fakeReturns := fake.getCapturePagesReturns
	fake.recordInvocation("GetCapturePages", []interface{}{arg1, arg2, arg3})
	fake.getCapturePagesMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePersistor) GetCapturePagesCallCount() int {
	fake.getCapturePagesMutex.RLock()
	defer fake.getCapturePagesMutex.RUnlock()
	return len(fake.getCapturePagesArgsForCall)
}

func (fake *FakePersistor) GetCapturePagesCalls(stub func(context.Context, persistence.TransactionHandler, *model.CapturePagesFilters) (*model.PaginatedCapturePages, error)) {
	fake.getCapturePagesMutex.Lock()
	defer fake.getCapturePagesMutex.Unlock()
	fake.GetCapturePagesStub = stub
}

func (fake *FakePersistor) GetCapturePagesArgsForCall(i int) (context.Context, persistence.TransactionHandler, *model.CapturePagesFilters) {
	fake.getCapturePagesMutex.RLock()
	defer fake.getCapturePagesMutex.RUnlock()
	argsForCall := fake.getCapturePagesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) GetCapturePagesReturns(result1 *model.PaginatedCapturePages, result2 error) {
	fake.getCapturePagesMutex.Lock()
	defer fake.getCapturePagesMutex.Unlock()
	fake.GetCapturePagesStub = nil
	fake.getCapturePagesReturns = struct {
		result1 *model.PaginatedCapturePages
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) GetCapturePagesReturnsOnCall(i int, result1 *model.PaginatedCapturePages, result2 error) {
	fake.getCapturePagesMutex.Lock()
	defer fake.getCapturePagesMutex.Unlock()
	fake.GetCapturePagesStub = nil
	if fake.getCapturePagesReturnsOnCall == nil {
		fake.getCapturePagesReturnsOnCall = make(map[int]struct {
			result1 *model.PaginatedCapturePages
			result2 error
		})
	}
	fake.getCapturePagesReturnsOnCall[i] = struct {
		result1 *model.PaginatedCapturePages
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createCapturePagesMutex.RLock()
	defer fake.createCapturePagesMutex.RUnlock()
	fake.getCapturePageTypeByIdMutex.RLock()
	defer fake.getCapturePageTypeByIdMutex.RUnlock()
	fake.getCapturePagesMutex.RLock()
	defer fake.getCapturePagesMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePersistor) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
