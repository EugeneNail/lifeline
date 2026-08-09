package resource

import "time"

// Journal represents the public journal fields returned to the client.
type Journal struct {
	Date      string    `json:"date"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
