package sessions

import "time"

// Store provides interface to create custom session stores
type Store interface {
	// Save helps in saving data to session store with a expiry
	Save(key string, data map[string]interface{}, expiry time.Duration) error
	// Get retrieves the data associated with the session key
	Get(key string) (data map[string]interface{}, err error)
	// Delete deletes a session from session store
	Delete(key string) error
	// Reset deletes all sessions from session store
	Reset() error
}