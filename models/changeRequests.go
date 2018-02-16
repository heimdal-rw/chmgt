package models

// CollectionChangeRequests is the name of the collection for change requests
var CollectionChangeRequests = "ChangeRequests"

// GetChangeRequests is a wrapper to get change request items
func (d *Datasource) GetChangeRequests(id string) ([]Item, error) {
	return d.GetItems(id, CollectionChangeRequests)
}

// InsertChangeRequest is a wrapper to insert a change request
func (d *Datasource) InsertChangeRequest(cr Item) error {
	return d.InsertItem(cr, CollectionChangeRequests)
}

// RemoveChangeRequest is a wrapper to remove a change request
func (d *Datasource) RemoveChangeRequest(cr Item) error {
	return d.RemoveItem(cr, CollectionChangeRequests)
}

// UpdateChangeRequest is a wrapper to update a change request
func (d *Datasource) UpdateChangeRequest(cr Item) error {
	return d.UpdateItem(cr, CollectionChangeRequests)
}
