package handlers

import "metrics-command/commands/handler"

type RatingHandler struct {
	handler handler.Handler
}

func NewRatingHandler(handler handler.Handler) *RatingHandler {
	return &RatingHandler{
		handler: handler,
	}
}
