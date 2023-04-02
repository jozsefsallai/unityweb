package unityweb

import "os"

// File represents the contents of a single file in a Unity Web data file.
type File struct {
	// Contents is a byte slice that holds the contents of the file.
	Contents []byte
}

// NewFile creates a new File object with the given contents.
func NewFile(contents []byte) *File {
	return &File{
		Contents: contents,
	}
}

// Dump will write the contents of the file to the given path. If a deeply
// nested path is provided, the function will recursively create all directories
// in the path (if they don't already exist), then write the file with the given
// file name.
func (fo *File) Dump(path string) error {
	err := mkdirAllForFile(path)
	if err != nil {
		return err
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.Write(fo.Contents); err != nil {
		return err
	}

	return nil
}
