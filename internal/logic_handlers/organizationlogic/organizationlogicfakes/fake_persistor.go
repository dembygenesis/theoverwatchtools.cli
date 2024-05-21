// Code generated by counterfeiter. DO NOT EDIT.
package organizationlogicfakes

import (
	"context"
	"sync"

	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
)

type FakePersistor struct {
	CreateOrganizationStub        func(context.Context, persistence.TransactionHandler, *model.Organization) (*model.Organization, error)
	createOrganizationMutex       sync.RWMutex
	createOrganizationArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.Organization
	}
	createOrganizationReturns struct {
		result1 *model.Organization
		result2 error
	}
	createOrganizationReturnsOnCall map[int]struct {
		result1 *model.Organization
		result2 error
	}
	DeleteOrganizationStub        func(context.Context, persistence.TransactionHandler, int) error
	deleteOrganizationMutex       sync.RWMutex
	deleteOrganizationArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 int
	}
	deleteOrganizationReturns struct {
		result1 error
	}
	deleteOrganizationReturnsOnCall map[int]struct {
		result1 error
	}
	GetOrganizationByNameStub        func(context.Context, persistence.TransactionHandler, string) (*model.Organization, error)
	getOrganizationByNameMutex       sync.RWMutex
	getOrganizationByNameArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 string
	}
	getOrganizationByNameReturns struct {
		result1 *model.Organization
		result2 error
	}
	getOrganizationByNameReturnsOnCall map[int]struct {
		result1 *model.Organization
		result2 error
	}
	GetOrganizationTypeByIdStub        func(context.Context, persistence.TransactionHandler, int) (*model.OrganizationType, error)
	getOrganizationTypeByIdMutex       sync.RWMutex
	getOrganizationTypeByIdArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 int
	}
	getOrganizationTypeByIdReturns struct {
		result1 *model.OrganizationType
		result2 error
	}
	getOrganizationTypeByIdReturnsOnCall map[int]struct {
		result1 *model.OrganizationType
		result2 error
	}
	GetOrganizationsStub        func(context.Context, persistence.TransactionHandler, *model.OrganizationFilters) (*model.PaginatedOrganizations, error)
	getOrganizationsMutex       sync.RWMutex
	getOrganizationsArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.OrganizationFilters
	}
	getOrganizationsReturns struct {
		result1 *model.PaginatedOrganizations
		result2 error
	}
	getOrganizationsReturnsOnCall map[int]struct {
		result1 *model.PaginatedOrganizations
		result2 error
	}
	RestoreOrganizationStub        func(context.Context, persistence.TransactionHandler, int) error
	restoreOrganizationMutex       sync.RWMutex
	restoreOrganizationArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 int
	}
	restoreOrganizationReturns struct {
		result1 error
	}
	restoreOrganizationReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePersistor) CreateOrganization(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 *model.Organization) (*model.Organization, error) {
	fake.createOrganizationMutex.Lock()
	ret, specificReturn := fake.createOrganizationReturnsOnCall[len(fake.createOrganizationArgsForCall)]
	fake.createOrganizationArgsForCall = append(fake.createOrganizationArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.Organization
	}{arg1, arg2, arg3})
	stub := fake.CreateOrganizationStub
	fakeReturns := fake.createOrganizationReturns
	fake.recordInvocation("CreateOrganization", []interface{}{arg1, arg2, arg3})
	fake.createOrganizationMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePersistor) CreateOrganizationCallCount() int {
	fake.createOrganizationMutex.RLock()
	defer fake.createOrganizationMutex.RUnlock()
	return len(fake.createOrganizationArgsForCall)
}

func (fake *FakePersistor) CreateOrganizationCalls(stub func(context.Context, persistence.TransactionHandler, *model.Organization) (*model.Organization, error)) {
	fake.createOrganizationMutex.Lock()
	defer fake.createOrganizationMutex.Unlock()
	fake.CreateOrganizationStub = stub
}

