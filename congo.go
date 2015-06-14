package congo

// User represents a person in the system.
// This can be an attendee, speaker, sponsor contact, etc.
type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string
}

// Event represents a specific conference event within a Series.
type Event struct {
	ID   int
	Name string
}

// Ticket represents the admission to an Event for a User.
type Ticket struct {
	UserID  int
	EventID int
}
