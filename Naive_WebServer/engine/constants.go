package engine

const (
	GET = iota
	POST
	PUT
	DELETE
	OPTION
)

const (
	HeaderKeyContentType = "Content-Type"
)

const (
	HeaderValueJson = "application/json"
	HeaderValueHTML = "text/html"
)
