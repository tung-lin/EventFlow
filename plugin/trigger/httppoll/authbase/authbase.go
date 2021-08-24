package authbase

type IAuthFactory interface {
	GetIdentifyName() string
	CreateAuth(config interface{}) IAuthPlugin
}

type IAuthPlugin interface {
	CreateHttpHeaders() []Header
	CheckParameter() error
}

type Header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
