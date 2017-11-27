package arn

// Mapping ...
type Mapping struct {
	Service   string `json:"service" editable:"true"`
	ServiceID string `json:"serviceId" editable:"true"`
	Created   string `json:"created"`
	CreatedBy string `json:"createdBy"`
}

// Name ...
func (mapping *Mapping) Name() string {
	switch mapping.Service {
	case "shoboi/anime":
		return "Shoboi"
	case "anilist/anime":
		return "AniList"
	case "myanimelist/anime":
		return "MyAnimeList"
	case "thetvdb/anime":
		return "TheTVDB"
	case "anidb/anime":
		return "AniDB"
	default:
		return ""
	}
}

// Link ...
func (mapping *Mapping) Link() string {
	switch mapping.Service {
	case "shoboi/anime":
		return "http://cal.syoboi.jp/tid/" + mapping.ServiceID
	case "anilist/anime":
		return "https://anilist.co/anime/" + mapping.ServiceID
	case "myanimelist/anime":
		return "https://myanimelist.net/anime/" + mapping.ServiceID
	case "thetvdb/anime":
		return "https://thetvdb.com/?tab=series&id=" + mapping.ServiceID
	case "anidb/anime":
		return "https://anidb.net/perl-bin/animedb.pl?show=anime&aid=" + mapping.ServiceID
	default:
		return ""
	}
}
