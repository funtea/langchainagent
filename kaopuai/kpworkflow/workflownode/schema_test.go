package workflownode

import (
	"context"
	"fmt"
	"testing"
)

func TestSchema(t *testing.T) {
	ctx := context.Background()
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
	startResult, _ := startNode.RunStart(ctx)
	fmt.Printf("startNode 结构体:%+v \r\n", startNode)
	fmt.Printf("startResult 结构体:%+v \r\n", startResult)

	//查找knowledge节点
	var knowledgeNode = Node{}
	for _, tmpNode := range nodeMap {
		if tmpNode.Type == TypeKnowledgeNode {
			knowledgeNode = tmpNode
		}
	}
	knowledge, err := NewKnowledgeNode(knowledgeNode.Id, nodeMap)
	if err != nil {
		return
	}
	knowledgeResult, err := knowledge.RunKnowledge(ctx, nodeMap, startResult)
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
	variable, err := NewVariableNode(variableNode.Id, nodeMap, startResult)
	if err != nil {
		return
	}

	variableResult, _, _ := variable.RunVariable(ctx, startResult)
	fmt.Printf("%+v", variableResult)

	//llm节点
	var llmNode = Node{}
	for _, tmpNode := range nodeMap {
		if tmpNode.Type == TypeLLMNode {
			llmNode = tmpNode
		}
	}
	llm, err := NewLLMNode(llmNode.Id, nodeMap)
	if err != nil {
		return
	}
	llmResult, _, err := llm.RunLLM(ctx, nodeMap, variableResult)
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
	condition, err := NewConditionNode(conditionNode.Id, nodeMap, llmResult)
	if err != nil {
		return
	}
	isSuccess, _, err := condition.RunCondition(ctx, llmResult, nodeMap)
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

	code, err := NewCodeNode(codeNode.Id, nodeMap, llmResult)
	if err != nil {
		return
	}
	fmt.Printf("code:%+v", code)

	codeResult, _, err := code.RunCode(ctx, variableResult)

	var endNode = Node{}
	for _, tmpNode := range nodeMap {
		if tmpNode.Type == TypeEndNode {
			endNode = tmpNode
		}
	}

	end, err := NewEndNode(endNode.Id, nodeMap, codeResult)
	endResult, endContent, err := end.RunEnd(ctx)
	if err != nil {
		return
	}
	fmt.Printf("\r\n %+v", endResult)
	fmt.Printf("\r\n %+v", endContent)
	return
}
