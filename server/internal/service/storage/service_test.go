package storage

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"blitzarx1/wisdom-fort/pkg/logger"
)

type mocks struct {
	mockStorage *MockkeyValStore
}

func initMocks(ctrl *gomock.Controller) *mocks {
	return &mocks{
		mockStorage: NewMockkeyValStore(ctrl),
	}
}

func TestService_StorageWithTTL(t *testing.T) {
	ctrl := gomock.NewController(t)
	m := initMocks(ctrl)

	key := "key"
	ttl := time.Second

	ctx := logger.WithCtx(context.Background(), logger.New(nil, "test"), "test_add_storage_with_ttl")
	s := build(ctx, func() keyValStore { return m.mockStorage })

	m.mockStorage.EXPECT().Set(key, uint(0))
	m.mockStorage.EXPECT().Delete(key)

	stID := s.AddStorageWithTTL(ctx, ttl)
	s.Set(stID, key, 0)

	time.Sleep(ttl + time.Second)
}
