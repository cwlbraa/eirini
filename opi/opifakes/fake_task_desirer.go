// Code generated by counterfeiter. DO NOT EDIT.
package opifakes

import (
	"sync"

	"code.cloudfoundry.org/eirini/opi"
)

type FakeTaskDesirer struct {
	DeleteStub        func(string) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 string
	}
	deleteReturns struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	DesireStub        func(*opi.Task) error
	desireMutex       sync.RWMutex
	desireArgsForCall []struct {
		arg1 *opi.Task
	}
	desireReturns struct {
		result1 error
	}
	desireReturnsOnCall map[int]struct {
		result1 error
	}
	DesireStagingStub        func(*opi.StagingTask) error
	desireStagingMutex       sync.RWMutex
	desireStagingArgsForCall []struct {
		arg1 *opi.StagingTask
	}
	desireStagingReturns struct {
		result1 error
	}
	desireStagingReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeTaskDesirer) Delete(arg1 string) error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Delete", []interface{}{arg1})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteReturns
	return fakeReturns.result1
}

func (fake *FakeTaskDesirer) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakeTaskDesirer) DeleteCalls(stub func(string) error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *FakeTaskDesirer) DeleteArgsForCall(i int) string {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTaskDesirer) DeleteReturns(result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTaskDesirer) DeleteReturnsOnCall(i int, result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	if fake.deleteReturnsOnCall == nil {
		fake.deleteReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTaskDesirer) Desire(arg1 *opi.Task) error {
	fake.desireMutex.Lock()
	ret, specificReturn := fake.desireReturnsOnCall[len(fake.desireArgsForCall)]
	fake.desireArgsForCall = append(fake.desireArgsForCall, struct {
		arg1 *opi.Task
	}{arg1})
	fake.recordInvocation("Desire", []interface{}{arg1})
	fake.desireMutex.Unlock()
	if fake.DesireStub != nil {
		return fake.DesireStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.desireReturns
	return fakeReturns.result1
}

func (fake *FakeTaskDesirer) DesireCallCount() int {
	fake.desireMutex.RLock()
	defer fake.desireMutex.RUnlock()
	return len(fake.desireArgsForCall)
}

func (fake *FakeTaskDesirer) DesireCalls(stub func(*opi.Task) error) {
	fake.desireMutex.Lock()
	defer fake.desireMutex.Unlock()
	fake.DesireStub = stub
}

func (fake *FakeTaskDesirer) DesireArgsForCall(i int) *opi.Task {
	fake.desireMutex.RLock()
	defer fake.desireMutex.RUnlock()
	argsForCall := fake.desireArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTaskDesirer) DesireReturns(result1 error) {
	fake.desireMutex.Lock()
	defer fake.desireMutex.Unlock()
	fake.DesireStub = nil
	fake.desireReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTaskDesirer) DesireReturnsOnCall(i int, result1 error) {
	fake.desireMutex.Lock()
	defer fake.desireMutex.Unlock()
	fake.DesireStub = nil
	if fake.desireReturnsOnCall == nil {
		fake.desireReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.desireReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTaskDesirer) DesireStaging(arg1 *opi.StagingTask) error {
	fake.desireStagingMutex.Lock()
	ret, specificReturn := fake.desireStagingReturnsOnCall[len(fake.desireStagingArgsForCall)]
	fake.desireStagingArgsForCall = append(fake.desireStagingArgsForCall, struct {
		arg1 *opi.StagingTask
	}{arg1})
	fake.recordInvocation("DesireStaging", []interface{}{arg1})
	fake.desireStagingMutex.Unlock()
	if fake.DesireStagingStub != nil {
		return fake.DesireStagingStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.desireStagingReturns
	return fakeReturns.result1
}

func (fake *FakeTaskDesirer) DesireStagingCallCount() int {
	fake.desireStagingMutex.RLock()
	defer fake.desireStagingMutex.RUnlock()
	return len(fake.desireStagingArgsForCall)
}

func (fake *FakeTaskDesirer) DesireStagingCalls(stub func(*opi.StagingTask) error) {
	fake.desireStagingMutex.Lock()
	defer fake.desireStagingMutex.Unlock()
	fake.DesireStagingStub = stub
}

func (fake *FakeTaskDesirer) DesireStagingArgsForCall(i int) *opi.StagingTask {
	fake.desireStagingMutex.RLock()
	defer fake.desireStagingMutex.RUnlock()
	argsForCall := fake.desireStagingArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeTaskDesirer) DesireStagingReturns(result1 error) {
	fake.desireStagingMutex.Lock()
	defer fake.desireStagingMutex.Unlock()
	fake.DesireStagingStub = nil
	fake.desireStagingReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeTaskDesirer) DesireStagingReturnsOnCall(i int, result1 error) {
	fake.desireStagingMutex.Lock()
	defer fake.desireStagingMutex.Unlock()
	fake.DesireStagingStub = nil
	if fake.desireStagingReturnsOnCall == nil {
		fake.desireStagingReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.desireStagingReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeTaskDesirer) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.desireMutex.RLock()
	defer fake.desireMutex.RUnlock()
	fake.desireStagingMutex.RLock()
	defer fake.desireStagingMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeTaskDesirer) recordInvocation(key string, args []interface{}) {
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

var _ opi.TaskDesirer = new(FakeTaskDesirer)
