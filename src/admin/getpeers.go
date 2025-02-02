package admin

import (
	"encoding/hex"
	"net"
	"sort"
	"time"

	"github.com/yggdrasil-network/yggdrasil-go/src/address"
)

type GetPeersRequest struct {
}

type GetPeersResponse struct {
	Peers []PeerEntry `json:"peers"`
}

type PeerEntry struct {
	URI           string        `json:"remote,omitempty"`
	Up            bool          `json:"up"`
	Inbound       bool          `json:"inbound"`
	IPAddress     string        `json:"address,omitempty"`
	PublicKey     string        `json:"key"`
	Port          uint64        `json:"port"`
	Priority      uint64        `json:"priority"`
	RXBytes       DataUnit      `json:"bytes_recvd,omitempty"`
	TXBytes       DataUnit      `json:"bytes_sent,omitempty"`
	Uptime        float64       `json:"uptime,omitempty"`
	LastError     string        `json:"last_error,omitempty"`
	LastErrorTime time.Duration `json:"last_error_time,omitempty"`
}

func (a *AdminSocket) getPeersHandler(req *GetPeersRequest, res *GetPeersResponse) error {
	peers := a.core.GetPeers()
	res.Peers = make([]PeerEntry, 0, len(peers))
	for _, p := range peers {
		peer := PeerEntry{
			Port:     p.Port,
			Up:       p.Up,
			Inbound:  p.Inbound,
			Priority: uint64(p.Priority), // can't be uint8 thanks to gobind
			URI:      p.URI,
			RXBytes:  DataUnit(p.RXBytes),
			TXBytes:  DataUnit(p.TXBytes),
			Uptime:   p.Uptime.Seconds(),
		}
		if addr := address.AddrForKey(p.Key); addr != nil {
			peer.PublicKey = hex.EncodeToString(p.Key)
			peer.IPAddress = net.IP(addr[:]).String()
		}
		if p.LastError != nil {
			peer.LastError = p.LastError.Error()
			peer.LastErrorTime = time.Since(p.LastErrorTime)
		}
		res.Peers = append(res.Peers, peer)
	}
	sort.Slice(res.Peers, func(i, j int) bool {
		if res.Peers[i].Port == res.Peers[j].Port {
			return res.Peers[i].Priority < res.Peers[j].Priority
		}
		return res.Peers[i].Port < res.Peers[j].Port
	})
	return nil
}
