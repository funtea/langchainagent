package basenode

import (
	"context"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
)

/**
 *nodeMap        key节点id   value  节点
 *nodeOutputMap  key节点id   value  节点输出的变量值
 */
func NewKnowledgeNode(nodeId string, nodeMap map[string]Node) (variable *Node, err error) {
	node := nodeMap[nodeId]
	if node.Type != TypeKnowledgeNode {
		return variable, errors.New("Knowledge 节点类型错误")
	}

	return &node, nil
}

func (node *Node) ParseKnowledgeInput(nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (err error) {
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
func (node *Node) RunKnowledge(ctx context.Context, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, err error) {
	err = node.ParseKnowledgeInput(nodeMap, nodeOutputMap)
	if err != nil {
		return nodeOutputMap, err
	}

	var datasetList []string
	var topK, strategy int64
	var minScore float64
	for _, dataset := range node.Data.Inputs.DatasetParam {
		if dataset.Name == "datasetList" {
			datasetList = dataset.Input.Value.StringArrayContent
		} else if dataset.Name == "topK" {
			topK, err = strconv.ParseInt(dataset.Input.Value.LiteralContent, 10, 64)
			if err != nil {
				return
			}
		} else if dataset.Name == "minScore" {
			minScore, _ = strconv.ParseFloat(dataset.Input.Value.LiteralContent, 64)
			if err != nil {
				return
			}
		} else if dataset.Name == "strategy" {
			strategy, _ = strconv.ParseInt(dataset.Input.Value.LiteralContent, 10, 64)
			if err != nil {
				return
			}
		}
	}

	fmt.Printf("%+v \n", datasetList)
	fmt.Printf("%+v \n", topK)
	fmt.Printf("%+v \n", strategy)
	fmt.Printf("%+v \n", minScore)

	var knowledgeJson string = `{"outputList":[{output:"aaa"},{output:"bbb"}]}`
	outputMap := make(map[string]SchemaOutputs)
	var output = SchemaOutputs{
		Type:  "string",
		Name:  "outputList",
		Value: knowledgeJson,
	}
	outputMap["outputList"] = output
	if _, ok := nodeOutputMap[node.Id]; !ok {
		nodeOutputMap[node.Id] = make(map[string]SchemaOutputs)
	}
	nodeOutputMap[node.Id] = outputMap

	return nodeOutputMap, err
}
