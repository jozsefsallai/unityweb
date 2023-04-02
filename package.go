package unityweb

import (
	"fmt"
	"os"
	"path"
)

const (
	magicFieldBytes       = 16
	startOffsetFieldBytes = 4
)

// Package is an object representation of a single Unity Web package file.
type Package struct {
	// Magic is the magic header of the package file. It must always be 16 bytes
	// long and contain the null-terminated string "UnityWebData1.0".
	Magic [16]byte // UnityWebData1.0, 16 bytes with NUL

	// StartOffset represents the byte offset at which the very first file's
	// contents begins. This will also always be the value of the offset field
	// in the first metadata object. It is represented as a 4-byte unsigned
	// integer in little-endian format.
	StartOffset uint32 // 4 bytes

	// FileMetadata is a slice of FileMetadata objects, each representing meta
	// information about a single file in the package.
	FileMetadata []FileMetadata // x * (4 + 4 + 4 + y) bytes

	// Files is a slice of File objects, each representing a single file in the
	// package. The number of files in this slice will always be the same as the
	// number of FileMetadata objects in the FileMetadata slice.
	Files []File // rest of the bytes
}

// NewPackage initializes a new Package object with the correct magic header.
// The offset fields, as well as the slices in the object will be zero/empty.
func NewPackage() *Package {
	return &Package{
		Magic: [16]byte{'U', 'n', 'i', 't', 'y', 'W', 'e', 'b', 'D', 'a', 't', 'a', '1', '.', '0', 0},
	}
}

// AddFile appends the metadata and contents of given byte slice with a given
// file name to the Package object.
func (p *Package) AddFile(filename string, contents []byte) {
	file := NewFile(contents)
	metadata := NewFileMetadata(filename, contents)
	p.Files = append(p.Files, *file)
	p.FileMetadata = append(p.FileMetadata, *metadata)
}

// RecalculateOffsets recalculates the offset fields in the FileMetadata slice
// and the StartOffset field in the Package object. This method should be called
// only after all files have been added.
func (p *Package) RecalculateOffsets() {
	if len(p.FileMetadata) == 0 {
		return
	}

	metadataBlockSizes := uint32(0)
	for i := range p.FileMetadata {
		metadataBlockSizes += p.FileMetadata[i].BlockSize()
	}

	p.StartOffset = magicFieldBytes + startOffsetFieldBytes + metadataBlockSizes

	offset := p.StartOffset

	for i := range p.FileMetadata {
		p.FileMetadata[i].Offset = offset
		offset += p.FileMetadata[i].Size
	}
}

// ReadFromPackageFile takes in an os.File object representing a Unity Web data
// package file, and parses it into a Package object. The file must be opened
// with read permissions. This method can fail if the file doesn't start with
// the correct magic header or if there are file I/O errors (such as premature
// EOF).
func (p *Package) ReadFromPackageFile(file *os.File) error {
	if _, err := file.Read(p.Magic[:]); err != nil {
		return err
	}

	if string(p.Magic[:len(p.Magic)-1]) != "UnityWebData1.0" {
		return ErrInvalidMagicHeader
	}

	startOffsetField := make([]byte, startOffsetFieldBytes)
	if _, err := file.Read(startOffsetField); err != nil {
		return err
	}

	p.StartOffset = bytesToUint32LE(startOffsetField)

	bytesRead := uint32(magicFieldBytes + startOffsetFieldBytes)

	for bytesRead < p.StartOffset {
		metadata := FileMetadata{}
		if err := metadata.FromPackageFile(file); err != nil {
			return err
		}

		p.FileMetadata = append(p.FileMetadata, metadata)
		bytesRead += metadata.BlockSize()
	}

	for i := range p.FileMetadata {
		fileObject := File{}

		buffer := make([]byte, p.FileMetadata[i].Size)
		if _, err := file.Read(buffer); err != nil {
			return err
		}

		fileObject.Contents = buffer

		p.Files = append(p.Files, fileObject)
	}

	return nil
}

// Dump will write all files in the package to the given directory. The method
// will recursively create directories if they don't exist. This method can fail
// if there are file I/O errors (such as permission errors).
func (p *Package) Dump(outputDirectoryPath string) error {
	for i := range p.Files {
		fullpath := path.Join(outputDirectoryPath, string(p.FileMetadata[i].Filename))
		if err := p.Files[i].Dump(fullpath); err != nil {
			return err
		}
	}

	return nil
}

// Pack will create a byte slice representing the package file. This method
// will make sure to recalculate the offsets before performing the packing.
func (p *Package) Pack() []byte {
	p.RecalculateOffsets()

	buffer := make([]byte, 0)

	buffer = append(buffer, p.Magic[:]...)
	buffer = append(buffer, uint32LEToBytes(p.StartOffset)...)

	for i := range p.FileMetadata {
		buffer = append(buffer, p.FileMetadata[i].ToBytes()...)
	}

	for i := range p.Files {
		buffer = append(buffer, p.Files[i].Contents...)
	}

	return buffer
}

// PackToFile will pack the package object into a given output file. The method
// will recursively create directories if they don't exist. This method can fail
// if there are file I/O errors (such as permission errors).
func (p *Package) PackToFile(filepath string) error {
	buffer := p.Pack()

	err := mkdirAllForFile(filepath)
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	if _, err := file.Write(buffer); err != nil {
		return err
	}

	return nil
}

// FromPackageFile takes in a file path as an argument and returns a Package
// object representing the data inside the Unity Web data package file.
func FromPackageFile(filename string) (*Package, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	pkg := NewPackage()
	err = pkg.ReadFromPackageFile(f)

	if err != nil {
		return nil, err
	}

	return pkg, nil
}

// PackDirectory takes in a directory path as an argument and constructs a new
// Package object with all the necessary metadata and file contents. This method
// can fail if there are file I/O errors (such as permission errors). After
// calling this method, you may call the PackToFile method to create a new Unity
// Web data package file.
func PackDirectory(directoryPath string) (*Package, error) {
	files, err := listFilesRecursively(directoryPath)
	if err != nil {
		return nil, err
	}

	pkg := NewPackage()

	for idx := range files {
		path := fmt.Sprintf("%s/%s", directoryPath, files[idx])

		file, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		pkg.AddFile(files[idx], file)
	}

	return pkg, nil
}
