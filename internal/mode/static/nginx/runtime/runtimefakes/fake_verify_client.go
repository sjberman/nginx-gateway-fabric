// Code generated by counterfeiter. DO NOT EDIT.
package runtimefakes

import (
	"context"
	"sync"

	"github.com/nginxinc/nginx-gateway-fabric/internal/mode/static/nginx/runtime"
)

type FakeVerifyClient struct {
	EnsureConfigVersionStub        func(context.Context, int) error
	ensureConfigVersionMutex       sync.RWMutex
	ensureConfigVersionArgsForCall []struct {
		arg1 context.Context
		arg2 int
	}
	ensureConfigVersionReturns struct {
		result1 error
	}
	ensureConfigVersionReturnsOnCall map[int]struct {
		result1 error
	}
	GetConfigVersionStub        func() (int, error)
	getConfigVersionMutex       sync.RWMutex
	getConfigVersionArgsForCall []struct {
	}
	getConfigVersionReturns struct {
		result1 int
		result2 error
	}
	getConfigVersionReturnsOnCall map[int]struct {
		result1 int
		result2 error
	}
	WaitForCorrectVersionStub        func(context.Context, int, string, []byte, runtime.ReadFileFunc) error
	waitForCorrectVersionMutex       sync.RWMutex
	waitForCorrectVersionArgsForCall []struct {
		arg1 context.Context
		arg2 int
		arg3 string
		arg4 []byte
		arg5 runtime.ReadFileFunc
	}
	waitForCorrectVersionReturns struct {
		result1 error
	}
	waitForCorrectVersionReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeVerifyClient) EnsureConfigVersion(arg1 context.Context, arg2 int) error {
	fake.ensureConfigVersionMutex.Lock()
	ret, specificReturn := fake.ensureConfigVersionReturnsOnCall[len(fake.ensureConfigVersionArgsForCall)]
	fake.ensureConfigVersionArgsForCall = append(fake.ensureConfigVersionArgsForCall, struct {
		arg1 context.Context
		arg2 int
	}{arg1, arg2})
	stub := fake.EnsureConfigVersionStub
	fakeReturns := fake.ensureConfigVersionReturns
	fake.recordInvocation("EnsureConfigVersion", []interface{}{arg1, arg2})
	fake.ensureConfigVersionMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVerifyClient) EnsureConfigVersionCallCount() int {
	fake.ensureConfigVersionMutex.RLock()
	defer fake.ensureConfigVersionMutex.RUnlock()
	return len(fake.ensureConfigVersionArgsForCall)
}

func (fake *FakeVerifyClient) EnsureConfigVersionCalls(stub func(context.Context, int) error) {
	fake.ensureConfigVersionMutex.Lock()
	defer fake.ensureConfigVersionMutex.Unlock()
	fake.EnsureConfigVersionStub = stub
}

