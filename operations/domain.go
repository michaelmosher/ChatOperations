package operations

type ActionRepository interface {
	FindById(id int) (Action, error)
	FindAll() ([]Action, error)
}

type ServerRepository interface {
	FindById(id int) (Server, error)
	FindAll() ([]Server, error)
}

type RequestRepository interface {
	Store(o Request) (int64, error)
	FindById(id int) (Request, error)
}

type Action struct {
	Id      int64
	Title   string
	Command string
}

type Server struct {
	Id	        int64
	Title       string
	Address     string
	Environment string
}

type Request struct {
	Id           int64
	Requester    string
	Server       Server
	Action       Action
	Responder    string
	Approved     bool
	Response_url string
}
