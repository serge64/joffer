package storagepg

import (
	"fmt"
	"strings"
	"testing"
)

func TestDB(t *testing.T, databaseURL string) (*Store, func(...string)) {
	t.Helper()

	store, err := New(databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	return store, func(tables ...string) {
		if len(tables) > 0 {
			query := fmt.Sprintf("TRUNCATE %s CASCADE", strings.Join(tables, ", "))
			store.db.Exec(query)
		}
		store.Close()
	}
}
