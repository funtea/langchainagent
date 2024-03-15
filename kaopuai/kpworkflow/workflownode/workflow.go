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
func NewWorkflowNode(nodeId string, nodeMap map[string]Node) (workflow *Node, err error) {
	node := nodeMap[nodeId]
	if node.Type != TypeWorkflowNode {
		return nil, errors.New("workflow节点类型错误")
	}

	return &node, nil
}

func (node *Node) RunWorkflow(ctx context.Context, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, resultJson string, err error) {
	nodeOutputMapResult = nodeOutputMap

	err = node.ParseWorkflowInput(nodeMap, nodeOutputMap)
	if err != nil {
		return nodeOutputMapResult, "", err
	}

	params, err := getEdgeParams(node)
	if err != nil {
		return nodeOutputMapResult, "", err
	}

	schemaJson, err := getSchemaJson(node)
	if err != nil {
		return nil, "", err
	}

	dealWorkFlowResult, err := RunEdges(ctx, schemaJson, params)
	if err != nil {
		return nodeOutputMapResult, "", err
	}

	var schemaOutputs = SchemaOutputs{
		Name:  "output",
		Value: dealWorkFlowResult.EndNode.TestResult.AnswerContent + dealWorkFlowResult.EndNode.TestResult.OutputVariable,
	}

	if _, ok := nodeOutputMap[node.Id]; !ok {
		nodeOutputMap[node.Id] = make(map[string]SchemaOutputs)
	}
	nodeOutputMap[node.Id]["output"] = schemaOutputs
	return
}

