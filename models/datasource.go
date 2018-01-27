package models

import (
	"gopkg.in/mgo.v2"
)

type Datasource struct {
	Session      *mgo.Session
	DatabaseName string
	DSN          string
}

func NewDatasource(dsn, dbname string) (*Datasource, error) {
	datasource := new(Datasource)
	if err := datasource.Connect(); err != nil {
		return nil, err
	}

	datasource.DatabaseName = dbname
	datasource.DSN = dsn

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
