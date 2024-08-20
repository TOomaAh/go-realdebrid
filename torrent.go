package gorealdebrid

import (
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Torrent struct {
	ID       string     `json:"id"`
	Filename string     `json:"filename"`
	Hash     string     `json:"hash"`
	Bytes    int        `json:"bytes"`
	Host     string     `json:"host"`
	Split    int        `json:"split"`
	Progress float64    `json:"progress"`
	Status   string     `json:"status"`
	Added    time.Time  `json:"added"`
	Links    []string   `json:"links"`
	Ended    *time.Time `json:"ended,omitempty"`
	Speed    *int       `json:"speed,omitempty"`
	Seeders  *int       `json:"seeders,omitempty"`
}

type TorrentFilter string

const (
	TorrentFilterActive TorrentFilter = "active"
)

type TorrentOptions struct {
	Offset int           `json:"offset,omitempty"`
	Page   int           `json:"page,omitempty"`
	Limit  int           `json:"limit,omitempty"`
	Filter TorrentFilter `json:"filter,omitempty"`
}

type Link struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	MimeType string `json:"mimeType"`
	Link     string `json:"link"`
	Host     string `json:"host"`
	Download string `json:"download"`

	Chunks   int64 `json:"chunks"`
	Crc      int64 `json:"crc"`
	FileSize int64 `json:"fileSize"`

	Streamable int `json:"streamable"`
}

type AddTorrent struct {
	Id  string `json:"id"`
	Uri string `json:"uri"`
}

// Get user torrents list
func (c *RealDebridClient) GetTorrents(options *TorrentOptions) ([]Torrent, error) {

	q := url.Values{}

	if options != nil {
		if options.Offset != 0 {
			q.Add("offset", strconv.Itoa(options.Offset))
		}
		if options.Page != 0 {
			q.Add("page", strconv.Itoa(options.Page))
		}
		if options.Limit != 0 {
			q.Add("limit", strconv.Itoa(options.Limit))
		}
		if options.Filter != "" {
			q.Add("filter", string(options.Filter))
		}
	}
	var req *http.Request
	var err error
	if len(q) == 0 {
		req, err = c.newRequest(http.MethodGet, "/torrents", nil, "", nil)
	} else {
		req, err = c.newRequest(http.MethodGet, "/torrents?", nil, q.Encode(), nil)
	}

	if err != nil {
		return nil, err
	}

	var torrents []Torrent

	err = c.do(req, &torrents)

	if err != nil {
		return nil, err
	}

	return torrents, nil
}

// Get all informations on the asked torrent
func (c *RealDebridClient) GetTorrent(id string) (*Torrent, error) {
	req, err := c.newRequest(http.MethodGet, "/torrents/info/"+id, nil, "", nil)

	if err != nil {
		return nil, err
	}

	var torrent Torrent

	err = c.do(req, &torrent)

	if err != nil {
		return nil, err
	}

	return &torrent, nil
}

// Add a torrent file to download, return a 201 HTTP code.
func (c *RealDebridClient) AddTorrent(file io.Reader) (*AddTorrent, error) {
	var addTorrent AddTorrent
	headers := http.Header{}
	headers.Set("Content-Type", "application/x-bittorrent")
	req, err := c.newRequest(http.MethodPost, "/torrents/addTorrent", headers, "", file)

	if err != nil {
		return nil, err
	}

	err = c.do(req, &addTorrent)

	if err != nil {
		return nil, err
	}

	return &addTorrent, nil

}

func (t *RealDebridClient) DebridTorrent(link string) (*Link, error) {
	var l Link
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	body := url.Values{}
	body.Set("link", link)

	req, err := t.newRequest(http.MethodPost, "/unrestrict/link", header, "", strings.NewReader(body.Encode()))

	if err != nil {
		return nil, err
	}

	err = t.do(req, &l)

	if err != nil {
		return nil, err
	}

	return &l, nil
}

func (t *RealDebridClient) AcceptTorrent(id string) error {
	header := http.Header{}
	header.Set("Content-Type", "application/x-www-form-urlencoded")

	body := url.Values{}
	body.Set("files", "all")

	req, err := t.newRequest(http.MethodPost, "/torrents/selectFiles/"+id, header, "", strings.NewReader(body.Encode()))

	if err != nil {
		return err
	}

	err = t.do(req, nil)

	if err != nil {
		return err
	}

	return nil
}
