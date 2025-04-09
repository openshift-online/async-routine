// Code generated by MockGen. DO NOT EDIT.
// Source: async_routine.go

// Package async is a generated GoMock package.
package async

import (
	reflect "reflect"
	time "time"

	gomock "go.uber.org/mock/gomock"
)

// MockAsyncRoutine is a mock of AsyncRoutine interface.
type MockAsyncRoutine struct {
	ctrl     *gomock.Controller
	recorder *MockAsyncRoutineMockRecorder
}

// MockAsyncRoutineMockRecorder is the mock recorder for MockAsyncRoutine.
type MockAsyncRoutineMockRecorder struct {
	mock *MockAsyncRoutine
}

// NewMockAsyncRoutine creates a new mock instance.
func NewMockAsyncRoutine(ctrl *gomock.Controller) *MockAsyncRoutine {
	mock := &MockAsyncRoutine{ctrl: ctrl}
	mock.recorder = &MockAsyncRoutineMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAsyncRoutine) EXPECT() *MockAsyncRoutineMockRecorder {
	return m.recorder
}

// CreatedAt mocks base method.
func (m *MockAsyncRoutine) CreatedAt() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatedAt")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// CreatedAt indicates an expected call of CreatedAt.
func (mr *MockAsyncRoutineMockRecorder) CreatedAt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatedAt", reflect.TypeOf((*MockAsyncRoutine)(nil).CreatedAt))
}

// FinishedAt mocks base method.
func (m *MockAsyncRoutine) FinishedAt() *time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinishedAt")
	ret0, _ := ret[0].(*time.Time)
	return ret0
}

// FinishedAt indicates an expected call of FinishedAt.
func (mr *MockAsyncRoutineMockRecorder) FinishedAt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinishedAt", reflect.TypeOf((*MockAsyncRoutine)(nil).FinishedAt))
}

// GetData mocks base method.
func (m *MockAsyncRoutine) GetData() map[string]string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetData")
	ret0, _ := ret[0].(map[string]string)
	return ret0
}

// GetData indicates an expected call of GetData.
func (mr *MockAsyncRoutineMockRecorder) GetData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetData", reflect.TypeOf((*MockAsyncRoutine)(nil).GetData))
}

// Name mocks base method.
func (m *MockAsyncRoutine) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockAsyncRoutineMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockAsyncRoutine)(nil).Name))
}

// OpId mocks base method.
func (m *MockAsyncRoutine) OpId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpId")
	ret0, _ := ret[0].(string)
	return ret0
}

// OpId indicates an expected call of OpId.
func (mr *MockAsyncRoutineMockRecorder) OpId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpId", reflect.TypeOf((*MockAsyncRoutine)(nil).OpId))
}

// OriginatorOpId mocks base method.
func (m *MockAsyncRoutine) OriginatorOpId() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OriginatorOpId")
	ret0, _ := ret[0].(string)
	return ret0
}

// OriginatorOpId indicates an expected call of OriginatorOpId.
func (mr *MockAsyncRoutineMockRecorder) OriginatorOpId() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OriginatorOpId", reflect.TypeOf((*MockAsyncRoutine)(nil).OriginatorOpId))
}

// StartedAt mocks base method.
func (m *MockAsyncRoutine) StartedAt() *time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartedAt")
	ret0, _ := ret[0].(*time.Time)
	return ret0
}

// StartedAt indicates an expected call of StartedAt.
func (mr *MockAsyncRoutineMockRecorder) StartedAt() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartedAt", reflect.TypeOf((*MockAsyncRoutine)(nil).StartedAt))
}

// Status mocks base method.
func (m *MockAsyncRoutine) Status() RoutineStatus {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(RoutineStatus)
	return ret0
}

// Status indicates an expected call of Status.
func (mr *MockAsyncRoutineMockRecorder) Status() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockAsyncRoutine)(nil).Status))
}

// hasExceededTimebox mocks base method.
func (m *MockAsyncRoutine) hasExceededTimebox() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "hasExceededTimebox")
	ret0, _ := ret[0].(bool)
	return ret0
}

// hasExceededTimebox indicates an expected call of hasExceededTimebox.
func (mr *MockAsyncRoutineMockRecorder) hasExceededTimebox() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "hasExceededTimebox", reflect.TypeOf((*MockAsyncRoutine)(nil).hasExceededTimebox))
}

// id mocks base method.
func (m *MockAsyncRoutine) id() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "id")
	ret0, _ := ret[0].(string)
	return ret0
}

// id indicates an expected call of id.
func (mr *MockAsyncRoutineMockRecorder) id() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "id", reflect.TypeOf((*MockAsyncRoutine)(nil).id))
}

// isFinished mocks base method.
func (m *MockAsyncRoutine) isFinished() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "isFinished")
	ret0, _ := ret[0].(bool)
	return ret0
}

// isFinished indicates an expected call of isFinished.
func (mr *MockAsyncRoutineMockRecorder) isFinished() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isFinished", reflect.TypeOf((*MockAsyncRoutine)(nil).isFinished))
}

// isRunning mocks base method.
func (m *MockAsyncRoutine) isRunning() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "isRunning")
	ret0, _ := ret[0].(bool)
	return ret0
}

// isRunning indicates an expected call of isRunning.
func (mr *MockAsyncRoutineMockRecorder) isRunning() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "isRunning", reflect.TypeOf((*MockAsyncRoutine)(nil).isRunning))
}

// run mocks base method.
func (m *MockAsyncRoutine) run(manager AsyncRoutineManager) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "run", manager)
}

// run indicates an expected call of run.
func (mr *MockAsyncRoutineMockRecorder) run(manager interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "run", reflect.TypeOf((*MockAsyncRoutine)(nil).run), manager)
}
