package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Message represents the structure of a message document in MongoDB
type Message struct {
	// ID is the MongoDB document's primary key.
	ID primitive.ObjectID `bson:"_id,omitempty"`

	// SenderID refers to the User ID of the sender.
	// In Mongoose, this was `mongoose.Schema.Types.ObjectId, ref: "User"`.
	// In Go, we use `primitive.ObjectID` for the ID itself.
	// `bson:"senderId"`: Maps this field to the "senderId" in MongoDB.
	SenderID primitive.ObjectID `bson:"senderId"`

	// ReceiverID refers to the User ID of the receiver.
	// `bson:"receiverId"`: Maps to "receiverId" in MongoDB.
	ReceiverID primitive.ObjectID `bson:"receiverId"`

	// Text content of the message. Optional in Mongoose.
	// `bson:"text,omitempty"`: Maps to "text". `omitempty` is used as it can be empty.
	Text string `bson:"text,omitempty"`

	// Image URL associated with the message. Optional in Mongoose.
	// `bson:"image,omitempty"`: Maps to "image". `omitempty` is used as it can be empty.
	Image string `bson:"image,omitempty"`

	// CreatedAt field, automatically added by Mongoose `timestamps: true`.
	CreatedAt time.Time `bson:"createdAt"`

	// UpdatedAt field, automatically added by Mongoose `timestamps: true`.
	UpdatedAt time.Time `bson:"updatedAt"`
}
