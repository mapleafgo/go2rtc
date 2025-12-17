package rtp

import (
	"github.com/AlexxIT/go2rtc/pkg/core"
)

// Start starts the producer session
func (s *Session) Start() error {
	return nil
}

// Stop stops the producer session and cleans up resources
func (s *Session) Stop() error {
	s.server.DelSession(s.SSRC)
	return nil
}

// GetTrack returns a receiver for the specified media and codec
func (s *Session) GetTrack(media *core.Media, codec *core.Codec) (*core.Receiver, error) {
	core.Assert(media.Direction == core.DirectionRecvonly)

	// Check if receiver already exists
	for _, track := range s.Receivers {
		if track.Codec == codec {
			return track, nil
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	track := core.NewReceiver(media, codec)
	s.Receivers = append(s.Receivers, track)

	return track, nil
}
