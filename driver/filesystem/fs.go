package filesystem

import (
	"errors"
	"os"

	"goftp.io/server"
)

type Driver struct {
	*server.FileDriver
	configuration struct {
		RootPath string `toml:"root"`
	}
}

func NewDriver(fn func(v interface{}) error) (d *Driver, err error) {
	d = &Driver{}
	if err = fn(&d.configuration); err != nil {
		return nil, err
	}

	if d.configuration.RootPath == "" {
		return nil, errors.New("root path is required")
	}

	info, err := os.Stat(d.configuration.RootPath)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, errors.New("root path must be a directory")
	}

	d.FileDriver = &server.FileDriver{
		RootPath: d.configuration.RootPath,
		Perm:     &perm{},
	}

	return d, nil

}
