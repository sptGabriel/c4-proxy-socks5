package app

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/sptGabriel/socks5/app/packets"
	"github.com/sptGabriel/socks5/app/packets/serverpackets"
)

type Server struct {
	Logger                               *log.Logger
	AuthNoAuthenticationRequiredCallback func(conn *Conn) error
	AuthUsernamePasswordCallback         func(conn *Conn, username, password []byte) error
	connectHandlers                      []ConnectHandler
	closeHandlers                        []CloseHandler
}

type Conn struct {
	server     *Server
	rwc        net.Conn
	Data       interface{}
	externalIP string
}

func New() *Server {
	return &Server{
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (srv *Server) HandleConnect(h ConnectHandler) {
	srv.connectHandlers = append(srv.connectHandlers, h)
}

func (srv *Server) HandleConnectFunc(h func(c *Conn, host string) (newHost string, err error)) {
	srv.connectHandlers = append(srv.connectHandlers, FuncConnectHandler(h))
}

func (srv *Server) HandleClose(h CloseHandler) {
	srv.closeHandlers = append(srv.closeHandlers, h)
}

func (srv *Server) HandleCloseFunc(h func(c *Conn)) {
	srv.closeHandlers = append(srv.closeHandlers, FuncCloseHandler(h))
}

func (srv *Server) ListenAndServe(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				srv.Logger.Printf("socks5: Accept error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}
		tempDelay = 0

		c, err := srv.newConn(rw)
		if err != nil {
			srv.Logger.Printf("socks5: Server.newConn: %v", err)
			continue
		}
		go c.serve()
	}
}

func (srv *Server) newConn(c net.Conn) (*Conn, error) {
	conn := &Conn{
		server: srv,
		rwc:    c,
	}
	return conn, nil
}

func (c *Conn) RemoteAddr() string {
	return c.rwc.RemoteAddr().String()
}

func (c *Conn) LocalAddr() string {
	return c.rwc.LocalAddr().String()
}

func (c *Conn) handshakeNoAuth() error {
	_, err := c.rwc.Write([]byte{verSocks5, authNoAuthenticationRequired})
	return err
}

func (c *Conn) handshake() error {
	var head header
	if _, err := head.ReadFrom(c.rwc); err != nil {
		return err
	}

	if err := c.handshakeNoAuth(); err != ErrAuthenticationFailed {
		return err
	}

	c.rwc.Write([]byte{verSocks5, authNoAcceptableMethods})

	return ErrAuthenticationFailed
}

func writeCommandErrorReply(c net.Conn, rep byte) error {
	_, err := c.Write([]byte{
		verSocks5,
		rep,
		rsvReserved,
		atypIPv4Address,
		0, 0, 0, 0,
		0, 0,
	})
	return err
}

func (c *Conn) commandConnect(cmd *cmd) error {
	dest := cmd.DestAddress()

	destPort := cmd.DestPort()
	fmt.Print(destPort)

	serverConn, err := net.Dial("tcp", dest)
	if err != nil {
		log.Printf("Erro ao conectar ao servidor real: %v", err)
		return err
	}
	defer serverConn.Close()

	r := &cmdResp{
		ver: verSocks5,
		rep: repSucceeded,
		rsv: rsvReserved,
	}

	host, port, err := net.SplitHostPort(serverConn.LocalAddr().String())
	if err != nil {
		writeCommandErrorReply(c.rwc, repGeneralSocksServerFailure)
		return err
	}

	ip := net.ParseIP(host)
	if ipv4 := ip.To4(); ipv4 != nil {
		r.atyp = atypIPv4Address
		r.bnd_addr = ipv4[:net.IPv4len]
	} else {
		r.atyp = atypIPv6Address
		r.bnd_addr = ip[:net.IPv6len]
	}

	prt, err := strconv.Atoi(port)
	if err != nil {
		writeCommandErrorReply(c.rwc, repGeneralSocksServerFailure)
		return err
	}
	r.bnd_port = uint16(prt)

	if _, err = r.WriteTo(c.rwc); err != nil {
		writeCommandErrorReply(c.rwc, repGeneralSocksServerFailure)
		return err
	}

	baseClient, err := NewBaseClient(serverConn, c.rwc, destPort)
	if err != nil {
		return err
	}

	return baseClient.HandleConnection()
}

func (c *Conn) Decrypt(data []byte, loginCrypt *LoginCrypt, gameCrypt *GameCrypt, port string) {
	toDecrypt := make([]byte, len(data)-2)
	copy(toDecrypt, data[2:])

	if port == "2106" {
		err := loginCrypt.Decrypt(toDecrypt, 0, len(toDecrypt))
		if err != nil {
			log.Printf("Erro decriptando dados do servidor para porta 2106: %v", err)
		}

		opcode := toDecrypt[0] & 0xff
		switch opcode {
		case 0x00:
			fmt.Println("RECEIVED PACKET Init")
		case 0x01:
			fmt.Println("RECEIVED PACKET LoginFail")
		case 0x02:
			fmt.Println("RECEIVED PACKET AccountKicked")
		case 0x03:
			fmt.Println("RECEIVED PACKET LoginOk")
		case 0x04:
			fmt.Println("RECEIVED PACKET ServerList")
		case 0x06:
			fmt.Println("RECEIVED PACKET PlayFail")
		case 0x07:
			fmt.Println("RECEIVED PACKET PlayOk")
		case 0x0b:
			fmt.Println("RECEIVED PACKET GGAuth")
		default:
			if len(data) > 2 {
				fmt.Printf("Unknown game packet received. [0x%x 0x%x] len=%d\n", opcode, data[1], len(data))
			}
		}

		return
	}

	gameCrypt.Decrypt(toDecrypt, 0, len(toDecrypt))

	opcode := toDecrypt[0] & 0xff
	switch opcode {
	case 0x00:
		key2 := toDecrypt[2:6]
		key := binary.LittleEndian.Uint32(key2)
		gameCrypt.SetKey(key)
		fmt.Println("RECEIVED PACKET INIT LOGIN")
	case 0x01:
		fmt.Println("RECEIVED PACKET MoveToLocation")
	case 0x02:
		fmt.Println("RECEIVED PACKET NpcSay")
	case 0x03:
		fmt.Println("RECEIVED PACKET CharInfo")
	case 0x04:
		fmt.Println("RECEIVED PACKET UserInfo")
	case 0x05:
		fmt.Println("RECEIVED PACKET Attack")
	case 0x06:
		fmt.Println("RECEIVED PACKET Die")
	case 0x07:
		fmt.Println("RECEIVED PACKET Revive")
	case 0x0B:
		fmt.Println("RECEIVED PACKET SpawnItem")
	case 0x0C:
		fmt.Println("RECEIVED PACKET DropItem")
	case 0x0D:
		fmt.Println("RECEIVED PACKET GetItem")
	case 0x0E:
		fmt.Println("RECEIVED PACKET StatusUpdate")
	case 0x0F:
		fmt.Println("RECEIVED PACKET NpcHtmlMessage")
	case 0x12:
		fmt.Println("RECEIVED PACKET DeleteObject")
	case 0x13:
		fmt.Println("RECEIVED PACKET CharSelectionInfo")
	case 0x15:
		fmt.Println("RECEIVED PACKET CharSelected")
	case 0x16:
		fmt.Println("RECEIVED PACKET NpcInfo")
	case 0x7F:
		fmt.Println("RECEIVED PACKET TutorialShowHtml")
	case 0xA1:
		fmt.Println("RECEIVED PACKET TutorialShowQuestionMark")
	case 0xA2:
		fmt.Println("RECEIVED PACKET TutorialEnableClientEvent")
	case 0xA3:
		fmt.Println("RECEIVED PACKET TutorialCloseHtml")
	case 0xA6:
		fmt.Println("RECEIVED PACKET MyTargetSelected")
	case 0xA7:
		fmt.Println("RECEIVED PACKET PartyMemberPosition")
	case 0xB6:
		fmt.Println("RECEIVED PACKET PetDelete")
	case 0xBA:
		fmt.Println("RECEIVED PACKET VehicleStarted")
	case 0xF8:
		fmt.Println("RECEIVED PACKET SSQInfo")
	case 0xE4:
		fmt.Println("RECEIVED PACKET HennaInfo")
	case 0x1B:
		fmt.Println("RECEIVED PACKET ItemList")
	case 0x1A:
		subOpcode := toDecrypt[1] & 0xff
		switch subOpcode {
		case 0x87:
			fmt.Println("RECEIVED PACKET UserInfo")
			reader := packets.NewReader(toDecrypt[2:])
			userInfo, err := serverpackets.ReadUserInfo(reader)
			fmt.Print(err, userInfo)
		default:
			fmt.Printf("Unknown extended game packet received. [0xfe 0x%x] len=%d\n", subOpcode, len(toDecrypt))
			fmt.Println("Game Data encrypted", data)
			fmt.Println("Game Decrypted Data", toDecrypt)
		}
	case 0x4A:
		fmt.Println("RECEIVED PACKET CreatureSay")
	case 0x39:
		fmt.Println("RECEIVED PACKET AskJoinParty")
	case 0x3A:
		fmt.Println("RECEIVED PACKET JoinParty")
	case 0x64:
		fmt.Println("RECEIVED PACKET SystemMessage")
	case 0x76:
		fmt.Println("RECEIVED PACKET MagicSkillLaunched")
	case 0x48:
		fmt.Println("RECEIVED PACKET MagicSkillUse")
	case 0x60:
		fmt.Println("RECEIVED PACKET MoveToPawn")
	case 0xCE:
		fmt.Println("RECEIVED PACKET RelationChanged")
	case 0x2D:
		fmt.Println("RECEIVED PACKET SocialAction")
	case 0x29:
		fmt.Println("RECEIVED PACKET TargetSelected")
	case 0x2A:
		fmt.Println("RECEIVED PACKET TargetUnselected")
	case 0x2B:
		fmt.Println("RECEIVED PACKET AutoAttackStart")
	case 0x2C:
		fmt.Println("RECEIVED PACKET AutoAttackStop")
	case 0xf3:
		fmt.Println("RECEIVED PACKET EtcStatusUpdate")
	case 0x61:
		fmt.Println("RECEIVED PACKET ValidateLocation")
	case 0x80:
		fmt.Println("RECEIVED PACKET QuestList")
	case 0xE7:
		fmt.Println("RECEIVED PACKET SendMacroList")
	case 0xfe: // Verifica pacotes com opcode extendido
		subOpcode := toDecrypt[1] & 0xff
		switch subOpcode {
		case 0x1B:
			fmt.Println("RECEIVED PACKET ExSendManorList")

		default:
			fmt.Printf("Unknown extended game packet received. [0xfe 0x%x] len=%d\n", subOpcode, len(toDecrypt))
			fmt.Println("Game Data encrypted", data)
			fmt.Println("Game Decrypted Data", toDecrypt)
		}
	default:
		if len(toDecrypt) > 2 {
			fmt.Printf("Unknown game packet received. [0x%x 0x%x] len=%d\n", opcode, toDecrypt[1], len(toDecrypt))
			fmt.Println("Game Data encrypted", data)
			fmt.Println("Game Decrypted Data", toDecrypt)
		}
	}

	return
}

func (c *Conn) command() error {
	var cmd cmd
	if _, err := cmd.ReadFrom(c.rwc); err != nil {
		if err == ErrAddressTypeNotSupported {
			writeCommandErrorReply(c.rwc, repAddressTypeNotSupported)
		}
		return err
	}

	switch cmd.cmd {
	case cmdConnect:
		return c.commandConnect(&cmd)
	default:
		return writeCommandErrorReply(c.rwc, repComandNotSupported)
	}
}

func (c *Conn) serve() {
	defer func() {
		if err := recover(); err != nil {
			const size = 4096
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			c.server.Logger.Printf("socks5: panic serving %v: %v\n%s", c.rwc.RemoteAddr(), err, buf)
		}
		c.close()
	}()

	if err := c.handshake(); err != nil {
		c.server.Logger.Printf("socks5: Conn.serve: Handshake failed: %v", err)
		return
	}

	if err := c.command(); err != nil {
		c.server.Logger.Printf("socks5: Conn.serve: command execution failed: %v", err)
		return
	}
}

func (c *Conn) close() {
	for _, h := range c.server.closeHandlers {
		h.HandleClose(c)
	}

	if c.rwc != nil {
		c.rwc.Close()
		c.rwc = nil
	}

	if c.externalIP != "" {
		// _ = c.r.Disconnected(c.externalIP)
	}
}

// func getCryptKey(port string) []byte {
// 	if port == "2106" {
// 		return []byte{
// 			0x5f, 0x3b, 0x35, 0x2e,
// 			0x5d, 0x39, 0x34, 0x2d,
// 			0x33, 0x31, 0x3d, 0x3d,
// 			0x2d, 0x25, 0x78, 0x54,
// 			0x21, 0x5e, 0x5b, 0x24, 0x00,
// 		}

// 	}

// 	return []byte{0x94, 0x35, 0x00, 0x00, 0xa1, 0x6c, 0x54, 0x87}
// }
