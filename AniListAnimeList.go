package arn

// AniListAnimeList ...
type AniListAnimeList struct {
	ID          int    `json:"id"`
	DisplayName string `json:"display_name"`
	// AnimeTime       int           `json:"anime_time"`
	// MangaChap       int           `json:"manga_chap"`
	// About           string        `json:"about"`
	// ListOrder       int           `json:"list_order"`
	// AdultContent    bool          `json:"adult_content"`
	// ForumHomepage   int           `json:"forum_homepage"`
	// LegacyLists     bool          `json:"legacy_lists"`
	// Donator         int           `json:"donator"`
	// Following       bool          `json:"following"`
	// ImageURLLge     string        `json:"image_url_lge"`
	// ImageURLMed     string        `json:"image_url_med"`
	// ImageURLBanner  string        `json:"image_url_banner"`
	// TitleLanguage   string        `json:"title_language"`
	// ScoreType       int           `json:"score_type"`
	// CustomListAnime []string      `json:"custom_list_anime"`
	// CustomListManga []interface{} `json:"custom_list_manga"`
	// Stats           struct {
	// 	StatusDistribution struct {
	// 		Anime struct {
	// 			Watching    int `json:"watching"`
	// 			PlanToWatch int `json:"plan to watch"`
	// 			Completed   int `json:"completed"`
	// 			Dropped     int `json:"dropped"`
	// 			OnHold      int `json:"on-hold"`
	// 		} `json:"anime"`
	// 		Manga struct {
	// 			Reading    int `json:"reading"`
	// 			PlanToRead int `json:"plan to read"`
	// 			Completed  int `json:"completed"`
	// 			Dropped    int `json:"dropped"`
	// 			OnHold     int `json:"on-hold"`
	// 		} `json:"manga"`
	// 	} `json:"status_distribution"`
	// 	ScoreDistribution struct {
	// 		Anime struct {
	// 			Num10  int `json:"10"`
	// 			Num20  int `json:"20"`
	// 			Num30  int `json:"30"`
	// 			Num40  int `json:"40"`
	// 			Num50  int `json:"50"`
	// 			Num60  int `json:"60"`
	// 			Num70  int `json:"70"`
	// 			Num80  int `json:"80"`
	// 			Num90  int `json:"90"`
	// 			Num100 int `json:"100"`
	// 		} `json:"anime"`
	// 		Manga struct {
	// 			Num10  int `json:"10"`
	// 			Num20  int `json:"20"`
	// 			Num30  int `json:"30"`
	// 			Num40  int `json:"40"`
	// 			Num50  int `json:"50"`
	// 			Num60  int `json:"60"`
	// 			Num70  int `json:"70"`
	// 			Num80  int `json:"80"`
	// 			Num90  int `json:"90"`
	// 			Num100 int `json:"100"`
	// 		} `json:"manga"`
	// 	} `json:"score_distribution"`
	// 	FavouriteGenres struct {
	// 		Drama         int `json:"Drama"`
	// 		Action        int `json:"Action"`
	// 		Comedy        int `json:"Comedy"`
	// 		Adventure     int `json:"Adventure"`
	// 		Romance       int `json:"Romance"`
	// 		Psychological int `json:"Psychological"`
	// 		Fantasy       int `json:"Fantasy"`
	// 		Mystery       int `json:"Mystery"`
	// 		SciFi         int `json:"Sci-Fi"`
	// 		Thriller      int `json:"Thriller"`
	// 	} `json:"favourite_genres"`
	// 	ListScores struct {
	// 		Anime struct {
	// 			Mean              int `json:"mean"`
	// 			StandardDeviation int `json:"standard_deviation"`
	// 		} `json:"anime"`
	// 		Manga struct {
	// 			Mean              int `json:"mean"`
	// 			StandardDeviation int `json:"standard_deviation"`
	// 		} `json:"manga"`
	// 	} `json:"list_scores"`
	// } `json:"stats"`
	// AdvancedRating      bool          `json:"advanced_rating"`
	// AdvancedRatingNames []interface{} `json:"advanced_rating_names"`
	// Notifications       int           `json:"notifications"`
	// AiringNotifications int           `json:"airing_notifications"`
	// UpdatedAt           int           `json:"updated_at"`
	Lists struct {
		Completed   []*AniListAnimeListItem `json:"completed"`
		Watching    []*AniListAnimeListItem `json:"watching"`
		PlanToWatch []*AniListAnimeListItem `json:"plan_to_watch"`
		Dropped     []*AniListAnimeListItem `json:"dropped"`
		OnHold      []*AniListAnimeListItem `json:"on_hold"`
	} `json:"lists"`
	CustomLists interface{} `json:"custom_lists"`
	// CSS         struct {
	// 	URL    string `json:"url"`
	// 	Popup  bool   `json:"popup"`
	// 	Count  bool   `json:"count"`
	// 	Covers string `json:"covers"`
	// } `json:"css"`
}
