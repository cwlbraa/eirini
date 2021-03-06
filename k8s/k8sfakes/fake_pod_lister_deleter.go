// Code generated by counterfeiter. DO NOT EDIT.
package k8sfakes

import (
	"sync"

	"code.cloudfoundry.org/eirini/k8s"
	v1a "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FakePodListerDeleter struct {
	DeleteStub        func(string, *v1.DeleteOptions) error
	deleteMutex       sync.RWMutex
	deleteArgsForCall []struct {
		arg1 string
		arg2 *v1.DeleteOptions
	}
	deleteReturns struct {
		result1 error
	}
	deleteReturnsOnCall map[int]struct {
		result1 error
	}
	ListStub        func(v1.ListOptions) (*v1a.PodList, error)
	listMutex       sync.RWMutex
	listArgsForCall []struct {
		arg1 v1.ListOptions
	}
	listReturns struct {
		result1 *v1a.PodList
		result2 error
	}
	listReturnsOnCall map[int]struct {
		result1 *v1a.PodList
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePodListerDeleter) Delete(arg1 string, arg2 *v1.DeleteOptions) error {
	fake.deleteMutex.Lock()
	ret, specificReturn := fake.deleteReturnsOnCall[len(fake.deleteArgsForCall)]
	fake.deleteArgsForCall = append(fake.deleteArgsForCall, struct {
		arg1 string
		arg2 *v1.DeleteOptions
	}{arg1, arg2})
	fake.recordInvocation("Delete", []interface{}{arg1, arg2})
	fake.deleteMutex.Unlock()
	if fake.DeleteStub != nil {
		return fake.DeleteStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deleteReturns
	return fakeReturns.result1
}

func (fake *FakePodListerDeleter) DeleteCallCount() int {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	return len(fake.deleteArgsForCall)
}

func (fake *FakePodListerDeleter) DeleteCalls(stub func(string, *v1.DeleteOptions) error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = stub
}

func (fake *FakePodListerDeleter) DeleteArgsForCall(i int) (string, *v1.DeleteOptions) {
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	argsForCall := fake.deleteArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakePodListerDeleter) DeleteReturns(result1 error) {
	fake.deleteMutex.Lock()
	defer fake.deleteMutex.Unlock()
	fake.DeleteStub = nil
	fake.deleteReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePodListerDeleter) DeleteReturnsOnCall(i int, result1 error) {
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

func (fake *FakePodListerDeleter) List(arg1 v1.ListOptions) (*v1a.PodList, error) {
	fake.listMutex.Lock()
	ret, specificReturn := fake.listReturnsOnCall[len(fake.listArgsForCall)]
	fake.listArgsForCall = append(fake.listArgsForCall, struct {
		arg1 v1.ListOptions
	}{arg1})
	fake.recordInvocation("List", []interface{}{arg1})
	fake.listMutex.Unlock()
	if fake.ListStub != nil {
		return fake.ListStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePodListerDeleter) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *FakePodListerDeleter) ListCalls(stub func(v1.ListOptions) (*v1a.PodList, error)) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = stub
}

func (fake *FakePodListerDeleter) ListArgsForCall(i int) v1.ListOptions {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	argsForCall := fake.listArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakePodListerDeleter) ListReturns(result1 *v1a.PodList, result2 error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 *v1a.PodList
		result2 error
	}{result1, result2}
}

func (fake *FakePodListerDeleter) ListReturnsOnCall(i int, result1 *v1a.PodList, result2 error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = nil
	if fake.listReturnsOnCall == nil {
		fake.listReturnsOnCall = make(map[int]struct {
			result1 *v1a.PodList
			result2 error
		})
	}
	fake.listReturnsOnCall[i] = struct {
		result1 *v1a.PodList
		result2 error
	}{result1, result2}
}

func (fake *FakePodListerDeleter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deleteMutex.RLock()
	defer fake.deleteMutex.RUnlock()
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePodListerDeleter) recordInvocation(key string, args []interface{}) {
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

var _ k8s.PodListerDeleter = new(FakePodListerDeleter)
