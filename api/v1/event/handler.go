package event

import (
	"blankfactor/event-calendar/model"
	"context"
	"sort"
)

type handler struct{}

func NewHandler() *handler {
	return &handler{}
}

// OverlappingEvents
// @Summary calculate overlapping events.
// @Param key body model.OverlappingRequest true "request body"
// @Tags event
// @Accept json
// @Product json
// @Success 201 {object} model.Response{meta=model.Meta,records=[]model.OverlappingResponse}
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Router /event [post]
func (h *handler) OverlappingEvents(ctx context.Context, param interface{}) (interface{}, error) {
	request := param.(*model.OverlappingRequest)
	events := request.Events

	sort.Slice(events[:], func(i, j int) bool {
		for x := range events[i] {
			if events[i][x] == events[j][x] {
				continue
			}
			return events[i][x] < events[j][x]
		}
		return false
	})

	result := make([][]int, 0)
	result = append(result, events[0])

	for i := 1; i < len(events); i++ {
		current := events[i]
		predecessor := result[len(result)-1]
		if predecessor[1] >= current[0] {
			limit := Max(predecessor[1], current[1])
			result[len(result)-1] = []int{predecessor[0], limit}
		} else {
			result = append(result, events[i])
		}
	}

	response := new(model.OverlappingResponse)
	response.Events = result

	return model.NewResponse(0, 0, 1, []interface{}{response}), nil
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x

}
