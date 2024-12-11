package client

import (
	"encoding/json"
	"fmt"
	"goproject_Music/datastruct"
	"io"
	"net/http"
	"net/url"
	"time"
)

type client struct {
	Client *http.Client
	Host   string
	Path   string
}

func NewClient(host, path string) *client {
	c := &http.Client{Timeout: 30 * time.Second}
	return &client{Client: c, Host: host, Path: path}
}

func (c *client) GetSongFromClient(name, group string) (*datastruct.SongDetail, error) {
	detail := &datastruct.SongDetail{}
	u := url.URL{
		Scheme: "http",
		Host:   c.Host,
		Path:   c.Path,
	}

	q := u.Query()
	q.Set("name", name)
	q.Set("group", group)
	u.RawQuery = q.Encode()

	url := u.String()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	var b []byte
	b, err = c.doRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, detail)
	if err != nil {
		return nil, err
	}

	if (*detail == datastruct.SongDetail{}) {
		return nil, datastruct.ErrBadMusicGroup
	}

	return detail, err
}

func (c *client) doRequest(req *http.Request) ([]byte, error) {
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if http.StatusOK != resp.StatusCode {
		return nil, fmt.Errorf("%s, code: %d", body, resp.StatusCode)
	}

	return body, nil

}
