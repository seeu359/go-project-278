package handlers

import (
	"github.com/seeu359/go-project-278/links"
)

type Handler struct {
	q *links.Queries
}

func New(q *links.Queries) *Handler {
	return &Handler{q: q}
}
