package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/seeu359/go-project-278/links"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type CreateLinkRequest struct {
	URL       string `json:"original_url"`
	ShortName string `json:"short_name"`
}

type GetLinkResponse struct {
	ID        int64  `json:"id"`
	ShortName string `json:"short_name"`
	URL       string `json:"original_url"`
}

var LinkNotFoundError = errors.New("Link Not Found")

func (h *Handler) CreateLink(c *gin.Context) {
	var req CreateLinkRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	link := links.CreateLinkParams{Url: req.URL, ShortName: req.ShortName}
	_, err := h.q.CreateLink(c.Request.Context(), link)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handler) GetLinkById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	link, err := h.q.GetLink(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": LinkNotFoundError.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := GetLinkResponse{
		ID:        link.ID,
		ShortName: link.ShortName,
		URL:       link.Url,
	}
	c.JSON(http.StatusOK, response)
}
