package basenode

import (
	"fmt"
	"testing"
)

func TestSchema(t *testing.T) {
	var params = map[string]SchemaOutputs{
		"a": {
			Type:  "string",
			Name:  "a",
			Value: "345",
		},
		"b": {
			Type:  "integer",
			Name:  "b",
			Value: 789,
		},
		"z": {
			Type:  "boolean",
			Name:  "z",
			Value: true,
		},
		"x": {
			Type:  "float",
			Name:  "x",
			Value: 8.5,
		},
		"y1": {
			Type:  "string",
			Name:  "y1",
			Value: "nihao",
		},
		"y2": {
			Type:  "string",
			Name:  "y2",
			Value: "hello world",
		},
	}
	schemaTest, nodeMap, err := NewSchema("")
	if err != nil {
		return
	}

	fmt.Printf("schema 结构体:%+v \r\n", schemaTest)
	fmt.Printf("nodeMap 结构体:%+v \r\n", nodeMap)

	startNode, err := NewStartNode(nodeMap, params)
	if err != nil {
		return
	}
	startResult := startNode.RunStart()
	fmt.Printf("startNode 结构体:%+v \r\n", startNode)
	fmt.Printf("startResult 结构体:%+v \r\n", startResult)

	//查找knowledge节点
	var knowledgeNode = Node{}
	for _, tmpNode := range nodeMap {
		if tmpNode.Type == TypeKnowledgeNode {
			knowledgeNode = tmpNode
		}
	}
	knowledge, err := NewKnowledgeNode(&knowledgeNode, nodeMap, startResult)
	if err != nil {
		return
	}
	knowledgeResult, err := knowledge.RunKnowledge(nodeMap, startResult)
	if err != nil {
		return
	}
	fmt.Printf("%+v", knowledgeResult)

	//查找变量节点
	var variableNode = Node{}
	for _, tmpNode := range nodeMap {
		if tmpNode.Type == TypeVariableNode {
			variableNode = tmpNode
		}
	}
	variable, err := NewVariableNode(&variableNode, nodeMap, startResult)
	if err != nil {
		return
	}

	variableResult := variable.RunVariable(startResult)
	fmt.Printf("%+v", variableResult)

	//llm节点
	var llmNode = Node{}
	for _, tmpNode := range nodeMap {
		if tmpNode.Type == TypeLLMNode {
			llmNode = tmpNode
		}
	}
	llm, err := NewLLMNode(&llmNode, variableResult)
	if err != nil {
		return
	}
	llmResult, err := llm.RunLLM(nodeMap, variableResult)
	if err != nil {
		return
	}

	//condition节点
	var conditionNode = Node{}
	for _, tmpNode := range nodeMap {
		if tmpNode.Type == TypeConditionNode {
			conditionNode = tmpNode
		}
	}
	condition, err := NewConditionNode(&conditionNode, nodeMap, llmResult)
	if err != nil {
		return
	}
	isSuccess, _, err := condition.RunCondition(llmResult, nodeMap)
	if err != nil {
		return
	}
	fmt.Println(isSuccess)

	//查找code节点
	var codeNode = Node{}
	for _, tmpNode := range nodeMap {
		if tmpNode.Type == TypeCodeNode {
			codeNode = tmpNode
		}
	}

	code, err := NewCodeNode(&codeNode, nodeMap, llmResult)
	if err != nil {
		return
	}
	fmt.Printf("code:%+v", code)

	codeResult, err := code.RunCode(variableResult)

	var endNode = Node{}
	for _, tmpNode := range nodeMap {
		if tmpNode.Type == TypeEndNode {
			endNode = tmpNode
		}
	}

	end, err := NewEndNode(&endNode, nodeMap, codeResult)
	endResult, endContent, err := end.RunEnd()
	if err != nil {
		return
	}
	fmt.Printf("\r\n %+v", endResult)
	fmt.Printf("\r\n %+v", endContent)
	return
}
