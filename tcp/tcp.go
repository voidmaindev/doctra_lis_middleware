package tcp

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/voidmaindev/doctra_lis_middleware/log"
)

// HL7BufferSize is the size of the HL7 buffer.
const hl7BufferSize = 1 << 15

// TCP is the struct that represents the TCP connection.
type TCP struct {
	Log        *log.Logger
	Listener   net.Listener
	RcvChannel chan RcvData
	Conns      map[string]*ConnData
}

// RcvData is the struct that represents the received data.
type RcvData struct {
	Conn       net.Conn
	ConnString string
	Data       []byte
	Wg         *sync.WaitGroup
}

// PrevData is the struct that represents the previous data of the connection.
type PrevData struct {
	Data    string
	Started bool
}

// ConnData is the struct that represents the connection data.
type ConnData struct {
	Conn       net.Conn
	ConnString string
	PrevData   *PrevData
	Wg         *sync.WaitGroup
}

// NewTCP creates a new TCP connection.
func NewTCP(log *log.Logger, listener net.Listener) *TCP {
	tcp := &TCP{
		Log:      log,
		Listener: listener,
	}

	tcp.Conns = map[string]*ConnData{}
	tcp.RcvChannel = make(chan RcvData)

	return tcp
}

// Accept accepts a connection.
func (t *TCP) AcceptConnections() {
	for {
		conn, err := t.Listener.Accept()
		if err != nil {
			if err == net.ErrClosed || err == io.EOF || err.Error() == "EOF" {
				t.Log.Info("TCP listener closed")
				return //TODO: return or recreate listener and continue
			}
			t.Log.Err(err, "failed to accept a connection")
			continue
		}

		connString := getConnString(conn)

		t.Log.Info("accepted a connection from " + connString)
		t.Conns[connString] = newConnData(conn, connString)

		go t.ReadMessages(conn, t.Conns[connString])
	}
}

// newConnData creates a new connection data.
func newConnData(conn net.Conn, connString string) *ConnData {
	return &ConnData{
		Conn:       conn,
		ConnString: connString,
		PrevData:   &PrevData{},
		Wg:         &sync.WaitGroup{},
	}
}

// ReadMessages reads messages from the connection.
func (t *TCP) ReadMessages(conn net.Conn, connData *ConnData) {
	// connString := getConnString(conn)
	buf := make([]byte, hl7BufferSize) // Allocate buffer once

	defer func() {
		conn.Close()
		// delete(t.Conns, connData.connString)
		t.Log.Info(fmt.Sprintf("connection from %s closed", connData.ConnString))
	}()

	for {
		connData.Wg.Wait()
		n, err := conn.Read(buf)
		if err != nil {
			// if err == net.ErrClosed || err == io.EOF || err.Error() == "EOF" {
			t.Log.Err(err, fmt.Sprintf("connection from %s closed", connData.ConnString))
			// t.Log.Info(fmt.Sprintf("connection from %s closed", connString))
			return
			// }

			// t.Log.Err(err, "failed to read from connection")
			// continue
		}

		rcvData := RcvData{
			Conn:       conn,
			ConnString: connData.ConnString,
			Data:       buf[:n],
			Wg:         connData.Wg,
		}
		connData.Wg.Add(1)
		t.RcvChannel <- rcvData
	}
}

// getConnString gets the connection string.
func getConnString(conn net.Conn) string {
	ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	return ip
}
