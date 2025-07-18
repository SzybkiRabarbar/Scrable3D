// Code generated by MockGen. DO NOT EDIT.
// Source: internal/ctrl/words.go
//
// Generated by this command:
//
//	mockgen -source=internal/ctrl/words.go -destination=internal/mock/mock_ctrl_words.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockWordsController is a mock of WordsController interface.
type MockWordsController struct {
	ctrl     *gomock.Controller
	recorder *MockWordsControllerMockRecorder
	isgomock struct{}
}

// MockWordsControllerMockRecorder is the mock recorder for MockWordsController.
type MockWordsControllerMockRecorder struct {
	mock *MockWordsController
}

// NewMockWordsController creates a new mock instance.
func NewMockWordsController(ctrl *gomock.Controller) *MockWordsController {
	mock := &MockWordsController{ctrl: ctrl}
	mock.recorder = &MockWordsControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWordsController) EXPECT() *MockWordsControllerMockRecorder {
	return m.recorder
}

// CheckWord mocks base method.
func (m *MockWordsController) CheckWord(word string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckWord", word)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckWord indicates an expected call of CheckWord.
func (mr *MockWordsControllerMockRecorder) CheckWord(word any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckWord", reflect.TypeOf((*MockWordsController)(nil).CheckWord), word)
}

// WordsNumber mocks base method.
func (m *MockWordsController) WordsNumber() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WordsNumber")
	ret0, _ := ret[0].(int)
	return ret0
}

// WordsNumber indicates an expected call of WordsNumber.
func (mr *MockWordsControllerMockRecorder) WordsNumber() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WordsNumber", reflect.TypeOf((*MockWordsController)(nil).WordsNumber))
}
