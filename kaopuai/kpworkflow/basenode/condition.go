package basenode

import (
	"errors"
	"github.com/tidwall/gjson"
)

/**
 *nodeMap        key节点id   value  节点
 *nodeOutputMap  key节点id   value  节点输出的变量值
 */
func NewConditionNode(node *Node, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (code *Node, err error) {
	if node.Type != TypeConditionNode {
		return nil, errors.New("condition节点类型错误")
	}

	err = node.ParseConditionInput(nodeMap, nodeOutputMap)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (condition *Node) ParseConditionInput(nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (err error) {
	if len(condition.Data.Inputs.InputParameters) == 0 {
		return errors.New("condition节点 输入变量不能为空")
	}

	for key, inputParameter := range condition.Data.Inputs.InputParameters {
		if inputParameter.Input.Value.Type == "ref" {
			//引用类型
			nodeId := condition.Data.Inputs.InputParameters[key].Input.Value.Content.BlockID
			varName := condition.Data.Inputs.InputParameters[key].Input.Value.Content.Name
			refNode := nodeOutputMap[nodeId]

			if nodeMap[nodeId].Type == TypeCodeNode {
				//如果是code节点,  只取第一层数据
				scriptJsonAny := refNode["scriptResultJson"].Value
				scriptJson := scriptJsonAny.(string)
				codeGjson := gjson.Parse(scriptJson)
				tmpValue := codeGjson.Get(varName).String()
				condition.Data.Inputs.InputParameters[key].Input.Value.Content.Value = tmpValue
			} else {
				//todo 检查变量名是否存在  varName
				condition.Data.Inputs.InputParameters[key].Input.Value.Content.Value = refNode[varName].Value
			}
		} else if inputParameter.Input.Value.Type == "literal" {
			//直接使用节点设置的类型值,  自身就是值
			//node.Data.Inputs.InputParameters[0].Input.Value.LiteralContent
		}
	}

	return
}

func (condition *Node) RunCondition(nodeOutputMap map[string]map[string]SchemaOutputs) map[string]map[string]SchemaOutputs {

	return nodeOutputMap
}
