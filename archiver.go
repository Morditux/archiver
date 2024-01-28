package main

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Archiver struct {
	Source   string
	Days     int
	Delete   bool
	Compress bool
	Name     string
}

var ToDelete = []string{}

func NewArchiver(source string, days int, delete bool, compress bool, name string) *Archiver {
	return &Archiver{Source: source, Days: days, Delete: delete, Compress: compress, Name: name}
}

func (a *Archiver) Archive() error {

	// Create a new archive
	file, err := os.Create(a.Name)
	if err != nil {
		return err
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()
	// Walk the directory tree recursively
	err = filepath.Walk(a.Source, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		// Check if the file is older than the specified days
		if info.ModTime().Before(time.Now().AddDate(0, 0, -a.Days)) {

			// Create a new file header
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}
			// remove leading / from name
			zpath := path
			if len(path) > 0 && path[0] == '/' {
				zpath = path[1:]
			}
			println(zpath)
			header.Name = zpath

			// Set the compression method to the one specified
			if a.Compress {
				header.Method = zip.Deflate
			} else {
				header.Method = zip.Store
			}

			// Create a new file in the archive
			writer, err := w.CreateHeader(header)
			if err != nil {
				return err
			}

			// Open the file
			file, err := os.Open(path)
			if err != nil {
				return err
			}

			// Copy the file to the archive
			_, err = io.Copy(writer, file)
			if err != nil {
				file.Close()
				return err
			}

			file.Close()
			// Delete the file if specified
			if a.Delete {
				ToDelete = append(ToDelete, path)
				//err = os.Remove(path)
				//if err != nil {
				//	return err
				//}
			}
		}
		// Delete the file if specified
		if a.Delete {
			a.DeleteFiles()
		}
		return nil
	})

	return err
}

func (a *Archiver) DeleteFiles() error {
	for _, path := range ToDelete {
		err := os.Remove(path)
		if err != nil {
			return err
		}
	}
	return nil
}
