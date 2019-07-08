package pipelinetool

import (
	"fmt"
)

//'or' logic between ConditionGroup
//'and' logic inside ConditionGroup

//Condition has one or more ConditionGroup
type Condition []struct {
	ConditionGroups []ConditionGroup `yaml:"conditiongroup"`
}

//ConditionGroup defines condition by key-value pair
type ConditionGroup struct {
	Metadata      string      `yaml:"metadata"`
	Value         interface{} `yaml:"value"`
	ThresholdType string      `yaml:"thresholdtype"`
}

//IsMatchCondition judges status by Condition and ConditionGroup
func IsMatchCondition(conditions *Condition, parameters *map[string]interface{}) (isMatched bool, err error) {

	for _, condition := range *conditions {

		matched := true

		for _, conditionGroup := range condition.ConditionGroups {
			if value, existed := (*parameters)[conditionGroup.Metadata]; existed {

				var conditionValue interface{}

				switch value.(type) {
				case float64:
					v, ok := conditionGroup.Value.(float64)
					if !ok {
						return false, fmt.Errorf(fmt.Sprintf("cannot cast %v from %T to float64", conditionGroup.Value, conditionGroup.Value))
					}
					conditionValue = v
				case string:
					if conditionGroup.ThresholdType != "eq" && conditionGroup.ThresholdType != "neq" {
						return false, fmt.Errorf("the condition(gt, gte, lt, lte) accepts only integer or float values")
					}

					v, ok := conditionGroup.Value.(string)
					if !ok {
						return false, fmt.Errorf(fmt.Sprintf("cannot cast %v from %T to string", conditionGroup.Value, conditionGroup.Value))
					}
					conditionValue = v
				}

				if conditionGroup.ThresholdType == "" {
					conditionGroup.ThresholdType = "eq"
				}

				switch thresholdType := conditionGroup.ThresholdType; thresholdType {
				case "eq":
					matched = value == conditionValue
				case "neq":
					matched = value != conditionValue
				}

				if !matched {
					break
				}

			} else {
				matched = false
				break
			}
		}

		if matched {
			return true, nil
		}
	}

	return false, nil
}
