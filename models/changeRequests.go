package models

import (
	"errors"

	"gopkg.in/mgo.v2/bson"
)

type ChangeRequest struct {
	ID          bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
}

func (cr *ChangeRequest) SetID(id string) {
	cr.ID = bson.ObjectIdHex(id)
}

func (d *Datasource) GetChangeRequests(id string) ([]ChangeRequest, error) {
	c := d.Session.DB(d.DatabaseName).C("ChangeRequests")
	var (
		crs []ChangeRequest
		err error
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
		err = query.All(&crs)
	} else {
		err = c.Find(nil).All(&crs)
	}
	if err != nil {
		return nil, err
	}

	return crs, err
}

func (d *Datasource) InsertChangeRequest(cr *ChangeRequest) error {
	c := d.Session.DB(d.DatabaseName).C("ChangeRequests")
	info, err := c.Upsert(new(ChangeRequest), cr)
	if err != nil {
		return err
	}
	if info.UpsertedId != nil {
		cr.ID = info.UpsertedId.(bson.ObjectId)
	} else {
		return errors.New("datasource: unknown error inserting user")
	}
	return nil
}

func (d *Datasource) RemoveChangeRequest(cr *ChangeRequest) error {
	c := d.Session.DB(d.DatabaseName).C("ChangeRequests")
	return c.RemoveId(cr.ID)
}

func (d *Datasource) UpdateChangeRequest(cr *ChangeRequest) error {
	c := d.Session.DB(d.DatabaseName).C("ChangeRequests")
	return c.UpdateId(cr.ID, cr)
}
