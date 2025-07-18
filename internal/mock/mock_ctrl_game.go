// Code generated by MockGen. DO NOT EDIT.
// Source: internal/ctrl/game.go
//
// Generated by this command:
//
//	mockgen -source=internal/ctrl/game.go -destination=internal/mock/mock_ctrl_game.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"
	dto "scrable3/internal/dto"

	gomock "go.uber.org/mock/gomock"
)

// MockGameController is a mock of GameController interface.
type MockGameController struct {
	ctrl     *gomock.Controller
	recorder *MockGameControllerMockRecorder
	isgomock struct{}
}

// MockGameControllerMockRecorder is the mock recorder for MockGameController.
type MockGameControllerMockRecorder struct {
	mock *MockGameController
}

// NewMockGameController creates a new mock instance.
func NewMockGameController(ctrl *gomock.Controller) *MockGameController {
	mock := &MockGameController{ctrl: ctrl}
	mock.recorder = &MockGameControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGameController) EXPECT() *MockGameControllerMockRecorder {
	return m.recorder
}

// GetAvaibleChars mocks base method.
func (m *MockGameController) GetAvaibleChars(ctx *dto.WsContext) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAvaibleChars", ctx)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAvaibleChars indicates an expected call of GetAvaibleChars.
func (mr *MockGameControllerMockRecorder) GetAvaibleChars(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAvaibleChars", reflect.TypeOf((*MockGameController)(nil).GetAvaibleChars), ctx)
}

// GetCurrentFields mocks base method.
func (m *MockGameController) GetCurrentFields(ctx *dto.WsContext) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentFields", ctx)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentFields indicates an expected call of GetCurrentFields.
func (mr *MockGameControllerMockRecorder) GetCurrentFields(ctx any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentFields", reflect.TypeOf((*MockGameController)(nil).GetCurrentFields), ctx)
}

// ReceiveChars mocks base method.
func (m *MockGameController) ReceiveChars(ctx *dto.WsContext, p *dto.PlayData) ([]byte, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReceiveChars", ctx, p)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ReceiveChars indicates an expected call of ReceiveChars.
func (mr *MockGameControllerMockRecorder) ReceiveChars(ctx, p any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReceiveChars", reflect.TypeOf((*MockGameController)(nil).ReceiveChars), ctx, p)
}
