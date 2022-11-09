package infrastructure

import (
	"fmt"
	"math/rand"
	"os"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/domain"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

//AppBufSize buffer size for async hooks
var AppBufSize uint = 8192

//AppHook application log hook struct
type AppHook struct {
	repo       interfaces.IGenericRepository[uuid.UUID, domain.AppLog]
	mutex      sync.RWMutex
	InsertFunc func(entry *logrus.Entry, repository interfaces.IGenericRepository[uuid.UUID, domain.AppLog]) error
}

//AppAsyncHook application log async hook struct
type AppAsyncHook struct {
	*AppHook
	buf        chan *logrus.Entry
	flush      chan bool
	wg         sync.WaitGroup
	Ticker     *time.Ticker
	InsertFunc func(entry *logrus.Entry, repository interfaces.IGenericRepository[uuid.UUID, domain.AppLog]) error
}

var insertAppFunc = func(entry *logrus.Entry, repository interfaces.IGenericRepository[uuid.UUID, domain.AppLog]) error {
	if entry.Data["method"] == nil && !fromLogController(entry) {
		ent := newEntityFromApp(entry)
		_, err := repository.Insert(nil, *ent)
		if err != nil {
			return errors.Internal.Wrap(err, "error inserting log to db")
		}
	}
	return nil
}

var asyncAppInsertFunc = func(entry *logrus.Entry, repository interfaces.IGenericRepository[uuid.UUID, domain.AppLog]) error {
	if entry.Data["method"] != nil {
		ent := newEntityFromApp(entry)
		_, err := repository.Insert(nil, *ent)
		if err != nil {
			return errors.Internal.Wrap(err, "error inserting log to db")
		}
	}
	return nil
}

func newEntityFromApp(entry *logrus.Entry) *domain.AppLog {
	var actionID uuid.UUID
	var source string
	if entry.Data["actionID"] != nil {
		actionID = entry.Data["actionID"].(uuid.UUID)
	}
	if entry.Data["source"] != nil {
		source = entry.Data["source"].(string)
	}
	return &domain.AppLog{
		EntityUUID: domain.EntityUUID{
			ID: uuid.New(),
		},
		Level:    entry.Level.String(),
		Source:   source,
		ActionID: actionID,
		Message:  entry.Message,
	}
}

//NewAppHook create new gorm hook for application logs
//Params
//	repo - gorm generic repository with domain.AppLog instantiated
//Return
//	*AppHook - gorm hook
func NewAppHook(repo interfaces.IGenericRepository[uuid.UUID, domain.AppLog]) *AppHook {
	return &AppHook{
		repo:       repo,
		InsertFunc: insertAppFunc,
	}
}

//NewAsyncAppHook create new async gorm hook for application logs
//Params
//	repo - gorm generic repository with domain.AppLog instantiated
//Return
//	*AppHook - async gorm hook
func NewAsyncAppHook(repo interfaces.IGenericRepository[uuid.UUID, domain.AppLog]) *AppAsyncHook {
	hook := &AppAsyncHook{
		AppHook:    NewAppHook(repo),
		buf:        make(chan *logrus.Entry, AppBufSize),
		flush:      make(chan bool),
		Ticker:     time.NewTicker(time.Second),
		InsertFunc: asyncAppInsertFunc,
	}
	go hook.fire()
	return hook
}

func (h *AppHook) newEntry(entry *logrus.Entry) *logrus.Entry {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return &logrus.Entry{
		Logger:  entry.Logger,
		Data:    entry.Data,
		Time:    entry.Time,
		Level:   entry.Level,
		Message: entry.Message,
	}
}

//Fire run hook insert function
//Params
//	entry - logrus entry with fields to log
//Return
//	error - if error occurs return error, otherwise nil
func (h *AppHook) Fire(entry *logrus.Entry) error {
	newEntry := h.newEntry(entry)
	err := h.InsertFunc(newEntry, h.repo)
	if err != nil {
		return errors.Internal.Wrap(err, "failed to fire hook")
	}
	return nil
}

//Fire run async hook insert function
//Params
//	entry - logrus entry with fields to log
//Return
//	error - if error occurs return error, otherwise nil
func (h *AppAsyncHook) Fire(entry *logrus.Entry) error {
	h.wg.Add(1)
	h.buf <- h.newEntry(entry)
	return nil
}

// Flush empty description
func (h *AppAsyncHook) Flush() {
	h.Ticker = time.NewTicker(100 * time.Millisecond)
	h.wg.Wait()
	h.flush <- true
	<-h.flush
}

func (h *AppAsyncHook) fire() {
	for {
		var err error
		if err != nil {
			select {
			case <-h.Ticker.C:
				continue
			}
		}

		var numEntries int
		rnd := rand.Intn(10)
		time.Sleep(time.Duration(rnd) * time.Second)
	Loop:
		for {
			select {
			case entry := <-h.buf:
				err = h.InsertFunc(entry, h.repo)
				if err != nil {
					fmt.Fprintf(os.Stderr, "[LogrusHook] Can't insert entry (%v): %v\n", entry, err)
				}
				numEntries++
			case <-h.Ticker.C:
				if numEntries > 0 {
					break Loop
				}
			case <-h.flush:
				h.flush <- true
				return
			}
		}

		for i := 0; i < numEntries; i++ {
			h.wg.Done()
		}
	}
}

// Levels returns the available logging levels.
func (h *AppHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
