package global

import "gopkg.in/olahol/melody.v1"

var (
	Clients  = make(map[string]interface{})
	Sessions []*melody.Session
)
