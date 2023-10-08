package main

type ApiResponseBase struct {
	Error *string     `json:"error"`
	Data  interface{} `json:"data"`
}

type ApiServer struct {
	Id              int64             `json:"id"`
	Address         string            `json:"address"`
	AddressAliases  []string          `json:"address_aliases"`
	NormalizedDesc  string            `json:"normalized_description"`
	Mods            map[string]string `json:"mods"`
	FaviconId       uint64            `json:"favicon_id"`
	FirstSeen       uint32            `json:"first_seen"`
	LastScan        uint32            `json:"last_scan"`
	LastSuccessScan uint32            `json:"last_success_scan"`
	PlayerCount     uint64            `json:"player_count"`
	PlayerMax       uint64            `json:"player_max"`
	VersionId       uint64            `json:"version_id"`
	VersionName     string            `json:"version_name"`
}

type ApiSearchResult struct {
	Servers []ApiServer `json:"servers"`
}
