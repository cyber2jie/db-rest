package api

const (
	CodeSuccess = 10001
	CodeError   = 10002
)
const (
	pathPrefix        = "/api/internal"
	pathTokenGet      = pathPrefix + "/token"
	pathWorkSpaceList = pathPrefix + "/workspace/list"
)

var (
	white_list = []string{
		pathTokenGet,
		pathTokenGet + "/",
		pathWorkSpaceList,
		pathWorkSpaceList + "/",
	}
)

type Result struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResult(code int, msg string) Result {
	return Result{
		Code: code,
		Msg:  msg,
	}
}
func NewResultWithData(code int, msg string, data interface{}) Result {
	return Result{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}
