package tcpnet

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/zen099/onetube/server/base/grmon"
	"github.com/zen099/onetube/server/base/guid"
	. "github.com/zen099/onetube/server/base/logger"
)

type TcpContext struct {
	LocalAddr  string
	RemoteAddr string
	PacketNum  int64
	Params     *sync.Map
	ConnID     int64
	Conn       net.Conn
	Mutex      *sync.Mutex
	writeChan  chan interface{}
	closeChan  chan interface{}
}

func (c *TcpContext) ReadFull(b []byte) error {
	_, err := io.ReadFull(c.Conn, b)
	return err
}

func (c *TcpContext) WriteFull(b []byte) {
	c.writeChan <- b
}

func (c *TcpContext) Close() {
	c.closeChan <- 1
}

func (c *TcpContext) Set(k interface{}, v interface{}) {
	c.Params.Store(k, v)
}

func (c *TcpContext) Get(k interface{}) (interface{}, bool) {
	return c.Params.Load(k)
}

func (c *TcpContext) MustGet(k interface{}) interface{} {
	obj, ok := c.Params.Load(k)
	if !ok {
		s := fmt.Sprintf("k=[%v] does not exist in TcpContext", k)
		Log.Error(s)
		panic(s)
	}

	return obj
}

func (c *TcpContext) Remove(k interface{}) {
	c.Params.Delete(k)
}

type TcpServer struct {
	connectionLimitChan chan int
	workerLimitChan     chan int

	MaxConnection   int
	MaxWorker       int
	MaxWriteChanLen int
	Address         string
	OnAccept        func(ctx *TcpContext) error
	OnData          func(ctx *TcpContext) (interface{}, error)
	OnDataHandler   func(ctx *TcpContext, m interface{}) error
	OnClose         func(ctx *TcpContext) error
}

func NewTcpServer(addr string) *TcpServer {
	s := &TcpServer{
		Address: addr,
	}
	return s
}

func (s *TcpServer) Run() error {
	addr := s.Address
	maxconn := s.MaxConnection
	if maxconn <= 0 {
		maxconn = 100
	}
	maxworker := s.MaxWorker
	if maxworker <= 0 {
		maxworker = 200
	}
	writeChanLen := s.MaxWriteChanLen
	if writeChanLen <= 0 {
		writeChanLen = 8192
	}

	s.connectionLimitChan = make(chan int, maxconn)
	s.workerLimitChan = make(chan int, maxworker)

	Log.Infof("TCP server listening at %v", addr)
	lsn, err := net.Listen("tcp", addr)
	if err != nil {
		Log.Error(err)
		return err
	}

	defer lsn.Close()
	for {
		conn, err := lsn.Accept()
		if err != nil {
			Log.Warn(err)
			continue
		}
		Log.Debugf("s:[%v] accept connection from %v", addr, conn.RemoteAddr().String())

		gm := guid.GetGUIDManager()
		connid := gm.Save(conn)
		ctx := &TcpContext{
			LocalAddr:  addr,
			RemoteAddr: conn.RemoteAddr().String(),
			PacketNum:  0,
			Params:     &sync.Map{},
			ConnID:     connid,
			Conn:       conn,
			Mutex:      &sync.Mutex{},
			writeChan:  make(chan interface{}, writeChanLen),
			closeChan:  make(chan interface{}, 1),
		}

		connClose := func() {
			s.OnClose(ctx)
			conn.Close()
			gm.Remove(connid)
		}

		if s.OnAccept != nil {
			err = s.OnAccept(ctx)
			if err != nil {
				Log.Error(err)
				connClose()
				return err
			}
		}

		// 控制handler并发数量
		s.connectionLimitChan <- 1

		grm := grmon.GetGRMon()
		grm.Go("tcp_connection_handler", func() {

			defer func() { <-s.connectionLimitChan }()

			s.ConnectionHandler(ctx)
		})
	}
}

func (s *TcpServer) ConnectionHandler(ctx *TcpContext) {

	connid := ctx.ConnID
	conn := ctx.Conn
	writeChan := ctx.writeChan
	closeChan := ctx.closeChan
	remoteAddr := conn.RemoteAddr().String()

	connClose := func(all bool) {
		conn.Close()
		gm := guid.GetGUIDManager()
		gm.Remove(ctx.ConnID)
		if all {
			s.OnClose(ctx)
			close(writeChan)
			close(closeChan)
		}
	}

	grm := grmon.GetGRMon()
	grm.Go("tcp_read", func() {
		for {
			msg, err := s.OnData(ctx)
			if err != nil {
				Log.Warnf("err = %v, connid = %v, remoteAddr = %v will be closed",
					err, connid, remoteAddr)
				connClose(true)
				return
			}
			ctx.PacketNum++

			// 控制并发数量
			s.workerLimitChan <- 1
			grm.Go("tcp_worker", func() {

				defer func() { <-s.workerLimitChan }()

				err := s.OnDataHandler(ctx, msg)
				if err != nil {
					Log.Error(err)
					connClose(false)
					return
				}
			})
		}
	})

	grm.Go("tcp_write", func() {
		for {
			var data interface{}
			var ok bool
			select {
			case data, ok = <-writeChan:
				if !ok {
					Log.Warnf("resp channel closed, write routine will exit, connid = %v", connid)
					return
				}
			case <-closeChan:
				Log.Warnf("recieved from close channel, will close the connection, connid = %v", connid)
				connClose(false) // 此处只尝试关闭连接即可
				return
			}

			if data == nil { // 空数据则不处理
				Log.Warnf("tcp will write nil. connid = %v", connid)
				continue
			}

			databytes := data.([]byte)
			nlen := len(databytes)
			left := nlen
			for left > 0 {
				sentdata := databytes[nlen-left:]
				// Log.Debugf("length = %v, sentdata = %v", len(databytes), hex.EncodeToString(sentdata))
				n, err := conn.Write(sentdata)
				if err != nil {
					Log.Error(err)
					connClose(false)
					return
				}
				left = left - n
			}
		}
	})
}
