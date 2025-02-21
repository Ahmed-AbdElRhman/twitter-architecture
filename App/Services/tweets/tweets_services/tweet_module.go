package tweets_services

import "time"

type Tweet struct {
	ID        string    `json:"id" bson:"_id,omitempty"`
	Owner     int       `json:"owner" bson:"owner"`
	Tweet     string    `json:"tweet" bson:"tweet"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}
