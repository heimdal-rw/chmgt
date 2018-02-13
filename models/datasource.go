package models

import (
	"errors"

	"gopkg.in/mgo.v2"
)

var ErrNoRows = errors.New("datasource: no records returned")

type Datasource struct {
	Session      *mgo.Session
	DatabaseName string
	DSN          string
}

func NewDatasource(dsn, dbname string) (*Datasource, error) {
	datasource := new(Datasource)
	datasource.DatabaseName = dbname
	datasource.DSN = dsn

	if err := datasource.Connect(); err != nil {
		return nil, err
	}

	userIdx := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	err := datasource.Session.DB(dbname).C(CollectionUsers).EnsureIndex(userIdx)
	if err != nil {
		return nil, err
	}

	return datasource, nil
}

func (d *Datasource) Connect() error {
	var err error
	d.Session, err = mgo.Dial(d.DSN)
	if err != nil {
		return err
	}

	return nil
}

func (d *Datasource) Close() {
	d.Session.Close()
}
