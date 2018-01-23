package throttlingtool

import "EventFlow/common/interface/pluginbase"

var policyFactories map[string]pluginbase.IThrottlingPolicyFactory

func init() {
	policyFactories = make(map[string]pluginbase.IThrottlingPolicyFactory)
}

func AddPolicyFactory(factory pluginbase.IThrottlingPolicyFactory) {
	if factory != nil {
		policyFactories[factory.GetIdentifyName()] = factory
	}
}

func CreatePolicy(mode string, config interface{}) pluginbase.IThrottlingPolicyPlugin {

	if policyFactory, existed := policyFactories[mode]; existed {
		return policyFactory.CreatePolicy(config)
	}

	return nil
}
