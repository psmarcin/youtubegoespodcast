package app

import (
	"context"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CacheRepository struct {
	mock.Mock
}

func (c *CacheRepository) GetKey(ctx context.Context, key string) (string, error) {
	args := c.Called(mock.Anything, key)
	return args.String(0), args.Error(1)
}

func (c *CacheRepository) SetKey(ctx context.Context, key string, value string, ex time.Duration) error {
	args := c.Called(mock.Anything, key, value, ex)
	return args.Error(0)
}

type ImportTo struct {
	Test string `json:"test"`
}

func TestCacheService_Get(t *testing.T) {
	cacheMock := new(CacheRepository)
	to := ImportTo{}

	ctx, _ := tracer.Start(context.Background(), "x")

	type fields struct {
		cache cacheRepository
	}
	type args struct {
		key string
		to  interface{}
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  bool
		before   func()
		toExpect ImportTo
	}{
		{
			name: "should return 1 value",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key: "1",
				to:  &to,
			},
			wantErr: false,
			before: func() {
				to = ImportTo{}
				cacheMock.On("GetKey", ctx, "1").Return(`{"test": "1"}`, nil)
			},
			toExpect: ImportTo{Test: "1"},
		},
		{
			name: "should return empty value",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key: "2",
				to:  &to,
			},
			wantErr: false,
			before: func() {
				to = ImportTo{}
				cacheMock.On("GetKey", ctx, "2").Return(`{}`, nil)
			},
			toExpect: ImportTo{Test: ""},
		},
		{
			name: "should return error on empty string value",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key: "3",
				to:  &to,
			},
			wantErr: true,
			before: func() {
				to = ImportTo{}
				cacheMock.On("GetKey", ctx, "3").Return(``, nil)
			},
			toExpect: ImportTo{Test: ""},
		},
		{
			name: "should return error on GetKey error",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key: "4",
				to:  &to,
			},
			wantErr: true,
			before: func() {
				to = ImportTo{}
				cacheMock.On("GetKey", ctx, "4").Return(``, errors.New("connection problem"))
			},
			toExpect: ImportTo{Test: ""},
		},
		{
			name: "should not update to because it's not set",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key: "5",
				to:  nil,
			},
			wantErr: false,
			before: func() {
				to = ImportTo{}
				cacheMock.On("GetKey", ctx, "5").Return(`{"test": "5"}`, nil)
			},
			toExpect: ImportTo{Test: ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CacheService{
				cache: tt.fields.cache,
			}
			tt.before()
			if err := c.Get(ctx, tt.args.key, tt.args.to); (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.EqualValues(t, to, tt.toExpect)
		})
	}
}

func TestCacheService_Set(t *testing.T) {
	cacheMock := new(CacheRepository)

	type fields struct {
		cache cacheRepository
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		before  func()
	}{
		{
			name: "should set cache for empty value",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key:   "1",
				value: nil,
			},
			wantErr: false,
			before: func() {
				cacheMock.On("SetKey", context.Background(), "1", "null", CacheTTL).Return(nil)
			},
		},
		{
			name: "should not return error for string value",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key:   "2",
				value: "aaa",
			},
			wantErr: false,
			before: func() {
				cacheMock.On("SetKey", context.Background(), "2", "\"aaa\"", CacheTTL).Return(nil)
			},
		},
		{
			name: "should not return error for struct value",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key: "3",
				value: ImportTo{
					Test: "3",
				},
			},
			wantErr: false,
			before: func() {
				cacheMock.On("SetKey", context.Background(), "3", "{\"test\":\"3\"}", CacheTTL).Return(nil)
			},
		},
		{
			name: "should not return error on cache repository setKey error, just log",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key:   "3",
				value: nil,
			},
			wantErr: false,
			before: func() {
				cacheMock.On("SetKey", context.Background(), "3", "null", CacheTTL).Return(errors.New("connection error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CacheService{
				cache: tt.fields.cache,
			}
			tt.before()
			if err := c.Set(context.Background(), tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCacheService_MarshalAndSet(t *testing.T) {
	cacheMock := new(CacheRepository)

	type fields struct {
		cache cacheRepository
	}
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		before  func()
	}{
		{
			name: "should set cache for empty value",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key:   "1",
				value: nil,
			},
			wantErr: false,
			before: func() {
				cacheMock.On("SetKey", context.Background(), "1", "\"null\"", CacheTTL).Return(nil)
			},
		},
		{
			name: "should not return error for string value",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key:   "2",
				value: "aaa",
			},
			wantErr: false,
			before: func() {
				cacheMock.On("SetKey", context.Background(), "2", "\"\\\"aaa\\\"\"", CacheTTL).Return(nil)
			},
		},
		{
			name: "should not return error for struct value",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key: "3",
				value: ImportTo{
					Test: "3",
				},
			},
			wantErr: false,
			before: func() {
				cacheMock.On("SetKey", context.Background(), "3", "\"{\\\"test\\\":\\\"3\\\"}\"", CacheTTL).Return(nil)
			},
		},
		{
			name: "should not return error on cache repository setKey error, just log",
			fields: fields{
				cache: cacheMock,
			},
			args: args{
				key:   "4",
				value: nil,
			},
			wantErr: false,
			before: func() {
				cacheMock.On("SetKey", context.Background(), "4", "\"null\"", CacheTTL).Return(errors.New("connection error"))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CacheService{
				cache: tt.fields.cache,
			}
			tt.before()
			if err := c.MarshalAndSet(context.Background(), tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MarshalAndSet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
