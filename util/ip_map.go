package util

import (
	"database/sql"
	"errors"
	"strconv"
	"strings"
)

func LoadIPMap(db *sql.DB) (map[string]map[string]map[string]map[string]map[string]interface{}, error) {
	if db == nil {
		return nil, errors.New("db connection is nil")
	}
	ipLib := map[string]map[string]map[string]map[string]map[string]interface{}{}

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

func insertInIPMap(ipLib map[string]map[string]map[string]map[string]map[string]interface{}, userID int64, ip string) {
	ipCut := strings.Split(ip, `.`)
	if len(ipCut) != 4 {
		return
	}

	// если записей по него еще нет то создай map
	if _, ok := ipLib[strconv.FormatInt(userID, 10)]; !ok {
		ipLib[strconv.FormatInt(userID, 10)] = map[string]map[string]map[string]map[string]interface{}{}
	}

	// вносим первую часть ip
	if _, ok := ipLib[strconv.FormatInt(userID, 10)][ipCut[0]]; !ok {
		ipLib[strconv.FormatInt(userID, 10)][ipCut[0]] = map[string]map[string]map[string]interface{}{}
	}

	// вторую
	if _, ok := ipLib[strconv.FormatInt(userID, 10)][ipCut[0]][ipCut[1]]; !ok {
		ipLib[strconv.FormatInt(userID, 10)][ipCut[0]][ipCut[1]] = map[string]map[string]interface{}{}
	}

	// третью
	if _, ok := ipLib[strconv.FormatInt(userID, 10)][ipCut[0]][ipCut[1]][ipCut[2]]; !ok {
		ipLib[strconv.FormatInt(userID, 10)][ipCut[0]][ipCut[1]][ipCut[2]] = map[string]interface{}{}
	}

	// четвертую
	if _, ok := ipLib[strconv.FormatInt(userID, 10)][ipCut[0]][ipCut[1]][ipCut[2]][ipCut[3]]; !ok {
		ipLib[strconv.FormatInt(userID, 10)][ipCut[0]][ipCut[1]][ipCut[2]][ipCut[3]] = struct{}{}
	}

	return
}

func CheckDupesInIPMap(ipLib map[string]map[string]map[string]map[string]map[string]interface{}, first, second string) bool {
	if first == second {
		return true
	}

	// есть ли вообще записи об этом id
	firstPartOne, ok := ipLib[first]
	if !ok {
		return false
	}

	secondPartOne, ok := ipLib[second]
	if !ok {
		return false
	}

	intersectionOne := checkFirstIntersection(firstPartOne, secondPartOne)

	if len(intersectionOne) <= 0 {
		return false
	}

	totalOverlap := 0

	// цикл по первой части ip
	for _, value1 := range intersectionOne {
		firstPartTwo, ok := ipLib[first][value1]
		if !ok {
			continue
		}
		secondPartTwo, ok := ipLib[second][value1]
		if !ok {
			continue
		}

		intersectionTwo := checkSecondIntersection(firstPartTwo, secondPartTwo)
		if len(intersectionTwo) <= 0 {
			continue
		}

		// цикл по второй части ip
		for _, value2 := range intersectionTwo {

			firstPartThree, ok := ipLib[first][value1][value2]
			if !ok {
				continue
			}
			secondPartThree, ok := ipLib[second][value1][value2]
			if !ok {
				continue
			}

			intersectionThree := checkThirdIntersection(firstPartThree, secondPartThree)
			if len(intersectionThree) <= 0 {
				continue
			}

			// цикл по третьей части ip
			for _, value3 := range intersectionThree {
				firstPartFour, ok := ipLib[first][value1][value2][value3]
				if !ok {
					continue
				}
				secondPartFour, ok := ipLib[second][value1][value2][value3]
				if !ok {
					continue
				}

				intersectionFour := checkFourIntersection(firstPartFour, secondPartFour)
				totalOverlap += len(intersectionFour)
				if totalOverlap >= 2 {
					return true
				}
			}

		}
	}

	return false
}

func checkFirstIntersection(first map[string]map[string]map[string]map[string]interface{}, second map[string]map[string]map[string]map[string]interface{}) []string {
	var intersection []string
	for v := range second {
		if _, ok := first[v]; ok {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

func checkSecondIntersection(first map[string]map[string]map[string]interface{}, second map[string]map[string]map[string]interface{}) []string {
	var intersection []string
	for v := range second {
		if _, ok := first[v]; ok {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

func checkThirdIntersection(first map[string]map[string]interface{}, second map[string]map[string]interface{}) []string {
	var intersection []string
	for v := range second {
		if _, ok := first[v]; ok {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

func checkFourIntersection(first map[string]interface{}, second map[string]interface{}) []string {
	var intersection []string
	for v := range second {
		if _, ok := first[v]; ok {
			intersection = append(intersection, v)
		}
	}
	return intersection
}
