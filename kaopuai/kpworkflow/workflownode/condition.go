package workflownode

import (
	"context"
	"errors"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
)

/**
 *nodeMap        key节点id   value  节点
 *nodeOutputMap  key节点id   value  节点输出的变量值
 */
func NewConditionNode(nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (condition *Node, err error) {
	node := nodeMap[nodeId]
	if node.Type != TypeConditionNode {
		return nil, errors.New("condition节点类型错误")
	}

	err = (&node).ParseConditionInput(nodeMap, nodeOutputMap)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (condition *Node) ParseConditionInput(nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (err error) {

	return
}

func (condition *Node) RunCondition(ctx context.Context, nodeOutputMap map[string]map[string]SchemaOutputs, nodeMap map[string]Node) (isSuccess bool, resultJson string, err error) {
	if len(condition.Data.Inputs.Branches) == 0 {
		return false, "", errors.New("condition branches 为空")
	}

	if len(condition.Data.Inputs.Branches[0].Condition.Conditions) == 0 {
		return false, "", errors.New("condition Conditions 为空")
	}

	//目前只有一个branches
	logic := condition.Data.Inputs.Branches[0].Condition.Logic              //2and 1or
	conditionList := condition.Data.Inputs.Branches[0].Condition.Conditions //所有判断条件
	var conditionResult []bool
	for _, conditionEntity := range conditionList {
		tmpConditionResult := dealCondition(conditionEntity, nodeOutputMap, nodeMap)
		conditionResult = append(conditionResult, tmpConditionResult)
	}

	if logic == 2 {
		//and
		isSuccess = true
		for _, tmpConditionResult := range conditionResult {
			if tmpConditionResult == false {
				isSuccess = false
				break
			}
		}
	} else {
		//or
		isSuccess = false
		for _, tmpConditionResult := range conditionResult {
			if tmpConditionResult == true {
				isSuccess = true
				break
			}
		}
	}

	resultJson = condition.getResultJson(isSuccess)
	return
}

func (condition *Node) getResultJson(isSuccess bool) (resultJson string) {
	if isSuccess {
		resultJson = "{\"result\":\"pass to if branch\"}"
	} else {
		resultJson = "{\"result\":\"pass to else branch\"}"
	}
	return
}

// left:string integer boolean number(float)
// integer number拥有：equal\not equal\is empty\is not empty\greater than\greater than or equal\less than\less than or equal
// bool拥有:equal\not equal\is empty\is not empty\is true\is false
// string拥有:equal\not equal\longer than\longer than or equal\shorter than\shorter than or equal\contain\not contain\is empty\is not empty
// todo 还未处理code返回值
func dealCondition(condition SchemaConditions, nodeOutputMap map[string]map[string]SchemaOutputs, nodeMap map[string]Node) (isTrue bool) {
	var left, right any
	leftBlockID := condition.Left.Input.Value.Content.BlockID
	leftTmpNodeVariableName := condition.Left.Input.Value.Content.Name

	rightBlockID := condition.Right.Input.Value.Content.BlockID
	rightTmpNodeVariableName := condition.Right.Input.Value.Content.Name
	leftBlock, isLeftJson := getNodeOutput(leftBlockID, leftTmpNodeVariableName, nodeOutputMap, nodeMap)
	rightBlock, isRightJson := getNodeOutput(rightBlockID, rightTmpNodeVariableName, nodeOutputMap, nodeMap)

	var tmpLeftValue, tmpRightValue any
	tmpLeftValue = getAnyValueByName(isLeftJson, leftBlock.Value, leftTmpNodeVariableName)
	tmpRightValue = getAnyValueByName(isRightJson, rightBlock.Value, rightTmpNodeVariableName)

	if condition.Operator == 1 {
		//equal 是否和输入框内容相等
		if condition.Left.Input.Type != condition.Right.Input.Type {
			return false
		}

		left, right = getLeftRightEqual(condition, tmpLeftValue, tmpRightValue, isLeftJson, isRightJson)
		if left == right {
			return true
		}
	} else if condition.Operator == 2 {
		//not equal  不等
		if condition.Left.Input.Type != condition.Right.Input.Type {
			return true
		}

		left, right = getLeftRightEqual(condition, tmpLeftValue, tmpRightValue, isLeftJson, isRightJson)
		if left != right {
			return true
		}
	} else if condition.Operator == 3 {
		//longer than
		isTrue = dealLongerThan(condition, tmpLeftValue, tmpRightValue, condition.Operator)
		return isTrue
	} else if condition.Operator == 4 {
		//longer than or equal
		isTrue = dealLongerThan(condition, tmpLeftValue, tmpRightValue, condition.Operator)
		return isTrue
	} else if condition.Operator == 5 {
		//shorter than
		isTrue = dealLongerThan(condition, tmpLeftValue, tmpRightValue, condition.Operator)
		return isTrue
	} else if condition.Operator == 6 {
		//shorter than or equal
		isTrue = dealLongerThan(condition, tmpLeftValue, tmpRightValue, condition.Operator)
		return isTrue
	} else if condition.Operator == 7 {
		//contain
		isTrue = dealLongerThan(condition, tmpLeftValue, tmpRightValue, condition.Operator)
		return isTrue
	} else if condition.Operator == 8 {
		//not contain
		isTrue = dealLongerThan(condition, tmpLeftValue, tmpRightValue, condition.Operator)
		return isTrue
	} else if condition.Operator == 9 {
		//is empty
		isTrue = isEmpty(condition, tmpLeftValue, isLeftJson)
		return isTrue
	} else if condition.Operator == 10 {
		//is not empty
		isTrue = isEmpty(condition, tmpLeftValue, isLeftJson)
		return !isTrue
	} else if condition.Operator == 11 {
		//is true
		isTrue = dealIsTrue(condition, tmpLeftValue)
		return isTrue
	} else if condition.Operator == 12 {
		//is false
		isTrue = dealIsTrue(condition, tmpLeftValue)
		return !isTrue
	} else if condition.Operator == 13 {
		//greater than  大于
		if condition.Left.Input.Type != condition.Right.Input.Type {
			return false
		}

		isTrue = dealGreaterThan(condition, tmpLeftValue, tmpRightValue, condition.Operator)
		return isTrue
	} else if condition.Operator == 14 {
		//greater than or equal   大于等于
		//todo interger  float类型不一样，可以比较大于等于
		if condition.Left.Input.Type != condition.Right.Input.Type {
			return false
		}

		isTrue = dealGreaterThan(condition, tmpLeftValue, tmpRightValue, condition.Operator)
		return isTrue
	} else if condition.Operator == 15 {
		//less than   小于
		if condition.Left.Input.Type != condition.Right.Input.Type {
			return false
		}

		isTrue = dealGreaterThan(condition, tmpLeftValue, tmpRightValue, condition.Operator)
		return isTrue
	} else if condition.Operator == 16 {
		//less than or equal  小于等于
		if condition.Left.Input.Type != condition.Right.Input.Type {
			return false
		}

		isTrue = dealGreaterThan(condition, tmpLeftValue, tmpRightValue, condition.Operator)
		return isTrue
	}
	return false
}

func getNodeOutput(nodeId, nodeName string, nodeOutputMap map[string]map[string]SchemaOutputs, nodeMap map[string]Node) (SchemaOutputs, bool) {
	if nodeMap[nodeId].Type == TypeCodeNode ||
		nodeMap[nodeId].Type == TypeLLMNode ||
		nodeMap[nodeId].Type == TypeKnowledgeNode {
		return nodeOutputMap[nodeId]["outputList"], true
	} else {
		return nodeOutputMap[nodeId][nodeName], false
	}

}

func getLeftRightEqual(condition SchemaConditions, tmpLeftValue, tmpRightValue any, isLeftJson, isRightJson bool) (left, right any) {

	if condition.Left.Input.Value.Type == "ref" && condition.Right.Input.Value.Type == "literal" {
		//左侧使用引用， 右侧直接使用input框内容  nodeOutputMap[leftBlockID][leftTmpNodeVariableName]
		left, right = dealLeftRightLiteral(condition, tmpLeftValue, isLeftJson)
		//left = (nodeOutputMap[blockID][tmpNodeVariableName].Value).(bool)
	} else if condition.Left.Input.Value.Type == "ref" && condition.Right.Input.Value.Type == "ref" {
		//左引用  右侧也是引用
		left, right = dealLeftRightRef(condition, tmpLeftValue, tmpRightValue, isLeftJson, isRightJson)
	}
	return
}

func dealLeftRightLiteral(condition SchemaConditions, tmpLeftValue any, isLeftJson bool) (left, right any) {
	if isLeftJson {
		left = (tmpLeftValue).(string)
	} else if condition.Left.Input.Type == "boolean" {
		left = (tmpLeftValue).(bool)
	} else if condition.Left.Input.Type == "string" {
		left = (tmpLeftValue).(string)
	} else if condition.Left.Input.Type == "integer" {
		left = (tmpLeftValue).(int64)
	} else if condition.Left.Input.Type == "float" {
		left = (tmpLeftValue).(float64)
	}

	if condition.Right.Input.Type == "boolean" {
		rightValue := condition.Right.Input.Value.LiteralContent
		right, _ = strconv.ParseBool(rightValue)
	} else if condition.Right.Input.Type == "string" {
		rightValue := condition.Right.Input.Value.LiteralContent
		right = rightValue
	} else if condition.Right.Input.Type == "integer" {
		rightValue := condition.Right.Input.Value.LiteralContent
		right, _ = strconv.ParseInt(rightValue, 10, 64)
	} else if condition.Right.Input.Type == "float" {
		rightValue := condition.Right.Input.Value.LiteralContent
		right, _ = strconv.ParseFloat(rightValue, 64)
	}
	return
}

func dealLeftRightRef(condition SchemaConditions, tmpLeftValue, tmpRightValue any, isLeftJson, isRightJson bool) (left, right any) {
	if isLeftJson {
		left = (tmpLeftValue).(string)
	} else if condition.Left.Input.Type == "boolean" {
		left = (tmpLeftValue).(bool)
	} else if condition.Left.Input.Type == "string" {
		left = (tmpLeftValue).(string)
	} else if condition.Left.Input.Type == "integer" {
		left = (tmpLeftValue).(int64)
	} else if condition.Left.Input.Type == "float" {
		left = (tmpLeftValue).(float64)
	}

	if isRightJson {
		right = (tmpRightValue).(string)
	} else if condition.Right.Input.Type == "boolean" {
		right = (tmpRightValue).(bool)
	} else if condition.Right.Input.Type == "string" {
		right = (tmpRightValue).(string)
	} else if condition.Right.Input.Type == "integer" {
		right = (tmpRightValue).(int64)
	} else if condition.Right.Input.Type == "float" {
		right = (tmpRightValue).(float64)
	}
	return
}

func isEmpty(condition SchemaConditions, tmpLeftValue any, isLeftJson bool) bool {
	if tmpLeftValue == nil {
		//	如果是空  则为空
		return true
	}

	if tmpLeftValue == "" {
		//	如果是空字符串  则为空
		return true
	}

	if isLeftJson {
		left := (tmpLeftValue).(string)
		if len(left) == 0 {
			return true
		}
	} else if condition.Left.Input.Type == "boolean" {
		left := (tmpLeftValue).(bool)
		if left != true && left != false {
			return true
		}
	} else if condition.Left.Input.Type == "string" {
		left := (tmpLeftValue).(string)
		if len(left) == 0 {
			return true
		}
	} else if condition.Left.Input.Type == "integer" {

	} else if condition.Left.Input.Type == "float" {

	}

	return false
}

func dealGreaterThan(condition SchemaConditions, tmpLeftValue, tmpRightValue any, operator int64) (isTrue bool) {

	if condition.Left.Input.Value.Type == "ref" && condition.Right.Input.Value.Type == "literal" {
		//左侧使用引用， 右侧直接使用input框内容
		isTrue = dealLeftRightGreaterThanLiteral(condition, tmpLeftValue, operator)
	} else if condition.Left.Input.Value.Type == "ref" && condition.Right.Input.Value.Type == "ref" {
		//左引用  右侧也是引用
		isTrue = dealLeftRightGreaterThanRef(condition, tmpLeftValue, tmpRightValue, operator)
	}

	return isTrue
}

// integer float 才能使用greater than
func dealLeftRightGreaterThanLiteral(condition SchemaConditions, tmpLeftValue any, operator int64) bool {
	if condition.Left.Input.Type != "integer" && condition.Left.Input.Type != "float" {
		return false
	}

	if condition.Left.Input.Type == "integer" {
		left := (tmpLeftValue).(int64)

		rightValue := condition.Right.Input.Value.LiteralContent
		right, _ := strconv.ParseInt(rightValue, 10, 64)

		if operator == 13 {
			if left > right {
				return true
			}
		} else if operator == 14 {
			if left >= right {
				return true
			}
		} else if operator == 15 {
			if left < right {
				return true
			}
		} else if operator == 16 {
			if left <= right {
				return true
			}
		}

	} else if condition.Left.Input.Type == "float" {
		left := (tmpLeftValue).(float64)

		rightValue := condition.Right.Input.Value.LiteralContent
		right, _ := strconv.ParseFloat(rightValue, 64)
		if operator == 13 {
			if left > right {
				return true
			}
		} else if operator == 14 {
			if left >= right {
				return true
			}
		} else if operator == 15 {
			if left < right {
				return true
			}
		} else if operator == 16 {
			if left <= right {
				return true
			}
		}
	}
	return false
}

func dealLeftRightGreaterThanRef(condition SchemaConditions, tmpLeftValue, tmpRightValue any, operator int64) bool {
	if condition.Left.Input.Type == "integer" {
		left := (tmpLeftValue).(int64)
		right := (tmpRightValue).(int64)
		if operator == 13 {
			if left > right {
				return true
			}
		} else if operator == 14 {
			if left >= right {
				return true
			}
		} else if operator == 15 {
			if left < right {
				return true
			}
		} else if operator == 16 {
			if left <= right {
				return true
			}
		}
	} else if condition.Left.Input.Type == "float" {
		left := (tmpLeftValue).(float64)
		right := (tmpRightValue).(float64)
		if operator == 13 {
			if left > right {
				return true
			}
		} else if operator == 14 {
			if left >= right {
				return true
			}
		} else if operator == 15 {
			if left < right {
				return true
			}
		} else if operator == 16 {
			if left <= right {
				return true
			}
		}
	}
	return false
}

func dealIsTrue(condition SchemaConditions, tmpLeftValue any) bool {
	if tmpLeftValue == nil {
		//	如果是空  则为空
		return false
	}

	if tmpLeftValue == "" {
		//	如果是空字符串  则为空
		return false
	}

	if condition.Left.Input.Type == "boolean" {
		left := (tmpLeftValue).(bool)
		if left == true {
			return true
		}
	}

	return false
}

func dealLongerThan(condition SchemaConditions, tmpLeftValue, tmpRightValue any, operator int64) (isTrue bool) {

	if condition.Left.Input.Value.Type == "ref" && condition.Right.Input.Value.Type == "literal" {
		//左侧使用引用， 右侧直接使用input框内容
		isTrue = dealLeftRightLongerThanLiteral(condition, tmpLeftValue, operator)
	} else if condition.Left.Input.Value.Type == "ref" && condition.Right.Input.Value.Type == "ref" {
		//左引用  右侧也是引用
		isTrue = dealLeftRightLongerThanRef(condition, tmpLeftValue, tmpRightValue, operator)
	}

	return isTrue
}

// string 才能使用greater than
func dealLeftRightLongerThanLiteral(condition SchemaConditions, tmpLeftValue any, operator int64) bool {
	if condition.Left.Input.Type != "string" {
		return false
	}

	if condition.Left.Input.Type == "integer" {
		left := (tmpLeftValue).(string)

		right := condition.Right.Input.Value.LiteralContent

		if operator == 3 {
			if len(left) > len(right) {
				return true
			}
		} else if operator == 4 {
			if len(left) >= len(right) {
				return true
			}
		} else if operator == 5 {
			if len(left) < len(right) {
				return true
			}
		} else if operator == 6 {
			if len(left) <= len(right) {
				return true
			}
		} else if operator == 7 {
			isContain := strings.Contains(left, right)
			return isContain
		} else if operator == 8 {
			isContain := strings.Contains(left, right)
			return !isContain
		}

	}
	return false
}

func dealLeftRightLongerThanRef(condition SchemaConditions, tmpLeftValue, tmpRightValue any, operator int64) bool {
	if condition.Left.Input.Type == "string" {
		left := (tmpLeftValue).(string)
		right := (tmpRightValue).(string)
		if operator == 3 {
			if len(left) > len(right) {
				return true
			}
		} else if operator == 4 {
			if len(left) >= len(right) {
				return true
			}
		} else if operator == 5 {
			if len(left) < len(right) {
				return true
			}
		} else if operator == 6 {
			if len(left) <= len(right) {
				return true
			}
		} else if operator == 7 {
			isContain := strings.Contains(left, right)
			return isContain
		} else if operator == 8 {
			isContain := strings.Contains(left, right)
			return !isContain
		}
	}
	return false
}

func getAnyValueByName(isLeftJson bool, leftBlockValue any, leftTmpNodeVariableName string) (tmpLeftValue any) {
	if isLeftJson {
		scriptJsonAny := leftBlockValue
		scriptJson := scriptJsonAny.(string)
		codeGjson := gjson.Parse(scriptJson)
		tmpLeftValue = codeGjson.Get(leftTmpNodeVariableName).String()
	} else {
		tmpLeftValue = leftBlockValue
	}
	return
}
