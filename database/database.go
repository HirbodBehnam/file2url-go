package database

import (
	"time"

	"github.com/google/uuid"
)

// Interface provides an interface to store the files in a database of choice
type Interface interface {
	// Store must store a File in database and return a unique ID mapped to file
	Store(File) (id uuid.UUID, err error)
	// Load must load a file from database. If it doesn't exist, returns false as exists
	Load(id uuid.UUID) (file File, exists bool)
	// Close must close the database
	Close() error
}

// File holds the info needed for a file
// See tg.InputDocumentFileLocation for more info about fields
type File struct {
	// When was this file added
	AddedTime     time.Time
	FileReference []byte
	// Simply the file name
	Name       string
	MimeType   string
	ID         int64
	AccessHash int64
	Size       int64
}
