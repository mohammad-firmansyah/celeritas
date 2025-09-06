package celeritas

import "net/http"

func (c *Celeritas) SessionLoad(next http.Handler) http.Handler {
	c.InfoLog.Println("Session Load Called")
	return c.Session.LoadAndSave(next)
}
