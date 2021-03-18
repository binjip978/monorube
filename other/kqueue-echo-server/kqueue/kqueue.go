package kqueue

import (
	"fmt"
	"kq-echo/socket"
	"syscall"
)

type EventLoop struct {
	KQueueFD int
	SocketFD int
}

type Handler = func(s *socket.Socket)

func NewEventLoop(s *socket.Socket) (*EventLoop, error) {
	kQueue, err := syscall.Kqueue()
	if err != nil {
		return nil, fmt.Errorf("failed to create kqueue fd: %v", err)
	}

	kEvent := syscall.Kevent_t{
		Ident:  uint64(s.FileDescriptor),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD | syscall.EV_ENABLE,
		Fflags: 0,
		Data:   0,
		Udata:  nil,
	}

	kEventRegistered, err := syscall.Kevent(kQueue, []syscall.Kevent_t{kEvent}, nil, nil)
	if err != nil || kEventRegistered == -1 {
		return nil, fmt.Errorf("failed to register change event: %v", err)
	}

	return &EventLoop{
		KQueueFD: kQueue,
		SocketFD: s.FileDescriptor,
	}, nil
}

func (el *EventLoop) Handle(handler Handler) {
	for {
		newEvents := make([]syscall.Kevent_t, 10)
		n, err := syscall.Kevent(el.KQueueFD, nil, newEvents, nil)
		if err != nil {
			continue
		}

		for i := 0; i < n; i++ {
			currEvent := newEvents[i]
			eFD := int(currEvent.Ident)

			if currEvent.Flags&syscall.EV_EOF != 0 {
				syscall.Close(eFD)
			} else if eFD == el.SocketFD {
				// new connection
				sConn, _, err := syscall.Accept(eFD)
				if err != nil {
					continue
				}

				sEvent := syscall.Kevent_t{
					Ident:  uint64(sConn),
					Filter: syscall.EVFILT_READ,
					Flags:  syscall.EV_ADD,
					Fflags: 0,
					Data:   0,
					Udata:  nil,
				}

				sEventReg, err := syscall.Kevent(el.KQueueFD, []syscall.Kevent_t{sEvent}, nil, nil)
				if err != nil || sEventReg == -1 {
					continue
				}
			} else if currEvent.Filter&syscall.EVFILT_READ != 0 {
				go handler(&socket.Socket{FileDescriptor: eFD})
			}
		}
	}

}
