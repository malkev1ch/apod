// Code generated by MockGen. DO NOT EDIT.
// Source: picture.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	model "github.com/malkev1ch/apod/internal/model"
)

// MockPictureService is a mock of PictureService interface.
type MockPictureService struct {
	ctrl     *gomock.Controller
	recorder *MockPictureServiceMockRecorder
}

// MockPictureServiceMockRecorder is the mock recorder for MockPictureService.
type MockPictureServiceMockRecorder struct {
	mock *MockPictureService
}

// NewMockPictureService creates a new mock instance.
func NewMockPictureService(ctrl *gomock.Controller) *MockPictureService {
	mock := &MockPictureService{ctrl: ctrl}
	mock.recorder = &MockPictureServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPictureService) EXPECT() *MockPictureServiceMockRecorder {
	return m.recorder
}

// GetAll mocks base method.
func (m *MockPictureService) GetAll(ctx context.Context) ([]model.Picture, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]model.Picture)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockPictureServiceMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockPictureService)(nil).GetAll), ctx)
}

// GetByDate mocks base method.
func (m *MockPictureService) GetByDate(ctx context.Context, date time.Time) (model.Picture, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByDate", ctx, date)
	ret0, _ := ret[0].(model.Picture)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByDate indicates an expected call of GetByDate.
func (mr *MockPictureServiceMockRecorder) GetByDate(ctx, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByDate", reflect.TypeOf((*MockPictureService)(nil).GetByDate), ctx, date)
}
