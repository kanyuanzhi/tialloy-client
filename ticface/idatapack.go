package ticface

type IDataPack interface {
	GetHeadLen() uint32
	Pack(message IMessage) ([]byte, error)
	Unpack(binaryData []byte) (IMessage, error)
}
