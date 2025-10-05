package tcpclient

import (
	"aiplag-agent/common/config"
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

func SendCommand(cmd byte, payload string) (byte, error) {
	address, err := config.UsedTCPAddress()
	if err != nil {
		return 0, err
	}
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return 0, fmt.Errorf("connect: %w", err)
	}
	defer conn.Close()

	// Build message
	data := append([]byte{cmd}, []byte(payload)...)
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint16(len(data)))
	buf.Write(data)

	// Send
	if _, err := conn.Write(buf.Bytes()); err != nil {
		return 0, fmt.Errorf("write: %w", err)
	}

	// Read response
	respHeader := make([]byte, 2)
	if _, err := conn.Read(respHeader); err != nil {
		return 0, fmt.Errorf("failed to read response header: %v", err)
	}
	respLength := binary.BigEndian.Uint16(respHeader)
	if respLength != 1 {
		return 0, fmt.Errorf("unexpected response length: %d", respLength)
	}

	respPayload := make([]byte, respLength)
	if _, err := conn.Read(respPayload); err != nil {
		return 0, fmt.Errorf("failed to read response payload: %v", err)
	}

	return respPayload[0], nil
}
