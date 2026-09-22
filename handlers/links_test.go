package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/seeu359/go-project-278/handlers"
	"github.com/seeu359/go-project-278/links"
	"github.com/seeu359/go-project-278/testutil"
	"github.com/seeu359/go-project-278/testutil/factory"
)

func TestHandler_GetLinks(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testutil.NewPostgres(t)
	q := links.New(pool)
	handler := handlers.New(q)
	require.NoError(t, factory.CreateLinks(t, q))

	router := gin.Default()
	router.GET("/api/links", handler.GetLinks)

	req := httptest.NewRequest(http.MethodGet, "/api/links", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var links []handlers.GetLinkResponse
	err := json.Unmarshal(w.Body.Bytes(), &links)

	require.NoError(t, err)
	require.Len(t, links, 2)
}

func TestHandler_GetLink(t *testing.T) {
	gin.SetMode(gin.TestMode)
	pool := testutil.NewPostgres(t)
	q := links.New(pool)
	handler := handlers.New(q)
	require.NoError(t, factory.CreateLinks(t, q))

	router := gin.Default()
	router.GET("/api/links/1", handler.GetLinkById)

	req := httptest.NewRequest(http.MethodGet, "/api/links/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var link handlers.GetLinkResponse
	err := json.Unmarshal(w.Body.Bytes(), &link)

	require.NoError(t, err)
	assert.Equal(t, link.ShortName, "short")
	assert.Equal(t, link.URL, "https://ya.ru")
}