func (node *Node) ParseWorkflowInput(nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (err error) {
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
				nodeMap[nodeId].Type == TypeWorkflowNode {

				//if inputParameter.Input.Type == "object" || inputParameter.Input.Type == "list" {
				//	return errors.New(fmt.Sprintf("参数%s 不能为object或者list", varName))
				//}

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

func getEdgeParams(node *Node) (map[string]SchemaOutputs, error) {
	params := make(map[string]SchemaOutputs)
	if len(node.Data.Inputs.InputDefs) == 0 {
		//允许空参数
		return params, nil
	}

	for _, inputDef := range node.Data.Inputs.InputDefs {
		var tmpValue any
		for _, inputParameter := range node.Data.Inputs.InputParameters {
			if inputDef.Name == inputParameter.Name {
				tmpValue = inputParameter.Input.Value.Content.Value
			}
		}

		var tmpSchemaOutputs = SchemaOutputs{
			Name:        inputDef.Name,
			Type:        inputDef.Type,
			Required:    inputDef.Required,
			Description: inputDef.Description,
			Value:       tmpValue,
		}

		params[inputDef.Name] = tmpSchemaOutputs
	}

	return params, nil
}

func getSchemaJson(node *Node) (schemaJson string, err error) {
	var workflowId = node.Data.Inputs.WorkflowId
	var spaceId = node.Data.Inputs.SpaceId
	var inputsType = node.Data.Inputs.Type

	//todo workflowId,spaceId,inputsType 通过mysql获取schemaJson
	fmt.Printf("%+v,%+v,%+v", workflowId, spaceId, inputsType)
	schemaJson = "{\"nodes\":[{\"id\":\"100001\",\"type\":\"1\",\"meta\":{\"position\":{\"x\":192,\"y\":0}},\"data\":{\"outputs\":[{\"type\":\"string\",\"name\":\"a\",\"required\":true,\"description\":\"参数a\"},{\"type\":\"integer\",\"name\":\"b\",\"required\":true,\"description\":\"参数b\"},{\"type\":\"boolean\",\"name\":\"z\",\"required\":false,\"description\":\"zzz\"},{\"type\":\"float\",\"name\":\"x\",\"required\":true,\"description\":\"xxx\"},{\"type\":\"string\",\"name\":\"y1\",\"required\":true},{\"type\":\"string\",\"name\":\"y2\",\"required\":true}],\"nodeMeta\":{\"title\":\"Start\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Start.png\",\"description\":\"The starting node of the workflow, used to set the information needed to initiate the workflow.\",\"subTitle\":\"\"}}},{\"id\":\"900001\",\"type\":\"2\",\"meta\":{\"position\":{\"x\":4343,\"y\":878}},\"data\":{\"nodeMeta\":{\"title\":\"End\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-End.png\",\"description\":\"The final node of the workflow, used to return the result information after the workflow runs.\",\"subTitle\":\"\"},\"inputs\":{\"terminatePlan\":\"useAnswerContent\",\"inputParameters\":[{\"name\":\"c\",\"input\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"list\",\"name\":\"key4511\",\"listSchema\":{\"type\":\"integer\"}}],\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"120710\",\"name\":\"key4.0.key45.key451\"}}}}],\"content\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"我的content\"}}}}},{\"id\":\"117411\",\"type\":\"11\",\"meta\":{\"position\":{\"x\":1362,\"y\":922}},\"data\":{\"nodeMeta\":{\"title\":\"Variable\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Variable.png\",\"description\":\"Used for reading and writing variables in your bot. The variable name must match the variable name in Bot.\",\"subTitle\":\"Variable\"},\"inputs\":{\"mode\":\"set\",\"inputParameters\":[{\"name\":\"botVariable\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"bbb\"}}}]},\"outputs\":[{\"type\":\"boolean\",\"name\":\"isSuccess\"}]}},{\"id\":\"120710\",\"type\":\"5\",\"meta\":{\"position\":{\"x\":3719.884210526316,\"y\":-78.4263157894737}},\"data\":{\"nodeMeta\":{\"title\":\"Code\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Code.png\",\"description\":\"Write code to process input variables to generate return values.\",\"subTitle\":\"Code\"},\"inputs\":{\"inputParameters\":[{\"name\":\"input1\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}},{\"name\":\"input2\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"3232\"}}},{\"name\":\"input3\",\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}},{\"name\":\"input4\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}],\"code\":\"function main( params ){\\n    const ret = {\\n        \\\"key0\\\": params.input1 + params.input2,\\n        \\\"key1\\\": \\\"hi\\\",\\n        \\\"key2\\\": [\\\"hello\\\", \\\"world\\\"],\\n        \\\"key3\\\": {\\n            \\\"key31\\\": \\\"hi\\\"\\n        },\\n        \\\"key4\\\": [{\\n            \\\"key41\\\": true,\\n            \\\"key42\\\": 1,\\n            \\\"key43\\\": 12.88,\\n            \\\"key44\\\": [\\\"hello\\\"],\\n            \\\"key45\\\": {\\n                \\\"key451\\\": {\\n                    \\\"key4511\\\":[9,4,2]\\n                }\\n            }\\n        }]\\n    };\\n\\n    return ret;\\n}\",\"language\":5},\"outputs\":[{\"type\":\"string\",\"name\":\"key0\"},{\"type\":\"string\",\"name\":\"key1\"},{\"type\":\"list\",\"name\":\"key2\",\"listSchema\":{\"type\":\"string\"}},{\"type\":\"object\",\"name\":\"key3\",\"objectSchema\":[{\"type\":\"string\",\"name\":\"key31\"}]},{\"type\":\"list\",\"name\":\"key4\",\"listSchema\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"boolean\",\"name\":\"key41\"},{\"type\":\"integer\",\"name\":\"key42\"},{\"type\":\"float\",\"name\":\"key43\"},{\"type\":\"list\",\"name\":\"key44\",\"listSchema\":{\"type\":\"string\"}},{\"type\":\"object\",\"name\":\"key45\",\"objectSchema\":[{\"type\":\"object\",\"name\":\"key451\",\"objectSchema\":[{\"type\":\"list\",\"name\":\"key4511\",\"listSchema\":{\"type\":\"integer\"}}]}]}]}}]}},{\"id\":\"171236\",\"type\":\"8\",\"meta\":{\"position\":{\"x\":2992,\"y\":751}},\"data\":{\"nodeMeta\":{\"title\":\"Condition\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Condition.png\",\"description\":\"Connect two downstream branches. If the set conditions are met, run only the ‘if’ branch; otherwise, run only the ‘else’ branch.\",\"subTitle\":\"Condition\"},\"inputs\":{\"branches\":[{\"condition\":{\"logic\":2,\"conditions\":[{\"operator\":1,\"left\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}},\"right\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"literal\",\"literalContent\":\"false\"}}}},{\"operator\":2,\"left\":{\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}},\"right\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}}},{\"operator\":9,\"left\":{\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}},{\"operator\":10,\"left\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"z\"}}}}},{\"operator\":13,\"left\":{\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"x\"}}}},\"right\":{\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}}},{\"operator\":11,\"left\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}}},{\"operator\":12,\"left\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}}},{\"operator\":14,\"left\":{\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"x\"}}}},\"right\":{\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}},{\"operator\":15,\"left\":{\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"x\"}}}},\"right\":{\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}},{\"operator\":16,\"left\":{\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"x\"}}}},\"right\":{\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}},{\"operator\":10,\"left\":{\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"181217\",\"name\":\"outputList.output\"}}}}}]}}]}}},{\"id\":\"196739\",\"type\":\"3\",\"meta\":{\"position\":{\"x\":2111.5,\"y\":624.011385199241}},\"data\":{\"nodeMeta\":{\"title\":\"LLM\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-LLM.png\",\"description\":\"Invoke the large language model, generate responses using variables and prompt words.\",\"subTitle\":\"LLM\"},\"inputs\":{\"inputParameters\":[{\"name\":\"input1\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"a\"}}},{\"name\":\"input2\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"1\"}}},{\"name\":\"input3\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"false\"}}},{\"name\":\"input4\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"1.5\"}}}],\"llmParam\":[{\"name\":\"modleName\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"GPT-3.5 (16K)\"}}},{\"name\":\"modelType\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"literal\",\"literalContent\":\"113\"}}},{\"name\":\"prompt\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"input1：{{.input1}}\\ninput2：{{.input2}}\\ninput3：{{.input3}}\\ninput4：{{.input4}}\"}}},{\"name\":\"temperature\",\"input\":{\"type\":\"float\",\"value\":{\"type\":\"literal\",\"literalContent\":\"0.7\"}}}]},\"outputs\":[{\"type\":\"string\",\"name\":\"output1\",\"description\":\"是input1的值\"},{\"type\":\"integer\",\"name\":\"output2\",\"description\":\"是input2的值\"},{\"type\":\"boolean\",\"name\":\"output3\",\"description\":\"是input3的值\"},{\"type\":\"float\",\"name\":\"output4\",\"description\":\"是input4的值\"},{\"type\":\"list\",\"name\":\"output5\",\"listSchema\":{\"type\":\"string\"},\"description\":\"是input1的集合\"},{\"type\":\"list\",\"name\":\"output6\",\"listSchema\":{\"type\":\"integer\"},\"description\":\"是input2的集合\"},{\"type\":\"list\",\"name\":\"output7\",\"listSchema\":{\"type\":\"boolean\"},\"description\":\"是input3的集合\"},{\"type\":\"list\",\"name\":\"output8\",\"listSchema\":{\"type\":\"float\"},\"description\":\"是input4的集合\"}],\"version\":\"2\"}},{\"id\":\"181217\",\"type\":\"6\",\"meta\":{\"position\":{\"x\":726,\"y\":755.1052631578947}},\"data\":{\"nodeMeta\":{\"title\":\"Knowledge\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Knowledge.png\",\"description\":\"In the selected knowledge, the best matching information is recalled based on the input variable and returned as an Array.\",\"subTitle\":\"Knowledge\"},\"outputs\":[{\"type\":\"list\",\"name\":\"outputList\",\"listSchema\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"string\",\"name\":\"output\"}]}}],\"inputs\":{\"inputParameters\":[{\"name\":\"Query\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}}],\"datasetParam\":[{\"name\":\"datasetList\",\"input\":{\"type\":\"list\",\"listSchema\":{\"type\":\"string\"},\"value\":{\"type\":\"literal\",\"stringArrayContent\":[\"7338252484214980609\"]}}},{\"name\":\"topK\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"literal\",\"literalContent\":\"3\"}}},{\"name\":\"minScore\",\"input\":{\"type\":\"number\",\"value\":{\"type\":\"literal\",\"literalContent\":\"0.5\"}}},{\"name\":\"strategy\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"literal\",\"literalContent\":\"1\"}}}]}}}],\"edges\":[{\"sourceNodeID\":\"120710\",\"targetNodeID\":\"900001\"},{\"sourceNodeID\":\"171236\",\"targetNodeID\":\"120710\",\"sourcePortID\":\"false\"},{\"sourceNodeID\":\"171236\",\"targetNodeID\":\"900001\",\"sourcePortID\":\"true\"},{\"sourceNodeID\":\"196739\",\"targetNodeID\":\"171236\"},{\"sourceNodeID\":\"100001\",\"targetNodeID\":\"181217\"},{\"sourceNodeID\":\"181217\",\"targetNodeID\":\"117411\"},{\"sourceNodeID\":\"117411\",\"targetNodeID\":\"196739\"}]}"

	return schemaJson, nil
}
