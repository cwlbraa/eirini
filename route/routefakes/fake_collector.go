// Code generated by counterfeiter. DO NOT EDIT.
package routefakes

import (
	"sync"

	"code.cloudfoundry.org/eirini/route"
)

type FakeCollector struct {
	CollectStub        func() ([]route.Message, error)
	collectMutex       sync.RWMutex
	collectArgsForCall []struct {
	}
	collectReturns struct {
		result1 []route.Message
		result2 error
	}
	collectReturnsOnCall map[int]struct {
		result1 []route.Message
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeCollector) Collect() ([]route.Message, error) {
	fake.collectMutex.Lock()
	ret, specificReturn := fake.collectReturnsOnCall[len(fake.collectArgsForCall)]
	fake.collectArgsForCall = append(fake.collectArgsForCall, struct {
	}{})
	fake.recordInvocation("Collect", []interface{}{})
	fake.collectMutex.Unlock()
	if fake.CollectStub != nil {
		return fake.CollectStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.collectReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeCollector) CollectCallCount() int {
	fake.collectMutex.RLock()
	defer fake.collectMutex.RUnlock()
	return len(fake.collectArgsForCall)
}

func (fake *FakeCollector) CollectCalls(stub func() ([]route.Message, error)) {
	fake.collectMutex.Lock()
	defer fake.collectMutex.Unlock()
	fake.CollectStub = stub
}

func (fake *FakeCollector) CollectReturns(result1 []route.Message, result2 error) {
	fake.collectMutex.Lock()
	defer fake.collectMutex.Unlock()
	fake.CollectStub = nil
	fake.collectReturns = struct {
		result1 []route.Message
		result2 error
	}{result1, result2}
}

func (fake *FakeCollector) CollectReturnsOnCall(i int, result1 []route.Message, result2 error) {
	fake.collectMutex.Lock()
	defer fake.collectMutex.Unlock()
	fake.CollectStub = nil
	if fake.collectReturnsOnCall == nil {
		fake.collectReturnsOnCall = make(map[int]struct {
			result1 []route.Message
			result2 error
		})
	}
	fake.collectReturnsOnCall[i] = struct {
		result1 []route.Message
		result2 error
	}{result1, result2}
}

func (fake *FakeCollector) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.collectMutex.RLock()
	defer fake.collectMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeCollector) recordInvocation(key string, args []interface{}) {
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

var _ route.Collector = new(FakeCollector)
