package model

type OverlappingRequest struct {
	Events [][]int `json:"events"  validate:"required"`
}

type OverlappingResponse struct {
	Events [][]int `json:"overlapping_events"`
}
