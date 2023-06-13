// Code generated by MockGen. DO NOT EDIT.
// Source: metadata/internal/controller/metadata/controller.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/mamalmaleki/go-movie/metadata/pkg/model"
)

// MockmetadataRepository is a mock of metadataRepository interface.
type MockmetadataRepository struct {
	ctrl     *gomock.Controller
	recorder *MockmetadataRepositoryMockRecorder
}

// MockmetadataRepositoryMockRecorder is the mock recorder for MockmetadataRepository.
type MockmetadataRepositoryMockRecorder struct {
	mock *MockmetadataRepository
}

// NewMockmetadataRepository creates a new mock instance.
func NewMockmetadataRepository(ctrl *gomock.Controller) *MockmetadataRepository {
	mock := &MockmetadataRepository{ctrl: ctrl}
	mock.recorder = &MockmetadataRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmetadataRepository) EXPECT() *MockmetadataRepositoryMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockmetadataRepository) Get(ctx context.Context, id string) (*model.Metadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, id)
	ret0, _ := ret[0].(*model.Metadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockmetadataRepositoryMockRecorder) Get(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockmetadataRepository)(nil).Get), ctx, id)
}
