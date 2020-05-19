package seafile

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync"

	"goftp.io/server"
)

type Driver struct {
	token      string
	baseURL    *url.URL
	httpClient *http.Client
	once       sync.Once

	configuration struct {
		Username string `toml:"username"`
		Password string `toml:"password"`
		APIBase  string `toml:"api"`
	}

	libraries  []Library
	librarymap map[string]int
}

func libraryPrefix(path string) (libraryName, realPath string) {
	path = strings.Trim(path, "/")
	i := strings.Index(path, "/")
	if i == -1 {
		return path, ""
	}
	return path[:i], path[i:]
}

func (d *Driver) Stat(path string) (server.FileInfo, error) {
	if path == "/" {
		return &fileInfo{name: "/", isDir: true}, nil
	}

	libraryName, subPath := libraryPrefix(path)
	lib, err := d.getLibrary(context.TODO(), libraryName)
	if err != nil {
		return nil, err
	}

	if subPath == "" || subPath == "/" {
		return lib, nil
	}

	// TODO: stat files not just dirs?
	return d.getDirectoryDetail(context.TODO(), lib.ID, subPath)
}

func (d *Driver) ListDir(path string, fn func(server.FileInfo) error) error {
	if path == "/" {
		_, err := d.getLibraries(context.TODO(), fn)
		return err
	}

	libraryName, subPath := libraryPrefix(path)
	lib, err := d.getLibrary(context.TODO(), libraryName)
	if err != nil {
		return err
	}

	entries, err := d.listDirectoryEntries(context.TODO(), lib.ID, subPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if err := fn(entry); err != nil {
			return err
		}
	}

	return nil
}

func (d *Driver) MakeDir(path string) error {
	libraryName, subPath := libraryPrefix(path)
	if libraryName == "" {
		return errors.New("Must upload outside of root")
	}

	lib, err := d.getLibrary(context.TODO(), libraryName)
	if err != nil {
		return err
	}

	if subPath == "/" || subPath == "" {
		return errors.New("Make a new directory")
	}

	if err := d.createDirectory(context.TODO(), lib.ID, subPath); err != nil {
		return err
	}

	return nil
}

func (d *Driver) PutFile(path string, stream io.Reader, abool bool) (int64, error) {
	libraryName, subPath := libraryPrefix(path)
	if libraryName == "" {
		return 0, errors.New("Must upload outside of root")
	}

	lib, err := d.getLibrary(context.TODO(), libraryName)
	if err != nil {
		return 0, err
	}

	link, err := d.createUploadLink(context.TODO(), lib.ID, subPath)
	if err != nil {
		return 0, err
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", subPath)
	if err != nil {
		return 0, err
	}
	sz, err := io.Copy(part, stream)
	if err != nil {
		return 0, err
	}
	if err := writer.WriteField("parent_dir", filepath.Dir(subPath)); err != nil {
		return 0, err
	}

	err = writer.Close()
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequest(http.MethodPost, link+"?ret-json=1", body)
	if err != nil {
		return 0, err
	}
	req.Header.Add("Authorization", "Token "+d.token)
	req.Header.Add("Accept", "application/json; indent=4")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return 0, errors.New(resp.Status)
	}
	return sz, nil
}

func NewDriver(fn func(v interface{}) error) (d *Driver, err error) {
	d = &Driver{}
	if err = fn(&d.configuration); err != nil {
		return nil, err
	}

	if d.configuration.Username == "" || d.configuration.Password == "" || d.configuration.APIBase == "" {
		return nil, errors.New("Configuration for seafile is invalid, required username, password and api")
	}

	d.baseURL, err = url.Parse(d.configuration.APIBase)
	if err != nil {
		return nil, fmt.Errorf(`failure while parsing SeafileURL: %w`, err)
	}

	if strings.Contains(d.baseURL.Path, "api") {
		return nil, errors.New(`Please provide only the url for your seafile installtion, leave off api any api2`)
	}

	return d, nil

}
