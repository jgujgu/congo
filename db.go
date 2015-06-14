package congo

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/gogo/protobuf/proto"
	"github.com/gopheracademy/congo/internal"
)

//go:generate protoc --gogo_out=. internal/congo.proto

// DB represents the data store for all congo data
type DB struct {
	path string
	db   *bolt.DB
}

// NewDB returns a new instance of DB at the given file path.
func NewDB(path string) *DB {
	return &DB{path: path}
}

// Path returns the path the database was initialized with.
func (db *DB) Path() string { return db.path }

// Open opens and initializes the database.
func (db *DB) Open() error {
	// Create path if not exists.
	if err := os.MkdirAll(db.path, 0777); err != nil {
		return err
	}

	// Open underlying bolt database.
	d, err := bolt.Open(filepath.Join(db.path, "db"), 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	db.db = d

	// Initialize top level buckets.
	if err := db.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("users"))
		return nil
	}); err != nil {
		db.Close()
		return err
	}

	return nil
}

// Close shuts down the database.
func (db *DB) Close() error {
	if db.db != nil {
		return db.db.Close()
	}
	return nil
}

//func (db *DB) User(id int) (*User, error)
//func (db *DB) UserByUsername(username string) (*User, error)
//func (db *DB) Users() ([]*User, error)

// CreateUser creates a user and returns a new instance with its ID.
func (db *DB) CreateUser(u *User) (int, error) {
	// Create new user.
	var id uint64
	if err := db.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		id, _ = b.NextSequence()

		// Marshal user object into bytes.
		buf, err := proto.Marshal(&internal.User{
			ID:        proto.Uint64(id),
			FirstName: proto.String(u.FirstName),
			LastName:  proto.String(u.LastName),
			Email:     proto.String(u.Email),
		})
		if err != nil {
			return err
		}

		// Insert encoded user into id key.
		if err := b.Put(u64tob(id), buf); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return 0, err
	}

	return int(id), nil
}

//func (db *DB) Event(id int) (*Event, error)
//func (db *DB) Events() ([]*Event, error)
//func (db *DB) CreateEvent(*Event) error

//func (db *DB) CreateTicket(*User, *Event) (*Ticket, error)

func u64tob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, v)
	return b
}

func btou64(b []byte) uint64 { return binary.BigEndian.Uint64(b) }
