package basenode

import (
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"strings"
)

/**
 *nodeMap        key节点id   value  节点
 *nodeOutputMap  key节点id   value  节点输出的变量值
 */
func NewEndNode(node *Node, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (variable *Node, err error) {
	if node.Type != TypeEndNode {
		return variable, errors.New("变量节点类型错误")
	}

	err = node.ParseEndInput(nodeMap, nodeOutputMap)
	if err != nil {
		return nil, err
	}
	return node, nil
}

/**
 *解析end节点的输入参数
 */
func (node *Node) ParseEndInput(nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (err error) {
	if len(node.Data.Inputs.InputParameters) == 0 {
		return errors.New("end节点 输入变量不能为空")
	}

	for key, inputParameter := range node.Data.Inputs.InputParameters {
		if inputParameter.Input.Value.Type == "ref" {
			//引用类型
			nodeId := node.Data.Inputs.InputParameters[key].Input.Value.Content.BlockID
			varName := node.Data.Inputs.InputParameters[key].Input.Value.Content.Name
			refNode := nodeOutputMap[nodeId]

			if nodeMap[nodeId].Type == TypeCodeNode ||
				nodeMap[nodeId].Type == TypeLLMNode ||
				nodeMap[nodeId].Type == TypeKnowledgeNode {
				//如果是code节点,  只取第一层数据
				scriptJsonAny := refNode["outputList"].Value
				scriptJson := scriptJsonAny.(string)
				codeGjson := gjson.Parse(scriptJson)
				tmpValue := codeGjson.Get(varName).String()
				node.Data.Inputs.InputParameters[key].Input.Value.Content.Value = tmpValue
			} else {
				//todo 检查变量名是否存在  varName
				node.Data.Inputs.InputParameters[key].Input.Value.Content.Value = refNode[varName].Value
			}
		} else if inputParameter.Input.Value.Type == "literal" {
			//直接使用节点设置的类型值,  自身就是值
			//node.Data.Inputs.InputParameters[0].Input.Value.LiteralContent
		}
	}

	return
}

// Run方法
func (end *Node) RunEnd() (outputVariable, answerContent string, err error) {
	if len(end.Data.Inputs.InputParameters) == 0 {
		return
	}

	var tmpCodeVariable string

	tmpCodeVariable = "{"
	for _, inputParam := range end.Data.Inputs.InputParameters {
		tmpCodeVariable += `"` + inputParam.Name + `":`
		if inputParam.Input.Type == "string" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += `"` + fmt.Sprintf("%s", inputParam.Input.Value.Content.Value) + `",`
			} else {
				tmpCodeVariable += `"` + fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `",`
			}
		} else if inputParam.Input.Type == "integer" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += fmt.Sprintf("%d", inputParam.Input.Value.Content.Value) + `,`
			} else {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `,`
			}
		} else if inputParam.Input.Type == "boolean" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += fmt.Sprintf("%t", inputParam.Input.Value.Content.Value) + `,`
			} else {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `,`
			}
		} else if inputParam.Input.Type == "float" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += fmt.Sprintf("%f", inputParam.Input.Value.Content.Value) + `,`
			} else {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `,`
			}
		} else if inputParam.Input.Type == "object" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.Content.Value) + `,`
			} else {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `,`
			}
		}

	}
	tmpCodeVariable = strings.TrimRight(tmpCodeVariable, ",")
	tmpCodeVariable += "}"

	if end.Data.Inputs.TerminatePlan == "useAnswerContent" {
		answerContent = end.Data.Inputs.Content.Value.LiteralContent
	}

	outputVariable = tmpCodeVariable
	return
}
