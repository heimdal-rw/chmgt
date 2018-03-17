package models

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/heimdal-rw/chmgt/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// ErrNoRows is the error to return when no records were found
var ErrNoRows = errors.New("datasource: no records returned")

// ErrObjID is the error to return when an Object ID is invalid
var ErrObjID = errors.New("invalid object id")

// CollectionChangeRequests is the name of the collection for change requests
var CollectionChangeRequests = "ChangeRequests"

// CollectionUsers is the name of the collection for users
var CollectionUsers = "Users"

// Item encompases objects saved in the database
type Item map[string]interface{}

// SetID turns a string into a valid MongoDB ID and sets it on the object
func (i Item) SetID(id string) error {
	if !bson.IsObjectIdHex(id) {
		return ErrObjID
	}
	i["_id"] = bson.ObjectIdHex(id)
	return nil
}

// MarshalJSON returns a json formatted string from the Item object
func (i *Item) MarshalJSON() ([]byte, error) {
	var j interface{}
	data, err := bson.Marshal(i)
	if err != nil {
		return nil, err
	}
	bson.Unmarshal(data, &j)
	return json.Marshal(&j)
}

// UnmarshalJSON returns an Item object from a json formatted string
func (i *Item) UnmarshalJSON(p []byte) error {
	var j map[string]interface{}
	json.Unmarshal(p, &j)
	data, err := bson.Marshal(&j)
	if err != nil {
		return err
	}
	return bson.Unmarshal(data, i)
}

// Datasource is an object containing the database info and connection
type Datasource struct {
	Session      *mgo.Session
	DatabaseName string
	DSN          string
}

// NewDatasource builds and connects to a database instance, then returns
// a Datasource object
func NewDatasource(config *config.Config) (*Datasource, error) {
	datasource := new(Datasource)
	datasource.DatabaseName = config.Database.Name
	datasource.DSN = config.DatabaseConnection()

	if err := datasource.Connect(); err != nil {
		return nil, err
	}

	if config.Database.AuthDB != "" {
		if err := datasource.Session.DB(config.Database.AuthDB).Login(
			config.Database.Username,
			config.Database.Password,
		); err != nil {
			return nil, err
		}
	}

	userIdx := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	err := datasource.Session.DB(datasource.DatabaseName).C(CollectionUsers).EnsureIndex(userIdx)
	if err != nil {
		return nil, err
	}

	return datasource, nil
}

// Connect creates a connection to the database
func (d *Datasource) Connect() error {
	var err error
	d.Session, err = mgo.Dial(d.DSN)
	if err != nil {
		return err
	}

	return nil
}

// Close terminates a connection to the database
func (d *Datasource) Close() {
	d.Session.Close()
}

func (d *Datasource) GetPassword(id bson.ObjectId) (string, error) {
	c := d.Session.DB(d.DatabaseName).C("Users")

	var items []Item

	err := c.FindId(id).All(&items)
	if err != nil {
		return "", err
	}

	if items == nil {
		return "", ErrNoRows
	}

	if val, ok := items[0]["password"].(string); ok {
		return val, nil
	}
	return "", nil
}

// GetItems queries the database for specified items
func (d *Datasource) GetItems(id, collection string) ([]Item, error) {
	c := d.Session.DB(d.DatabaseName).C(collection)
	var (
		items []Item
		err   error
	)
	if id != "" {
		if !bson.IsObjectIdHex(id) {
			return nil, ErrObjID
		}
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

	// Remove the password from the items if it exists
	for _, item := range items {
		delete(item, "password")
	}

	if items == nil {
		items = make([]Item, 0)
	}

	return items, err
}

// InsertItem inserts an object into the database
func (d *Datasource) InsertItem(item Item, collection string) error {
	c := d.Session.DB(d.DatabaseName).C(collection)
	item["_id"] = bson.NewObjectId()
	err := c.Insert(item)
	if err != nil {
		return err
	}
	return nil
}

// RemoveItem removes the specified object from the database
func (d *Datasource) RemoveItem(item Item, collection string) error {
	c := d.Session.DB(d.DatabaseName).C(collection)
	return c.RemoveId(item["_id"])
}

// UpdateItem updates the specified object in the database
func (d *Datasource) UpdateItem(item Item, collection string) error {
	if strings.EqualFold("Users", collection) {
		if _, ok := item["password"]; !ok {
			password, err := d.GetPassword(item["_id"].(bson.ObjectId))
			if err != nil {
				return err
			}
			item["password"] = password
		}
	}
	c := d.Session.DB(d.DatabaseName).C(collection)
	return c.UpdateId(item["_id"], item)
}

// ValidateUser checks the user credentials in the database
func (d *Datasource) ValidateUser(username, password string) (bool, error) {
	c := d.Session.DB(d.DatabaseName).C(CollectionUsers)
	query := bson.M{
		"username": username,
		"password": password,
	}
	num, err := c.Find(query).Count()
	if err != nil {
		return false, err
	}
	// Since usernames are unique, there should be only one record
	if num != 1 {
		return false, nil
	}
	return true, nil
}
