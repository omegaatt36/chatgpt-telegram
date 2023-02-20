// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/omegaatt36/chatgpt-telegram/appmodule/telegram/usecase (interfaces: TelegramUseCase)

// Package usecase is a generated GoMock package.
package usecase

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockTelegramUseCase is a mock of TelegramUseCase interface.
type MockTelegramUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockTelegramUseCaseMockRecorder
}

// MockTelegramUseCaseMockRecorder is the mock recorder for MockTelegramUseCase.
type MockTelegramUseCaseMockRecorder struct {
	mock *MockTelegramUseCase
}

// NewMockTelegramUseCase creates a new mock instance.
func NewMockTelegramUseCase(ctrl *gomock.Controller) *MockTelegramUseCase {
	mock := &MockTelegramUseCase{ctrl: ctrl}
	mock.recorder = &MockTelegramUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTelegramUseCase) EXPECT() *MockTelegramUseCaseMockRecorder {
	return m.recorder
}

// SendAsLiveOutput mocks base method.
func (m *MockTelegramUseCase) SendAsLiveOutput(arg0 int64, arg1 <-chan string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendAsLiveOutput", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendAsLiveOutput indicates an expected call of SendAsLiveOutput.
func (mr *MockTelegramUseCaseMockRecorder) SendAsLiveOutput(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAsLiveOutput", reflect.TypeOf((*MockTelegramUseCase)(nil).SendAsLiveOutput), arg0, arg1)
}
