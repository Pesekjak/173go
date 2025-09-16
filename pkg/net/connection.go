package net

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/Pesekjak/173go/pkg/base"
	"github.com/Pesekjak/173go/pkg/chat"
	"github.com/Pesekjak/173go/pkg/log"
	"github.com/Pesekjak/173go/pkg/prot"
	"github.com/Pesekjak/173go/pkg/system"
)

type Connection struct {
	tcp    *net.TCPConn
	reader *bufio.Reader
	writer *bufio.Writer

	listening bool
	closed    bool

	mu sync.Mutex

	logger *log.Logger

	report chan system.Message
}

func NewConnection(tcp *net.TCPConn, logger *log.Logger, report chan system.Message) *Connection {
	return &Connection{
		tcp:    tcp,
		reader: bufio.NewReader(tcp),
		writer: bufio.NewWriter(tcp),

		listening: false,
		closed:    false,

		logger: logger,

		report: report,
	}
}

func (c *Connection) isActive() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.listening && !c.closed
}

func (c *Connection) isClosed() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closed
}

func (c *Connection) startListening(handler prot.PacketHandler) error {
	c.mu.Lock()
	if c.listening {
		c.mu.Unlock()
		return fmt.Errorf("this connection is already listening")
	}
	c.listening = true
	c.mu.Unlock()

	c.logger.Info("handling new connection from ", c)

	for {
		packetID, err := c.reader.ReadByte()

		if err != nil {
			c.mu.Lock()
			closed := c.closed
			c.mu.Unlock()

			// the connection was closed somewhere else
			if closed {
				return nil
			}

			if err == io.EOF {
				c.logger.Info("connection closed by ", c)
				c.Close(nil)
				return nil
			}

			c.logger.Severe("error reading packet ID: ", err)
			c.Close(err)
			return err
		}

		var packet prot.PacketIn
		packet, err = prot.ReadPacket(packetID, c.reader)
		if err != nil {
			c.logger.Severe("error reading packet: ", err)
			c.Close(err)
			return err
		}

		err = packet.Handle(handler)
		if err != nil {
			c.logger.Severe("error handling packet: ", err)
			c.Close(err)
			return err
		}
	}
}

func (c *Connection) WritePacket(packet prot.PacketOut, flush bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return io.ErrClosedPipe
	}

	c.logger.Info("send packet: ", packet)
	if err := prot.WritePacket(packet, c.writer); err != nil {
		return err
	}
	if flush {
		if err := c.writer.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Connection) Close(err error) {
	if err == nil {
		c.CloseWith(nil, "")
	}
	c.CloseWith(err, base.ConvertToString(err))
}

func (c *Connection) CloseWith(err error, reason string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.closed {
		return
	}
	c.closed = true

	if reason == "" && err != nil {
		reason = base.ConvertToString(err)
	}

	if reason != "" {
		// try to kick with error message for context before closing the connection
		// write directly to prevent deadlock, we hold the connection mutex
		_ = prot.WritePacket(&prot.PacketOutKick{Reason: reason}, c.writer)
		_ = c.writer.Flush()
	}

	if err != nil {
		c.logger.Severe("closed connection ", c, ": ", err)
	} else if reason != "" {
		c.logger.Info("closed connection ", c, ": ", chat.StripColorCodes(reason))
	} else {
		c.logger.Info("closed connection ", c)
	}

	err = c.tcp.Close()
	if err != nil {
		c.logger.Severe("failed to close connection ", c, ": ", err)
	}
}

func (c *Connection) String() string {
	return c.tcp.RemoteAddr().String()
}
