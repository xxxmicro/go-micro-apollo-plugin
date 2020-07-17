package apollo

import (
	"errors"
	"github.com/micro/go-micro/v2/config/source"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/zouyx/agollo/v3/storage"
	"time"
)

type watcher struct {
	name      string
	exit      chan bool
	eventChan chan *storage.ChangeEvent
}

func (w *watcher) OnChange(changeEvent *storage.ChangeEvent) {
	log.Info("change listener.")
	log.Info(changeEvent.Changes)
	log.Info(changeEvent.Namespace)
	w.eventChan <- changeEvent
}

func (w *watcher) Next() (*source.ChangeSet, error) {
	select {
	case event := <-w.eventChan:
		log.Info(event.Changes)
		content := event.Changes["content"].NewValue
		log.Info("content")
		log.Info(content)

		cs := &source.ChangeSet{
			Timestamp: time.Now(),
			Format:    "yaml",
			Source:    w.name,
			Data:      []byte(content),
		}
		cs.Checksum = cs.Sum()
		return cs, nil
	case <-w.exit:
		return nil, errors.New("watcher stopped")
	}
}

func (w *watcher) Stop() error {
	select {
	case <-w.exit:
	default:
	}
	return nil
}

func newWatcher(name string) (*watcher, error) {
	return &watcher{
		name:      name,
		exit:      make(chan bool),
		eventChan: make(chan *storage.ChangeEvent),
	}, nil
}
