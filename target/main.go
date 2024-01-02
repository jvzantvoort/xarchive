package target

import (
	"fmt"
	"os"
	"path/filepath"
)

type Target struct {
	Path string
	Name string
	Size int64
	Hash string
}

func FatalErr(err error) {
	if err == nil {
		return
	}
	fmt.Printf("Error %s\n", err)
	os.Exit(1)

}

func NewTarget(path string) (*Target, error) {
	var err error
	retv := &Target{}
	retv.Path = path

	abspath, err := filepath.Abs(retv.Path)
	if err != nil {
		return retv, err
	}

	retv.Path = abspath
	retv.Size, err = GetFileSize(retv.Path)
	FatalErr(err)
	retv.Hash, err = GetSHA256(retv.Path)
	FatalErr(err)
	retv.Name = filepath.Base(path)

	return retv, nil

}
