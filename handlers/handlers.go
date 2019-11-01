package handlers

import (
	"database/sql"
	"net/http"
)

type DBConnect struct {
	DB *sql.DB
}

func (connect DBConnect) IsOk(w http.ResponseWriter, r *http.Request) {
	if connect.DB == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err := connect.DB.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("ok"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
