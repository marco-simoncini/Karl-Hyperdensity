package sidecar

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

type SocketTransport struct {
	socketPath string
	conn       net.Conn
	reader     *bufio.Reader
}

func NewSocketTransport(socketPath string) *SocketTransport {
	return &SocketTransport{socketPath: socketPath}
}

func (t *SocketTransport) Connect() error {
	conn, err := net.DialTimeout("unix", t.socketPath, 2*time.Second)
	if err != nil {
		return err
	}
	t.conn = conn
	t.reader = bufio.NewReader(conn)
	_, err = t.reader.ReadBytes('\n')
	return err
}

func (t *SocketTransport) Close() error {
	if t.conn == nil {
		return nil
	}
	return t.conn.Close()
}

func (t *SocketTransport) Execute(command string) (map[string]any, error) {
	if t.conn == nil {
		return nil, fmt.Errorf("qmp transport not connected")
	}
	payload := map[string]any{"execute": command}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	if _, err := t.conn.Write(append(encoded, '\n')); err != nil {
		return nil, err
	}
	line, err := t.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	var decoded map[string]any
	if err := json.Unmarshal(line, &decoded); err != nil {
		return nil, err
	}
	if errValue, exists := decoded["error"]; exists {
		return nil, fmt.Errorf("qmp command failed: %v", errValue)
	}
	if returnValue, exists := decoded["return"].(map[string]any); exists {
		return returnValue, nil
	}
	return decoded, nil
}
