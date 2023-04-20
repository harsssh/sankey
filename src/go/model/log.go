package model

type Request struct {
	Method string
	URI    string
}

type Log struct {
	UA		string
	Request Request
}

func (r Request) String() string {
	return r.Method + " " + r.URI
}
