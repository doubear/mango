package contracts

// Authenable is authentication manager.
type Authenable interface {
	UserID() string
	SetUserID(string)
}
