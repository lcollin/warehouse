package mocks

import helpers "github.com/lcollin/warehouse/helpers"
import mock "github.com/stretchr/testify/mock"
import models "github.com/lcollin/warehouse/models"

// SubOrderI is an autogenerated mock type for the SubOrderI type
type SubOrderI struct {
	mock.Mock
}

// Delete provides a mock function with given fields: _a0
func (_m *SubOrderI) Delete(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAll provides a mock function with given fields: _a0, _a1
func (_m *SubOrderI) GetAll(_a0 int, _a1 int) ([]*models.SubOrder, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []*models.SubOrder
	if rf, ok := ret.Get(0).(func(int, int) []*models.SubOrder); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.SubOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetByID provides a mock function with given fields: _a0
func (_m *SubOrderI) GetByID(_a0 string) (*models.SubOrder, error) {
	ret := _m.Called(_a0)

	var r0 *models.SubOrder
	if rf, ok := ret.Get(0).(func(string) *models.SubOrder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SubOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetBySubOrderID provides a mock function with given fields: _a0
func (_m *SubOrderI) GetBySubOrderID(_a0 string) (*models.SubOrder, error) {
	ret := _m.Called(_a0)

	var r0 *models.SubOrder
	if rf, ok := ret.Get(0).(func(string) *models.SubOrder); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SubOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Insert provides a mock function with given fields: _a0
func (_m *SubOrderI) Insert(_a0 *models.SubOrder) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.SubOrder) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Update provides a mock function with given fields: _a0, _a1
func (_m *SubOrderI) Update(_a0 *models.SubOrder, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.SubOrder, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

var _ helpers.SubOrderI = (*SubOrderI)(nil)
