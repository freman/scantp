package driver

import (
	"io"

	"goftp.io/server"
)

// Driver trims the interface of server.Driver cos we're not supporting a bunch of things
type Driver interface {
	Stat(string) (server.FileInfo, error)
	ListDir(string, func(server.FileInfo) error) error
	MakeDir(string) error
	PutFile(string, io.Reader, bool) (int64, error)
}
