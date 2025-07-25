package models

import (
	"time" // Required for `time.Time` fields for timestamps

	"go.mongodb.org/mongo-driver/bson/primitive" // Required for `primitive.ObjectID` to handle MongoDB's unique identifiers
)

// User represents the structure of a user document in MongoDB.
// Each field is defined with its Go type and a `bson` tag.
// The `bson` tag tells the MongoDB Go driver how to map the Go struct field
// to the corresponding field name in the MongoDB document.
type User struct {
	// ID is the MongoDB document's primary key.
	// `primitive.ObjectID` is the Go type for MongoDB's `_id`.
	// `bson:"_id,omitempty"`:
	//   - `_id`: Maps this field to MongoDB's `_id` field.
	//   - `omitempty`: This option means the field will be omitted from the BSON document
	//     if its value is the zero value for its type (e.g., an empty ObjectID).
	//     This is useful when creating new documents where MongoDB generates the _id.
	ID primitive.ObjectID `bson:"_id,omitempty"`

	// Email field, required and unique in your Mongoose schema.
	// `bson:"email"`: Maps this field to the "email" field in MongoDB.
	Email string `bson:"email"`

	// FullName field, required in your Mongoose schema.
	// `bson:"fullName"`: Maps to "fullName" in MongoDB.
	FullName string `bson:"fullName"`

	// Password field, required and minlength 6 in your Mongoose schema.
	// This field will store the hashed password.
	// `bson:"password"`: Maps to "password" in MongoDB.
	Password string `bson:"password"`

	// ProfilePic field, optional with a default empty string in Mongoose.
	// `bson:"profilePic,omitempty"`: Maps to "profilePic". `omitempty` is used
	//   because it's an optional field and might be an empty string.
	ProfilePic string `bson:"profilePic,omitempty"`

	// CreatedAt field, automatically added by Mongoose `timestamps: true`.
	// `time.Time` is the Go type for timestamps.
	// `bson:"createdAt"`: Maps to "createdAt" in MongoDB.
	CreatedAt time.Time `bson:"createdAt"`

	// UpdatedAt field, automatically added by Mongoose `timestamps: true`.
	// `bson:"updatedAt"`: Maps to "updatedAt" in MongoDB.
	UpdatedAt time.Time `bson:"updatedAt"`
}