package ticnet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/kanyuanzhi/tialloy-client/global"
	"github.com/kanyuanzhi/tialloy-client/ticface"
)

type DataPack struct {
}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	//DataLen uint32(4字节) + MsgID uint32(4字节)
	return 8
}

func (d *DataPack) Pack(message ticface.IMessage) ([]byte, error) {
	dataBuf := bytes.NewBuffer([]byte{})

	if err := binary.Write(dataBuf, binary.BigEndian, message.GetDataLen()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuf, binary.BigEndian, message.GetMsgID()); err != nil {
		return nil, err
	}

	if err := binary.Write(dataBuf, binary.BigEndian, message.GetData()); err != nil {
		return nil, err
	}

	return dataBuf.Bytes(), nil
}

func (d *DataPack) Unpack(binaryData []byte) (ticface.IMessage, error) {
	dataBuf := bytes.NewBuffer(binaryData)

	message := &Message{} //只解压head的信息，得到dataLen和msgID

	if err := binary.Read(dataBuf, binary.BigEndian, &message.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuf, binary.BigEndian, &message.MsgID); err != nil {
		return nil, err
	}

	if global.Object.TcpMaxPacketSize > 0 && message.DataLen > global.Object.TcpMaxPacketSize {
		return nil, errors.New("received data size is larger than max packet size")
	}

	return message, nil
}
