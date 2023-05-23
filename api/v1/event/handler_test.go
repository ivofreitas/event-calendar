package event

import (
	"blankfactor/event-calendar/log"
	"blankfactor/event-calendar/model"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	h *handler
)

func init() {
	h = NewHandler()
}

func TestOverlappingEvents(t *testing.T) {
	testCases := []struct {
		Name        string
		Request     *model.OverlappingRequest
		Response    *model.Response
		ResponseErr error
	}{
		{
			Name:    "Test Case 1",
			Request: &model.OverlappingRequest{Events: [][]int{{1, 2}, {3, 5}, {4, 7}, {6, 8}, {9, 10}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{1, 2}, {3, 8}, {9, 10}}}},
			},
		},
		{
			Name:    "Test Case 2",
			Request: &model.OverlappingRequest{Events: [][]int{{1, 3}, {2, 8}, {9, 10}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{1, 8}, {9, 10}}}},
			},
		},
		{
			Name:    "Test Case 3",
			Request: &model.OverlappingRequest{Events: [][]int{{1, 10}, {10, 20}, {20, 30}, {30, 40}, {40, 50}, {50, 60}, {60, 70}, {70, 80}, {80, 90}, {90, 100}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{1, 100}}}},
			},
		},
		{
			Name:    "Test Case 4",
			Request: &model.OverlappingRequest{Events: [][]int{{1, 10}, {11, 20}, {21, 30}, {31, 40}, {41, 50}, {51, 60}, {61, 70}, {71, 80}, {81, 90}, {91, 100}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{1, 10}, {11, 20}, {21, 30}, {31, 40}, {41, 50}, {51, 60}, {61, 70}, {71, 80}, {81, 90}, {91, 100}}}},
			},
		},
		{
			Name:    "Test Case 5",
			Request: &model.OverlappingRequest{Events: [][]int{{100, 105}, {1, 104}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{1, 105}}}},
			},
		},
		{
			Name:    "Test Case 6",
			Request: &model.OverlappingRequest{Events: [][]int{{-50, 20}, {70, 95}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{-50, 20}, {70, 95}}}},
			},
		},
		{
			Name:    "Test Case 7",
			Request: &model.OverlappingRequest{Events: [][]int{{-5, -4}, {-4, -3}, {-3, -2}, {-2, -1}, {-1, 0}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{-5, 0}}}},
			},
		},
		{
			Name:    "Test Case 8",
			Request: &model.OverlappingRequest{Events: [][]int{{43, 49}, {9, 12}, {12, 54}, {45, 90}, {91, 93}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{9, 90}, {91, 93}}}},
			},
		},
		{
			Name:    "Test Case 9",
			Request: &model.OverlappingRequest{Events: [][]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{0, 0}}}},
			},
		},
		{
			Name:    "Test Case 10",
			Request: &model.OverlappingRequest{Events: [][]int{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 1}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{0, 1}}}},
			},
		},
		{
			Name:    "Test Case 11",
			Request: &model.OverlappingRequest{Events: [][]int{{1, 22}, {-20, 30}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{-20, 30}}}},
			},
		},
		{
			Name:    "Test Case 12",
			Request: &model.OverlappingRequest{Events: [][]int{{20, 21}, {22, 23}, {0, 1}, {3, 4}, {23, 24}, {25, 27}, {5, 6}, {7, 19}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{0, 1}, {3, 4}, {5, 6}, {7, 19}, {20, 21}, {22, 24}, {25, 27}}}},
			},
		},
		{
			Name:    "Test Case 13",
			Request: &model.OverlappingRequest{Events: [][]int{{2, 3}, {4, 5}, {6, 7}, {8, 9}, {1, 10}}},
			Response: &model.Response{
				Meta: model.Meta{
					Offset:      0,
					Limit:       0,
					RecordCount: 1,
				},
				Records: []interface{}{&model.OverlappingResponse{Events: [][]int{{1, 10}}}},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			ctx = log.InitParams(ctx)

			response, err := h.OverlappingEvents(ctx, tc.Request)
			if tc.ResponseErr != nil {
				assert.Error(t, err)
				responseError := err.(*model.ResponseError)
				assert.Equal(t, tc.ResponseErr.Error(), responseError.DeveloperMessage)
			}

			assert.NoError(t, err)
			assert.NotNil(t, response)
			assert.Equal(t, tc.Response, response)
		})
	}
}
