// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	config "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/config"
	database "github.com/go-park-mail-ru/2019_1_Escapade/internal/pkg/database"

	mock "github.com/stretchr/testify/mock"

	sql "database/sql"

	time "time"
)

// Interface is an autogenerated mock type for the Interface type
type Interface struct {
	mock.Mock
}

// Begin provides a mock function with given fields:
func (_m *Interface) Begin() (database.TransactionI, error) {
	ret := _m.Called()

	var r0 database.TransactionI
	if rf, ok := ret.Get(0).(func() database.TransactionI); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(database.TransactionI)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Close provides a mock function with given fields:
func (_m *Interface) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Exec provides a mock function with given fields: query, args
func (_m *Interface) Exec(query string, args ...interface{}) (sql.Result, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 sql.Result
	if rf, ok := ret.Get(0).(func(string, ...interface{}) sql.Result); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(sql.Result)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Open provides a mock function with given fields: cdb
func (_m *Interface) Open(cdb config.Database) error {
	ret := _m.Called(cdb)

	var r0 error
	if rf, ok := ret.Get(0).(func(config.Database) error); ok {
		r0 = rf(cdb)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Ping provides a mock function with given fields:
func (_m *Interface) Ping() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Query provides a mock function with given fields: query, args
func (_m *Interface) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sql.Rows
	if rf, ok := ret.Get(0).(func(string, ...interface{}) *sql.Rows); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Rows)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...interface{}) error); ok {
		r1 = rf(query, args...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryRow provides a mock function with given fields: query, args
func (_m *Interface) QueryRow(query string, args ...interface{}) *sql.Row {
	var _ca []interface{}
	_ca = append(_ca, query)
	_ca = append(_ca, args...)
	ret := _m.Called(_ca...)

	var r0 *sql.Row
	if rf, ok := ret.Get(0).(func(string, ...interface{}) *sql.Row); ok {
		r0 = rf(query, args...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*sql.Row)
		}
	}

	return r0
}

// SetConnMaxLifetime provides a mock function with given fields: d
func (_m *Interface) SetConnMaxLifetime(d time.Duration) {
	_m.Called(d)
}

// SetMaxIdleConns provides a mock function with given fields: n
func (_m *Interface) SetMaxIdleConns(n int) {
	_m.Called(n)
}

// SetMaxOpenConns provides a mock function with given fields: n
func (_m *Interface) SetMaxOpenConns(n int) {
	_m.Called(n)
}
