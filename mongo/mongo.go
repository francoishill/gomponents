package mongo

import (
	"io"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"gopkg.in/mgo.v2"
)

type Mongo interface {
	Close()
	Refresh()
	RefreshIfConnectionError(err error) error
	Collection(name string) *mgo.Collection
	EnsureIndexes(collection *mgo.Collection, indexes []mgo.Index) error
	IsErrNotFound(err error) bool
	IsDupErr(err error) bool
}

type defaultMongo struct {
	busyRefreshing bool
	refreshLock    sync.RWMutex

	db      *mgo.Database
	session *mgo.Session
}

//Close the session
func (m *defaultMongo) Close() {
	m.session.Close()
}

//Refresh the session
func (m *defaultMongo) Refresh() {
	m.refreshLock.Lock()
	defer m.refreshLock.Unlock()

	if m.busyRefreshing {
		return
	}

	m.busyRefreshing = true
	logrus.Error("Database is down or disconnected, will now refresh session")
	m.session.Refresh()
	time.Sleep(time.Second) //To avoid refreshing too much
	m.busyRefreshing = false
}

//inspiration from https://github.com/ti/mdb/blob/master/mdb.go
func (m *defaultMongo) RefreshIfConnectionError(err error) error {
	if !isNetworkError(err) {
		return err
	}
	m.Refresh() //throttling happens inside that method
	return err
}

func (m *defaultMongo) Collection(name string) *mgo.Collection {
	return m.db.C(name)
}

func (m *defaultMongo) EnsureIndexes(collection *mgo.Collection, indexes []mgo.Index) error {
	for _, index := range indexes {
		if err := collection.EnsureIndex(index); err != nil {
			return errors.Wrapf(err, "Failed to set mongo index (%v) on %s Collection", index, collection.Name)
		}
	}
	return nil
}

func (m *defaultMongo) IsErrNotFound(err error) bool { return err == mgo.ErrNotFound }
func (m *defaultMongo) IsDupErr(err error) bool      { return mgo.IsDup(err) }

func isNetworkError(err error) bool {
	if err == nil {
		return false
	}
	if err == io.EOF {
		return true
	}
	if _, ok := err.(*net.OpError); ok {
		return true
	}
	e := strings.ToLower(err.Error())
	if strings.HasPrefix(e, "closed") || strings.HasSuffix(e, "closed") {
		return true
	}
	return false
}

func DefaultMongo(connectionString string) *defaultMongo {
	urlInfo, err := url.Parse(connectionString)
	if err != nil {
		logrus.Panicf("Failed to parse mongo connection string as URL, error: %s", err.Error())
	}

	dbName := strings.TrimLeft(urlInfo.Path, "/ ")
	userName, pwd := "", ""
	if user := urlInfo.User; user != nil {
		userName = user.Username()
		pwd, _ = urlInfo.User.Password()
	}
	info := &mgo.DialInfo{
		Addrs:    []string{urlInfo.Host},
		Timeout:  10 * time.Second,
		Database: dbName,
		Source:   dbName,
		Username: userName,
		Password: pwd,
	}

	logrus.Debug("Connecting mongo")
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		logrus.Panicf("Failed to connect to mongo, error: %s", err.Error())
	}
	logrus.Debug("Connected Mongo")

	return &defaultMongo{
		db:      session.DB(dbName),
		session: session,
	}
}
