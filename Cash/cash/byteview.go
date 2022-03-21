package cash

type ByteView struct {
	bytes []byte
}

func (b ByteView) Len() int {
	return len(b.bytes)
}

func (b ByteView) String() string {
	return string(b.bytes)
}

func (b ByteView) Clone() []byte {
	newBytes := make([]byte, b.Len())
	copy(newBytes, b.bytes)
	return newBytes
}
