package rtp

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/AlexxIT/go2rtc/pkg/core"
	"github.com/pion/rtp"
)

// Session represents an RTP session for a specific SSRC
type Session struct {
	core.Connection
	core.Listener
	server *Server

	SSRC     uint32
	LastSeen time.Time

	conn net.PacketConn
	addr net.Addr // remote address

	mu sync.RWMutex
}

// WriteRTP sends an RTP packet to the remote endpoint
func (s *Session) WriteRTP(packet *rtp.Packet) (int, error) {
	if s.conn != nil {
		b, err := packet.Marshal()
		if err != nil {
			return 0, err
		}
		return s.conn.WriteTo(b, s.addr)
	}
	return 0, fmt.Errorf("session connection not available")
}

// ReadRTP processes an incoming RTP packet
func (s *Session) ReadRTP(packet *rtp.Packet) {
	// Update last seen time
	s.LastSeen = time.Now()

	// Directly forward to all receivers
	s.mu.RLocker()
	receivers := make([]*core.Receiver, len(s.Receivers))
	copy(receivers, s.Receivers)
	s.mu.RUnlock()

	for _, r := range receivers {
		r.Input(packet)
	}
}
