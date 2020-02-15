package sessions

import (
	"github.com/flannel-dev-lab/cyclops/cookie"
	"github.com/google/uuid"
	"net/http"
	"time"
)

type Session struct {
	Store  Store
	Cookie cookie.CyclopsCookie
}

func (session *Session) generateSessionId() string {
	uuidObj := uuid.New()
	return uuidObj.String()
}

// Set will set a session to the database with a map of values and expiry. Once written, it will send the session_id as cookie
func (session *Session) Set(w http.ResponseWriter, value map[string]interface{}, expiry time.Duration) (err error) {
	sessionId := session.generateSessionId()

	err = session.Store.Save(sessionId, value, expiry)
	if err != nil {
		return err
	}

	session.Cookie.Name = "session_id"
	session.Cookie.Value = sessionId
	session.Cookie.Expires = expiry

	session.Cookie.SetCookie(w)
	return nil
}

// Get Retrieves a particular session with key "session_id"
func (session *Session) Get(r *http.Request) (data map[string]interface{}, err error) {
	cookieData, err := session.Cookie.GetCookie(r, "session_id")
	if err != nil {
		return data, err
	}

	return session.Store.Get(cookieData.Value)
}

func (session *Session) getSessionDataWithKey(r *http.Request) (data map[string]interface{}, key string, err error) {
	cookieData, err := session.Cookie.GetCookie(r, "session_id")
	if err != nil {
		return data, key, err
	}

	sessionData, err := session.Store.Get(cookieData.Value)
	if err != nil {
		return data, key, err
	}

	return sessionData, cookieData.Value, nil
}

// Delete only deletes a particular session from database
func (session *Session) Delete(w http.ResponseWriter, r *http.Request) (err error) {
	cookieData, err := session.Cookie.GetCookie(r, "session_id")
	if err != nil {
		return err
	}

	session.Cookie.Delete(w, cookieData)

	return session.Store.Delete(cookieData.Value)
}

// Reset deletes all session values from the database
func (session *Session) Reset(r *http.Request) (err error) {
	return session.Store.Reset()
}

// Update updates the existing session with the data passes
func (session *Session) Update(r *http.Request, w http.ResponseWriter, data map[string]interface{}, expiry time.Duration) (err error) {
	existingSessionData, key, err := session.getSessionDataWithKey(r)
	if err != nil {
		if err == http.ErrNoCookie {
			// Sets a new cookie if not found
			err = session.Set(w, data, expiry)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	for key, value := range data {
		existingSessionData[key] = value
	}

	return session.Store.Save(key, existingSessionData, expiry)
}
