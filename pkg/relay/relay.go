package relay

import (
	"log"
	"net"
	"slices"

	"github.com/orangeseeds/holepunching/pkg/utils"
)

type NATRecord struct {
	Private string
	Public  string
}

type MsgType = int

const (
	SYNC MsgType = 0
	CONN         = 1
	ACK          = 2
	ERR          = 3
)

type Message struct {
	Src  net.UDPAddr
	Dest net.UDPAddr
	Type MsgType
}

type Relay struct {
	Port    string
	Records []NATRecord
	Peers   []net.Addr
	Conn    net.PacketConn
}

func (r *Relay) ListenUDP() {
	conn, err := net.ListenPacket("udp", r.Port)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer conn.Close()
	r.Conn = conn

	for {
		buf := make([]byte, 512)
		i, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println(err)
			continue
		}

		r.AddUniquePeers(addr)
		// r.AddUniqueRecord(NATRecord{
		// 	Private: addr.String(),
		// 	Public:  ,
		// })

		var m Message
		err = utils.Decode(buf[0:i], &m)
		if err != nil {
			log.Println(err)
			continue
		}

		switch m.Type {
		case CONN, ACK:
			log.Println(m)
			destInPeers := slices.ContainsFunc(r.Peers, func(p net.Addr) bool {
				if p.String() == m.Dest.String() {
					return true
				}
				return false
			})
			if destInPeers {
				_, err := conn.WriteTo(buf, &m.Dest)
				if err != nil {
					log.Println(err)
					continue
				}
			} else {
				msg, err := utils.Encode(Message{
					Type: ERR,
				})
				if err != nil {
					log.Println(err)
					continue
				}
				conn.WriteTo(msg.Bytes(), addr)
			}
		default:
			log.Println("Invalid Message Type")
			continue
		}
	}
}

func (r *Relay) AddUniquePeers(addr net.Addr) {
	exists := slices.ContainsFunc(r.Peers, func(a net.Addr) bool {
		if a.String() == addr.String() {
			return true
		}
		return false
	})
	if !exists {
		r.Peers = append(r.Peers, addr)
	}
}

func (r *Relay) AddUniqueRecord(record NATRecord) {
	exists := slices.ContainsFunc(r.Records, func(n NATRecord) bool {
		if n.Public == record.Public || n.Private == record.Private {
			return true
		}
		return false
	})
	if !exists {
		r.Records = append(r.Records, record)
	}
}

// func (r *Relay) FloodTable() {
// 	for _, addr := range r.Peers {
// 		buf, err := utils.Encode(r.Records)
// 		if err != nil {
// 			log.Println(err)
// 			continue
// 		}
// 		r.Conn.WriteTo(buf.Bytes(), addr)
// 	}
// 	// }
// }
