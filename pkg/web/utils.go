package web

import (
	"log"
	"net/http"
	"strconv"
)

func ParseBool(r *http.Request, s string) bool {
	b, err := strconv.ParseBool(r.URL.Query().Get(s))
	if err != nil {
		log.Println(err)
		return false
	}
	return b
}
