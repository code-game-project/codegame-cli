package external

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Bananenpro/cli"
)

func InstallProgram(name, filename, url, version, path string) (string, error) {
	cli.BeginLoading("Installing %s v%s...", name, version)

	exeName := fmt.Sprintf("%s_%s", name, strings.ReplaceAll(version, ".", "-"))
	if runtime.GOOS == "windows" {
		exeName = exeName + ".exe"
	}

	if _, err := os.Stat(filepath.Join(path, exeName)); err == nil {
		return exeName, nil
	}

	downloadFile := fmt.Sprintf("%s-%s-%s.tar.gz", name, runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		downloadFile = fmt.Sprintf("%s-%s-%s.zip", name, runtime.GOOS, runtime.GOARCH)
	}

	res, err := http.Get(fmt.Sprintf("%s/releases/download/v%s/%s", strings.TrimSuffix(url, "/"), version, downloadFile))
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return "", err
	}

	binaries, err := os.ReadDir(path)
	if err == nil {
		for _, b := range binaries {
			info, err := b.Info()
			if err == nil && info.Name() != exeName && strings.HasPrefix(info.Name(), fmt.Sprintf("%s_%s", filename, strings.Join(strings.Split(version, ".")[:2], "-"))) {
				os.Remove(filepath.Join(path, info.Name()))
			}
		}
	}

	if runtime.GOOS == "windows" {
		err = unzipFile(res.Body, filename+".exe", filepath.Join(path, exeName))
	} else {
		err = untargzFile(res.Body, filename, filepath.Join(path, exeName))
	}
	if err == nil {
		cli.FinishLoading()
	}
	return exeName, err
}

// untargzFile first decompresses source with gzip, then extracts the file with fileName into outputFileName.
func untargzFile(source io.Reader, fileName, outputFileName string) error {
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

// unzipFile first decompresses source with gzip, then extracts the file with fileName into outputFileName.
func unzipFile(source io.Reader, fileName, outputFileName string) error {
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
