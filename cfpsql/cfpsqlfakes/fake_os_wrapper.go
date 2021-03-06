// Code generated by counterfeiter. DO NOT EDIT.
package cfpsqlfakes

import (
	"sync"

	"github.com/jaecktec/cf-psql-plugin/cfpsql"
)

type FakeOsWrapper struct {
	LookupEnvStub        func(key string) (string, bool)
	lookupEnvMutex       sync.RWMutex
	lookupEnvArgsForCall []struct {
		key string
	}
	lookupEnvReturns struct {
		result1 string
		result2 bool
	}
	lookupEnvReturnsOnCall map[int]struct {
		result1 string
		result2 bool
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeOsWrapper) LookupEnv(key string) (string, bool) {
	fake.lookupEnvMutex.Lock()
	ret, specificReturn := fake.lookupEnvReturnsOnCall[len(fake.lookupEnvArgsForCall)]
	fake.lookupEnvArgsForCall = append(fake.lookupEnvArgsForCall, struct {
		key string
	}{key})
	fake.recordInvocation("LookupEnv", []interface{}{key})
	fake.lookupEnvMutex.Unlock()
	if fake.LookupEnvStub != nil {
		return fake.LookupEnvStub(key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.lookupEnvReturns.result1, fake.lookupEnvReturns.result2
}

func (fake *FakeOsWrapper) LookupEnvCallCount() int {
	fake.lookupEnvMutex.RLock()
	defer fake.lookupEnvMutex.RUnlock()
	return len(fake.lookupEnvArgsForCall)
}

func (fake *FakeOsWrapper) LookupEnvArgsForCall(i int) string {
	fake.lookupEnvMutex.RLock()
	defer fake.lookupEnvMutex.RUnlock()
	return fake.lookupEnvArgsForCall[i].key
}

func (fake *FakeOsWrapper) LookupEnvReturns(result1 string, result2 bool) {
	fake.LookupEnvStub = nil
	fake.lookupEnvReturns = struct {
		result1 string
		result2 bool
	}{result1, result2}
}

func (fake *FakeOsWrapper) LookupEnvReturnsOnCall(i int, result1 string, result2 bool) {
	fake.LookupEnvStub = nil
	if fake.lookupEnvReturnsOnCall == nil {
		fake.lookupEnvReturnsOnCall = make(map[int]struct {
			result1 string
			result2 bool
		})
	}
	fake.lookupEnvReturnsOnCall[i] = struct {
		result1 string
		result2 bool
	}{result1, result2}
}

func (fake *FakeOsWrapper) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.lookupEnvMutex.RLock()
	defer fake.lookupEnvMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeOsWrapper) recordInvocation(key string, args []interface{}) {
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

var _ cfpsql.OsWrapper = new(FakeOsWrapper)
