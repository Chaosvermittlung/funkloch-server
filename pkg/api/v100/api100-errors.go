package api100

type ErrorResponse struct {
	Httpstatus   string `json:"httpstatus"`
	Errorcode    string `json:"errorcode"`
	Errormessage string `json:"errormessage"`
}

type APIErrorcode int

const (
	ERROR_WRONGCREDENTIALS APIErrorcode = 1 + iota
	ERROR_DBQUERYFAILED
	ERROR_MALFORMEDAUTH
	ERROR_NOHASH
	ERROR_NOTOKEN
	ERROR_JSONERROR
	ERROR_FILEERROR
	ERROR_FILEHASH
	ERROR_USERNOTAUTHORIZED
	ERROR_INVALIDPARAMETER
	ERROR_NOTFOUND
)

func (e *APIErrorcode) String() string {
	switch *e {
	case ERROR_WRONGCREDENTIALS:
		return "Wrong Username or Password"
	case ERROR_DBQUERYFAILED:
		return "Database Query failed"
	case ERROR_MALFORMEDAUTH:
		return "Authorization request is malformed"
	case ERROR_NOHASH:
		return "Could not generate Hash from Password"
	case ERROR_NOTOKEN:
		return "Could not generate Token"
	case ERROR_JSONERROR:
		return "JSON Marshal error"
	case ERROR_FILEERROR:
		return "file read/write error"
	case ERROR_FILEHASH:
		return "file md5-hash incorrect, file possibly corrupted"
	case ERROR_USERNOTAUTHORIZED:
		return "User not authorized"
	case ERROR_INVALIDPARAMETER:
		return "Invalid parameter"
	case ERROR_NOTFOUND:
		return "Resource not found"
	default:
		return "unknown error"
	}
}
