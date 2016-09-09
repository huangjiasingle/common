package db

import (
	"io"

	"common/errors"
	"common/etc"
	l "common/log"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var log = l.New("common/db")

var (
	driver = etc.String("datastore", "driver")
	dsn    = etc.String("datastore", "dsn")
)

// --------------
// Engine

type Engine struct {
	*xorm.Engine
}

func New() (*Engine, error) {
	engine, err := xorm.NewEngine(driver, dsn)
	if err != nil {
		return nil, errors.As(err)
	}

	// cache
	// cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	// engine.SetDefaultCacher(cacher)

	return &Engine{engine}, nil
}

func (engine *Engine) Debug() {
	engine.ShowSQL(true)
}

func (engine *Engine) Close() error {
	return engine.Close()
}

// -------------------
// Common

type Closer interface {
	io.Closer
}

func Close(db Closer) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Warn(errors.As(err))
		}
	}
}
