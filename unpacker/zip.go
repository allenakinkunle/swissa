package unpacker

import (
	"archive/zip"
	"fmt"
	"io"
)

// ZipUnpacker unpacks files in a ZIP file
type ZipUnpacker struct {
	reader *zip.Reader
}

// NewZipUnpacker creates a new instance of ZipUnpacker and returns the pointer to it
func NewZipUnpacker(r io.ReaderAt, size int64) (*ZipUnpacker, error) {
	reader, err := zip.NewReader(r, size)

	if err != nil {
		return nil, fmt.Errorf("could not create ZipUnpacker, %w", err)
	}

	return &ZipUnpacker{
		reader: reader,
	}, nil
}
