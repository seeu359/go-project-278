package handlers_test

import (
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seeu359/go-project-278/testutil"
)

var TestPool *pgxpool.Pool

func TestMain(m *testing.M) {
	pool, cleanup, err := testutil.StartPostgres()
	gin.SetMode(gin.TestMode)

	if err != nil {
		log.Fatal("Error during setup postgres")
	}

	TestPool = pool
	code := m.Run()
	cleanup()
	os.Exit(code)
}
