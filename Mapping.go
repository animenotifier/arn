package arn

// Register a list of supported services.
func init() {
	DataLists["mapping-services"] = []*Option{
		&Option{"anidb/anime", "anidb/anime"},
		&Option{"anilist/anime", "anilist/anime"},
		&Option{"kitsu/anime", "kitsu/anime"},
		&Option{"kitsu/character", "kitsu/character"},
		&Option{"myanimelist/anime", "myanimelist/anime"},
		&Option{"myanimelist/character", "myanimelist/character"},
		&Option{"myanimelist/producer", "myanimelist/producer"},
		&Option{"shoboi/anime", "shoboi/anime"},
		&Option{"thetvdb/anime", "thetvdb/anime"},
		&Option{"trakt/anime", "trakt/anime"},
	}
}

// Mapping ...
type Mapping struct {
	Service   string `json:"service" editable:"true" datalist:"mapping-services"`
	ServiceID string `json:"serviceId" editable:"true"`
}

// Name ...
func (mapping *Mapping) Name() string {
	switch mapping.Service {
	case "anidb/anime":
		return "AniDB"
	case "anilist/anime":
		return "AniList"
	case "imdb/anime":
		return "IMDb"
	case "kitsu/anime":
		return "Kitsu"
	case "myanimelist/anime":
		return "MAL"
	case "shoboi/anime":
		return "Shoboi"
	case "thetvdb/anime":
		return "TVDB"
	case "trakt/anime":
		return "Trakt"
	default:
		return mapping.Service
	}
}

// Link ...
func (mapping *Mapping) Link() string {
	switch mapping.Service {
	case "kitsu/anime":
		return "https://kitsu.io/anime/" + mapping.ServiceID
	case "shoboi/anime":
		return "http://cal.syoboi.jp/tid/" + mapping.ServiceID
	case "anilist/anime":
		return "https://anilist.co/anime/" + mapping.ServiceID
	case "imdb/anime":
		return "https://www.imdb.com/title/" + mapping.ServiceID
	case "myanimelist/anime":
		return "https://myanimelist.net/anime/" + mapping.ServiceID
	case "thetvdb/anime":
		return "https://thetvdb.com/?tab=series&id=" + mapping.ServiceID
	case "anidb/anime":
		return "https://anidb.net/perl-bin/animedb.pl?show=anime&aid=" + mapping.ServiceID
	case "trakt/anime":
		return "https://trakt.tv/shows/" + mapping.ServiceID
	default:
		return ""
	}
}
