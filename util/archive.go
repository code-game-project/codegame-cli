package util

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

// UntargzFile first decompresses source with gzip, then extracts the file with fileName into outputFileName.
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
		if !info.IsDir() && info.Name() == fileName {
			file, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(file, tarReader)
			return err
		}
	}

	return errors.New("file not found")
}

// UnzipFile first decompresses source with gzip, then extracts the file with fileName into outputFileName.
func UnzipFile(source io.Reader, fileName, outputFileName string) error {
	data, err := ioutil.ReadAll(source)
	if err != nil {
		return err
	}

	reader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return err
	}

	for _, f := range reader.File {
		if !f.FileInfo().IsDir() && f.FileInfo().Name() == fileName {
			out, err := os.OpenFile(outputFileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, f.FileInfo().Mode())
			if err != nil {
				return err
			}
			defer out.Close()
			in, err := f.Open()
			if err != nil {
				return err
			}
			defer in.Close()

			_, err = io.Copy(out, in)
			return err
		}
	}

	return errors.New("file not found")
}
