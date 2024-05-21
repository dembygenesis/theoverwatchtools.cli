// Code generated by counterfeiter. DO NOT EDIT.
package apifakes

import (
	"context"
	"sync"

	"github.com/dembygenesis/local.tools/internal/model"
)

type FakeUserService struct {
	ListUsersStub        func(context.Context, *model.UserFilters) (*model.PaginatedUsers, error)
	listUsersMutex       sync.RWMutex
	listUsersArgsForCall []struct {
		arg1 context.Context
		arg2 *model.UserFilters
	}
	listUsersReturns struct {
		result1 *model.PaginatedUsers
		result2 error
	}
	listUsersReturnsOnCall map[int]struct {
		result1 *model.PaginatedUsers
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeUserService) ListUsers(arg1 context.Context, arg2 *model.UserFilters) (*model.PaginatedUsers, error) {
	fake.listUsersMutex.Lock()
	ret, specificReturn := fake.listUsersReturnsOnCall[len(fake.listUsersArgsForCall)]
	fake.listUsersArgsForCall = append(fake.listUsersArgsForCall, struct {
		arg1 context.Context
		arg2 *model.UserFilters
	}{arg1, arg2})
	stub := fake.ListUsersStub
	fakeReturns := fake.listUsersReturns
	fake.recordInvocation("ListUsers", []interface{}{arg1, arg2})
	fake.listUsersMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeUserService) ListUsersCallCount() int {
	fake.listUsersMutex.RLock()
	defer fake.listUsersMutex.RUnlock()
	return len(fake.listUsersArgsForCall)
}

func (fake *FakeUserService) ListUsersCalls(stub func(context.Context, *model.UserFilters) (*model.PaginatedUsers, error)) {
	fake.listUsersMutex.Lock()
	defer fake.listUsersMutex.Unlock()
	fake.ListUsersStub = stub
}

func (fake *FakeUserService) ListUsersArgsForCall(i int) (context.Context, *model.UserFilters) {
	fake.listUsersMutex.RLock()
	defer fake.listUsersMutex.RUnlock()
	argsForCall := fake.listUsersArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeUserService) ListUsersReturns(result1 *model.PaginatedUsers, result2 error) {
	fake.listUsersMutex.Lock()
	defer fake.listUsersMutex.Unlock()
	fake.ListUsersStub = nil
	fake.listUsersReturns = struct {
		result1 *model.PaginatedUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeUserService) ListUsersReturnsOnCall(i int, result1 *model.PaginatedUsers, result2 error) {
	fake.listUsersMutex.Lock()
	defer fake.listUsersMutex.Unlock()
	fake.ListUsersStub = nil
	if fake.listUsersReturnsOnCall == nil {
		fake.listUsersReturnsOnCall = make(map[int]struct {
			result1 *model.PaginatedUsers
			result2 error
		})
	}
	fake.listUsersReturnsOnCall[i] = struct {
		result1 *model.PaginatedUsers
		result2 error
	}{result1, result2}
}

func (fake *FakeUserService) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listUsersMutex.RLock()
	defer fake.listUsersMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeUserService) recordInvocation(key string, args []interface{}) {
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