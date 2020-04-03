package server

type (
	Server struct{}
)

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
