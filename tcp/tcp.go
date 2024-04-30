package tcp

import (
	"fmt"
	"net"

	"github.com/voidmaindev/doctra_lis_middleware/log"
)

// HL7BufferSize is the size of the HL7 buffer.
const hl7BufferSize = 1 << 15

// TCP is the struct that represents the TCP connection.
type TCP struct {
	Log        *log.Logger
	Listener   net.Listener
	RcvChannel chan RcvData
	Conns      map[string]net.Conn
}

type RcvData struct {
	Conn       net.Conn
	ConnString string
	Data       []byte
}

// NewTCP creates a new TCP connection.
func NewTCP(log *log.Logger, listener net.Listener) *TCP {
	tcp := &TCP{
		Log:      log,
		Listener: listener,
	}

	tcp.Conns = map[string]net.Conn{}
	tcp.RcvChannel = make(chan RcvData)

	return tcp
}

// Accept accepts a connection.
func (t *TCP) AcceptConnections() {
	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			if err == net.ErrClosed || err.Error() == "EOF" {
				t.Log.Info("TCP listener closed")
				return
			}
			t.Log.Err(err, "failed to accept a connection")
		}

		connString := getConnString(conn)

		t.Log.Info("accepted a connection from " + connString)
		t.Conns[connString] = conn

		go t.ReadMessages(conn)
	}
}

// ReadMessages reads messages from the connection.
func (t *TCP) ReadMessages(conn net.Conn) {
	connString := getConnString(conn)
	defer conn.Close()
	defer delete(t.Conns, connString)

	for {
		buf := make([]byte, hl7BufferSize)
		n, err := conn.Read(buf)
		if err != nil {
			if err == net.ErrClosed || err.Error() == "EOF" {
				t.Log.Info(fmt.Sprintf("connection from %s closed", connString))
				return
			}

			t.Log.Err(err, "failed to read from connection")
			continue
		}

		t.RcvChannel <- RcvData{Conn: conn, ConnString: connString, Data: buf[:n]}
	}
}

// Close closes the connection.
func getConnString(conn net.Conn) string {
	ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	return ip
}

