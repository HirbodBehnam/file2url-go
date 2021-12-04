package database

// Interface provides an interface to store the files in a database of choice
type Interface interface {
	// Store must store a File in database and return a unique ID mapped to file
	Store(File) (id string)
	// Load must load a file from database. If it doesn't exist, returns false as exists
	Load(id string) (file File, exists bool)
}

// File holds the info needed for a file
// See tg.InputDocumentFileLocation for more info about fields
type File struct {
	FileReference []byte
	// Simply the file name
	Name       string
	MimeType   string
	ID         int64
	AccessHash int64
	Size       int64
}
