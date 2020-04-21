package seafile

import (
	"os"
	"time"
)

type fileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
	isDir   bool
}

func (f fileInfo) Name() string {
	return f.name
}
func (f fileInfo) Size() int64 {
	return f.size
}
func (f fileInfo) Mode() os.FileMode {
	return f.mode
}
func (f fileInfo) ModTime() time.Time {
	return f.modTime
}
func (f fileInfo) IsDir() bool {
	return f.isDir
}
func (f fileInfo) Sys() interface{} {
	return nil
}

func (f fileInfo) Owner() string {
	return "nobody"
}

func (f fileInfo) Group() string {
	return "nobody"
}
