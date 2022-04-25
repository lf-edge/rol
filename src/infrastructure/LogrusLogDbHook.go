package infrastructure

import (
	"fmt"
	"math/rand"
	"os"
	"rol/app/interfaces"
	"rol/domain"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var AppBufSize uint = 8192

//AppHook application log hook struct
type AppHook struct {
	repo       interfaces.IGenericRepository[domain.AppLog]
	mutex      sync.RWMutex
	InsertFunc func(entry *logrus.Entry, repository interfaces.IGenericRepository[domain.AppLog]) error
}

//AppAsyncHook application log async hook struct
type AppAsyncHook struct {
	*AppHook
	buf        chan *logrus.Entry
	flush      chan bool
	wg         sync.WaitGroup
	Ticker     *time.Ticker
	InsertFunc func(entry *logrus.Entry, repository interfaces.IGenericRepository[domain.AppLog]) error
}

var insertAppFunc = func(entry *logrus.Entry, repository interfaces.IGenericRepository[domain.AppLog]) error {
	if entry.Data["method"] == nil && !fromLogController(entry) {
		ent := newEntityFromApp(entry)
		_, err := repository.Insert(nil, *ent)
		return err
	}
	return nil
}

var asyncAppInsertFunc = func(entry *logrus.Entry, repository interfaces.IGenericRepository[domain.AppLog]) error {
	if entry.Data["method"] != nil {
		ent := newEntityFromApp(entry)
		_, err := repository.Insert(nil, *ent)
		return err
	}
	return nil
}

func newEntityFromApp(entry *logrus.Entry) *domain.AppLog {
	var actionId uuid.UUID
	var source string
	if entry.Data["actionId"] != nil {
		actionId = entry.Data["actionId"].(uuid.UUID)
	}
	if entry.Data["source"] != nil {
		source = entry.Data["source"].(string)
	}
	return &domain.AppLog{
		Entity: domain.Entity{
			ID: uuid.New(),
		},
		Level:    entry.Level.String(),
		Source:   source,
		ActionID: actionId,
		Message:  entry.Message,
	}
}

//NewAppHook create new gorm hook for application logs
//Params
//	repo - gorm generic repository with domain.AppLog instantiated
//Return
//	*AppHook - gorm hook
func NewAppHook(repo interfaces.IGenericRepository[domain.AppLog]) *AppHook {
	return &AppHook{
		repo:       repo,
		InsertFunc: insertAppFunc,
	}
}

//NewAppHook create new async gorm hook for application logs
//Params
//	repo - gorm generic repository with domain.AppLog instantiated
//Return
//	*AppHook - async gorm hook
func NewAsyncAppHook(repo interfaces.IGenericRepository[domain.AppLog]) *AppAsyncHook {
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

func (hook *AppHook) newEntry(entry *logrus.Entry) *logrus.Entry {
	hook.mutex.RLock()
	defer hook.mutex.RUnlock()

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
func (hook *AppHook) Fire(entry *logrus.Entry) error {
	newEntry := hook.newEntry(entry)
	return hook.InsertFunc(newEntry, hook.repo)
}

//Fire run async hook insert function
//Params
//	entry - logrus entry with fields to log
//Return
//	error - if error occurs return error, otherwise nil
func (hook *AppAsyncHook) Fire(entry *logrus.Entry) error {
	hook.wg.Add(1)
	hook.buf <- hook.newEntry(entry)
	return nil
}

// Flush empty description
func (hook *AppAsyncHook) Flush() {
	hook.Ticker = time.NewTicker(100 * time.Millisecond)
	hook.wg.Wait()
	hook.flush <- true
	<-hook.flush
}

func (hook *AppAsyncHook) fire() {
	for {
		var err error
		if err != nil {
			select {
			case <-hook.Ticker.C:
				continue
			}
		}

		var numEntries int
		rnd := rand.Intn(10)
		time.Sleep(time.Duration(rnd) * time.Second)
	Loop:
		for {
			select {
			case entry := <-hook.buf:
				err = hook.InsertFunc(entry, hook.repo)
				if err != nil {
					fmt.Fprintf(os.Stderr, "[LogrusHook] Can't insert entry (%v): %v\n", entry, err)
				}
				numEntries++
			case <-hook.Ticker.C:
				if numEntries > 0 {
					break Loop
				}
			case <-hook.flush:
				hook.flush <- true
				return
			}
		}

		for i := 0; i < numEntries; i++ {
			hook.wg.Done()
		}
	}
}

// Levels returns the available logging levels.
func (hook *AppHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
