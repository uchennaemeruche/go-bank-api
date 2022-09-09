package token

import "time"

/** Todo:
* Create Token - Takes username and duration and returns a toke string or an error.
* Verify Token - Takes a token string and returns the pointer to the Payload or an error
**/

// Defines functions that different Token providers will implement.
type Maker interface {
	CreateToken(username string, duration time.Duration) (token string, payload *Payload, err error)
	VerifyToken(token string) (*Payload, error)
}
