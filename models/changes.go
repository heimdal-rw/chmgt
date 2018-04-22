package models

import (
	"time"

	"github.com/divideandconquer/go-merge/merge"
	"gopkg.in/mgo.v2/bson"
)

// TBLCHANGES is the name of the collection for users
var TBLCHANGES = "Changes"

// Signatures is a slice of Signature objects
type Signatures []Signature

// Signature is an object containing user ID and time signed
type Signature struct {
	UserID bson.ObjectId `bson:"userid,omitempty" json:"userid"`
	Signed *time.Time    `bson:"signed" json:"signed"`
}

// Changes is a slice of Change objects
type Changes []Change

// Change is an object contianing change information
type Change struct {
	ID                 bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Created            Signature     `bson:"created" json:"created"`
	Requested          Signature     `bson:"requested" json:"requested"`
	Scheduled          *time.Time    `bson:"scheduled" json:"scheduled"`
	ScheduledDateOnly  bool          `bson:"scheduleddateonly" json:"scheduleddateonly"`
	Assigned           bson.ObjectId `bson:"assigned,omitempty" json:"assigned"`
	Executed           Signature     `bson:"executed" json:"executed"`
	Title              string        `bson:"title" json:"title"`
	Content            string        `bson:"content" json:"content"`
	Signatures         Signatures    `bson:"signatures" json:"signatures"`
	OverrideSignatures Signatures    `bson:"overridesignatures" json:"overridesignatures"`
}

// GetChanges is a handler to return the specified user objects
func (d *Datasource) GetChanges(id string) (Changes, error) {
	var changes Changes
	c := d.Database.C(TBLCHANGES)

	if id == "" {
		if err := c.Find(nil).All(&changes); err != nil {
			return nil, err
		}
	} else {
		if !bson.IsObjectIdHex(id) {
			return nil, ErrObjID
		}
		query := c.FindId(bson.ObjectIdHex(id))
		num, err := query.Count()
		if err != nil {
			return nil, err
		}
		if num <= 0 {
			return nil, ErrNotFound
		}
		if err := query.All(&changes); err != nil {
			return nil, err
		}
	}
	return changes, nil
}

// InsertChange inserts a new change object
func (d *Datasource) InsertChange(change Change) (bson.ObjectId, error) {
	change.ID = bson.NewObjectId()

	if err := d.Database.C(TBLCHANGES).Insert(&change); err != nil {
		return "", err
	}
	return change.ID, nil
}

// UpdateChange updates a specified change object
func (d *Datasource) UpdateChange(change Change) error {
	var currChange Change
	c := d.Database.C(TBLCHANGES)

	// get the current user data
	if err := c.FindId(change.ID).One(&currChange); err != nil {
		return err
	}

	// merge the new data with the current
	result := merge.Merge(currChange, change)

	return c.UpdateId(change.ID, result.(*Change))
}

// RemoveChange removes a specified change object
func (d *Datasource) RemoveChange(id string) error {
	if !bson.IsObjectIdHex(id) {
		return ErrObjID
	}
	return d.Database.C(TBLCHANGES).RemoveId(bson.ObjectIdHex(id))
}
