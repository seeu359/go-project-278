package handlers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/seeu359/go-project-278/handlers"
	"github.com/seeu359/go-project-278/testutil"
	"github.com/seeu359/go-project-278/testutil/factory"
)

type RespError struct {
	Error string `json:"error"`
}

func TestHandler_GetLinks(t *testing.T) {
	q := testutil.Tx(t, TestPool)

	require.NoError(t, factory.CreateLinks(t, q))

	router := gin.Default()
	router.GET("/api/links", handlers.New(q).GetLinks)

	req := httptest.NewRequest(http.MethodGet, "/api/links", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var links []handlers.GetLinkResponse
	err := json.Unmarshal(w.Body.Bytes(), &links)

	require.NoError(t, err)
	require.Len(t, links, 2)
}

func TestHandler_GetLink(t *testing.T) {
	q := testutil.Tx(t, TestPool)

	id, err := factory.CreateLink(t, q)
	require.NoError(t, err)

	router := gin.Default()
	router.GET("/api/links/:id", handlers.New(q).GetLinkById)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/links/%d", id), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var link handlers.GetLinkResponse
	err = json.Unmarshal(w.Body.Bytes(), &link)

	require.NoError(t, err)
	assert.Equal(t, "short", link.ShortName)
	assert.Equal(t, "https://ya.ru", link.URL)
}

func TestHandler_GetLinkErrorNotFound(t *testing.T) {
	q := testutil.Tx(t, TestPool)
	id, err := factory.CreateLink(t, q)
	require.NoError(t, err)

	router := gin.Default()
	router.GET("/api/links/:id", handlers.New(q).GetLinkById)

	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/links/%d", id+100), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var respErr RespError
	err = json.Unmarshal(w.Body.Bytes(), &respErr)
	assert.Equal(t, "Link Not Found", respErr.Error)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
