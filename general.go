package remtools

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

/*
GeneralContext is an implementation of Context for general platform.
*/
type GeneralContext struct {
	trashBoxPath string
}

/*
Path returns the trash box path.
*/
func (context *GeneralContext) Path() string {
	return context.trashBoxPath
}

/*
InitPath builds the trash box path.
*/
func (context *GeneralContext) InitPath() {
	home, _ := homedir.Dir()
	path := filepath.Join(home, ".Trash")
	context.trashBoxPath = path
	if !Exists(path) {
		os.Mkdir(path, 0755)
	}
}

/*
Move function moves file to trash box.
*/
func (context *GeneralContext) Move(path string) bool {
	toPath := filepath.Join(context.Path(), filepath.Base(path))
	err := os.Link(path, toPath)
	if err != nil {
		os.Remove(path)
		return true
	}
	return copy(path, toPath)
}

/*
Copy function copies file to trash box.
*/
func (context *GeneralContext) Copy(path string) bool {
	toPath := filepath.Join(context.Path(), filepath.Base(path))
	return copy(path, toPath)
}

func removeFile(dir string, info os.FileInfo) bool {
	path := filepath.Join(dir, info.Name())
	if !info.IsDir() {
		return os.Remove(path) == nil
	}
	return os.RemoveAll(path) == nil
}

/*
EmptyTrash remove all files in the trash box.
*/
func (context *GeneralContext) EmptyTrash() bool {
	path := context.Path()
	infos, _ := ioutil.ReadDir(path)
	ok := true
	for _, info := range infos {
		ok = ok && removeFile(path, info)
	}
	return ok
}
