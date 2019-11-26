package util

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

func LoadIPMap(db *sql.DB) (map[string]map[string]map[string]map[string]map[string]string, error) {
	if db == nil {
		return nil, errors.New("db connection is nil")
	}
	ipLib := map[string]map[string]map[string]map[string]map[string]string{}

	// postgres
	rows, err := db.Query("SELECT user_id, ip_addr FROM conn_log")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userID int64
	var ip string
	for rows.Next() {
		err = rows.Scan(&userID, &ip)
		if err != nil {
			return nil, err
		}
		insertInIPMap(ipLib, userID, ip)
	}

	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return ipLib, nil
}

func insertInIPMap(ipLib map[string]map[string]map[string]map[string]map[string]string, userID int64, ip string) {
	ipCut := strings.Split(ip, `.`)
	if len(ipCut) != 4 {
		return
	}

	// если записей по него еще нет то создай map
	if _, ok := ipLib[strconv.FormatInt(userID, 10)]; !ok {
		ipLib[strconv.FormatInt(userID, 10)] = map[string]map[string]map[string]map[string]string{}
	}

	// вносим первую часть ip
	if _, ok := ipLib[strconv.FormatInt(userID, 10)][ipCut[0]]; !ok {
		ipLib[strconv.FormatInt(userID, 10)][ipCut[0]] = map[string]map[string]map[string]string{}
	}

	// вторую
	if _, ok := ipLib[strconv.FormatInt(userID, 10)][ipCut[0]][ipCut[1]]; !ok {
		ipLib[strconv.FormatInt(userID, 10)][ipCut[0]][ipCut[1]] = map[string]map[string]string{}
	}

	// третью
	if _, ok := ipLib[strconv.FormatInt(userID, 10)][ipCut[0]][ipCut[1]][ipCut[2]]; !ok {
		ipLib[strconv.FormatInt(userID, 10)][ipCut[0]][ipCut[1]][ipCut[2]] = map[string]string{}
	}

	// четвертую
	if _, ok := ipLib[strconv.FormatInt(userID, 10)][ipCut[0]][ipCut[1]][ipCut[2]][ipCut[3]]; !ok {
		ipLib[strconv.FormatInt(userID, 10)][ipCut[0]][ipCut[1]][ipCut[2]][ipCut[3]] = ""
	}

	return
}

func CheckDupesInIPMap(ipLib map[string]map[string]map[string]map[string]map[string]string, first, second string) bool {
	return false
}
