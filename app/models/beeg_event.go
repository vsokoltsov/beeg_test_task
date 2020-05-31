package models

// BeegEvent represents event structure
type BeegEvent struct {
	ID    int    `db:"id" json:"id"`
	Label string `db:"label" json:"label"`
	Count int    `db:"count" json:"count"`
}
