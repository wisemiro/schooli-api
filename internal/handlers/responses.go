package handlers

import (
	"net/http"
	"time"
)

const (
	SuccessMessage string = "success"
)

type SuccessResponse struct {
	TimeStamp time.Time   `json:"time_stamp"`
	Message   string      `json:"message"`
	Status    int         `json:"status"`
	Data      interface{} `json:"data"`
}

func NewStatusOkResponse(message string, data any) *SuccessResponse {
	return &SuccessResponse{
		TimeStamp: time.Now(),
		Message:   message,
		Status:    http.StatusOK,
		Data:      data,
	}
}

func NewStatusCreatedResponse(message string, data any) *SuccessResponse {
	return &SuccessResponse{
		TimeStamp: time.Now(),
		Message:   message,
		Status:    http.StatusCreated,
		Data:      data,
	}
}

func NewDeleteResponse(message string, data any) *SuccessResponse {
	return &SuccessResponse{
		TimeStamp: time.Now(),
		Message:   message,
		Status:    http.StatusNoContent,
		Data:      data,
	}
}
