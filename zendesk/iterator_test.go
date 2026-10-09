package zendesk

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test for NewPaginationOptions function
func TestNewPaginationOptions(t *testing.T) {
	opts := NewPaginationOptions()

	assert.Equal(t, 100, opts.PageSize)
	assert.Equal(t, true, opts.IsCBP)
}

// Test for HasMore function
func TestHasMore(t *testing.T) {
	iter := &Iterator[int]{
		hasMore: true,
	}

	result := iter.HasMore()

	assert.Equal(t, true, result)
}

// Mock functions for GetNext testing
func mockObpFunc(ctx context.Context, opts *OBPOptions) ([]int, Page, error) {
	nextPage := "2"
	return []int{1, 2, 3}, Page{NextPage: &nextPage, Count: 3}, nil
}

func mockCbpFunc(ctx context.Context, opts *CBPOptions) ([]int, CursorPaginationMeta, error) {
	return []int{1, 2, 3}, CursorPaginationMeta{HasMore: true, AfterCursor: "3"}, nil
}

// Test for GetNext function
func TestGetNext(t *testing.T) {
	ctx := context.Background()

	iter := &Iterator[int]{
		pageSize:  2,
		hasMore:   true,
		isCBP:     false,
		pageIndex: 1,
		ctx:       ctx,
		obpFunc:   mockObpFunc,
		cbpFunc:   mockCbpFunc,
	}

	results, err := iter.GetNext()

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, results)
	assert.Equal(t, true, iter.HasMore())
	assert.Equal(t, 2, iter.pageIndex, "pageIndex should advance for the next OBP page")
}

func mockObpFuncError(ctx context.Context, opts *OBPOptions) ([]int, Page, error) {
	return nil, Page{}, errors.New("obp failed")
}

func mockCbpFuncError(ctx context.Context, opts *CBPOptions) ([]int, CursorPaginationMeta, error) {
	return nil, CursorPaginationMeta{}, errors.New("cbp failed")
}

func TestGetNextOBPFailureStopsIteration(t *testing.T) {
	iter := &Iterator[int]{
		hasMore: true,
		ctx:     context.Background(),
		obpFunc: mockObpFuncError,
	}

	results, err := iter.GetNext()

	assert.Error(t, err)
	assert.Nil(t, results)
	assert.False(t, iter.HasMore(), "HasMore must be false after an error so callers stop iterating")
}

func TestGetNextCBPFailureStopsIteration(t *testing.T) {
	iter := &Iterator[int]{
		hasMore: true,
		isCBP:   true,
		ctx:     context.Background(),
		cbpFunc: mockCbpFuncError,
	}

	results, err := iter.GetNext()

	assert.Error(t, err)
	assert.Nil(t, results)
	assert.False(t, iter.HasMore(), "HasMore must be false after an error so callers stop iterating")
}

func TestGetNextCBPAdvancesCursor(t *testing.T) {
	iter := &Iterator[int]{
		pageSize: 2,
		hasMore:  true,
		isCBP:    true,
		ctx:      context.Background(),
		cbpFunc:  mockCbpFunc,
	}

	results, err := iter.GetNext()

	assert.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, results)
	assert.Equal(t, "3", iter.pageAfter, "pageAfter should carry the cursor to the next page")
	assert.True(t, iter.HasMore())
}

func TestGetNextOBPStopsAtLastPage(t *testing.T) {
	iter := &Iterator[int]{
		hasMore: true,
		ctx:     context.Background(),
		obpFunc: func(ctx context.Context, opts *OBPOptions) ([]int, Page, error) {
			return []int{42}, Page{}, nil
		},
	}

	results, err := iter.GetNext()

	assert.NoError(t, err)
	assert.Equal(t, []int{42}, results)
	assert.False(t, iter.HasMore(), "no next page means iteration should stop")
}
