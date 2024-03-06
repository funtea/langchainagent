package basenode

import (
	"errors"
	"github.com/tidwall/gjson"
)

type Variable struct {
	Id   string       `json:"id"`
	Type string       `json:"type"`
	Data VariableData `json:"data"`
}

type VariableData struct {
	//NodeMeta NodeMeta          `json:"nodeMeta"`
	Inputs  VariableInputs    `json:"inputs"`
	Outputs []VariableOutputs `json:"outputs"`
}

type VariableInputs struct {
	Mode            string                    `json:"mode"`
	InputParameters []VariableInputParameters `json:"inputParameters"`
}

type VariableOutputs struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type VariableInputParameters struct {
	Name  string    `json:"name"`
	Input InputData `json:"input"`
}

type InputData struct {
	Type  string    `json:"type"`
	Value ValueData `json:"value"`
}

type ValueData struct {
	Type    string      `json:"type"`
	Content ContentData `json:"content"`
}

type ContentData struct {
	Source  string `json:"source"`
	BlockID string `json:"blockID"`
	Name    string `json:"name"`
}

/**
 *nodeMap        key节点id   value  节点
 *nodeOutputMap  key节点id   value  节点输出的变量值
 */
func NewVariableNode(node *Node, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (variable *Node, err error) {
	if node.Type != TypeVariableNode {
		return variable, errors.New("变量节点类型错误")
	}

	err = node.ParseVariableInput(nodeMap, nodeOutputMap)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (node *Node) ParseVariableInput(nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (err error) {
	if node.Data.Inputs.Mode == "set" {
		//向机器人设置变量
	}

	if len(node.Data.Inputs.InputParameters) == 0 {
		return errors.New("variable 输入变量不能为空")
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

	return nil
}

// Run方法
func (variable *Node) RunVariable(nodeOutputMap map[string]map[string]SchemaOutputs) map[string]map[string]SchemaOutputs {
	if len(variable.Data.Outputs) == 0 {
		return nodeOutputMap
	}
	for _, output := range variable.Data.Outputs {
		//todo 由于机器人还未搭建，这里设置默认值
		output.Value = true

		if _, ok := nodeOutputMap[variable.Id]; !ok {
			nodeOutputMap[variable.Id] = make(map[string]SchemaOutputs)
		}
		nodeOutputMap[variable.Id][output.Name] = output
	}
	return nodeOutputMap
}
