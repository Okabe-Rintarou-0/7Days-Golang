package cash

type PeerClient interface {
	Get(key string) (ByteView, error)
	Put(key string, value ByteView) error
	Del(key string) (ByteView, error)
}

type PeerPicker interface {
	PickPeer(key string, namespace string) PeerClient
}
