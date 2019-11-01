package handlers

import (
	"database/sql"
	"net/http"

	"github.com/ant0ine/go-json-rest/rest"
)

type DBConnect struct {
	DB *sql.DB
}

func (connect DBConnect) IsOk(w rest.ResponseWriter, r *rest.Request) {
	if connect.DB == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := connect.DB.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.(http.ResponseWriter).Write([]byte("ok"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
