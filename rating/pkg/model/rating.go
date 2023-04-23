package model

// RecordID defines a record id. Together with RecordType
// identifies unique records across all types.
type RecordID string

// RecordType defines a record type. Together with RecordID
// identifies unique records across all types
type RecordType string

// RecordTypeMovie Existing record types.
const (
	RecordTypeMovie       = RecordType("movie")
	RatingEventTypePut    = "put"
	RatingEventTypeDelete = "delete"
)

// UserID defines a user id.
type UserID string

// RatingValue defines a value of a rating record.
type RatingValue int

// Rating defines an individual rating created by a user for some records.
type Rating struct {
	RecordID   RecordID    `json:"recordID"`
	RecordType RecordType  `json:"recordType"`
	UserID     UserID      `json:"userID"`
	Value      RatingValue `json:"value"`
}

// RatingEventType defines the type of rating event.
type RatingEventType string

// RatingEvent defines an event containing rating information.
type RatingEvent struct {
	UserID     UserID          `json:"userId"`
	RecordID   RecordID        `json:"recordId"`
	RecordType RecordType      `json:"recordType"`
	Value      RatingValue     `json:"value"`
	ProviderID string          `json:"providerId"`
	EventType  RatingEventType `json:"eventType"`
}
