package mocks

import gateways "github.com/lcollin/warehouse/gateways"
import mock "github.com/stretchr/testify/mock"
import models "github.com/lcollin/warehouse/models"
import uuid "github.com/pborman/uuid"

// Warehouse is an autogenerated mock type for the Warehouse type
type Warehouse struct {
	mock.Mock
}

// DeleteItem provides a mock function with given fields: id
func (_m *Warehouse) DeleteItem(id uuid.UUID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteOrder provides a mock function with given fields: id
func (_m *Warehouse) DeleteOrder(id uuid.UUID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteSubOrder provides a mock function with given fields: id
func (_m *Warehouse) DeleteSubOrder(id uuid.UUID) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(uuid.UUID) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetAllItems provides a mock function with given fields: offset, limit
func (_m *Warehouse) GetAllItems(offset int, limit int) ([]*models.Item, error) {
	ret := _m.Called(offset, limit)

	var r0 []*models.Item
	if rf, ok := ret.Get(0).(func(int, int) []*models.Item); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Item)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllOrders provides a mock function with given fields: offset, limit
func (_m *Warehouse) GetAllOrders(offset int, limit int) ([]*models.Order, error) {
	ret := _m.Called(offset, limit)

	var r0 []*models.Order
	if rf, ok := ret.Get(0).(func(int, int) []*models.Order); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Order)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetAllSubOrders provides a mock function with given fields: offset, limit
func (_m *Warehouse) GetAllSubOrders(offset int, limit int) ([]*models.SubOrder, error) {
	ret := _m.Called(offset, limit)

	var r0 []*models.SubOrder
	if rf, ok := ret.Get(0).(func(int, int) []*models.SubOrder); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.SubOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetItemByID provides a mock function with given fields: id
func (_m *Warehouse) GetItemByID(id uuid.UUID) (*models.Item, error) {
	ret := _m.Called(id)

	var r0 *models.Item
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Item); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Item)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderByID provides a mock function with given fields: id
func (_m *Warehouse) GetOrderByID(id uuid.UUID) (*models.Order, error) {
	ret := _m.Called(id)

	var r0 *models.Order
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.Order); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrderLabel provides a mock function with given fields: id
func (_m *Warehouse) GetOrderLabel(id uuid.UUID) (string, error) {
	ret := _m.Called(id)

	var r0 string
	if rf, ok := ret.Get(0).(func(uuid.UUID) string); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrdersByRoasterID provides a mock function with given fields: id, offset, limit
func (_m *Warehouse) GetOrdersByRoasterID(id uuid.UUID, offset int, limit int) ([]*models.Order, error) {
	ret := _m.Called(id, offset, limit)

	var r0 []*models.Order
	if rf, ok := ret.Get(0).(func(uuid.UUID, int, int) []*models.Order); ok {
		r0 = rf(id, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Order)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, int, int) error); ok {
		r1 = rf(id, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetOrdersByUserID provides a mock function with given fields: id, offset, limit
func (_m *Warehouse) GetOrdersByUserID(id uuid.UUID, offset int, limit int) ([]*models.Order, error) {
	ret := _m.Called(id, offset, limit)

	var r0 []*models.Order
	if rf, ok := ret.Get(0).(func(uuid.UUID, int, int) []*models.Order); ok {
		r0 = rf(id, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.Order)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID, int, int) error); ok {
		r1 = rf(id, offset, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetSubOrderByID provides a mock function with given fields: id
func (_m *Warehouse) GetSubOrderByID(id uuid.UUID) (*models.SubOrder, error) {
	ret := _m.Called(id)

	var r0 *models.SubOrder
	if rf, ok := ret.Get(0).(func(uuid.UUID) *models.SubOrder); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SubOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uuid.UUID) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewItem provides a mock function with given fields: newItem
func (_m *Warehouse) NewItem(newItem *models.Item) (*models.Item, error) {
	ret := _m.Called(newItem)

	var r0 *models.Item
	if rf, ok := ret.Get(0).(func(*models.Item) *models.Item); ok {
		r0 = rf(newItem)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Item)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Item) error); ok {
		r1 = rf(newItem)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewOrder provides a mock function with given fields: newOrder
func (_m *Warehouse) NewOrder(newOrder *models.Order) (*models.Order, error) {
	ret := _m.Called(newOrder)

	var r0 *models.Order
	if rf, ok := ret.Get(0).(func(*models.Order) *models.Order); ok {
		r0 = rf(newOrder)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Order) error); ok {
		r1 = rf(newOrder)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewSubOrder provides a mock function with given fields: newSubOrder
func (_m *Warehouse) NewSubOrder(newSubOrder *models.SubOrder) (*models.SubOrder, error) {
	ret := _m.Called(newSubOrder)

	var r0 *models.SubOrder
	if rf, ok := ret.Get(0).(func(*models.SubOrder) *models.SubOrder); ok {
		r0 = rf(newSubOrder)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SubOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.SubOrder) error); ok {
		r1 = rf(newSubOrder)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateItem provides a mock function with given fields: update
func (_m *Warehouse) UpdateItem(update *models.Item) (*models.Item, error) {
	ret := _m.Called(update)

	var r0 *models.Item
	if rf, ok := ret.Get(0).(func(*models.Item) *models.Item); ok {
		r0 = rf(update)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Item)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Item) error); ok {
		r1 = rf(update)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateOrder provides a mock function with given fields: update
func (_m *Warehouse) UpdateOrder(update *models.Order) (*models.Order, error) {
	ret := _m.Called(update)

	var r0 *models.Order
	if rf, ok := ret.Get(0).(func(*models.Order) *models.Order); ok {
		r0 = rf(update)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Order)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.Order) error); ok {
		r1 = rf(update)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateSubOrder provides a mock function with given fields: update
func (_m *Warehouse) UpdateSubOrder(update *models.SubOrder) (*models.SubOrder, error) {
	ret := _m.Called(update)

	var r0 *models.SubOrder
	if rf, ok := ret.Get(0).(func(*models.SubOrder) *models.SubOrder); ok {
		r0 = rf(update)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.SubOrder)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*models.SubOrder) error); ok {
		r1 = rf(update)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

var _ gateways.Warehouse = (*Warehouse)(nil)
