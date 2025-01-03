package models

type UserData struct {
	FirstName string `json:"first_name" bson:"first_name"`
	LastName  string `json:"last_name" bson:"last_name"`
	Email     string `json:"email" bson:"email"`
	Password  string `json:"password" bson:"password"`
	// Geolocation data
	// Country string `json:"-"`
	// City    string `json:"-"`
	// Not used currently, because the database controller
	// doesn't have the corresponding field.
	// CreatedAt time.Time `json:"-"`
	// UpdatedAt time.Time `json:"-"`
}

// TODO: Create a single struct which contains user's data,
// together with it's geolocation so we can easily insert it
// into a database and retrieve the data as well,
// that would help for things like password validation etc.
// Try using sqlx ORM together with pgx driver.
