// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/rugwirobaker/hermes (interfaces: SendService,Pubsub,AppStore,MessageStore)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	hermes "github.com/rugwirobaker/hermes"
)

// MockSendService is a mock of SendService interface.
type MockSendService struct {
	ctrl     *gomock.Controller
	recorder *MockSendServiceMockRecorder
}

// MockSendServiceMockRecorder is the mock recorder for MockSendService.
type MockSendServiceMockRecorder struct {
	mock *MockSendService
}

// NewMockSendService creates a new mock instance.
func NewMockSendService(ctrl *gomock.Controller) *MockSendService {
	mock := &MockSendService{ctrl: ctrl}
	mock.recorder = &MockSendServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSendService) EXPECT() *MockSendServiceMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockSendService) Send(arg0 context.Context, arg1 *hermes.SMS) (*hermes.Report, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0, arg1)
	ret0, _ := ret[0].(*hermes.Report)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Send indicates an expected call of Send.
func (mr *MockSendServiceMockRecorder) Send(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockSendService)(nil).Send), arg0, arg1)
}

// MockPubsub is a mock of Pubsub interface.
type MockPubsub struct {
	ctrl     *gomock.Controller
	recorder *MockPubsubMockRecorder
}

// MockPubsubMockRecorder is the mock recorder for MockPubsub.
type MockPubsubMockRecorder struct {
	mock *MockPubsub
}

// NewMockPubsub creates a new mock instance.
func NewMockPubsub(ctrl *gomock.Controller) *MockPubsub {
	mock := &MockPubsub{ctrl: ctrl}
	mock.recorder = &MockPubsubMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPubsub) EXPECT() *MockPubsubMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockPubsub) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockPubsubMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPubsub)(nil).Close))
}

// Done mocks base method.
func (m *MockPubsub) Done(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Done", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Done indicates an expected call of Done.
func (mr *MockPubsubMockRecorder) Done(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Done", reflect.TypeOf((*MockPubsub)(nil).Done), arg0, arg1)
}

// Publish mocks base method.
func (m *MockPubsub) Publish(arg0 context.Context, arg1 hermes.Event) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Publish", arg0, arg1)
}

// Publish indicates an expected call of Publish.
func (mr *MockPubsubMockRecorder) Publish(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockPubsub)(nil).Publish), arg0, arg1)
}

// Subscribe mocks base method.
func (m *MockPubsub) Subscribe(arg0 context.Context, arg1 string) (<-chan hermes.Event, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", arg0, arg1)
	ret0, _ := ret[0].(<-chan hermes.Event)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockPubsubMockRecorder) Subscribe(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockPubsub)(nil).Subscribe), arg0, arg1)
}

// MockAppStore is a mock of AppStore interface.
type MockAppStore struct {
	ctrl     *gomock.Controller
	recorder *MockAppStoreMockRecorder
}

// MockAppStoreMockRecorder is the mock recorder for MockAppStore.
type MockAppStoreMockRecorder struct {
	mock *MockAppStore
}

// NewMockAppStore creates a new mock instance.
func NewMockAppStore(ctrl *gomock.Controller) *MockAppStore {
	mock := &MockAppStore{ctrl: ctrl}
	mock.recorder = &MockAppStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAppStore) EXPECT() *MockAppStoreMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockAppStore) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockAppStoreMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockAppStore)(nil).Delete), arg0, arg1)
}

// FindByToken mocks base method.
func (m *MockAppStore) FindByToken(arg0 context.Context, arg1 string) (*hermes.App, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByToken", arg0, arg1)
	ret0, _ := ret[0].(*hermes.App)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByToken indicates an expected call of FindByToken.
func (mr *MockAppStoreMockRecorder) FindByToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByToken", reflect.TypeOf((*MockAppStore)(nil).FindByToken), arg0, arg1)
}

// Get mocks base method.
func (m *MockAppStore) Get(arg0 context.Context, arg1 string) (*hermes.App, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*hermes.App)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockAppStoreMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockAppStore)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockAppStore) List(arg0 context.Context) ([]*hermes.App, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]*hermes.App)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockAppStoreMockRecorder) List(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockAppStore)(nil).List), arg0)
}

// Register mocks base method.
func (m *MockAppStore) Register(arg0 context.Context, arg1 *hermes.App) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockAppStoreMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAppStore)(nil).Register), arg0, arg1)
}

// Update mocks base method.
func (m *MockAppStore) Update(arg0 context.Context, arg1 *hermes.App) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockAppStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockAppStore)(nil).Update), arg0, arg1)
}

// MockMessageStore is a mock of MessageStore interface.
type MockMessageStore struct {
	ctrl     *gomock.Controller
	recorder *MockMessageStoreMockRecorder
}

// MockMessageStoreMockRecorder is the mock recorder for MockMessageStore.
type MockMessageStoreMockRecorder struct {
	mock *MockMessageStore
}

// NewMockMessageStore creates a new mock instance.
func NewMockMessageStore(ctrl *gomock.Controller) *MockMessageStore {
	mock := &MockMessageStore{ctrl: ctrl}
	mock.recorder = &MockMessageStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageStore) EXPECT() *MockMessageStoreMockRecorder {
	return m.recorder
}

// Insert mocks base method.
func (m *MockMessageStore) Insert(arg0 context.Context, arg1 *hermes.Message) (*hermes.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", arg0, arg1)
	ret0, _ := ret[0].(*hermes.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Insert indicates an expected call of Insert.
func (mr *MockMessageStoreMockRecorder) Insert(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockMessageStore)(nil).Insert), arg0, arg1)
}

// List mocks base method.
func (m *MockMessageStore) List(arg0 context.Context, arg1 *hermes.ListOptions) ([]*hermes.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]*hermes.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockMessageStoreMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockMessageStore)(nil).List), arg0, arg1)
}

// MessageByID mocks base method.
func (m *MockMessageStore) MessageByID(arg0 context.Context, arg1 string) (*hermes.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MessageByID", arg0, arg1)
	ret0, _ := ret[0].(*hermes.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MessageByID indicates an expected call of MessageByID.
func (mr *MockMessageStoreMockRecorder) MessageByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MessageByID", reflect.TypeOf((*MockMessageStore)(nil).MessageByID), arg0, arg1)
}

// MessageByPhone mocks base method.
func (m *MockMessageStore) MessageByPhone(arg0 context.Context, arg1 string) (*hermes.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MessageByPhone", arg0, arg1)
	ret0, _ := ret[0].(*hermes.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MessageByPhone indicates an expected call of MessageByPhone.
func (mr *MockMessageStoreMockRecorder) MessageByPhone(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MessageByPhone", reflect.TypeOf((*MockMessageStore)(nil).MessageByPhone), arg0, arg1)
}

// MessageBySerial mocks base method.
func (m *MockMessageStore) MessageBySerial(arg0 context.Context, arg1 string) (*hermes.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MessageBySerial", arg0, arg1)
	ret0, _ := ret[0].(*hermes.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MessageBySerial indicates an expected call of MessageBySerial.
func (mr *MockMessageStoreMockRecorder) MessageBySerial(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MessageBySerial", reflect.TypeOf((*MockMessageStore)(nil).MessageBySerial), arg0, arg1)
}

// Update mocks base method.
func (m *MockMessageStore) Update(arg0 context.Context, arg1 *hermes.Message) (*hermes.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*hermes.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockMessageStoreMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockMessageStore)(nil).Update), arg0, arg1)
}
