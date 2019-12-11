package util

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"ipchecker/models"
)

const psLimit = 200000

func LoadIPMap(instance *models.ServiceInstance) (map[string]map[string]map[string]map[string]map[int64]interface{}, map[string]interface{}, map[string]interface{}, error) {
	if instance == nil || instance.DB == nil {
		return nil, nil, nil, errors.New("db connection is nil")
	}
	ipLib := map[string]map[string]map[string]map[string]map[int64]interface{}{}
	potential := map[string]interface{}{}
	dupes := map[string]interface{}{}

	hasMore := true
	var err error
	offset := 0
	// postgres
	for hasMore {
		hasMore, err = loadPart(instance, ipLib, potential, dupes, offset)
		if err != nil {
			return nil, nil, nil, err
		}
		offset += psLimit
	}

	return ipLib, potential, dupes, nil
}

func loadPart(instance *models.ServiceInstance, ipLib map[string]map[string]map[string]map[string]map[int64]interface{}, potential, dupes map[string]interface{}, offset int) (bool, error) {

	start := time.Now()

	rows, err := instance.DB.Query("SELECT user_id, ip_addr FROM conn_log ORDER BY ts ASC OFFSET $1 LIMIT $2", offset, psLimit)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	var userID int64
	var ip string

	index := 0
	for rows.Next() {
		err = rows.Scan(&userID, &ip)
		if err != nil {
			return false, err
		}
		insertInIPMap(ipLib, potential, dupes, userID, ip)
		index += 1
	}

	instance.Log.Printf("%d loaded %f", offset+index, time.Since(start).Seconds())

	err = rows.Err()
	if err != nil {
		return false, err
	}

	if index < psLimit {
		return false, nil
	}

	return true, nil
}

func insertInIPMap(ipLib map[string]map[string]map[string]map[string]map[int64]interface{}, potential, dupes map[string]interface{}, userID int64, ip string) {
	ipCut := strings.Split(ip, `.`)
	if len(ipCut) != 4 {
		return
	}

	ipExists := true

	// вносим первую часть ip
	if _, ok := ipLib[ipCut[0]]; !ok {
		ipLib[ipCut[0]] = map[string]map[string]map[string]map[int64]interface{}{}
		ipExists = false
	}

	// вторую
	if _, ok := ipLib[ipCut[0]][ipCut[1]]; !ok {
		ipLib[ipCut[0]][ipCut[1]] = map[string]map[string]map[int64]interface{}{}
		ipExists = false
	}

	// третью
	if _, ok := ipLib[ipCut[0]][ipCut[1]][ipCut[2]]; !ok {
		ipLib[ipCut[0]][ipCut[1]][ipCut[2]] = map[string]map[int64]interface{}{}
		ipExists = false
	}

	// четвертую
	if _, ok := ipLib[ipCut[0]][ipCut[1]][ipCut[2]][ipCut[3]]; !ok {
		ipLib[ipCut[0]][ipCut[1]][ipCut[2]][ipCut[3]] = map[int64]interface{}{}
		ipExists = false
	}

	// если такого ip еще не было то просто добавим в него userID
	if !ipExists {
		ipLib[ipCut[0]][ipCut[1]][ipCut[2]][ipCut[3]][userID] = struct{}{}
		return
	}

	// иначе надо проверить на дубликаты и потенциальные дубликаты

	// проверяем если наш id тут уже есть, тогда просто выходим
	for key := range ipLib[ipCut[0]][ipCut[1]][ipCut[2]][ipCut[3]] {
		if key == userID {
			return
		}
	}

	// если он тут новый то надо проверить его сочетания со всеми прочими id
	for key := range ipLib[ipCut[0]][ipCut[1]][ipCut[2]][ipCut[3]] {
		pair := generatePair(userID, key)
		getNext := false

		// проверяем есть ли эта пара среди дубликатов
		for d := range dupes {
			if d == pair {
				getNext = true
				break
			}
		}
		// если да то берем следующую пару
		if getNext {
			continue
		}

		// если есть среди потенциальных
		for p := range potential {
			if p == pair {
				getNext = true
				break
			}
		}
		//  то переводим в дубликаты и идем дальше
		if getNext {
			dupes[pair] = struct{}{}
			delete(potential, pair)
			continue
		}

		potential[pair] = struct{}{}
	}

	ipLib[ipCut[0]][ipCut[1]][ipCut[2]][ipCut[3]][userID] = struct{}{}
	return
}

func generatePair(first, second int64) string {
	if first < second {
		return strconv.FormatInt(first, 10) + `|` + strconv.FormatInt(second, 10)
	}
	return strconv.FormatInt(second, 10) + `|` + strconv.FormatInt(first, 10)
}

func CheckDupesInIPMap(dupes map[string]interface{}, first, second string) bool {
	if first == second {
		return true
	}

	f, err := strconv.ParseInt(first, 10, 64)
	if err != nil {
		return false
	}

	s, err := strconv.ParseInt(second, 10, 64)
	if err != nil {
		return false
	}

	if _, ok := dupes[generatePair(f, s)]; ok {
		return true
	}

	return false
}
