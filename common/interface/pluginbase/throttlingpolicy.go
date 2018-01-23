package pluginbase

type IThrottlingPolicyFactory interface {
	IsDefaultFactory() bool
	GetIdentifyName() string
	CreatePolicy(config interface{}) IThrottlingPolicy
}

type IThrottlingPolicy interface {
	FireAction(throttlingIdFromTrigger string) bool
}

type Policy struct {
	Mode    string
	Setting interface{}
}
