package infrastructure

import (
	"fmt"
	"math/rand"
	"os"
	"rol/app/interfaces"
	"rol/app/utils"
	"rol/domain"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var HttpBufSize uint = 8192

const httpApiName = "httplog"
const appApiName = "applog"

//HttpHook log hook struct
type HttpHook struct {
	repo       interfaces.IGenericRepository[domain.HTTPLog]
	mutex      sync.RWMutex
	InsertFunc func(entry *logrus.Entry, repository interfaces.IGenericRepository[domain.HTTPLog]) error
}

//HttpAsyncHook log hook struct
type HttpAsyncHook struct {
	*HttpHook
	buf        chan *logrus.Entry
	flush      chan bool
	wg         sync.WaitGroup
	Ticker     *time.Ticker
	InsertFunc func(entry *logrus.Entry, repository interfaces.IGenericRepository[domain.HTTPLog]) error
}

var httpInsertFunc = func(entry *logrus.Entry, repository interfaces.IGenericRepository[domain.HTTPLog]) error {
	if entry.Data["method"] != nil && !fromLogController(entry) {
		ent := newEntityFromHttp(entry)
		_, err := repository.Insert(nil, *ent)
		return err
	}
	return nil
}

var asyncHttpInsertFunc = func(entry *logrus.Entry, repository interfaces.IGenericRepository[domain.HTTPLog]) error {
	if entry.Data["method"] != nil {
		ent := newEntityFromHttp(entry)
		_, err := repository.Insert(nil, *ent)
		return err
	}
	return nil
}

func newEntityFromHttp(entry *logrus.Entry) *domain.HTTPLog {
	customHeaders := entry.Data["customHeaders"].(string)
	queryParams := entry.Data["queryParams"].(string)

	customHeadersInd := utils.CutIndexingString(customHeaders)
	queryParamsInd := utils.CutIndexingString(queryParams)

	return &domain.HTTPLog{
		Entity: domain.Entity{
			ID: entry.Data["requestID"].(uuid.UUID),
		},
		HTTPMethod:              entry.Data["method"].(string),
		ClientIP:                entry.Data["clientIP"].(string),
		RequestBody:             entry.Data["requestBody"].(string),
		Domain:                  entry.Data["domain"].(string),
		RelativePath:            entry.Data["path"].(string),
		QueryParams:             entry.Data["queryParams"].(string),
		QueryParamsInd:          queryParamsInd,
		RequestHeaders:          entry.Data["headers"].(string),
		Latency:                 entry.Data["latency"].(int),
		ResponseBody:            entry.Data["responseBody"].(string),
		ResponseHeaders:         entry.Data["responseHeaders"].(string),
		CustomRequestHeaders:    entry.Data["customHeaders"].(string),
		CustomRequestHeadersInd: customHeadersInd,
	}
}

//NewAppHook create new gorm hook for http logs
//Params
//	repo - gorm generic repository with domain.HttpHook instantiated
//Return
//	*HttpHook - gorm hook
func NewHttpHook(repo interfaces.IGenericRepository[domain.HTTPLog]) *HttpHook {
	return &HttpHook{
		repo:       repo,
		InsertFunc: httpInsertFunc,
	}
}

//NewAppHook create new async gorm hook for http logs
//Params
//	repo - gorm generic repository with domain.HttpHook instantiated
//Return
//	*HttpHook - async gorm hook
func NewAsyncHttpHook(repo *GormGenericRepository[domain.HTTPLog]) *HttpAsyncHook {
	hook := &HttpAsyncHook{
		HttpHook:   NewHttpHook(repo),
		buf:        make(chan *logrus.Entry, HttpBufSize),
		flush:      make(chan bool),
		Ticker:     time.NewTicker(time.Second),
		InsertFunc: asyncHttpInsertFunc,
	}
	go hook.fire()
	return hook
}

func (hook *HttpHook) newEntry(entry *logrus.Entry) *logrus.Entry {
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
func (hook *HttpHook) Fire(entry *logrus.Entry) error {
	newEntry := hook.newEntry(entry)
	return hook.InsertFunc(newEntry, hook.repo)
}

//Fire run async hook insert function
//Params
//	entry - logrus entry with fields to log
//Return
//	error - if error occurs return error, otherwise nil
func (hook *HttpAsyncHook) Fire(entry *logrus.Entry) error {
	hook.wg.Add(1)
	hook.buf <- hook.newEntry(entry)
	return nil
}

func (hook *HttpAsyncHook) Flush() {
	hook.Ticker = time.NewTicker(100 * time.Millisecond)
	hook.wg.Wait()
	hook.flush <- true
	<-hook.flush
}

func (hook *HttpAsyncHook) fire() {
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

//Levels returns the available logging levels.
func (hook *HttpHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}

func fromLogController(entry *logrus.Entry) bool {
	if entry.Data["path"] == nil {
		return false
	} else if strings.Contains(entry.Data["path"].(string), httpApiName) || strings.Contains(entry.Data["path"].(string), appApiName) {
		return true
	}
	return false
}
