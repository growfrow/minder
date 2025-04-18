// Code generated by MockGen. DO NOT EDIT.
// Source: ./service.go
//
// Generated by this command:
//
//	mockgen -package mock_service -destination=./mock/service.go -source=./service.go
//

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	v1 "github.com/mindersec/minder/pkg/api/protobuf/go/minder/v1"
	gomock "go.uber.org/mock/gomock"
)

// MockEntityService is a mock of EntityService interface.
type MockEntityService struct {
	ctrl     *gomock.Controller
	recorder *MockEntityServiceMockRecorder
	isgomock struct{}
}

// MockEntityServiceMockRecorder is the mock recorder for MockEntityService.
type MockEntityServiceMockRecorder struct {
	mock *MockEntityService
}

// NewMockEntityService creates a new mock instance.
func NewMockEntityService(ctrl *gomock.Controller) *MockEntityService {
	mock := &MockEntityService{ctrl: ctrl}
	mock.recorder = &MockEntityServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEntityService) EXPECT() *MockEntityServiceMockRecorder {
	return m.recorder
}

// DeleteEntityByID mocks base method.
func (m *MockEntityService) DeleteEntityByID(ctx context.Context, entityID, projectID uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteEntityByID", ctx, entityID, projectID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteEntityByID indicates an expected call of DeleteEntityByID.
func (mr *MockEntityServiceMockRecorder) DeleteEntityByID(ctx, entityID, projectID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteEntityByID", reflect.TypeOf((*MockEntityService)(nil).DeleteEntityByID), ctx, entityID, projectID)
}

// GetEntityByID mocks base method.
func (m *MockEntityService) GetEntityByID(ctx context.Context, entityID, projectID uuid.UUID) (*v1.EntityInstance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntityByID", ctx, entityID, projectID)
	ret0, _ := ret[0].(*v1.EntityInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntityByID indicates an expected call of GetEntityByID.
func (mr *MockEntityServiceMockRecorder) GetEntityByID(ctx, entityID, projectID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntityByID", reflect.TypeOf((*MockEntityService)(nil).GetEntityByID), ctx, entityID, projectID)
}

// GetEntityByName mocks base method.
func (m *MockEntityService) GetEntityByName(ctx context.Context, name string, projectID, providerID uuid.UUID, entityType v1.Entity) (*v1.EntityInstance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntityByName", ctx, name, projectID, providerID, entityType)
	ret0, _ := ret[0].(*v1.EntityInstance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntityByName indicates an expected call of GetEntityByName.
func (mr *MockEntityServiceMockRecorder) GetEntityByName(ctx, name, projectID, providerID, entityType any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntityByName", reflect.TypeOf((*MockEntityService)(nil).GetEntityByName), ctx, name, projectID, providerID, entityType)
}

// ListEntities mocks base method.
func (m *MockEntityService) ListEntities(ctx context.Context, projectID, providerID uuid.UUID, entityType v1.Entity, cursor string, limit int64) ([]*v1.EntityInstance, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEntities", ctx, projectID, providerID, entityType, cursor, limit)
	ret0, _ := ret[0].([]*v1.EntityInstance)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ListEntities indicates an expected call of ListEntities.
func (mr *MockEntityServiceMockRecorder) ListEntities(ctx, projectID, providerID, entityType, cursor, limit any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEntities", reflect.TypeOf((*MockEntityService)(nil).ListEntities), ctx, projectID, providerID, entityType, cursor, limit)
}
