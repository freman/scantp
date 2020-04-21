package seafile

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"goftp.io/server"
)

func (d *Driver) getLibraries(ctx context.Context, fn func(server.FileInfo) error) ([]Library, error) {
	req, err := d.newRequest(ctx, http.MethodGet, "api2/repos", nil)
	if err != nil {
		return nil, err
	}

	// TODO: could probably streamline this with streaming decoding
	if err := d.doRequest(req, jsonResponse(&d.libraries, httpStatusOk)); err != nil {
		return nil, err
	}

	d.librarymap = make(map[string]int, len(d.libraries))
	for i, v := range d.libraries {
		if fn != nil {
			if lerr := fn(v); err != nil {
				fn = nil
				err = lerr
			}
		}
		d.librarymap[v.FName] = i
	}

	return d.libraries, err
}

func (d *Driver) getLibrary(ctx context.Context, libraryName string) (*Library, error) {
	if d.librarymap != nil {
		if idx, found := d.librarymap[libraryName]; found {
			return &d.libraries[idx], nil
		}
	}

	if _, err := d.getLibraries(ctx, nil); err != nil {
		return nil, err
	}

	if idx, found := d.librarymap[libraryName]; found {
		return &d.libraries[idx], nil
	}

	return nil, os.ErrNotExist
}

func (d *Driver) getDirectoryDetail(ctx context.Context, libraryID, path string) (*directoryDetail, error) {
	query := url.Values{
		"path": []string{path},
	}.Encode()

	uri := "api/v2.1/repos/" + libraryID + "/dir/detail/?" + query
	req, err := d.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	var detail directoryDetail
	if err := d.doRequest(req, jsonResponse(&detail, httpStatusOk)); err != nil {
		return nil, err
	}

	return &detail, nil
}

func (d *Driver) listDirectoryEntries(ctx context.Context, libraryID, path string) ([]directoryEntry, error) {
	query := url.Values{
		"p": []string{path},
	}.Encode()

	uri := "api2/repos/" + libraryID + "/dir/?" + query
	req, err := d.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}

	// TODO: could probably streamline this with streaming json decoding and passing the path
	var entries []directoryEntry
	if err := d.doRequest(req, jsonResponse(&entries, httpStatusOk)); err != nil {
		return nil, err
	}

	return entries, nil
}

func (d *Driver) createUploadLink(ctx context.Context, libraryID, path string) (string, error) {
	query := url.Values{
		"p": []string{filepath.Dir(path)},
	}.Encode()

	uri := "api2/repos/" + libraryID + "/upload-link/?" + query
	req, err := d.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return "", err
	}

	var str string
	if err := d.doRequest(req, jsonResponse(&str, httpStatusOk)); err != nil {
		return "", err
	}

	return str, nil
}

func (d *Driver) createDirectory(ctx context.Context, libraryID, path string) error {
	form := url.Values{
		"operation": []string{"mkdir"},
	}.Encode()

	query := url.Values{
		"p": []string{path},
	}.Encode()

	uri := "api2/repos/" + libraryID + "/dir/" + query
	req, err := d.newRequest(ctx, http.MethodPost, uri, strings.NewReader(form))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var str string
	if err := d.doRequest(req, jsonResponse(&str, httpStatusCreated)); err != nil {
		return err
	}

	return nil
}

func (d *Driver) resolveURL(path string) string {
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	return d.baseURL.ResolveReference(u).String()
}

func (d *Driver) newUnauthenticatedRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error) {
	resolved := d.resolveURL(path)
	req, err := http.NewRequestWithContext(ctx, method, resolved, body)
	if err != nil {
		return req, err
	}

	req.Header.Add("Accept", "application/json; indent=4")
	return req, err
}

func (d *Driver) newRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	if d.token == "" {
		var err error
		d.once.Do(func() {
			if d.httpClient == nil {
				d.httpClient = &http.Client{
					Timeout: 10 * time.Second,
				}
			}

			if err = d.unauthenticatedPing(ctx); err != nil {
				return
			}
			if err = d.authenticate(ctx); err != nil {
				return
			}
			if err = d.ping(ctx); err != nil {
				return
			}
		})
		if err != nil {
			return nil, fmt.Errorf("login process failed: %w", err)
		}

		if d.token == "" {
			return nil, errors.New("no authentication token found, perhaps logging in failed")
		}
	}

	r, err := d.newUnauthenticatedRequest(ctx, method, url, body)
	if err != nil {
		return r, err
	}

	r.Header.Add("Authorization", "Token "+d.token)

	return r, err
}

func (d *Driver) unauthenticatedPing(ctx context.Context) error {
	req, err := d.newUnauthenticatedRequest(ctx, http.MethodGet, "api2/ping/", nil)
	if err != nil {
		return err
	}

	return d.expectPong(req)
}

func (d *Driver) ping(ctx context.Context) error {
	req, err := d.newRequest(ctx, http.MethodGet, "api2/auth/ping/", nil)
	if err != nil {
		return err
	}

	return d.expectPong(req)
}

func (d *Driver) expectPong(r *http.Request) error {
	var str string
	if err := d.doRequest(r, jsonResponse(&str, httpStatusOk)); err != nil {
		return err
	}

	if !strings.EqualFold(str, "pong") {
		return fmt.Errorf("expected pong got %q", str)
	}

	return nil
}

func (d *Driver) doRequest(r *http.Request, fn func(r *http.Response) error) error {
	resp, err := d.httpClient.Do(r)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return fn(resp)
	}
	return fmt.Errorf("API http status error %d %s", resp.StatusCode, resp.Status)
}

func (d *Driver) authenticate(ctx context.Context) error {
	req, err := d.newUnauthenticatedRequest(ctx, http.MethodPost, "api2/auth-token/", bytes.NewBufferString(url.Values{
		"username": []string{d.configuration.Username},
		"password": []string{d.configuration.Password},
	}.Encode()))

	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var tok struct {
		Token string `json:"token"`
	}

	if err := d.doRequest(req, jsonResponse(&tok, httpStatusOk)); err != nil {
		return err
	}

	if tok.Token == "" {
		return errors.New("returned empty token")
	}

	d.token = tok.Token

	return nil
}

func jsonResponse(v interface{}, statusCheck ...func(int) bool) func(r *http.Response) error {
	return func(r *http.Response) error {
		var ok bool
		for _, fn := range statusCheck {
			ok = ok || fn(r.StatusCode)
			if ok {
				break
			}
		}

		if !ok {
			return fmt.Errorf("unexpected http response %d %s", r.StatusCode, r.Status)
		}

		return json.NewDecoder(r.Body).Decode(v)
	}
}

func httpStatusOk(v int) bool {
	return v == http.StatusOK
}

func httpStatusCreated(v int) bool {
	return v == http.StatusCreated
}
