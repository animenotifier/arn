package arn

// Register a list of gaming servers.
func init() {
	DataLists["ffxiv-servers"] = []*Option{
		&Option{"", ""},
		&Option{"Adamantoise", "Adamantoise"},
		&Option{"Aegis", "Aegis"},
		&Option{"Alexander", "Alexander"},
		&Option{"Anima", "Anima"},
		&Option{"Asura", "Asura"},
		&Option{"Atomos", "Atomos"},
		&Option{"Bahamut", "Bahamut"},
		&Option{"Balmung", "Balmung"},
		&Option{"Behemoth", "Behemoth"},
		&Option{"Belias", "Belias"},
		&Option{"Brynhildr", "Brynhildr"},
		&Option{"Cactuar", "Cactuar"},
		&Option{"Carbuncle", "Carbuncle"},
		&Option{"Cerberus", "Cerberus"},
		&Option{"Chocobo", "Chocobo"},
		&Option{"Coeurl", "Coeurl"},
		&Option{"Diabolos", "Diabolos"},
		&Option{"Durandal", "Durandal"},
		&Option{"Excalibur", "Excalibur"},
		&Option{"Exodus", "Exodus"},
		&Option{"Faerie", "Faerie"},
		&Option{"Famfrit", "Famfrit"},
		&Option{"Fenrir", "Fenrir"},
		&Option{"Garuda", "Garuda"},
		&Option{"Gilgamesh", "Gilgamesh"},
		&Option{"Goblin", "Goblin"},
		&Option{"Gungnir", "Gungnir"},
		&Option{"Hades", "Hades"},
		&Option{"Hyperion", "Hyperion"},
		&Option{"Ifrit", "Ifrit"},
		&Option{"Ixion", "Ixion"},
		&Option{"Jenova", "Jenova"},
		&Option{"Kujata", "Kujata"},
		&Option{"Lamia", "Lamia"},
		&Option{"Leviathan", "Leviathan"},
		&Option{"Lich", "Lich"},
		&Option{"Louisoix", "Louisoix"},
		&Option{"Malboro", "Malboro"},
		&Option{"Mandragora", "Mandragora"},
		&Option{"Masamune", "Masamune"},
		&Option{"Mateus", "Mateus"},
		&Option{"Midgardsormr", "Midgardsormr"},
		&Option{"Moogle", "Moogle"},
		&Option{"Odin", "Odin"},
		&Option{"Omega", "Omega"},
		&Option{"Pandaemonium", "Pandaemonium"},
		&Option{"Phoenix", "Phoenix"},
		&Option{"Ragnarok", "Ragnarok"},
		&Option{"Ramuh", "Ramuh"},
		&Option{"Ridill", "Ridill"},
		&Option{"Sargatanas", "Sargatanas"},
		&Option{"Shinryu", "Shinryu"},
		&Option{"Shiva", "Shiva"},
		&Option{"Siren", "Siren"},
		&Option{"Tiamat", "Tiamat"},
		&Option{"Titan", "Titan"},
		&Option{"Tonberry", "Tonberry"},
		&Option{"Typhon", "Typhon"},
		&Option{"Ultima", "Ultima"},
		&Option{"Ultros", "Ultros"},
		&Option{"Unicorn", "Unicorn"},
		&Option{"Valefor", "Valefor"},
		&Option{"Yojimbo", "Yojimbo"},
		&Option{"Zalera", "Zalera"},
		&Option{"Zeromus", "Zeromus"},
		&Option{"Zodiark", "Zodiark"},
	}
}

// UserAccounts represents a user's accounts on external services.
type UserAccounts struct {
	Facebook struct {
		ID string `json:"id" private:"true"`
	} `json:"facebook"`

	Google struct {
		ID string `json:"id" private:"true"`
	} `json:"google"`

	Twitter struct {
		ID   string `json:"id" private:"true"`
		Nick string `json:"nick"`
	} `json:"twitter"`

	Osu struct {
		Nick     string  `json:"nick" editable:"true"`
		PP       float64 `json:"pp"`
		Accuracy float64 `json:"accuracy"`
		Level    float64 `json:"level"`
	} `json:"osu"`

	Overwatch struct {
		BattleTag   string `json:"battleTag" editable:"true"`
		SkillRating int    `json:"skillRating"`
		Tier        string `json:"tier"`
	} `json:"overwatch"`

	FinalFantasyXIV struct {
		Nick      string `json:"nick" editable:"true"`
		Server    string `json:"server" editable:"true" datalist:"ffxiv-servers"`
		Class     string `json:"class"`
		Level     int    `json:"level"`
		ItemLevel int    `json:"itemLevel"`
	} `json:"ffxiv"`

	AniList struct {
		Nick string `json:"nick" editable:"true"`
	} `json:"anilist"`

	AnimePlanet struct {
		Nick string `json:"nick" editable:"true"`
	} `json:"animeplanet"`

	MyAnimeList struct {
		Nick string `json:"nick" editable:"true"`
	} `json:"myanimelist"`

	Kitsu struct {
		Nick string `json:"nick" editable:"true"`
	} `json:"kitsu"`
}
