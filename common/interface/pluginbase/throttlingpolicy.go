package pluginbase

type IThrottlingPolicyFactory interface {
	GetIdentifyName() string
	CreatePolicy(config interface{}) IThrottlingPolicyPlugin
}

type IThrottlingPolicyPlugin interface {
	Throttling(throttlingIdFromTrigger string) bool
}

type Policy struct {
	Mode    string
	Setting interface{}
}
