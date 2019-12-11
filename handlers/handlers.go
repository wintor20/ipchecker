package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"

	"ipchecker/models"
	"ipchecker/util"
)

type DBConnect struct {
	DB        *sql.DB
	IpLib     map[string]map[string]map[string]map[string]map[string]interface{}
	Dupes     map[string]interface{}
	Potential map[string]interface{}
}

// IsOk checks db connection
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

// CheckIP process GET request to check if two user_ids have at least two matching ip address
func (connect DBConnect) CheckIP(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if connect.DB == nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ids := strings.Split(strings.TrimPrefix(r.URL.Path, `/`), `/`)
	if len(ids) != 2 {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var answer models.CheckIPAnswer
	answer.Dupes = util.CheckDupesInIPMap(connect.Dupes, ids[0], ids[1])
	prepareAnswer(answer, &w)
	return
}

func prepareAnswer(answer models.CheckIPAnswer, w *http.ResponseWriter) {
	b, err := json.Marshal(answer)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = (*w).Write(b)
	if err != nil {
		(*w).WriteHeader(http.StatusInternalServerError)
		return
	}
}
