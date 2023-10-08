package main

type DBServerRow struct {
	Id              int64  `db:"id"`
	Address         string `db:"address"`
	NormalizedDesc  string `db:"normalized_desc"`
	FaviconId       uint64 `db:"favicon_id"`
	FirstSeen       uint32 `db:"first_seen"`
	LastScan        uint32 `db:"last_scan"`
	LastSuccessScan uint32 `db:"last_success_scan"`
	ModList         string `db:"modlist"`
	PlayerCount     uint64 `db:"player_count"`
	PlayerMax       uint64 `db:"player_max"`
	VersionId       uint64 `db:"version_id"`
	VersionName     string `db:"version_name"`
}
