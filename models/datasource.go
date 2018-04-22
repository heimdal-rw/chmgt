package models

import (
	"github.com/heimdal-rw/chmgt/config"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Datasource is an object containing the database info and connection
type Datasource struct {
	session      *mgo.Session
	Database     *mgo.Database
	DatabaseName string
	DSN          string
}

// NewDatasource builds and connects to a database instance, then returns
// a Datasource object
func NewDatasource(config *config.Config) (*Datasource, error) {
	d := new(Datasource)
	d.DatabaseName = config.Database.Name
	d.DSN = config.DatabaseConnection()

	// make initial connection to Mongo
	var err error
	d.session, err = mgo.Dial(d.DSN)
	if err != nil {
		return nil, err
	}

	// authenticate to Mongo if configured
	if config.Database.AuthDB != "" {
		if err := d.session.DB(config.Database.AuthDB).Login(
			config.Database.Username,
			config.Database.Password,
		); err != nil {
			return nil, err
		}
	}

	// connect to specific database
	d.Database = d.session.DB(d.DatabaseName)

	// make sure we have unique usernames
	userIdx := mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	err = d.Database.C(TBLUSERS).EnsureIndex(userIdx)
	if err != nil {
		return nil, err
	}

	// make sure we have unique email addresses
	emailIdx := mgo.Index{
		Key:        []string{"email"},
		Unique:     true,
		DropDups:   false,
		Background: false,
		Sparse:     false,
	}
	err = d.Database.C(TBLUSERS).EnsureIndex(emailIdx)
	if err != nil {
		return nil, err
	}

	return d, nil
}

// Close terminates a connection to the database
func (d *Datasource) Close() {
	d.session.Close()
}

// ValidateUser checks the user credentials in the database
func (d *Datasource) ValidateUser(username, password string) (bool, error) {
	c := d.Database.C(TBLUSERS)
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
