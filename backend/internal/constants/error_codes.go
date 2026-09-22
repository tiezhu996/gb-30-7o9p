package constants

const (
	CodeSuccess         = 0
	CodeBadRequest      = 40000
	CodeUnauthorized    = 40100
	CodeForbidden       = 40300
	CodeNotFound        = 40400
	CodeConflict        = 40900
	CodeRateLimited     = 42900
	CodeValidationError = 42200
	CodeInternalError   = 50000
)

const (
	MsgInvalidParam  = "invalid request parameter"
	MsgUnauthorized  = "authentication required"
	MsgForbidden     = "permission denied"
	MsgNotFound      = "resource not found"
	MsgConflict      = "resource conflict"
	MsgRateLimited   = "too many requests"
	MsgInternalError = "internal server error"
)
