package conditiontool

import (
	"fmt"
	"reflect"
)

//'or' logic between ConditionGroup
//'and' logic inside ConditionGroup

//Condition has one or more ConditionGroup
type Condition []struct {
	ConditionGroups []ConditionGroup `yaml:"conditiongroup"`
}

//ConditionGroup defines condition by key-value pair
type ConditionGroup struct {
	Metadata   string      `yaml:"metadata"`
	Value      interface{} `yaml:"value"`
	Expression string      `yaml:"expression"`
}

var floatType = reflect.TypeOf(float64(0))
var stringType = reflect.TypeOf(string(""))

//IsMatchCondition judges status by Condition and ConditionGroup
func IsMatchCondition(conditions *Condition, parameters *map[string]interface{}) (isMatched bool, err error) {

	for _, condition := range *conditions {

		matched := true

		for _, conditionGroup := range condition.ConditionGroups {
			if dataValue, existed := (*parameters)[conditionGroup.Metadata]; existed {

				var configValue = conditionGroup.Value
				var configValueCasted interface{}
				var dataValueCasted interface{}
				var configValueType reflect.Type
				var dataValueType reflect.Type
				var v reflect.Value

				if conditionGroup.Expression == "" {
					conditionGroup.Expression = "eq"
				}

				if configValue != nil {
					if reflect.TypeOf(configValue).Kind() == reflect.String && conditionGroup.Expression != "eq" && conditionGroup.Expression != "neq" {
						return false, fmt.Errorf("the condition(gt, gte, lt, lte) accepts only integer or float values")
					}

					v = reflect.ValueOf(configValue)
					v = reflect.Indirect(v)
					configValueType = v.Type()

				} else {
					if conditionGroup.Expression != "eq" && conditionGroup.Expression != "neq" {
						return false, fmt.Errorf("the condition(gt, gte, lt, lte) accepts only integer or float values")
					}
				}

				if dataValue != nil {
					dataValueType = reflect.TypeOf(dataValue)
				}

				if configValueType != nil && configValueType.ConvertibleTo(floatType) {
					configValueCasted = v.Convert(floatType).Float()

					if dataValue == nil {
						return false, fmt.Errorf("cannot compare %v(%v) to nil", configValue, configValueType)
					}

					v = reflect.ValueOf(dataValue)
					v = reflect.Indirect(v)

					if !v.Type().ConvertibleTo(floatType) {
						return false, fmt.Errorf("cannot compare %v(%v) to %v(%v)", configValue, configValueType, v, dataValueType.Kind())
					}

					dataValueCasted = v.Convert(floatType).Float()

				} else if configValueType == nil || configValueType.ConvertibleTo(stringType) {

					if configValueType != nil {
						configValueCasted = v.Convert(stringType).String()
					} else {
						configValueCasted = nil
					}

					if dataValue != nil {
						v = reflect.ValueOf(dataValue)
						v = reflect.Indirect(v)

						if !v.Type().ConvertibleTo(stringType) {
							return false, fmt.Errorf("cannot compare %v(%T) to %v(%v)", configValue, configValue, v, dataValueType.Kind())
						}

						dataValueCasted = v.Convert(stringType).String()
					} else {
						dataValueCasted = nil
					}
				}

				switch expression := conditionGroup.Expression; expression {
				case "eq":
					matched = configValueCasted == dataValueCasted
				case "neq":
					matched = configValueCasted != dataValueCasted
				case "gt":
					matched = dataValueCasted.(float64) > configValueCasted.(float64)
				case "gte":
					matched = dataValueCasted.(float64) >= configValueCasted.(float64)
				case "lt":
					matched = dataValueCasted.(float64) < configValueCasted.(float64)
				case "lte":
					matched = dataValueCasted.(float64) <= configValueCasted.(float64)
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
