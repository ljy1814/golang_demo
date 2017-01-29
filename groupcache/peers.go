package groupcache

import "demo/groupcache/groupcachepb"

//定义进程如何互相发现以及交流

type Context interface{}

type ProtoGetter interface {
	Get(context Context, in *groupcachepb.GetRequest, out *groupcachepb.GetResponse) error
}

type PeerPicker interface {
	//返回nil, false则表示是自身,正常返回其它
	PickPeer(key string) (peer ProtoGetter, ok bool)
}

//实现用于不查找peer的PeerPicker
type NoPeers struct{}

func (NoPeers) PickPeer(key string) (peer ProtoGetter, ok bool) { return }

var (
	portPicker func(groupName string) PeerPicker
)

func RegisterPeerPicker(fn func() PeerPicker) {
	if portPicker != nil {
		panic("RegisterPeerPicker called more than once")
	}
	portPicker = func(_ string) PeerPicker { return fn() }
}

func RegisterPerGroupPeerPicker(fn func(groupName string) PeerPicker) {
	if portPicker != nil {
		panic("RegisterPeerPicker called more than once")
	}
	portPicker = fn
}

func getPeers(groupName string) PeerPicker {
	if portPicker == nil {
		return NoPeers{}
	}
	pk := portPicker(groupName)
	if pk == nil {
		pk = NoPeers{}
	}
	return pk
}
