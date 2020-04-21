package driver

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/freman/scantp/driver/filesystem"
	"github.com/freman/scantp/driver/seafile"
	"goftp.io/server"
)

func (factory *MultipleDriverFactory) AddPath(name, driverName string, md toml.MetaData, primative toml.Primitive) error {
	if factory.drivers == nil {
		factory.drivers = make(map[string]Driver)
	}

	if _, exists := factory.drivers[name]; exists {
		return fmt.Errorf("paths must be unique %q is already used", name)
	}

	subDriver, err := factory.subDriver(driverName, func(v interface{}) error {
		return md.PrimitiveDecode(primative, v)
	})
	if err != nil {
		return err
	}

	factory.drivers[name] = subDriver
	return nil
}

func (factory *MultipleDriverFactory) subDriver(name string, fn func(v interface{}) error) (Driver, error) {
	switch name {
	case `seafile`:
		return seafile.NewDriver(fn)
	case `fs`, `filesystem`, `local`:
		return filesystem.NewDriver(fn)
	}

	return nil, fmt.Errorf("unknown driver %q", name)
}

var _ server.Driver = &MultipleDriver{}

type MultipleDriver struct {
	drivers map[string]Driver
}

func (driver *MultipleDriver) Stat(path string) (server.FileInfo, error) {
	driverName, realPath := driverPrefix(path)
	if driverName == "" || driverName == "/" {
		return &fileInfo{
			name: "/",
		}, nil
	}

	if subDriver, isa := driver.drivers[driverName]; isa {
		if realPath == "" || realPath == "/" {
			return &fileInfo{
				name: driverName,
			}, nil
		}

		return subDriver.Stat(realPath)
	}

	return nil, errors.New("Not a file")
}

func (driver *MultipleDriver) ListDir(path string, callback func(server.FileInfo) error) error {
	driverName, realPath := driverPrefix(path)
	fmt.Println(">", driverName, realPath)
	if driverName == "" || driverName == "/" {
		for name := range driver.drivers {
			if err := callback(&fileInfo{
				name: name,
			}); err != nil {
				return err
			}
		}
		return nil
	}

	if subDriver, isa := driver.drivers[driverName]; isa {
		return subDriver.ListDir(realPath, callback)
	}

	return errors.New("Not path with that name configured")
}

func (driver *MultipleDriver) DeleteDir(string) error {
	return errors.New("Permission Denied")
}

func (driver *MultipleDriver) DeleteFile(string) error {
	return errors.New("Permission Denied")
}

func (driver *MultipleDriver) Rename(string, string) error {
	return errors.New("Permission Denied")
}

func (driver *MultipleDriver) MakeDir(path string) error {
	driverName, realPath := driverPrefix(path)
	if driverName == "" || driverName == "/" {
		return errors.New("Virtual file system, not writable")
	}

	if subDriver, isa := driver.drivers[driverName]; isa {
		return subDriver.MakeDir(realPath)
	}

	return errors.New("Not path with that name configured")
}

func (driver *MultipleDriver) GetFile(string, int64) (int64, io.ReadCloser, error) {
	return 0, nil, errors.New("Permission Denied")
}

// PutFile implements Driver
func (driver *MultipleDriver) PutFile(destPath string, data io.Reader, appendData bool) (int64, error) {
	driverName, realPath := driverPrefix(destPath)
	if driverName == "" || driverName == "/" {
		return 0, errors.New("Virtual file system, not writable")
	}

	if subDriver, isa := driver.drivers[driverName]; isa {
		return subDriver.PutFile(realPath, data, appendData)
	}

	return 0, errors.New("unknown driver")
}

type MultipleDriverFactory struct {
	drivers map[string]Driver
}

func (factory *MultipleDriverFactory) NewDriver() (server.Driver, error) {
	return &MultipleDriver{factory.drivers}, nil
}

func driverPrefix(path string) (driverName, realPath string) {
	path = strings.Trim(path, "/")
	i := strings.Index(path, "/")
	if i == -1 {
		return path, "/"
	}
	return path[:i], path[i:]
}
