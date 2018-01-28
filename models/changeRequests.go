package models

import (
	"encoding/json"
	"errors"

	"gopkg.in/mgo.v2/bson"
)

type ChangeRequest struct {
	ID          bson.ObjectId          `bson:"_id,omitempty" json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Extras      map[string]interface{} `bson:",inline" json:"-"`
}

func (cr *ChangeRequest) SetID(id string) {
	cr.ID = bson.ObjectIdHex(id)
}

func (cr *ChangeRequest) MarshalJSON() ([]byte, error) {
	var j interface{}
	data, err := bson.Marshal(cr)
	if err != nil {
		return nil, err
	}
	bson.Unmarshal(data, &j)
	return json.Marshal(&j)
}

func (cr *ChangeRequest) UnmarshalJSON(p []byte) error {
	var j map[string]interface{}
	json.Unmarshal(p, &j)
	data, err := bson.Marshal(&j)
	if err != nil {
		return err
	}
	return bson.Unmarshal(data, cr)
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
