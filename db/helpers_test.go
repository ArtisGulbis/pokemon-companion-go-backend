package db

import "testing"

// setupTest creates an in-memory database for testing
func setupTest(t *testing.T) *Database {
	database, err := New(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() { database.Close() })

	return database
}
