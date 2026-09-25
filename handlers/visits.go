package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/seeu359/go-project-278/links"
	"github.com/seeu359/go-project-278/utils"
)

type GetVisitResponse struct {
	ID        int64     `json:"id"`
	LinkID    int64     `json:"link_id"`
	CreatedAt time.Time `json:"created_at"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	Status    int       `json:"status"`
}

func (h *Handler) RedirectTo(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is empty"})
		return
	}
	l, err := h.q.GetLinkByCode(c, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": LinkNotFoundError.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	status := http.StatusMovedPermanently
	_, err = h.q.AddVisit(c, links.AddVisitParams{
		LinkID:    int32(l.ID),
		Ip:        c.ClientIP(),
		UserAgent: pgtype.Text{String: c.Request.UserAgent(), Valid: true},
		Status:    pgtype.Int4{Int32: int32(status), Valid: true},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.Redirect(status, l.Url)
}

func (h *Handler) GetVisits(c *gin.Context) {
	dbVisits, err := h.q.GetVisits(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	pagination := c.Query("range")
	if pagination != "" {
		dbVisits, err = getPaginated(pagination, dbVisits)
		if err != nil {
			c.JSON(
				http.StatusBadRequest,
				gin.H{
					"error": utils.InvalidPaginationParamsError.Error(),
				},
			)
			return
		}
	}
	resp := make([]GetVisitResponse, 0, len(dbVisits))
	for _, v := range dbVisits {
		resp = append(resp, GetVisitResponse{
			ID:        v.ID,
			LinkID:    int64(v.LinkID),
			CreatedAt: v.CreatedAt.Time,
			IP:        v.Ip,
			UserAgent: v.UserAgent.String,
			Status:    int(v.Status.Int32),
		})
	}
	c.JSON(http.StatusOK, resp)
}
