package workflownode

import (
	"context"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
)

/**
 *nodeMap        key节点id   value  节点
 *nodeOutputMap  key节点id   value  节点输出的变量值
 */
func NewPluginsNode(nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (condition *Node, err error) {
	node := nodeMap[nodeId]
	if node.Type != TypePluginsNode {
		return nil, errors.New("code节点类型错误")
	}

	return &node, nil
}

func (node *Node) ParsePluginsInput(nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (err error) {

	//处理参数
	if len(node.Data.Inputs.InputParameters) == 0 {
		//参数可以为空
		return nil
	}

	for key, inputParameter := range node.Data.Inputs.InputParameters {
		if inputParameter.Input.Value.Type == "ref" {
			//引用类型
			nodeId := node.Data.Inputs.InputParameters[key].Input.Value.Content.BlockID
			varName := node.Data.Inputs.InputParameters[key].Input.Value.Content.Name
			refNode := nodeOutputMap[nodeId]

			if nodeMap[nodeId].Type == TypeCodeNode ||
				nodeMap[nodeId].Type == TypeLLMNode ||
				nodeMap[nodeId].Type == TypeKnowledgeNode ||
				nodeMap[nodeId].Type == TypePluginsNode {

				//解析
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

func (node *Node) RunPlugins(ctx context.Context, nodeOutputMap map[string]map[string]SchemaOutputs, nodeMap map[string]Node) (nodeOutputMapResult map[string]map[string]SchemaOutputs, resultJson string, err error) {
	nodeOutputMapResult = nodeOutputMap
	err = node.ParsePluginsInput(nodeMap, nodeOutputMap)
	if err != nil {
		return nodeOutputMapResult, "", err
	}

	apiID, apiName, pluginID, pluginName, pluginVersion, tips, outDocLink := parseApiParam(node.Data.Inputs.ApiParam)
	resultJson, err = dealPlugins(node, apiID, apiName, pluginID, pluginName, pluginVersion, tips, outDocLink)

	//整理输出
	var schemaOutputs = SchemaOutputs{
		Name:  "outputList",
		Value: resultJson,
	}

	if _, ok := nodeOutputMapResult[node.Id]; !ok {
		nodeOutputMapResult[node.Id] = make(map[string]SchemaOutputs)
	}
	nodeOutputMapResult[node.Id]["outputList"] = schemaOutputs
	fmt.Println(resultJson)

	return
}

func parseApiParam(ApiParamList []SchemaInputParameters) (apiID, apiName, pluginID, pluginName, pluginVersion, tips, outDocLink string) {
	if len(ApiParamList) == 0 {
		return
	}

	for _, apiParam := range ApiParamList {
		if apiParam.Name == "apiID" {
			apiID = apiParam.Input.Value.LiteralContent
		} else if apiParam.Name == "apiName" {
			apiName = apiParam.Input.Value.LiteralContent
		} else if apiParam.Name == "pluginID" {
			pluginID = apiParam.Input.Value.LiteralContent
		} else if apiParam.Name == "pluginName" {
			pluginName = apiParam.Input.Value.LiteralContent
		} else if apiParam.Name == "pluginVersion" {
			pluginVersion = apiParam.Input.Value.LiteralContent
		} else if apiParam.Name == "tips" {
			tips = apiParam.Input.Value.LiteralContent
		} else if apiParam.Name == "outDocLink" {
			outDocLink = apiParam.Input.Value.LiteralContent
		}
	}

	return
}

func dealPlugins(node *Node, apiID, apiName, pluginID, pluginName, pluginVersion, tips, outDocLink string) (resultJson string, err error) {

	resultJson = "{\"error\":\"\",\"request_id\":\"xxxxxx\",\"status\":200,\"result\":{\"data\":{\"base_currency_name\":\"cn\",\"status\":\"success\",\"updated_date\":\"20240307\",\"amount\":\"11.2\",\"base_currency_code\":\"1111\",\"rates\":{\"CNY\":{\"currency_name\":\"aaa\",\"rate\":\"bbb\",\"rate_for_amount\":\"ccc\"}}}}}"
	return
}
