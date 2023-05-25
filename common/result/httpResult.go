package result

const (
	Ok           = 200
	Error        = 500
	InvalidParam = 400

	AccountHasExist = 10001
	AccountNotExist = 10002
	WrongPassword   = 10003
)

var m map[int]string

func init() {
	m = make(map[int]string)
	m[Ok] = "ok"
	m[Error] = "服务器内部错误"
	m[InvalidParam] = "非法参数"
	m[AccountHasExist] = "账号已经存在"
	m[AccountNotExist] = "账号不存在"
	m[WrongPassword] = "密码错误"

}

func GetMsg(code int) string {
	return m[code]
}

type Result struct {
	Data interface{} `json:"data"`
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
}

type UserResult struct {
	Token    string   `json:"token,omitempty"`
	IpList   []string `json:"ipList,omitempty"`
	NickName string   `json:"nickName,omitempty"`
}

func Fail(code int, err error) *Result {
	result := &Result{
		Code: code,
		Msg:  GetMsg(code),
	}
	if err != nil {
		result.Data = err.Error()
	}
	return result
}

func Success(code int, data interface{}) *Result {
	return &Result{
		Code: code,
		Msg:  GetMsg(code),
		Data: data,
	}
}
