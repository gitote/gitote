package zip

import (
	"archive/zip"
	"os"
	"strings"
)

// OpenFile is the generalized open call; most users will use Open
// instead. It opens the named zip file with specified flag
// (O_RDONLY etc.) if applicable. If successful,
// methods on the returned ZipArchive can be used for I/O.
// If there is an error, it will be of type *PathError.
func (z *ZipArchive) Open(name string, flag int, perm os.FileMode) error {
	// Create a new archive if it's specified and not exist.
	if flag&os.O_CREATE != 0 {
		f, err := os.Create(name)
		if err != nil {
			return err
		}
		zw := zip.NewWriter(f)
		if err = zw.Close(); err != nil {
			return err
		}
	}

	rc, err := zip.OpenReader(name)
	if err != nil {
		return err
	}

	z.ReadCloser = rc
	z.FileName = name
	z.Comment = rc.Comment
	z.NumFiles = len(rc.File)
	z.Flag = flag
	z.Permission = perm
	z.isHasChanged = false

	z.files = make([]*File, z.NumFiles)
	for i, f := range rc.File {
		z.files[i] = &File{}
		z.files[i].FileHeader, err = zip.FileInfoHeader(f.FileInfo())
		if err != nil {
			return err
		}
		z.files[i].Name = strings.Replace(f.Name, "\\", "/", -1)
		if f.FileInfo().IsDir() && !strings.HasSuffix(z.files[i].Name, "/") {
			z.files[i].Name += "/"
		}
	}
	return nil
}
