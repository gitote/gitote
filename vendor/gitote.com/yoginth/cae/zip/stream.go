package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// A StreamArchive represents a streamable archive.
type StreamArchive struct {
	*zip.Writer
}

// NewStreamArachive returns a new streamable archive with given io.Writer.
// It's caller's responsibility to close io.Writer and streamer after operation.
func NewStreamArachive(w io.Writer) *StreamArchive {
	return &StreamArchive{zip.NewWriter(w)}
}

// StreamFile streams a file or directory entry into StreamArchive.
func (s *StreamArchive) StreamFile(relPath string, fi os.FileInfo, data []byte) error {
	if fi.IsDir() {
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		fh.Name = relPath + "/"
		if _, err = s.Writer.CreateHeader(fh); err != nil {
			return err
		}
	} else {
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return err
		}
		fh.Name = filepath.Join(relPath, fi.Name())
		fh.Method = zip.Deflate
		fw, err := s.Writer.CreateHeader(fh)
		if err != nil {
			return err
		} else if _, err = fw.Write(data); err != nil {
			return err
		}
	}
	return nil
}

// StreamReader streams data from io.Reader to StreamArchive.
func (s *StreamArchive) StreamReader(relPath string, fi os.FileInfo, r io.Reader) (err error) {
	fh, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}
	fh.Name = filepath.Join(relPath, fi.Name())

	fw, err := s.Writer.CreateHeader(fh)
	if err != nil {
		return err
	}
	_, err = io.Copy(fw, r)
	return err
}
