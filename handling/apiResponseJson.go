package handling

import (
	"github.com/heimdal-rw/chmgt/models"
)

// APIResponseJSON creates the structure of an API response
type APIResponseJSON struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    []models.Item `json:"data"`
}
