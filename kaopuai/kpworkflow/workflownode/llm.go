package workflownode

import (
	"context"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
	"strconv"
)

/**
 *nodeMap        key节点id   value  节点
 *nodeOutputMap  key节点id   value  节点输出的变量值
 */
func NewLLMNode(nodeId string, nodeMap map[string]Node) (nodeResult *Node, err error) {
	node := nodeMap[nodeId]
	if node.Type != TypeLLMNode {
		return nodeResult, errors.New("LLM节点类型错误")
	}

	return &node, nil
}

func (node *Node) ParseLLMInput(nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (err error) {
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
func (node *Node) RunLLM(ctx context.Context, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, resultJson string, err error) {
	err = node.ParseLLMInput(nodeMap, nodeOutputMap)
	if err != nil {
		return nodeOutputMap, "", err
	}

	//输入参数

	//配置参数
	var modelName,
		modelType, // 113 GPT3.5 (16k) \ 124对应gpt4(8k) \  133  gpt4 16k
		promptTemplate string

	var temperature float64
	for _, llmParam := range node.Data.Inputs.LlmParam {
		if llmParam.Name == "modelName" {
			modelName = llmParam.Input.Value.LiteralContent
		} else if llmParam.Name == "modelType" {
			modelType = llmParam.Input.Value.LiteralContent
		} else if llmParam.Name == "prompt" {
			promptTemplate = llmParam.Input.Value.LiteralContent
		} else if llmParam.Name == "temperature" {
			temperature, _ = strconv.ParseFloat(llmParam.Input.Value.LiteralContent, 64)
		}
	}
	fmt.Println(modelName)

	//拼接prompt
	promptTemplate, templateVariable, replaceVariable := getPromptVariable(promptTemplate, node)

	llmResultJson, err := runLLM(promptTemplate, modelType, temperature, templateVariable, replaceVariable)
	if err != nil {
		return nil, llmResultJson, err
	}

	var schemaOutputs = SchemaOutputs{
		Name:  "outputList",
		Value: llmResultJson,
	}

	if _, ok := nodeOutputMap[node.Id]; !ok {
		nodeOutputMap[node.Id] = make(map[string]SchemaOutputs)
	}
	nodeOutputMap[node.Id]["outputList"] = schemaOutputs

	return nodeOutputMap, llmResultJson, err
}

// 模版  llm节点参数
func getPromptVariable(promptTemplate string, llmNode *Node) (promptTemplateWithPre string, templateVariable []string, replaceVariable map[string]any) {
	replaceVariable = make(map[string]any)

	promptTemplateWithPre = "你是一个ai机器人，你必须以json格式返回内容。\n"
	promptTemplateWithPre += "prompt：\n"
	promptTemplateWithPre += promptTemplate
	promptTemplateWithPre += "\n输出格式为json，应该包含以下参数：\n"

	if len(llmNode.Data.Inputs.InputParameters) == 0 {
		return
	}

	for _, inputParameter := range llmNode.Data.Inputs.InputParameters {
		var newString string
		if inputParameter.Input.Value.Type == "ref" {
			newString = fmt.Sprintf("%s", inputParameter.Input.Value.Content.Value)
		} else {
			newString = fmt.Sprintf("%s", inputParameter.Input.Value.LiteralContent)
		}

		templateVariable = append(templateVariable, inputParameter.Name)
		replaceVariable[inputParameter.Name] = newString
	}

	//拼接规则
	for _, output := range llmNode.Data.Outputs {
		var outputType string
		if output.Type != "list" {
			outputType = output.Type
		} else {
			if output.ListSchema.Type == "string" {
				outputType = "array<string>"
			} else if output.ListSchema.Type == "integer" {
				outputType = "array<int64>"
			} else if output.ListSchema.Type == "boolean" {
				outputType = "array<boolean>"
			} else if output.ListSchema.Type == "float" {
				outputType = "array<float64>"
			} else {
				//全部没有匹配到
				fmt.Println("llm节点，全部没有匹配到")
				outputType = "array<string>"
			}
		}
		promptTemplateWithPre += fmt.Sprintf("%s:%s类型，%s。", output.Name, outputType, output.Description)
	}

	return
}

// modelType  113 GPT3.5 (16k) \ 124对应gpt4(8k) \  133  gpt 16k
func runLLM(promptTemplate, modelType string, temperature float64, templateVariable []string, replaceVariable map[string]any) (string, error) {
	// We can construct an LLMChain from a PromptTemplate and an LLM.
	var Options openai.Option
	if modelType == "113" {
		//gpt3
		Options = openai.WithModel("gpt-3.5-turbo")
	} else if modelType == "124" {
		Options = openai.WithModel("gpt-4")
	} else if modelType == "133" {
		Options = openai.WithModel("gpt-4-32k")
	} else {
		Options = openai.WithModel("gpt-4")
	}

	llm, err := openai.New(Options)
	if err != nil {
		return "", err
	}

	// If a chain only needs one input we can use the run function to execute chain.
	ctx := context.Background()

	translatePrompt := prompts.NewPromptTemplate(
		promptTemplate,
		templateVariable,
	)

	fmt.Println(translatePrompt.Template)
	llmChain := chains.NewLLMChain(llm, translatePrompt, chains.WithTemperature(temperature))

	// Otherwise the call function must be used.
	outputValues, err := chains.Call(ctx, llmChain, replaceVariable)
	if err != nil {
		return "", err
	}

	out, ok := outputValues[llmChain.OutputKey].(string)
	if !ok {
		return "", fmt.Errorf("invalid chain return")
	}
	fmt.Println(out)

	return out, nil
}
