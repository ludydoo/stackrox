package handler

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
	"github.com/stackrox/rox/central/scannerdefinitions/file"
	"github.com/stackrox/rox/pkg/concurrency"
	"github.com/stackrox/rox/pkg/sync"
	"github.com/stackrox/rox/pkg/utils"
)

const (
	lastModifiedHeader    = "Last-Modified"
	ifModifiedSinceHeader = "If-Modified-Since"
)

// updater periodically updates a file by downloading the contents from the downloadURL.
type updater struct {
	file *file.Metadata

	client      *http.Client
	downloadURL string
	interval    time.Duration
	once        sync.Once
	stopSig     concurrency.Signal
}

// newUpdater creates a new updater.
func newUpdater(file *file.Metadata, client *http.Client, downloadURL string, interval time.Duration) *updater {
	return &updater{
		file:        file,
		client:      client,
		downloadURL: downloadURL,
		interval:    interval,
		stopSig:     concurrency.NewSignal(),
	}
}

// Stop stops the updater.
func (u *updater) Stop() {
	u.stopSig.Signal()
}

// Start starts the updater.
// The updater is only started once.
func (u *updater) Start() {
	u.once.Do(func() {
		// Run the first update in a blocking-manner.
		u.update()
		go u.runForever()
	})
}

func (u *updater) runForever() {
	t := time.NewTicker(u.interval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			u.update()
		case <-u.stopSig.Done():
			return
		}
	}
}

func (u *updater) update() {
	if err := u.doUpdate(); err != nil {
		log.Errorf("Scanner vulnerability updater for endpoint %q failed: %v", u.downloadURL, err)
	}
}

func (u *updater) doUpdate() error {
	req, err := http.NewRequest(http.MethodGet, u.downloadURL, nil)
	if err != nil {
		return errors.Wrap(err, "constructing request")
	}
	// No need to grab a read lock on u.file.LastModifiedTime
	// as a parallel write to this read is not possible.
	req.Header.Set(ifModifiedSinceHeader, u.file.GetLastModifiedTime().Format(http.TimeFormat))

	resp, err := u.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "executing request")
	}
	defer utils.IgnoreError(resp.Body.Close)

	if resp.StatusCode == http.StatusNotModified {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("invalid response from google storage; got code %d", resp.StatusCode)
	}

	lastModified, err := time.Parse(time.RFC1123, resp.Header.Get(lastModifiedHeader))
	if err != nil {
		return errors.Errorf("unable to determine upstream definitions file's modified time: %v", err)
	}

	return file.Write(u.file, resp.Body, lastModified)
}