func (fake *FakePersistor) CreateOrganizationArgsForCall(i int) (context.Context, persistence.TransactionHandler, *model.Organization) {
	fake.createOrganizationMutex.RLock()
	defer fake.createOrganizationMutex.RUnlock()
	argsForCall := fake.createOrganizationArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) CreateOrganizationReturns(result1 *model.Organization, result2 error) {
	fake.createOrganizationMutex.Lock()
	defer fake.createOrganizationMutex.Unlock()
	fake.CreateOrganizationStub = nil
	fake.createOrganizationReturns = struct {
		result1 *model.Organization
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) CreateOrganizationReturnsOnCall(i int, result1 *model.Organization, result2 error) {
	fake.createOrganizationMutex.Lock()
	defer fake.createOrganizationMutex.Unlock()
	fake.CreateOrganizationStub = nil
	if fake.createOrganizationReturnsOnCall == nil {
		fake.createOrganizationReturnsOnCall = make(map[int]struct {
			result1 *model.Organization
			result2 error
		})
	}
	fake.createOrganizationReturnsOnCall[i] = struct {
		result1 *model.Organization
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) DeleteOrganization(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 int) error {
	fake.deleteOrganizationMutex.Lock()
	ret, specificReturn := fake.deleteOrganizationReturnsOnCall[len(fake.deleteOrganizationArgsForCall)]
	fake.deleteOrganizationArgsForCall = append(fake.deleteOrganizationArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 int
	}{arg1, arg2, arg3})
	stub := fake.DeleteOrganizationStub
	fakeReturns := fake.deleteOrganizationReturns
	fake.recordInvocation("DeleteOrganization", []interface{}{arg1, arg2, arg3})
	fake.deleteOrganizationMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakePersistor) DeleteOrganizationCallCount() int {
	fake.deleteOrganizationMutex.RLock()
	defer fake.deleteOrganizationMutex.RUnlock()
	return len(fake.deleteOrganizationArgsForCall)
}

func (fake *FakePersistor) DeleteOrganizationCalls(stub func(context.Context, persistence.TransactionHandler, int) error) {
	fake.deleteOrganizationMutex.Lock()
	defer fake.deleteOrganizationMutex.Unlock()
	fake.DeleteOrganizationStub = stub
}

func (fake *FakePersistor) DeleteOrganizationArgsForCall(i int) (context.Context, persistence.TransactionHandler, int) {
	fake.deleteOrganizationMutex.RLock()
	defer fake.deleteOrganizationMutex.RUnlock()
	argsForCall := fake.deleteOrganizationArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) DeleteOrganizationReturns(result1 error) {
	fake.deleteOrganizationMutex.Lock()
	defer fake.deleteOrganizationMutex.Unlock()
	fake.DeleteOrganizationStub = nil
	fake.deleteOrganizationReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePersistor) DeleteOrganizationReturnsOnCall(i int, result1 error) {
	fake.deleteOrganizationMutex.Lock()
	defer fake.deleteOrganizationMutex.Unlock()
	fake.DeleteOrganizationStub = nil
	if fake.deleteOrganizationReturnsOnCall == nil {
		fake.deleteOrganizationReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteOrganizationReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePersistor) GetOrganizationByName(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 string) (*model.Organization, error) {
	fake.getOrganizationByNameMutex.Lock()
	ret, specificReturn := fake.getOrganizationByNameReturnsOnCall[len(fake.getOrganizationByNameArgsForCall)]
	fake.getOrganizationByNameArgsForCall = append(fake.getOrganizationByNameArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetOrganizationByNameStub
	fakeReturns := fake.getOrganizationByNameReturns
	fake.recordInvocation("GetOrganizationByName", []interface{}{arg1, arg2, arg3})
	fake.getOrganizationByNameMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePersistor) GetOrganizationByNameCallCount() int {
	fake.getOrganizationByNameMutex.RLock()
	defer fake.getOrganizationByNameMutex.RUnlock()
	return len(fake.getOrganizationByNameArgsForCall)
}

func (fake *FakePersistor) GetOrganizationByNameCalls(stub func(context.Context, persistence.TransactionHandler, string) (*model.Organization, error)) {
	fake.getOrganizationByNameMutex.Lock()
	defer fake.getOrganizationByNameMutex.Unlock()
	fake.GetOrganizationByNameStub = stub
}

func (fake *FakePersistor) GetOrganizationByNameArgsForCall(i int) (context.Context, persistence.TransactionHandler, string) {
	fake.getOrganizationByNameMutex.RLock()
	defer fake.getOrganizationByNameMutex.RUnlock()
	argsForCall := fake.getOrganizationByNameArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) GetOrganizationByNameReturns(result1 *model.Organization, result2 error) {
	fake.getOrganizationByNameMutex.Lock()
	defer fake.getOrganizationByNameMutex.Unlock()
	fake.GetOrganizationByNameStub = nil
	fake.getOrganizationByNameReturns = struct {
		result1 *model.Organization
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) GetOrganizationByNameReturnsOnCall(i int, result1 *model.Organization, result2 error) {
	fake.getOrganizationByNameMutex.Lock()
	defer fake.getOrganizationByNameMutex.Unlock()
	fake.GetOrganizationByNameStub = nil
	if fake.getOrganizationByNameReturnsOnCall == nil {
		fake.getOrganizationByNameReturnsOnCall = make(map[int]struct {
			result1 *model.Organization
			result2 error
		})
	}
	fake.getOrganizationByNameReturnsOnCall[i] = struct {
		result1 *model.Organization
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) GetOrganizationTypeById(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 int) (*model.OrganizationType, error) {
	fake.getOrganizationTypeByIdMutex.Lock()
	ret, specificReturn := fake.getOrganizationTypeByIdReturnsOnCall[len(fake.getOrganizationTypeByIdArgsForCall)]
	fake.getOrganizationTypeByIdArgsForCall = append(fake.getOrganizationTypeByIdArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 int
	}{arg1, arg2, arg3})
	stub := fake.GetOrganizationTypeByIdStub
	fakeReturns := fake.getOrganizationTypeByIdReturns
	fake.recordInvocation("GetOrganizationTypeById", []interface{}{arg1, arg2, arg3})
	fake.getOrganizationTypeByIdMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePersistor) GetOrganizationTypeByIdCallCount() int {
	fake.getOrganizationTypeByIdMutex.RLock()
	defer fake.getOrganizationTypeByIdMutex.RUnlock()
	return len(fake.getOrganizationTypeByIdArgsForCall)
}

func (fake *FakePersistor) GetOrganizationTypeByIdCalls(stub func(context.Context, persistence.TransactionHandler, int) (*model.OrganizationType, error)) {
	fake.getOrganizationTypeByIdMutex.Lock()
	defer fake.getOrganizationTypeByIdMutex.Unlock()
	fake.GetOrganizationTypeByIdStub = stub
}

func (fake *FakePersistor) GetOrganizationTypeByIdArgsForCall(i int) (context.Context, persistence.TransactionHandler, int) {
	fake.getOrganizationTypeByIdMutex.RLock()
	defer fake.getOrganizationTypeByIdMutex.RUnlock()
	argsForCall := fake.getOrganizationTypeByIdArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) GetOrganizationTypeByIdReturns(result1 *model.OrganizationType, result2 error) {
	fake.getOrganizationTypeByIdMutex.Lock()
	defer fake.getOrganizationTypeByIdMutex.Unlock()
	fake.GetOrganizationTypeByIdStub = nil
	fake.getOrganizationTypeByIdReturns = struct {
		result1 *model.OrganizationType
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) GetOrganizationTypeByIdReturnsOnCall(i int, result1 *model.OrganizationType, result2 error) {
	fake.getOrganizationTypeByIdMutex.Lock()
	defer fake.getOrganizationTypeByIdMutex.Unlock()
	fake.GetOrganizationTypeByIdStub = nil
	if fake.getOrganizationTypeByIdReturnsOnCall == nil {
		fake.getOrganizationTypeByIdReturnsOnCall = make(map[int]struct {
			result1 *model.OrganizationType
			result2 error
		})
	}
	fake.getOrganizationTypeByIdReturnsOnCall[i] = struct {
		result1 *model.OrganizationType
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) GetOrganizations(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 *model.OrganizationFilters) (*model.PaginatedOrganizations, error) {
	fake.getOrganizationsMutex.Lock()
	ret, specificReturn := fake.getOrganizationsReturnsOnCall[len(fake.getOrganizationsArgsForCall)]
	fake.getOrganizationsArgsForCall = append(fake.getOrganizationsArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.OrganizationFilters
	}{arg1, arg2, arg3})
	stub := fake.GetOrganizationsStub
	fakeReturns := fake.getOrganizationsReturns
	fake.recordInvocation("GetOrganizations", []interface{}{arg1, arg2, arg3})
	fake.getOrganizationsMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePersistor) GetOrganizationsCallCount() int {
	fake.getOrganizationsMutex.RLock()
	defer fake.getOrganizationsMutex.RUnlock()
	return len(fake.getOrganizationsArgsForCall)
}

func (fake *FakePersistor) GetOrganizationsCalls(stub func(context.Context, persistence.TransactionHandler, *model.OrganizationFilters) (*model.PaginatedOrganizations, error)) {
	fake.getOrganizationsMutex.Lock()
	defer fake.getOrganizationsMutex.Unlock()
	fake.GetOrganizationsStub = stub
}

func (fake *FakePersistor) GetOrganizationsArgsForCall(i int) (context.Context, persistence.TransactionHandler, *model.OrganizationFilters) {
	fake.getOrganizationsMutex.RLock()
	defer fake.getOrganizationsMutex.RUnlock()
	argsForCall := fake.getOrganizationsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) GetOrganizationsReturns(result1 *model.PaginatedOrganizations, result2 error) {
	fake.getOrganizationsMutex.Lock()
	defer fake.getOrganizationsMutex.Unlock()
	fake.GetOrganizationsStub = nil
	fake.getOrganizationsReturns = struct {
		result1 *model.PaginatedOrganizations
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) GetOrganizationsReturnsOnCall(i int, result1 *model.PaginatedOrganizations, result2 error) {
	fake.getOrganizationsMutex.Lock()
	defer fake.getOrganizationsMutex.Unlock()
	fake.GetOrganizationsStub = nil
	if fake.getOrganizationsReturnsOnCall == nil {
		fake.getOrganizationsReturnsOnCall = make(map[int]struct {
			result1 *model.PaginatedOrganizations
			result2 error
		})
	}
	fake.getOrganizationsReturnsOnCall[i] = struct {
		result1 *model.PaginatedOrganizations
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) RestoreOrganization(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 int) error {
	fake.restoreOrganizationMutex.Lock()
	ret, specificReturn := fake.restoreOrganizationReturnsOnCall[len(fake.restoreOrganizationArgsForCall)]
	fake.restoreOrganizationArgsForCall = append(fake.restoreOrganizationArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 int
	}{arg1, arg2, arg3})
	stub := fake.RestoreOrganizationStub
	fakeReturns := fake.restoreOrganizationReturns
	fake.recordInvocation("RestoreOrganization", []interface{}{arg1, arg2, arg3})
	fake.restoreOrganizationMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakePersistor) RestoreOrganizationCallCount() int {
	fake.restoreOrganizationMutex.RLock()
	defer fake.restoreOrganizationMutex.RUnlock()
	return len(fake.restoreOrganizationArgsForCall)
}

func (fake *FakePersistor) RestoreOrganizationCalls(stub func(context.Context, persistence.TransactionHandler, int) error) {
	fake.restoreOrganizationMutex.Lock()
	defer fake.restoreOrganizationMutex.Unlock()
	fake.RestoreOrganizationStub = stub
}

func (fake *FakePersistor) RestoreOrganizationArgsForCall(i int) (context.Context, persistence.TransactionHandler, int) {
	fake.restoreOrganizationMutex.RLock()
	defer fake.restoreOrganizationMutex.RUnlock()
	argsForCall := fake.restoreOrganizationArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) RestoreOrganizationReturns(result1 error) {
	fake.restoreOrganizationMutex.Lock()
	defer fake.restoreOrganizationMutex.Unlock()
	fake.RestoreOrganizationStub = nil
	fake.restoreOrganizationReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePersistor) RestoreOrganizationReturnsOnCall(i int, result1 error) {
	fake.restoreOrganizationMutex.Lock()
	defer fake.restoreOrganizationMutex.Unlock()
	fake.RestoreOrganizationStub = nil
	if fake.restoreOrganizationReturnsOnCall == nil {
		fake.restoreOrganizationReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.restoreOrganizationReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePersistor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.createOrganizationMutex.RLock()
	defer fake.createOrganizationMutex.RUnlock()
	fake.deleteOrganizationMutex.RLock()
	defer fake.deleteOrganizationMutex.RUnlock()
	fake.getOrganizationByNameMutex.RLock()
	defer fake.getOrganizationByNameMutex.RUnlock()
	fake.getOrganizationTypeByIdMutex.RLock()
	defer fake.getOrganizationTypeByIdMutex.RUnlock()
	fake.getOrganizationsMutex.RLock()
	defer fake.getOrganizationsMutex.RUnlock()
	fake.restoreOrganizationMutex.RLock()
	defer fake.restoreOrganizationMutex.RUnlock()
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