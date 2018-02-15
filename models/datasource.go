package models

import (
	"encoding/json"
	"errors"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var ErrNoRows = errors.New("datasource: no records returned")

type Item map[string]interface{}

func (i Item) SetID(id string) {
	i["_id"] = bson.ObjectIdHex(id)
}

func (i *Item) MarshalJSON() ([]byte, error) {
	var j interface{}
	data, err := bson.Marshal(i)
	if err != nil {
		return nil, err
	}
	bson.Unmarshal(data, &j)
	return json.Marshal(&j)
}

func (i *Item) UnmarshalJSON(p []byte) error {
	var j map[string]interface{}
	json.Unmarshal(p, &j)
	data, err := bson.Marshal(&j)
	if err != nil {
		return err
	}
	return bson.Unmarshal(data, i)
}

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

func (d *Datasource) GetItems(id, collection string) ([]Item, error) {
	c := d.Session.DB(d.DatabaseName).C(collection)
	var (
		items []Item
		err   error
	)
	if id != "" {
		query := c.FindId(bson.ObjectIdHex(id))
		num, err := query.Count()
		if err != nil {
			return nil, err
		}
		if num <= 0 {
			return nil, ErrNoRows
		}
		err = query.All(&items)
	} else {
		err = c.Find(nil).All(&items)
	}
	if err != nil {
		return nil, err
	}

	return items, err
}

func (d *Datasource) InsertItem(item Item, collection string) error {
	c := d.Session.DB(d.DatabaseName).C(collection)
	item["_id"] = bson.NewObjectId()
	err := c.Insert(item)
	if err != nil {
		return err
	}
	return nil
}

func (d *Datasource) RemoveItem(user Item, collection string) error {
	c := d.Session.DB(d.DatabaseName).C(collection)
	return c.RemoveId(user["_id"])
}

func (d *Datasource) UpdateItem(user Item, collection string) error {
	c := d.Session.DB(d.DatabaseName).C(collection)
	return c.UpdateId(user["_id"], user)
}
