package cmd

import (
	"log"
	"os"
)

func UploadTo(r Repository) error {

}

func Mkdir(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		log.Print("Mkdir : ", path)
		return os.Mkdir(path, 0755)
	}
	return nil
}

func uploadToGit(path string, fname string, monte []string, srcdir string) error {
	log.Print("Upload to Git repository (", path, ")")

	if err := Mkdir(path); err != nil {
		return err
	}

	if _, err := os.Stat(srcdir); err != nil {
		return err
	}

	fp, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

}

type Repository struct {
	Type        string
	Path        string
	DirPattern  string
	FilePattern string
}
