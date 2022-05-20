package external

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
)

func UntargzFile(source io.Reader, fileName, outputFileName string) error {
	archive, err := gzip.NewReader(source)
	if err != nil {
		return err
	}
	defer archive.Close()

	tarReader := tar.NewReader(archive)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		info := header.FileInfo()
		if info.Name() == fileName {
			file, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(file, tarReader)
			if err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}
