package util

type Connection struct {
	conn *net.TCPConn
}

type ConnectionPool struct {
	connections []*Connection
	mutex       sync.Mutex
}

func NewConnectionPool() *ConnectionPool {
	return *ConnectionPool{
		connections: make([]*Connection, 0),
	}
}

type Server struct {
	Listener *net.TCPListener
	Pool     *ConnectionPool
}

func NewServer(addr string) *Server {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	return &Server{
		Listener: listener,
		Pool:     NewConnectionPool(),
	}
}
