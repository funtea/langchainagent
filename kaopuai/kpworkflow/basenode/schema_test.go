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

	//查找变量节点
	var codeNode = Node{}
	for _, tmpNode := range nodeMap {
		if tmpNode.Type == TypeCodeNode {
			codeNode = tmpNode
		}
	}

	code, err := NewCodeNode(&codeNode, nodeMap, variableResult)
	if err != nil {
		return
	}
	fmt.Printf("code:%+v", code)

	codeResult := code.RunCode(variableResult)

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
