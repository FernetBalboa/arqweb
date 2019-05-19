// Code generated by MockGen. DO NOT EDIT.
// Source: category_storage.go

// Package mock is a generated GoMock package.
package mock

import (
	domain "github.com/fernetbalboa/arqweb/src/api/domain"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockCategoryStorage is a mock of CategoryStorage interface
type MockCategoryStorage struct {
	ctrl     *gomock.Controller
	recorder *MockCategoryStorageMockRecorder
}

// MockCategoryStorageMockRecorder is the mock recorder for MockCategoryStorage
type MockCategoryStorageMockRecorder struct {
	mock *MockCategoryStorage
}

// NewMockCategoryStorage creates a new mock instance
func NewMockCategoryStorage(ctrl *gomock.Controller) *MockCategoryStorage {
	mock := &MockCategoryStorage{ctrl: ctrl}
	mock.recorder = &MockCategoryStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCategoryStorage) EXPECT() *MockCategoryStorageMockRecorder {
	return m.recorder
}

// SaveCategory mocks base method
func (m *MockCategoryStorage) SaveCategory(category *domain.Category) (*domain.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveCategory", category)
	ret0, _ := ret[0].(*domain.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveCategory indicates an expected call of SaveCategory
func (mr *MockCategoryStorageMockRecorder) SaveCategory(category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveCategory", reflect.TypeOf((*MockCategoryStorage)(nil).SaveCategory), category)
}

// GetCategories mocks base method
func (m *MockCategoryStorage) GetCategories() ([]domain.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategories")
	ret0, _ := ret[0].([]domain.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategories indicates an expected call of GetCategories
func (mr *MockCategoryStorageMockRecorder) GetCategories() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategories", reflect.TypeOf((*MockCategoryStorage)(nil).GetCategories))
}

// SearchCategory mocks base method
func (m *MockCategoryStorage) SearchCategory(filters *domain.CategoryFilter) ([]*domain.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchCategory", filters)
	ret0, _ := ret[0].([]*domain.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchCategory indicates an expected call of SearchCategory
func (mr *MockCategoryStorageMockRecorder) SearchCategory(filters interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchCategory", reflect.TypeOf((*MockCategoryStorage)(nil).SearchCategory), filters)
}

// EditCategory mocks base method
func (m *MockCategoryStorage) EditCategory(newVersionCategory *domain.Category) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditCategory", newVersionCategory)
	ret0, _ := ret[0].(error)
	return ret0
}

// EditCategory indicates an expected call of EditCategory
func (mr *MockCategoryStorageMockRecorder) EditCategory(newVersionCategory interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditCategory", reflect.TypeOf((*MockCategoryStorage)(nil).EditCategory), newVersionCategory)
}
