package unityweb

import (
	"os"
	"path/filepath"
	"strings"
)

func bytesToUint32LE(b []byte) uint32 {
	return uint32(b[3])<<24 | uint32(b[2])<<16 | uint32(b[1])<<8 | uint32(b[0])
}

func uint32LEToBytes(i uint32) []byte {
	return []byte{byte(i), byte(i >> 8), byte(i >> 16), byte(i >> 24)}
}

func mkdirAllForFile(fp string) error {
	slashIndex := strings.LastIndex(fp, string(os.PathSeparator))
	if slashIndex == -1 {
		return nil
	}

	return os.MkdirAll(fp[:slashIndex], 0755)
}

func listFilesRecursively(p string) ([]string, error) {
	paths := make([]string, 0)

	err := filepath.Walk(p, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		pathWithoutBase := strings.TrimPrefix(path, p+string(os.PathSeparator))
		paths = append(paths, pathWithoutBase)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}
