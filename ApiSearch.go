package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

func ApiSearch(w http.ResponseWriter, r *http.Request) (error, interface{}) {
	add2conds := func(conds *string, to_add string) {
		if *conds == "" {
			*conds += "WHERE " + to_add
		} else {
			*conds += " AND " + to_add
		}
	}
	conds := ""
	condvars := make([]interface{}, 0)
	query := r.URL.Query()

	if query.Has("version") {
		add2conds(&conds, "version_id=?")
		condvars = append(condvars, query.Get("version"))
	}

	if query.Has("mods") {
		mods := strings.Split(query.Get("mods"), ",")
		for i := range mods {
			mods[i] = escape_like(mods[i])
		}
		sort.Slice(mods, func(i, j int) bool {
			return strings.Compare(mods[i], mods[j]) < 0
		})
		add2conds(&conds, "modlist LIKE ?")
		condvars = append(condvars, "%|"+strings.Join(mods, ";%|")+";%")
	}

	if query.Get("text") != "" {
		add2conds(&conds, "normalized_desc LIKE ?")
		condvars = append(condvars, "%"+escape_like(query.Get("text"))+"%")
	}

	if query.Get("page") != "" {
		a, _ := strconv.Atoi(query.Get("page"))
		log.Println(a)
		if a >= 1 {
			a -= 1
		} else {
			a = 0
		}
		log.Println(a)
		condvars = append(condvars, a*100, 100)
	} else {
		condvars = append(condvars, 0, 100)
	}
	rows, err := runtime.db.Queryx(`SELECT
	id, address, normalized_desc, favicon_id, first_seen,
	last_scan, last_success_scan, modlist, player_count,
	player_max, version_id, version_name
	FROM servers `+conds+" ORDER BY `first_seen` DESC LIMIT ?,?", condvars...)

	if err != nil {
		return err, nil
	}
	defer rows.Close()

	var ret ApiSearchResult
	ret.Servers = make([]ApiServer, 0)

	for rows.Next() {
		var tmp DBServerRow

		if err := rows.StructScan(&tmp); err != nil {
			return err, nil
		}

		ret.Servers = append(ret.Servers, DBServerRow2Api(tmp))
	}

	return nil, ret
}
