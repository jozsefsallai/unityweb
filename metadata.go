package unityweb

import "os"

const (
	offsetFieldBytes       = 4
	sizeFieldBytes         = 4
	filenameSizeFieldBytes = 4
)

// FileMetadata contains information about a single file in a Unity Web data
// file.
type FileMetadata struct {
	// Offset represents at which byte offset the file's contents begin. Until
	// recalculated, this offset will always be zero.
	Offset uint32 // 4 bytes

	// Size represents the size of the file's contents in bytes.
	Size uint32 // 4 bytes

	// FilenameLength represents the length of the file's name in bytes.
	FilenameLength uint32 // 4 bytes

	// Filename is a byte slice containing the file's name. It is exactly the
	// same length as FilenameLength.
	Filename []byte // x bytes
}

// NewFileMetadata takes in a filename as an argument and its byte buffer, and
// initializes a FileMetadata object with the given information. The offset in
// the created object will be zero until recalculated in the Package object.
func NewFileMetadata(filename string, contents []byte) *FileMetadata {
	return &FileMetadata{
		Offset:         0, // we don't know the offset ahead of time
		Size:           uint32(len(contents)),
		FilenameLength: uint32(len(filename)),
		Filename:       []byte(filename),
	}
}

// BlockSize calculates the size of the FileMetadata block. The size consists of
// the size of each field in the struct (4 bytes each), plus the length of the
// filename.
func (m *FileMetadata) BlockSize() uint32 {
	return m.FilenameLength + offsetFieldBytes + sizeFieldBytes + filenameSizeFieldBytes
}

// FromPackageFile will populate the FileMetadata object with data read from the
// given package file.
func (m *FileMetadata) FromPackageFile(file *os.File) error {
	offsetField := make([]byte, offsetFieldBytes)

	if _, err := file.Read(offsetField); err != nil {
		return err
	}

	m.Offset = bytesToUint32LE(offsetField)

	sizeField := make([]byte, sizeFieldBytes)

	if _, err := file.Read(sizeField); err != nil {
		return err
	}

	m.Size = bytesToUint32LE(sizeField)

	filenameSizeField := make([]byte, filenameSizeFieldBytes)

	if _, err := file.Read(filenameSizeField); err != nil {
		return err
	}

	m.FilenameLength = bytesToUint32LE(filenameSizeField)

	m.Filename = make([]byte, m.FilenameLength)

	if _, err := file.Read(m.Filename); err != nil {
		return err
	}

	return nil
}

// ToBytes will return a byte slice containing the FileMetadata object's data.
// These bytes can then be written to a package file.
func (m *FileMetadata) ToBytes() []byte {
	bytes := make([]byte, m.BlockSize())

	offsetBytes := uint32LEToBytes(m.Offset)
	sizeBytes := uint32LEToBytes(m.Size)
	filenameSizeBytes := uint32LEToBytes(m.FilenameLength)

	copy(bytes, offsetBytes)
	copy(bytes[offsetFieldBytes:], sizeBytes)
	copy(bytes[offsetFieldBytes+sizeFieldBytes:], filenameSizeBytes)
	copy(bytes[offsetFieldBytes+sizeFieldBytes+filenameSizeFieldBytes:], m.Filename)

	return bytes
}
