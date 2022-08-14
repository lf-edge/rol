package infrastructure

import (
	"fmt"
	"math/rand"
	"os"
	"rol/app/errors"
	"rol/app/interfaces"
	"rol/app/utils"
	"rol/domain"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

//HTTPBufSize buffer size for async hooks
var HTTPBufSize uint = 8192

const httpAPIName = "httplog"
const appAPIName = "applog"

//HTTPHook log hook struct
type HTTPHook struct {
	repo       interfaces.IGenericRepository[domain.HTTPLog]
	mutex      sync.RWMutex
	InsertFunc func(entry *logrus.Entry, repository interfaces.IGenericRepository[domain.HTTPLog]) error
}

//HTTPAsyncHook log hook struct
type HTTPAsyncHook struct {
	*HTTPHook
	buf        chan *logrus.Entry
	flush      chan bool
	wg         sync.WaitGroup
	Ticker     *time.Ticker
	InsertFunc func(entry *logrus.Entry, repository interfaces.IGenericRepository[domain.HTTPLog]) error
}

var httpInsertFunc = func(entry *logrus.Entry, repository interfaces.IGenericRepository[domain.HTTPLog]) error {
	if entry.Data["method"] != nil && !fromLogController(entry) {
		ent := newEntityFromHTTP(entry)
		_, err := repository.Insert(nil, *ent)
		if err != nil {
			return errors.Internal.Wrap(err, "error inserting http log to db")
		}
	}
	return nil
}

var asyncHTTPInsertFunc = func(entry *logrus.Entry, repository interfaces.IGenericRepository[domain.HTTPLog]) error {
	if entry.Data["method"] != nil {
		ent := newEntityFromHTTP(entry)
		_, err := repository.Insert(nil, *ent)
		return errors.Internal.Wrap(err, "error inserting http log to db")
	}
	return nil
}

func newEntityFromHTTP(entry *logrus.Entry) *domain.HTTPLog {
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

//NewHTTPHook create new gorm hook for http logs
//Params
//	repo - gorm generic repository with domain.HTTPHook instantiated
//Return
//	*HTTPHook - gorm hook
func NewHTTPHook(repo interfaces.IGenericRepository[domain.HTTPLog]) *HTTPHook {
	return &HTTPHook{
		repo:       repo,
		InsertFunc: httpInsertFunc,
	}
}

//NewAsyncHTTPHook create new async gorm hook for http logs
//Params
//	repo - gorm generic repository with domain.HTTPHook instantiated
//Return
//	*HTTPHook - async gorm hook
func NewAsyncHTTPHook(repo *GormGenericRepository[domain.HTTPLog]) *HTTPAsyncHook {
	hook := &HTTPAsyncHook{
		HTTPHook:   NewHTTPHook(repo),
		buf:        make(chan *logrus.Entry, HTTPBufSize),
		flush:      make(chan bool),
		Ticker:     time.NewTicker(time.Second),
		InsertFunc: asyncHTTPInsertFunc,
	}
	go hook.fire()
	return hook
}

func (h *HTTPHook) newEntry(entry *logrus.Entry) *logrus.Entry {
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
func (h *HTTPHook) Fire(entry *logrus.Entry) error {
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
func (h *HTTPAsyncHook) Fire(entry *logrus.Entry) error {
	h.wg.Add(1)
	h.buf <- h.newEntry(entry)
	return nil
}

//Flush flush async hook
func (h *HTTPAsyncHook) Flush() {
	h.Ticker = time.NewTicker(100 * time.Millisecond)
	h.wg.Wait()
	h.flush <- true
	<-h.flush
}

func (h *HTTPAsyncHook) fire() {
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

//Levels returns the available logging levels.
func (h *HTTPHook) Levels() []logrus.Level {
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
	} else if strings.Contains(entry.Data["path"].(string), httpAPIName) || strings.Contains(entry.Data["path"].(string), appAPIName) {
		return true
	}
	return false
}
