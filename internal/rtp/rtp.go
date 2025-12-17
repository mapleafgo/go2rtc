package rtp

import (
	"strings"

	"github.com/AlexxIT/go2rtc/internal/app"
	"github.com/AlexxIT/go2rtc/internal/streams"
	"github.com/AlexxIT/go2rtc/pkg/core"
	"github.com/AlexxIT/go2rtc/pkg/rtp"
)

func Init() {
	var cfg struct {
		Mod struct {
			Listen string `yaml:"listen"`
		} `yaml:"rtp"`
	}

	// default config
	cfg.Mod.Listen = ":8000"

	// load config from YAML
	app.LoadConfig(&cfg)

	if cfg.Mod.Listen == "" {
		return
	}

	// create RTP server (endpoint) for receiving RTP streams
	Server = rtp.NewServer(cfg.Mod.Listen)

	// Register RTP handler for incoming streams
	streams.HandleFunc("rtp", NewProducer)
}

// Server Global RTP server for handling all RTP traffic
var Server *rtp.Server

func NewProducer(url string) (core.Producer, error) {
	s := &rtp.Session{
		Connection: core.Connection{
			ID:         core.NewID(),
			FormatName: "rtp",
		},
	}
	ssrc, rawQuery, _ := strings.Cut(url, "#")

	if rawQuery != "" {
		query := streams.ParseQuery(rawQuery)
		_ = query.Get("tcp") == "true"
	}

	s.SSRC = uint32(core.Atoi(ssrc))

	Server.AddSession(s)

	return s, nil
}
