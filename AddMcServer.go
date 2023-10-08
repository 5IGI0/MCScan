package main

import (
	"crypto/sha1"
	"log"
	"math/big"
	"net"
	"strings"
	"sync"
)

var addMutex sync.Mutex

func add_server_IP2I(ip string) *uint32 {
	if pip := net.ParseIP(strings.Split(ip, ":")[0]); pip != nil && strings.Contains(ip, ".") {
		i := big.NewInt(0)
		i.SetBytes(pip)
		ret := uint32(i.Uint64())
		return &ret
	}
	return nil
}

// server_id == 0 -> new
func add_server(serv_addr string, server_id int64) {
	serv_addr = normalize_addr(serv_addr)

	addMutex.Lock()
	defer addMutex.Unlock()
	log.Println("scanning", serv_addr, "...")

	if server_id == 0 {
		var count int
		assert_err(runtime.db.Get(&count, "SELECT COUNT(*) FROM `server_aliases` WHERE str_addr=?", serv_addr))
		if count != 0 {
			return
		}
	}

	ret, err := GetServerStatus(serv_addr)
	if err != nil {
		if server_id != 0 {
			runtime.db.MustExec("UPDATE `servers` SET `last_scan`=? WHERE `id`=?", get_timestamp(), server_id)
		}
		log.Println("error while getting", serv_addr, "'s status:", err)
		return
	}

	var favicon_id uint64 = 0
	if ret.Favicon != nil {
		favicon_hash := sha1.Sum(ret.Favicon)
		favicon_id =
			uint64(favicon_hash[0]) |
				(uint64(favicon_hash[1]) << 8) |
				(uint64(favicon_hash[2]) << 16) |
				(uint64(favicon_hash[3]) << 24) |
				(uint64(favicon_hash[4]) << 32) |
				(uint64(favicon_hash[5]) << 40) |
				(uint64(favicon_hash[6]) << 48) |
				(uint64(favicon_hash[7]) << 56)
		runtime.db.Exec("INSERT INTO `favicons`(`id`,`raw_favicon`) VALUES(?,?)",
			favicon_id, ret.Favicon)
	}

	if server_id == 0 {
		server_id, _ = runtime.db.MustExec(
			"INSERT INTO `servers`(`address`,`normalized_desc`,`favicon_id`,`first_seen`,`last_scan`,"+
				"`last_success_scan`,`modlist`,`player_count`,`player_max`,`version_id`,`version_name`) "+
				"VALUES(?,?,?,?,?,?,?,?,?,?,?)",
			ret.Address, ret.NormalizedDesc, favicon_id, get_timestamp(), get_timestamp(), get_timestamp(),
			ret.Mods, ret.StatusObj.Players.Online, ret.StatusObj.Players.Max, ret.StatusObj.Version.Protocol, delete_mc_symbols(ret.StatusObj.Version.Name)).LastInsertId()
	} else {
		runtime.db.MustExec(
			"UPDATE `servers` SET `address`=?,`normalized_desc`=?,`favicon_id`=?,`last_scan`=?,"+
				"`last_success_scan`=?,`modlist`=?,`player_count`=?,`player_max`=?,`version_id`=?,`version_name`=? WHERE id=?",
			ret.Address, ret.NormalizedDesc, favicon_id, get_timestamp(), get_timestamp(),
			ret.Mods, ret.StatusObj.Players.Online, ret.StatusObj.Players.Max, ret.StatusObj.Version.Protocol, delete_mc_symbols(ret.StatusObj.Version.Name),
			server_id)
	}

	{
		var args []interface{}
		tmp := "("
		for _, alias := range ret.AddrAliases {
			tmp += "?,"
			args = append(args, alias)
		}
		tmp = tmp[:len(tmp)-1] + ")"
		runtime.db.MustExec("DELETE FROM `server_aliases` WHERE `str_addr` IN "+tmp, args...)
		args = append(args, server_id)
		runtime.db.MustExec("DELETE FROM `servers` WHERE `address` IN "+tmp+" AND `id` != ?", args...)
		statment, err := runtime.db.Prepare("INSERT INTO `server_aliases`(`server_id`,`str_addr`,`int_addr`) VALUES(?,?,?)")
		assert_err(err)

		for _, alias := range ret.AddrAliases {
			_, err = statment.Exec(server_id, alias, add_server_IP2I(alias))
			assert_err(err)
		}
	}
	{
		runtime.db.MustExec("DELETE FROM `mods` WHERE `server_id`=?", server_id)
		statment, err := runtime.db.Prepare("INSERT INTO `mods`(`server_id`,`modid`,`version`) VALUES(?,?,?)")
		assert_err(err)

		for _, mod := range ret.StatusObj.ModInfo.ModList {
			_, err = statment.Exec(server_id, mod.ModId, mod.Version)
			assert_err(err)
		}
	}

	runtime.db.MustExec("INSERT INTO `status_history`(`server_id`,`at`,`status`,`aliases`) VALUES(?,?,?,?)",
		server_id, get_timestamp(), ret.StatusRaw, strings.Join(ret.AddrAliases, "|"))
}
