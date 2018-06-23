package api

type spotifyTokenResponse struct {
	AccessToken    string `json:"access_token"`
	TokenType      string `json:"token_type"`
	Scope          string `json:"scope"`
	ExpirationTime int    `json:"expires_in"`
	RefreshToken   string `json:"refresh_token"`
}

type spotifyAlbum struct {
	AlbumType            string                `json:"album_type"`
	Artists              []spotifySimpleArtist `json:"artists"`
	AvailableMarkets     []string              `json:"available_markets"`
	Copyrights           []copyright           `json:"copyrights"`
	ExternalIds          []string              `json:"-"`
	ExternalUrls         []externalURL         `json:"-"`
	Genres               []string              `json:"genres"`
	Href                 string                `json:"href"`
	ID                   string                `json:"id"`
	Images               []image               `json:"images"`
	Label                string                `json:"label"`
	Name                 string                `json:"name"`
	Popularity           int                   `json:"popularity"`
	ReleaseDate          string                `json:"release_date"`
	ReleaseDatePrecision string                `json:"release_date_precision"`
	Restrictions         []string              `json:"-"`
	Tracks               []string              `json:"-"`
	ObjectType           string                `json:"type"`
	URI                  string                `json:"uri"`
}

type spotifySimpleArtist struct {
	ExternalUrls []externalURL `json:"-"`
	Href         string        `json:"href"`
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	ObjectType   string        `json:"type"`
	URI          string        `json:"uri"`
}

type spotifySimplePlaylist struct {
	Collaborative bool          `json:"collaborative"`
	Href          string        `json:"href"`
	ExternalUrls  []externalURL `json:"-"`
	ID            string        `json:"id"`
	Images        []image       `json:"images"`
	Name          string        `json:"name"`
	Owner         spotifyUser   `json:"owner"`
	Public        bool          `json:"public"`
	SnapshotID    string        `json:"snapshot_id"`
	Tracks        spotifyTracks `json:"tracks"`
	ObjectType    string        `json:"type"`
	URI           string        `json:"uri"`
}

type spotifyPlaylistTrack struct {
	AddedAt   string       `json:"added_at"`
	AddedBy   spotifyUser  `json:"added_by"`
	LocalFile bool         `json:"is_local"`
	Track     spotifyTrack `json:"track"`
}

type spotifyTrack struct {
	Album            spotifySimpleAlbum    `json:"album"`
	Artists          []spotifySimpleArtist `json:"artists"`
	AvailableMarkets []string              `json:"available_markets"`
	DiscNumber       int                   `json:"disc_number"`
	DurationMS       int                   `json:"duration_ms"`
	Explicit         bool                  `json:"explicit"`
	ExternalID       []string              `json:"-"`
	ExternalUrls     []externalURL         `json:"-"`
	Href             string                `json:"href"`
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	Popularity       int                   `json:"popularity"`
	PreviewURL       string                `json:"preview_url"`
	TrackNumber      int                   `json:"track_number"`
	ObjectType       string                `json:"type"`
	URI              string                `json:"uri"`
}

type spotifySimpleAlbum struct {
	AlbumType            string                `json:"album_type"`
	Artists              []spotifySimpleArtist `json:"artists"`
	AvailableMarkets     []string              `json:"available_markets"`
	ExternalUrls         []externalURL         `json:"-"`
	Href                 string                `json:"href"`
	ID                   string                `json:"id"`
	Images               []image               `json:"images"`
	Name                 string                `json:"name"`
	ReleaseDate          string                `json:"release_date"`
	ReleaseDatePrecision string                `json:"release_date_precision"`
	Restrictions         []string              `json:"-"`
	ObjectType           string                `json:"type"`
	URI                  string                `json:"uri"`
}

type spotifyUser struct {
	DisplayName  string        `json:"display_name"`
	ExternalUrls []externalURL `json:"-"`
	Followers    string        `json:"-"`
	Href         string        `json:"href"`
	ID           string        `json:"id"`
	Images       []image       `json:"images"`
	ObjectType   string        `json:"type"`
	URI          string        `json:"uri"`
}

type spotifyTracks struct {
	TracksURI      string `json:"href"`
	NumberOfTracks int    `json:"total"`
}

type image struct {
	Height int    `json:"height"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
}

type externalURL struct {
	Type  string
	Value string
}

type copyright struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type spotifyBasePage struct {
	Href     string `json:"href"`
	Limit    int    `json:"limit"`
	Next     string `json:"next"`
	Offset   int    `json:"offset"`
	Previous string `json:"previous"`
	Total    int    `json:"total"`
}

type spotifyAlbumPage struct {
	spotifyBasePage
	Albums []savedAlbum `json:"items"`
}

type spotifyPlaylistPage struct {
	spotifyBasePage
	Playlists []spotifySimplePlaylist `json:"items"`
}

type spotifyTrackPage struct {
	spotifyBasePage
	PlaylistTracks []spotifyPlaylistTrack `json:"items"`
}

type savedAlbum struct {
	AddedAt string       `json:"added_at"`
	Album   spotifyAlbum `json:"album"`
}
