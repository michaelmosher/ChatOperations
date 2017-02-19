package operations

type Request struct {
	Id           int64
	Requester    string
	Server       string
	Action       string
	Responder    string
	Approved     bool
	Response_url string
}
