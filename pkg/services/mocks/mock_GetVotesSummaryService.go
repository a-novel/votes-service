// Code generated by mockery v2.33.2. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	models "github.com/a-novel/votes-service/pkg/models"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// GetVotesSummaryService is an autogenerated mock type for the GetVotesSummaryService type
type GetVotesSummaryService struct {
	mock.Mock
}

type GetVotesSummaryService_Expecter struct {
	mock *mock.Mock
}

func (_m *GetVotesSummaryService) EXPECT() *GetVotesSummaryService_Expecter {
	return &GetVotesSummaryService_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, targetID, target
func (_m *GetVotesSummaryService) Get(ctx context.Context, targetID uuid.UUID, target string) (*models.VotesSummary, error) {
	ret := _m.Called(ctx, targetID, target)

	var r0 *models.VotesSummary
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) (*models.VotesSummary, error)); ok {
		return rf(ctx, targetID, target)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID, string) *models.VotesSummary); ok {
		r0 = rf(ctx, targetID, target)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.VotesSummary)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID, string) error); ok {
		r1 = rf(ctx, targetID, target)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetVotesSummaryService_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type GetVotesSummaryService_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - targetID uuid.UUID
//   - target string
func (_e *GetVotesSummaryService_Expecter) Get(ctx interface{}, targetID interface{}, target interface{}) *GetVotesSummaryService_Get_Call {
	return &GetVotesSummaryService_Get_Call{Call: _e.mock.On("Get", ctx, targetID, target)}
}

func (_c *GetVotesSummaryService_Get_Call) Run(run func(ctx context.Context, targetID uuid.UUID, target string)) *GetVotesSummaryService_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID), args[2].(string))
	})
	return _c
}

func (_c *GetVotesSummaryService_Get_Call) Return(_a0 *models.VotesSummary, _a1 error) *GetVotesSummaryService_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *GetVotesSummaryService_Get_Call) RunAndReturn(run func(context.Context, uuid.UUID, string) (*models.VotesSummary, error)) *GetVotesSummaryService_Get_Call {
	_c.Call.Return(run)
	return _c
}

// NewGetVotesSummaryService creates a new instance of GetVotesSummaryService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewGetVotesSummaryService(t interface {
	mock.TestingT
	Cleanup(func())
}) *GetVotesSummaryService {
	mock := &GetVotesSummaryService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
