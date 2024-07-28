package husky

import "fmt"

type Code interface {
	Code() int64
	Msg() string
	I18n() string
	Meta() map[string]string
	WithMsg(format string, args ...any) Code
	WithMeta(key string, value string) Code
	error
}

type _Code struct {
	code int64
	msg  string
	i18n string
	meta map[string]string
}

func (c *_Code) Code() int64 {
	return c.code
}
func (c *_Code) Msg() string {
	return c.msg
}
func (c *_Code) I18n() string {
	return c.i18n
}
func (c *_Code) Meta() map[string]string {
	return c.meta
}
func (c *_Code) WithMsg(format string, args ...any) Code {
	return NewCode(c.code, fmt.Sprintf(format, args...), c.i18n, c.meta)
}
func (c *_Code) WithMeta(key string, value string) Code {
	meta := make(map[string]string)
	if c.meta != nil {
		for k, v := range c.meta {
			meta[k] = v
		}
	}
	meta[key] = value
	return NewCode(c.code, c.msg, c.i18n, meta)
}

func (c *_Code) Error() string {
	return fmt.Sprintf("Code(%d), Msg(%s), I18n(%s, %+v)", c.code, c.msg, c.i18n, c.meta)
}

func NewCode(code int64, msg string, i18n string, meta map[string]string) Code {
	return &_Code{
		code: code,
		msg:  msg,
		i18n: i18n,
		meta: meta,
	}
}

var (
	Ok                      = NewCode(0, "OK", "OK", nil)
	UnknownAbnormality      = NewCode(1, "未知异常", "UnknownAbnormality", nil)
	ServiceException        = NewCode(2, "服务异常", "ServiceException", nil)
	InProgress              = NewCode(3, "正在处理中", "InProgress", nil)
	ParameterError          = NewCode(4, "参数错误", "ParameterError", nil)
	AccessRestricted        = NewCode(5, "访问受限", "AccessRestricted", nil)
	AccessTooFrequent       = NewCode(6, "访问过于频繁", "AccessTooFrequent", nil)
	AccessResourceNotExists = NewCode(7, "访问资源不存在", "AccessResourceNotExists", nil)
)
