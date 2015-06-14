package congo_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gopheracademy/congo"
)

// Ensure the database can create a user.
func TestDB_CreateUser(t *testing.T) {
	db := OpenDB()
	defer db.Close()

	// Create user.
	id, err := db.CreateUser(&congo.User{
		FirstName: "bob",
		LastName:  "smith",
		Email:     "bob@smith.com",
	})
	if err != nil {
		t.Fatal(err)
	} else if id != 1 {
		t.Fatalf("unexpected id: %d", id)
	}
}

// DB represents a test wrapper for congo.DB.
type DB struct {
	*congo.DB
}

// NewDB returns a new instance of DB in a temporary path.
func NewDB() *DB {
	path, _ := ioutil.TempDir("", "congo-")
	return &DB{congo.NewDB(path)}
}

// OpenDB returns an open instance of DB.
func OpenDB() *DB {
	db := NewDB()
	if err := db.Open(); err != nil {
		panic(err)
	}
	return db
}

// Close closes and deletes the database.
func (db *DB) Close() {
	db.DB.Close()
	os.RemoveAll(db.Path())
}
