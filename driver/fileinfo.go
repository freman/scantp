package driver

import (
	"os"
	"time"
)

type fileInfo struct {
	name string
}

func (f fileInfo) Name() string {
	return f.name
}
func (f fileInfo) Size() int64 {
	return 0
}
func (f fileInfo) Mode() os.FileMode {
	return 0555
}
func (f fileInfo) ModTime() time.Time {
	return time.Time{}
}
func (f fileInfo) IsDir() bool {
	return true
}
func (f fileInfo) Sys() interface{} {
	return nil
}

func (f fileInfo) Owner() string {
	return "vfs"
}

func (f fileInfo) Group() string {
	return "vfs"
}
