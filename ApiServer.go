package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func ApiGetServerByAddr(w http.ResponseWriter, r *http.Request) (error, interface{}) {
	var id int64

	if err := runtime.db.Get(&id, "SELECT `server_id` FROM `server_aliases` WHERE `str_addr`=?", normalize_addr(mux.Vars(r)["addr"])); err != nil {
		return err, nil
	}

	mux.Vars(r)["id"] = fmt.Sprint(id)
	return ApiGetServerById(w, r)
}

func InternalGetServerById(id int, want_aliases bool) (error, interface{}) {
	var dbrow DBServerRow

	row := runtime.db.QueryRowx(`SELECT
	address, normalized_desc, favicon_id, first_seen,
	last_scan, last_success_scan, modlist, player_count,
	player_max, version_id, version_name
	FROM servers WHERE id=?`, id)

	if err := row.StructScan(&dbrow); err != nil {
		return err, nil
	}

	dbrow.Id = int64(id)
	ret := DBServerRow2Api(dbrow)

	if want_aliases {
		rows, err := runtime.db.Queryx("SELECT `str_addr` FROM `server_aliases` WHERE `server_id`=?", id)
		if err != nil {
			return err, nil
		}
		defer rows.Close()
		for rows.Next() {
			var tmp string
			if err = rows.Scan(&tmp); err != nil {
				return err, nil
			}
			ret.AddressAliases = append(ret.AddressAliases, tmp)
		}
	}

	return nil, ret
}

func ApiGetServerById(w http.ResponseWriter, r *http.Request) (error, interface{}) {
	reqvars := mux.Vars(r)
	reqqers := r.URL.Query()
	srv_id, _ := strconv.Atoi(reqvars["id"])
	return InternalGetServerById(
		srv_id,
		reqqers["full"] != nil || reqqers["aliases"] != nil)
}

func DBServerRow2Api(server DBServerRow) ApiServer {
	var ret ApiServer

	ret.Id = server.Id
	ret.Address = server.Address
	ret.NormalizedDesc = server.NormalizedDesc
	ret.FaviconId = server.FaviconId
	ret.FirstSeen = server.FirstSeen
	ret.LastScan = server.LastScan
	ret.LastSuccessScan = server.LastSuccessScan
	{
		ret.Mods = make(map[string]string)
		tmp := strings.Split(server.ModList, "||")
		for i := 1; i < len(tmp)-1; i++ {
			t := strings.Split(tmp[i], ";")
			ret.Mods[t[0]] = t[1]
		}
	}
	ret.PlayerCount = server.PlayerCount
	ret.PlayerMax = server.PlayerMax
	ret.VersionId = server.VersionId
	ret.VersionName = server.VersionName

	return ret
}
