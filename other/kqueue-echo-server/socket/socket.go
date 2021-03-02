package socket

import (
	"fmt"
	"net"
	"syscall"
)

type Socket struct {
	FileDescriptor int
}

func (s Socket) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	nb, err := syscall.Read(s.FileDescriptor, b)
	if err != nil {
		nb = 0
	}

	return nb, err
}

func (s Socket) Write(b []byte) (int, error) {
	nb, err := syscall.Write(s.FileDescriptor, b)
	if err != nil {
		nb = 0
	}

	return nb, err
}

func (s Socket) Close() error {
	return syscall.Close(s.FileDescriptor)
}

func (s Socket) String() string {
	return fmt.Sprintf("%d", s.FileDescriptor)
}

func Listen(ip string, port int) (*Socket, error) {
	socket := &Socket{}
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to create socket: %v", err)
	}
	socket.FileDescriptor = fd

	socketAddress := &syscall.SockaddrInet4{Port: port}
	copy(socketAddress.Addr[:], net.ParseIP(ip))

	err = syscall.Bind(socket.FileDescriptor, socketAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to bind socket: %v", err)
	}

	err = syscall.Listen(socket.FileDescriptor, syscall.SOMAXCONN)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on socket %v", err)
	}

	return socket, nil
}