func (fake *FakeVerifyClient) EnsureConfigVersionArgsForCall(i int) (context.Context, int) {
	fake.ensureConfigVersionMutex.RLock()
	defer fake.ensureConfigVersionMutex.RUnlock()
	argsForCall := fake.ensureConfigVersionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeVerifyClient) EnsureConfigVersionReturns(result1 error) {
	fake.ensureConfigVersionMutex.Lock()
	defer fake.ensureConfigVersionMutex.Unlock()
	fake.EnsureConfigVersionStub = nil
	fake.ensureConfigVersionReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeVerifyClient) EnsureConfigVersionReturnsOnCall(i int, result1 error) {
	fake.ensureConfigVersionMutex.Lock()
	defer fake.ensureConfigVersionMutex.Unlock()
	fake.EnsureConfigVersionStub = nil
	if fake.ensureConfigVersionReturnsOnCall == nil {
		fake.ensureConfigVersionReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.ensureConfigVersionReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeVerifyClient) GetConfigVersion() (int, error) {
	fake.getConfigVersionMutex.Lock()
	ret, specificReturn := fake.getConfigVersionReturnsOnCall[len(fake.getConfigVersionArgsForCall)]
	fake.getConfigVersionArgsForCall = append(fake.getConfigVersionArgsForCall, struct {
	}{})
	stub := fake.GetConfigVersionStub
	fakeReturns := fake.getConfigVersionReturns
	fake.recordInvocation("GetConfigVersion", []interface{}{})
	fake.getConfigVersionMutex.Unlock()
	if stub != nil {
		return stub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeVerifyClient) GetConfigVersionCallCount() int {
	fake.getConfigVersionMutex.RLock()
	defer fake.getConfigVersionMutex.RUnlock()
	return len(fake.getConfigVersionArgsForCall)
}

func (fake *FakeVerifyClient) GetConfigVersionCalls(stub func() (int, error)) {
	fake.getConfigVersionMutex.Lock()
	defer fake.getConfigVersionMutex.Unlock()
	fake.GetConfigVersionStub = stub
}

func (fake *FakeVerifyClient) GetConfigVersionReturns(result1 int, result2 error) {
	fake.getConfigVersionMutex.Lock()
	defer fake.getConfigVersionMutex.Unlock()
	fake.GetConfigVersionStub = nil
	fake.getConfigVersionReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeVerifyClient) GetConfigVersionReturnsOnCall(i int, result1 int, result2 error) {
	fake.getConfigVersionMutex.Lock()
	defer fake.getConfigVersionMutex.Unlock()
	fake.GetConfigVersionStub = nil
	if fake.getConfigVersionReturnsOnCall == nil {
		fake.getConfigVersionReturnsOnCall = make(map[int]struct {
			result1 int
			result2 error
		})
	}
	fake.getConfigVersionReturnsOnCall[i] = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakeVerifyClient) WaitForCorrectVersion(arg1 context.Context, arg2 int, arg3 string, arg4 []byte, arg5 runtime.ReadFileFunc) error {
	var arg4Copy []byte
	if arg4 != nil {
		arg4Copy = make([]byte, len(arg4))
		copy(arg4Copy, arg4)
	}
	fake.waitForCorrectVersionMutex.Lock()
	ret, specificReturn := fake.waitForCorrectVersionReturnsOnCall[len(fake.waitForCorrectVersionArgsForCall)]
	fake.waitForCorrectVersionArgsForCall = append(fake.waitForCorrectVersionArgsForCall, struct {
		arg1 context.Context
		arg2 int
		arg3 string
		arg4 []byte
		arg5 runtime.ReadFileFunc
	}{arg1, arg2, arg3, arg4Copy, arg5})
	stub := fake.WaitForCorrectVersionStub
	fakeReturns := fake.waitForCorrectVersionReturns
	fake.recordInvocation("WaitForCorrectVersion", []interface{}{arg1, arg2, arg3, arg4Copy, arg5})
	fake.waitForCorrectVersionMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeVerifyClient) WaitForCorrectVersionCallCount() int {
	fake.waitForCorrectVersionMutex.RLock()
	defer fake.waitForCorrectVersionMutex.RUnlock()
	return len(fake.waitForCorrectVersionArgsForCall)
}

func (fake *FakeVerifyClient) WaitForCorrectVersionCalls(stub func(context.Context, int, string, []byte, runtime.ReadFileFunc) error) {
	fake.waitForCorrectVersionMutex.Lock()
	defer fake.waitForCorrectVersionMutex.Unlock()
	fake.WaitForCorrectVersionStub = stub
}

func (fake *FakeVerifyClient) WaitForCorrectVersionArgsForCall(i int) (context.Context, int, string, []byte, runtime.ReadFileFunc) {
	fake.waitForCorrectVersionMutex.RLock()
	defer fake.waitForCorrectVersionMutex.RUnlock()
	argsForCall := fake.waitForCorrectVersionArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeVerifyClient) WaitForCorrectVersionReturns(result1 error) {
	fake.waitForCorrectVersionMutex.Lock()
	defer fake.waitForCorrectVersionMutex.Unlock()
	fake.WaitForCorrectVersionStub = nil
	fake.waitForCorrectVersionReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeVerifyClient) WaitForCorrectVersionReturnsOnCall(i int, result1 error) {
	fake.waitForCorrectVersionMutex.Lock()
	defer fake.waitForCorrectVersionMutex.Unlock()
	fake.WaitForCorrectVersionStub = nil
	if fake.waitForCorrectVersionReturnsOnCall == nil {
		fake.waitForCorrectVersionReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.waitForCorrectVersionReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeVerifyClient) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.ensureConfigVersionMutex.RLock()
	defer fake.ensureConfigVersionMutex.RUnlock()
	fake.getConfigVersionMutex.RLock()
	defer fake.getConfigVersionMutex.RUnlock()
	fake.waitForCorrectVersionMutex.RLock()
	defer fake.waitForCorrectVersionMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeVerifyClient) recordInvocation(key string, args []interface{}) {
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