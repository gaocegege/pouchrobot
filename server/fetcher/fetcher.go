package fetcher

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/pouchcontainer/pouchrobot/server/gh"
)

// FETCHINTERVAL refers the interval of fetch action
const FETCHINTERVAL = 3 * time.Minute

// Fetcher is a worker to periodically get elements from github.
type Fetcher struct {
	client *gh.Client
}

// New initializes a brand new fetch.
func New(client *gh.Client) *Fetcher {
	return &Fetcher{
		client: client,
	}
}

// Run starts periodical work
func (f *Fetcher) Run() {
	logrus.Info("start to run fetcher")
	for {
		f.CheckPRsConflict()
		time.Sleep(FETCHINTERVAL)
	}
}
