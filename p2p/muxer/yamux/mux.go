package sm_yamux

import (
	"io/ioutil"
	"net"

	mux "github.com/libp2p/go-libp2p-core/mux"
	yamux "github.com/libp2p/go-yamux"
)

var DefaultTransport *Multiplexer

func init() {
	config := yamux.DefaultConfig()
	// We've bumped this to 16MiB as this critically limits throughput.
	//
	// 1MiB means a best case of 10MiB/s (83.89Mbps) on a connection with
	// 100ms latency. The default gave us 2.4MiB *best case* which was
	// totally unacceptable.
	config.MaxStreamWindowSize = uint32(16 * 1024 * 1024)
	// don't spam
	config.LogOutput = ioutil.Discard
	// We always run over a security transport that buffers internally
	// (i.e., uses a block cipher).
	config.ReadBufSize = 0
	DefaultTransport = (*Multiplexer)(config)
}

// Multiplexer implements mux.Multiplexer that constructs
// yamux-backed muxed connections.
type Multiplexer yamux.Config

func (t *Multiplexer) NewConn(nc net.Conn, isServer bool) (mux.MuxedConn, error) {
	var s *yamux.Session
	var err error
	if isServer {
		s, err = yamux.Server(nc, t.Config())
	} else {
		s, err = yamux.Client(nc, t.Config())
	}
	return (*conn)(s), err
}

func (t *Multiplexer) Config() *yamux.Config {
	return (*yamux.Config)(t)
}

var _ mux.Multiplexer = &Multiplexer{}
