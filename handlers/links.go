package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/seeu359/go-project-278/links"
	"github.com/seeu359/go-project-278/utils"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type CreateLinkRequest struct {
	URL       string `json:"original_url" binding:"required,url"`
	ShortName string `json:"short_name" binding:"omitempty,min=3,max=32"`
}

type UpdateLinkRequest struct {
	CreateLinkRequest
}

type GetLinkResponse struct {
	ID        int64  `json:"id"`
	ShortName string `json:"short_name"`
	URL       string `json:"original_url"`
	ShortURL  string `json:"short_url"`
}

var LinkNotFoundError = errors.New("Link Not Found")
var ShortNameAlreadyInUseError = errors.New("short name already in use")

func (h *Handler) CreateLink(c *gin.Context) {
	var req CreateLinkRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	if req.ShortName == "" {
		shortName, err := links.RandomShortName(10)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		req.ShortName = shortName
	}

	link := links.CreateLinkParams{Url: req.URL, ShortName: req.ShortName}
	_, err := h.q.CreateLink(c.Request.Context(), link)
	if err != nil {
		if utils.IsUniqueViolationError(err) {
			c.JSON(http.StatusUnprocessableEntity,
				gin.H{
					"errors": gin.H{
						"short_name": ShortNameAlreadyInUseError.Error(),
					},
				},
			)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusCreated)
}

func getPaginated[T any](pagination string, entities []T) ([]T, error) {
	p, err := utils.FormatedPagination(pagination)
	fmt.Println("pagination", p)
	if err != nil {
		return nil, err

	}
	if p.Start >= len(entities) {
		return nil, nil
	}
	if p.End > len(entities) {
		p.End = len(entities)
	} else {
		p.End += 1
	}
	return entities[p.Start:p.End], nil
}

func (h *Handler) getDbLink(id int64) (links.Link, error) {
	ctx := context.Background()
	link, err := h.q.GetLink(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return link, LinkNotFoundError
		}
		return link, err
	}
	return link, nil
}

func (h *Handler) GetLinkById(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	dbLink, err := h.getDbLink(id)
	if err != nil {
		if errors.Is(err, LinkNotFoundError) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	shortURL, err := links.GetShortURL(dbLink.ShortName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	response := GetLinkResponse{
		ID:        dbLink.ID,
		ShortName: dbLink.ShortName,
		URL:       dbLink.Url,
		ShortURL:  shortURL,
	}
	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetLinks(c *gin.Context) {
	dbLinks, err := h.q.GetLinks(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pagination := c.Query("range")
	if pagination != "" {
		dbLinks, err = getPaginated(pagination, dbLinks)
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
	respLinks := make([]GetLinkResponse, 0, len(dbLinks))
	for _, l := range dbLinks {
		shortURL, err := links.GetShortURL(l.ShortName)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		respLinks = append(respLinks, GetLinkResponse{ID: l.ID, ShortName: l.ShortName, URL: l.Url, ShortURL: shortURL})
	}
	c.JSON(http.StatusOK, respLinks)
}

func (h *Handler) UpdateLink(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id passed"})
		return
	}

	var req UpdateLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	_, err = h.q.UpdateLink(
		c.Request.Context(),
		links.UpdateLinkParams{
			ID:        id,
			ShortName: req.ShortName,
			Url:       req.URL,
		},
	)
	if err != nil {
		if utils.IsUniqueViolationError(err) {
			c.JSON(http.StatusUnprocessableEntity,
				gin.H{
					"errors": gin.H{
						"short_name": ShortNameAlreadyInUseError.Error(),
					},
				},
			)
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (h *Handler) DeleteLink(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "incorrect id passed"})
		return
	}
	_, err = h.q.DeleteLink(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Status(http.StatusNoContent)
}
