package storagepg_test

import (
	"os"
	"testing"
)

var (
	databaseURL string
)

func TestMain(m *testing.M) {
	databaseURL = os.Getenv("DATABASE_TEST_URL")
	if databaseURL == "" {
		databaseURL = "host=localhost dbname=postgres_test user=postgres password=postgres sslmode=disable"
	}

	os.Exit(m.Run())
}
