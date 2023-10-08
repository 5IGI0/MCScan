package main

import (
	"encoding/base64"
	"encoding/json"
	"sort"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type MCServerStatus struct {
	Version struct {
		Protocol uint64 `json:"protocol"`
		Name     string `json:"name"`
	} `json:"version"`
	Players struct {
		Sample []struct {
			Name string `json:"name"`
			Id   string `json:"id"`
		} `json:"sample"`
		Online uint64 `json:"online"`
		Max    uint64 `json:"max"`
	} `json:"players"`
	Description any    `json:"description"`
	Favicon     string `json:"favicon"`
	ModInfo     struct {
		Type    string `json:"type"`
		ModList []struct {
			ModId   string `json:"modid"`
			Version string `json:"version"`
		} `json:"modList"`
	} `json:"modinfo"`
}

func AnalyzeServerStatus(status string, ret *ServerQueryResult) error {
	var statusobj MCServerStatus
	if err := json.Unmarshal([]byte(status), &statusobj); err != nil {
		return err
	}

	ret.StatusRaw = status
	ret.StatusObj = statusobj
	{
		if str, ok := statusobj.Description.(string); ok {
			ret.NormalizedDesc = str
		} else {
			var chat ChatComponent
			if err := mapstructure.Decode(statusobj.Description, &chat); err != nil {
				return err
			}
			ret.NormalizedDesc = NormalizeChatComponent(chat)
		}
	}
	ret.NormalizedDesc = strings.TrimSpace(delete_mc_symbols(ret.NormalizedDesc))

	sort.Slice(statusobj.ModInfo.ModList, func(i, j int) bool {
		return strings.Compare(statusobj.ModInfo.ModList[i].ModId, statusobj.ModInfo.ModList[j].ModId) < 0
	})

	ret.Mods = statusobj.ModInfo.Type + "||"
	for _, mod := range statusobj.ModInfo.ModList {
		ret.Mods += mod.ModId + ";" + mod.Version + "||"
	}

	if strings.HasPrefix(statusobj.Favicon, "data:image/png;base64,") {
		ret.Favicon = make([]byte, base64.RawStdEncoding.DecodedLen(len(statusobj.Favicon)-22))
		n, _ := base64.RawStdEncoding.Decode(ret.Favicon, []byte(statusobj.Favicon[22:]))
		ret.Favicon = ret.Favicon[:n]
	}

	return nil
}
