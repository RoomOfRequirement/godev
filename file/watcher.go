package file

import (
	"errors"
	"os"
	"time"
)

// Watcher struct
type Watcher struct {
	filepath      string
	lastEditTime  time.Time
	checkInterval time.Duration
	callback      func()

	closeChan chan struct{}
	closed    bool
}

// NewWatcher ...
func NewWatcher(filepath string, callback func(), checkInterval time.Duration) (*Watcher, error) {
	fi, err := os.Stat(filepath)
	if err != nil {
		return nil, err
	}
	if callback == nil {
		return nil, errors.New("no callback when file change")
	}
	if checkInterval == 0 {
		return nil, errors.New("invalid check interval")
	}
	w := &Watcher{
		filepath:      filepath,
		lastEditTime:  fi.ModTime(),
		checkInterval: checkInterval,
		callback:      callback,
		closeChan:     make(chan struct{}, 1),
		closed:        false,
	}
	go w.watch()
	return w, nil
}

func (w *Watcher) watch() {
	for {
		select {
		case <-w.closeChan:
			return
		case <-time.After(w.checkInterval):
			if w.closed {
				return
			}
			fi, err := os.Stat(w.filepath)
			// no err -> err, certainly changed
			if err != nil {
				// log.Println(err)
				w.callback()
			} else {
				// changed
				if t := fi.ModTime(); t.After(w.lastEditTime) {
					w.lastEditTime = t
					w.callback()
				}
				// unchanged
			}
		}
	}
}

// Close ...
func (w *Watcher) Close() {
	close(w.closeChan)
	w.closed = true
}
