package server

type (
	Server struct{}
)

// H.264 High Profile and VP9 (profile 0)

func New() *Server {
	//TODO implement create new server
	return &Server{}
}

func (s *Server) Start() {
	//TODO start new server

}

func (s *Server) Close() {
	//TODO Stop accepting new connections and Close all connections

}
