package main

import (
	"net"
	"os"
	"syscall"
)

// netSocket is a file descriptor for a system socket.
type netSocket struct {
	fd int
}

func (ns netSocket) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	n, err := syscall.Read(ns.fd, p)
	if err != nil {
		n = 0
	}
	return n, err
}

func (ns netSocket) Write(p []byte) (int, error) {
	n, err := syscall.Write(ns.fd, p)
	if err != nil {
		n = 0
	}
	return n, err
}

// Creates a new netSocket for the next pending connection request.
func (ns *netSocket) Accept() (*netSocket, error) {
	// syscall.ForkLock doc states lock not needed for blocking accept.
	nfd, _, err := syscall.Accept(ns.fd)
	if err == nil {
		syscall.CloseOnExec(nfd)
	}
	if err != nil {
		return nil, err
	}
	return &netSocket{nfd}, nil
}

func (ns *netSocket) Close() error {
	return syscall.Close(ns.fd)
}

// Creates a new socket file descriptor, binds it and listens on it.
func newNetSocket(ip net.IP, port int) (*netSocket, error) {
	// ForkLock docs state that socket syscall requires the lock.
	syscall.ForkLock.Lock()
	// AF_INET = Address Family for IPv4
	// SOCK_STREAM = virtual circuit service
	// 0: the protocol for SOCK_STREAM, there's only 1.
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, os.NewSyscallError("socket", err)
	}
	syscall.ForkLock.Unlock()

	// Allow reuse of recently-used addresses.
	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		syscall.Close(fd)
		return nil, os.NewSyscallError("setsockopt", err)
	}

	// Bind the socket to a port
	sa := &syscall.SockaddrInet4{Port: port}
	copy(sa.Addr[:], ip)
	if err = syscall.Bind(fd, sa); err != nil {
		return nil, os.NewSyscallError("bind", err)
	}

	// Listen for incoming connections.
	if err = syscall.Listen(fd, syscall.SOMAXCONN); err != nil {
		return nil, os.NewSyscallError("listen", err)
	}

	return &netSocket{fd: fd}, nil
}

// type request struct {
// 	method string // GET, POST, etc.
// 	header textproto.MIMEHeader
// 	body   []byte
// 	uri    string // The raw URI from the request
// 	proto  string // "HTTP/1.1"
// }

// func parseRequest(c *netSocket) (*request, error) {
// 	b := bufio.NewReader(*c)
// 	tp := textproto.NewReader(b)
// 	req := new(request)

// 	// First line: parse "GET /index.html HTTP/1.0"
// 	var s string
// 	s, _ = tp.ReadLine()
// 	sp := strings.Split(s, " ")
// 	req.method, req.uri, req.proto = sp[0], sp[1], sp[2]

// 	// Parse headers
// 	mimeHeader, _ := tp.ReadMIMEHeader()
// 	req.header = mimeHeader

// 	// Parse body
// 	if req.method == "GET" || req.method == "HEAD" {
// 		return req, nil
// 	}
// 	if len(req.header["Content-Length"]) == 0 {
// 		return nil, errors.New("no content length")
// 	}
// 	length, err := strconv.Atoi(req.header["Content-Length"][0])
// 	if err != nil {
// 		return nil, err
// 	}
// 	body := make([]byte, length)
// 	if _, err = io.ReadFull(b, body); err != nil {
// 		return nil, err
// 	}
// 	req.body = body
// 	return req, nil
// }
