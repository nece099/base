package app

const (
	// public errors
	EC_SUCCESS = "0"
	EC_FAILED  = "-1"

	EC_PARAMS_ERROR                   = "100001"
	EC_DB_ERROR                       = "100002"
	EC_INVALID_REQUEST                = "100003"
	EC_INVALID_APPID                  = "100004"
	EC_INCORRECT_USERNAME_OR_PASSWORD = "100005"
	EC_INVALID_OP                     = "100006"

	// auth server
	EC_USER_NOT_EXISTS            = "A00001"
	EC_UNSUPPORT_LOGIN_TYPE       = "A00002"
	EC_SESSION_NOT_EXISTS         = "A00003"
	EC_INVALID_LOGIN_TOKEN        = "A00004"
	EC_TOKEN_EXPIRED              = "A00005"
	EC_ROLE_NOT_EXISTS            = "A00006"
	EC_PRIVILEGE_GROUP_NOT_EXISTS = "A00007"
	EC_PRIVILEGE_NOT_EXISTS       = "A00008"
	EC_POLICY_RULE                = "A00009"
	EC_POLICY_TYPE                = "A00010"
	EC_USER_ROLE_NOT_EXISTS       = "A00011"
)

var ErrMsgMap = map[string]string{
	// public errors
	EC_SUCCESS: "成功",
	EC_FAILED:  "失败",

	EC_PARAMS_ERROR:                   "参数错误 : %v",
	EC_DB_ERROR:                       "数据错误 : %v",
	EC_INVALID_REQUEST:                "无效请求",
	EC_INVALID_APPID:                  "无效AppID: %v",
	EC_INCORRECT_USERNAME_OR_PASSWORD: "用户名或密码错误",
	EC_INVALID_OP:                     "无效操作",

	// auth server
	EC_USER_NOT_EXISTS:            "用户不存在: %v",
	EC_UNSUPPORT_LOGIN_TYPE:       "不支持登陆类型: %v",
	EC_SESSION_NOT_EXISTS:         "SESSION不存在: %v",
	EC_INVALID_LOGIN_TOKEN:        "无效登录TOKEN: %v",
	EC_TOKEN_EXPIRED:              "TOKEN已过期",
	EC_ROLE_NOT_EXISTS:            "角色不存在: %v",
	EC_PRIVILEGE_GROUP_NOT_EXISTS: "权限组不存在: %v",
	EC_PRIVILEGE_NOT_EXISTS:       "权限不存在: %v",
	EC_POLICY_RULE:                "策略规则错误: %v",
	EC_POLICY_TYPE:                "策略类型错误: %v",
	EC_USER_ROLE_NOT_EXISTS:       "用户角色不存在: %v",
}
