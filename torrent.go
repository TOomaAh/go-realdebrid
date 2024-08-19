package gorealdebrid

import (
	"net/http"
	"net/url"
	"strconv"
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

	err = c.get(req, &torrents)

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

	err = c.get(req, &torrent)

	if err != nil {
		return nil, err
	}

	return &torrent, nil
}
