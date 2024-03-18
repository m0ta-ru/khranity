package native

import (
	"io"
	"os"
	"fmt"
	"path"
	"bufio"
	"strings"
	"archive/tar"
	"path/filepath"
	"compress/gzip"

	"khranity/app/ignore"
)

var ignoreFile = ".khranityignore"

func tarAppend(source, target string, ignores []string) error {
	ignoreFile = path.Join(source, ignoreFile)
	
	tarfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer tarfile.Close()

	gzball := gzip.NewWriter(tarfile)
	tarball := tar.NewWriter(gzball)
	defer tarball.Close()
	defer gzball.Close()

	info, err := os.Stat(source)
	if err != nil {
		return nil
	}

	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}

	return filepath.Walk(source,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			ignoreData, err := ignore.CompileIgnoreFileAndLines(ignoreFile, ignores...)
			if err == nil {
				if ignoreData.MatchesPath(path) {
					return nil // skip
				}
			} else {
				ignoreData := ignore.CompileIgnoreLines(ignores...)
				if ignoreData.MatchesPath(path) {
					return nil // skip
				}
			}

			var link string
			if info.Mode() & os.ModeSymlink == os.ModeSymlink {
				if link, err = os.Readlink(path); err != nil {
					return err
				}
			}
			
			header, err := tar.FileInfoHeader(info, link/*info.Name()*/)
			if err != nil {
				return err
			}

			if baseDir != "" {
				header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			}

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if !info.Mode().IsRegular() {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			
			_, err = io.Copy(tarball, file)
			return err
		})
}

func tarExtract(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()

	gzReader, err := gzip.NewReader(bufio.NewReader(reader))
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	//tarReader := tar.NewReader(reader)

	//count := 0
	for {
		header, err := tarReader.Next()
		// fmt.Printf("%v: %v (%v)\n", count, header, err)
		// count++
		if err == io.EOF {
			break
		} else if err != nil {
			// TODO with "archive/tar: invalid tar header"
			return err
			//return nil
		}

		path := filepath.Join(target, header.Name)
		info := header.FileInfo()
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func Compress(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}

	filename := filepath.Base(source)
	target = filepath.Join(target, fmt.Sprintf("%s.gz", filename))
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	return err
}

func Extract(source, target string) error {
	reader, err := os.Open(source)
	if err != nil {
		return err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer archive.Close()

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)
	return err
}

// func Append(source, target string) error {
// 	err := tarAppend(source, target)
// 	if (err != nil) {
// 		return err
// 	}
// 	return nil
// }

// func Extract(source, target string) error {
// 	err := tarExtract(source, target)
// 	if (err != nil) {
// 		return err
// 	}
// 	return nil
// }
