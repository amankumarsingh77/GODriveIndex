package config

type PlayerConfig struct {
	Player          string `json:"player"`
	VideojsVersion  string `json:"videojs_version"`
	PlyrIOVersion   string `json:"plyr_io_version"`
	JwplayerVersion string `json:"jwplayer_version"`
}
