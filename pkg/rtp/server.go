package rtp

import (
	"net"
	"strconv"
	"sync"

	"github.com/pion/rtp"
)

type Server struct {
	address  string
	conn     net.PacketConn
	sessions map[uint32]*Session
	mu       sync.Mutex
}

func NewServer(address string) *Server {
	return &Server{
		address:  address,
		sessions: map[uint32]*Session{},
	}
}

func (s *Server) Port() int {
	if s.conn != nil {
		return s.conn.LocalAddr().(*net.UDPAddr).Port
	}

	_, a, _ := net.SplitHostPort(s.address)
	i, _ := strconv.Atoi(a)
	return i
}

func (s *Server) AddSession(session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.sessions) == 0 {
		var err error
		if s.conn, err = net.ListenPacket("udp", s.address); err != nil {
			return
		}
		go s.handle()
	}

	// Set server reference for session
	session.server = s

	// Update session connection
	session.conn = s.conn

	// Add to sessions map
	s.sessions[session.SSRC] = session
}

func (s *Server) DelSession(ssrc uint32) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.sessions, ssrc)

	// If no more sessions, close the connection
	if len(s.sessions) == 0 && s.conn != nil {
		_ = s.conn.Close()
		s.conn = nil
	}
}

func (s *Server) GetSession(ssrc uint32) (session *Session) {
	s.mu.Lock()
	defer s.mu.Unlock()
	session = s.sessions[ssrc]
	return
}

// handle incoming RTP packets
func (s *Server) handle() error {
	b := make([]byte, 2048)

	for {
		n, addr, err := s.conn.ReadFrom(b)
		if err != nil {
			return err
		}

		// Minimum size for RTP packet (header only)
		if n < 12 {
			continue
		}

		// Parse RTP packet
		packet := &rtp.Packet{}
		if err := packet.Unmarshal(b[:n]); err != nil {
			continue // not a valid RTP packet
		}

		// Check if we have a session for this SSRC
		session := s.GetSession(packet.SSRC)
		if session == nil {
			// No session registered for this SSRC, skip
			continue
		}

		// Update session remote address
		if session.addr == nil {
			session.addr = addr
		} else if session.addr.String() != addr.String() {
			// SSRC changed source address, skip
			continue
		}

		// Process packet in session
		session.ReadRTP(packet)
	}
}
