package file_store

import (
	"io"
	"os"
	"path"
)

type localAdapter struct {
	tag string
}

func newLocalAdapter(tag string) *localAdapter {
	return &localAdapter{tag: tag}
}

func (o *localAdapter) Store(rootPath, bucket, key string, size int64, fs io.ReadCloser) error {

	fullPath := path.Join(rootPath, bucket, key)
	previewFs, err := os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	buf := make([]byte, 1024*1024*3)
	for {

		n, err := fs.Read(buf)

		if err == io.EOF && n <= 0 {
			break
		}

		if n > 0 {
			_, _ = previewFs.Write(buf[:n])
		}

	}

	return nil
}
