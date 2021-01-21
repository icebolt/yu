package node

import (
	"bytes"
	"context"
	"crypto/rand"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"io"
	"os"
	"yu/config"
	. "yu/node"
	"yu/storage/kv"
)

var MasterNodeInfoKey = []byte("master-node-info")

type MasterNode struct {
	info   *MasterNodeInfo
	metadb kv.KV
}

func NewMasterNode(cfg *config.Conf) (*MasterNode, error) {
	metadb, err := kv.NewKV(&cfg.NodeDB)
	if err != nil {
		return nil, err
	}
	p2pHost, err := makeP2pHost(&cfg.NodeConf)
	if err != nil {
		return nil, err
	}
	data, err := metadb.Get(MasterNodeInfoKey)
	if err != nil {
		return nil, err
	}
	var info *MasterNodeInfo
	if data == nil {
		info = &MasterNodeInfo{
			Name:        cfg.NodeConf.NodeName,
			WorkerNodes: cfg.NodeConf.WorkerNodes,
		}
	} else {
		info, err = DecodeMasterNodeInfo(data)
		if err != nil {
			return nil, err
		}
	}

	info.P2pID = p2pHost.ID().String()

	return &MasterNode{
		info,
		metadb,
	}, nil
}

func makeP2pHost(cfg *config.NodeConf) (host.Host, error) {
	r, err := loadNodeKeyReader(cfg)
	if err != nil {
		return nil, err
	}
	priv, _, err := crypto.GenerateKeyPairWithReader(cfg.NodeKeyType, cfg.NodeKeyBits, r)
	if err != nil {
		return nil, err
	}

	return libp2p.New(
		context.Background(),
		libp2p.Identity(priv),
		libp2p.ListenAddrStrings(cfg.P2pAddrs...),
	)
}

func loadNodeKeyReader(cfg *config.NodeConf) (io.Reader, error) {
	if cfg.NodeKey != "" {
		return bytes.NewBufferString(cfg.NodeKey), nil
	}
	if cfg.NodeKeyFile != "" {
		return os.Open(cfg.NodeKeyFile)
	}
	return rand.Reader, nil
}