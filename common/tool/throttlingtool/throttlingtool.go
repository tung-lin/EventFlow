package throttlingtool

import "EventFlow/common/interface/pluginbase"

var defaultPolicyFactory pluginbase.IThrottlingPolicyFactory
var policyFactories map[string]pluginbase.IThrottlingPolicyFactory
var policyRules map[string]pluginbase.IThrottlingPolicy

func init() {
	policyFactories = make(map[string]pluginbase.IThrottlingPolicyFactory)
	policyRules = make(map[string]pluginbase.IThrottlingPolicy)
}

func CreatePolicy(mode string, config interface{}) pluginbase.IThrottlingPolicy {

	policyFactory, existed := policyFactories[mode]

	if !existed {
		policyFactory = defaultPolicyFactory
	}

	if policyFactory != nil {
		return policyFactory.CreatePolicy(config)
	}

	return nil
}

func AddPolicyFactory(factory pluginbase.IThrottlingPolicyFactory) {
	if factory != nil {
		policyFactories[factory.GetIdentifyName()] = factory

		if factory.IsDefaultFactory() {
			defaultPolicyFactory = factory
		}
	}
}

func RemovePolicy(throttlingIdFromTrigger string) {
	delete(policyRules, throttlingIdFromTrigger)
}
