package basenode

import (
	"context"
	"errors"
	"fmt"
)

func RunEdges(ctx context.Context, schemaJson string, params map[string]SchemaOutputs) (edgesResult string, err error) {
	schema, nodeMap, err := NewSchema(schemaJson)
	if err != nil {
		return
	}

	edgeList := schema.Edges
	if len(edgeList) == 0 {
		return "", errors.New("无路径节点")
	}

	//获取开始，结束点
	startNodeId, _, err := FindFirstAndEndNode(schema.Nodes)
	if err != nil {
		return "", err
	}

	//获取节点路径map
	edgeSourceNodeIdMap := GetEdgeMap(edgeList)

	nodeList := schema.Nodes
	fmt.Printf("%+v", nodeList)

	edgesResult, err = DealWorkFlow(ctx, startNodeId, edgeSourceNodeIdMap, nodeMap, params)

	return
}

func GetEdgeMap(edgeList []Edge) map[string]Edge {
	var edgeMap = make(map[string]Edge)
	for _, edge := range edgeList {
		if edge.SourcePortID == "true" {
			edgeMap[edge.SourceNodeID+"true"] = edge
		} else if edge.SourcePortID == "false" {
			edgeMap[edge.SourceNodeID+"false"] = edge
		} else {
			edgeMap[edge.SourceNodeID] = edge
		}

	}

	return edgeMap
}

func FindFirstAndEndNode(nodeList []Node) (startNodeId, endNodeId string, err error) {
	if len(nodeList) == 0 {
		return "", "", errors.New("无节点")
	}

	for _, node := range nodeList {
		if node.Type == TypeStartNode {
			startNodeId = node.Id
		}
		if node.Type == TypeEndNode {
			endNodeId = node.Id
		}
	}
	return
}

// 处理工作流主函数
func DealWorkFlow(ctx context.Context, startNodeId string, edgeSourceNodeIdMap map[string]Edge, nodeMap map[string]Node, params map[string]SchemaOutputs) (edgesResult string, err error) {
	var i int64 = 0
	//初始化节点id，默认开始节点
	var thisNodeId = startNodeId
	var thisNode Node
	var nodeOutputMap = make(map[string]map[string]SchemaOutputs)

	var outputVariable, answerContent string
	for {
		thisNode = nodeMap[thisNodeId]

		if thisNode.Type == TypeStartNode {
			//开始节点
			nodeOutputMap, err = runStart(ctx, nodeMap, params)
			if err != nil {
				return "", err
			}
		}

		if thisNode.Type == TypeKnowledgeNode {
			//知识点节点
			nodeOutputMap, err = runKnowledge(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return "", err
			}
		}

		if thisNode.Type == TypeWorkflowNode {
			//工作流节点
			nodeOutputMap, err = runWorkflow(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return "", err
			}
		}

		if thisNode.Type == TypeVariableNode {
			//变量节点
			nodeOutputMap, err = runVariable(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return "", err
			}
		}

		if thisNode.Type == TypeLLMNode {
			//llm节点
			nodeOutputMap, err = runLLMNode(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return "", err
			}
		}

		if thisNode.Type == TypeConditionNode {
			//condition节点
			isSuccess, err := runConditionNode(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return "", err
			}

			//todo nextNodeId
			if isSuccess == true {
				thisNodeId = edgeSourceNodeIdMap[thisNodeId+"true"].TargetNodeID
			} else {
				thisNodeId = edgeSourceNodeIdMap[thisNodeId+"false"].TargetNodeID
			}
			continue
		}

		if thisNode.Type == TypeCodeNode {
			//code节点
			nodeOutputMap, err = runCodeNode(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return "", err
			}
		}

		if thisNode.Type == TypePluginsNode {
			//变量节点
			nodeOutputMap, err = runPluginsNode(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return "", err
			}
		}

		//赋值下一轮节点id
		if thisNode.Type == TypeEndNode {
			outputVariable, answerContent, err = runEndNode(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return "", err
			}
			break
		}

		thisNodeId = edgeSourceNodeIdMap[thisNodeId].TargetNodeID
		i++
	}

	fmt.Println(outputVariable, answerContent)

	edgesResult = answerContent + outputVariable
	return
}

func runStart(ctx context.Context, nodeMap map[string]Node, params map[string]SchemaOutputs) (nodeOutputMap map[string]map[string]SchemaOutputs, err error) {
	startNode, err := NewStartNode(nodeMap, params)
	if err != nil {
		return
	}
	nodeOutputMap = startNode.RunStart(ctx)
	fmt.Printf("startNode 结构体:%+v \r\n", startNode)
	fmt.Printf("startResult 结构体:%+v \r\n", nodeOutputMap)
	return
}

func runKnowledge(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, err error) {
	nodeOutputMapResult = nodeOutputMap
	knowledge, err := NewKnowledgeNode(nodeId, nodeMap)
	if err != nil {
		return
	}
	nodeOutputMapResult, err = knowledge.RunKnowledge(ctx, nodeMap, nodeOutputMap)
	if err != nil {
		return
	}
	fmt.Printf("knowledge %+v", nodeOutputMapResult)
	return
}

func runVariable(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, err error) {
	nodeOutputMapResult = nodeOutputMap
	variable, err := NewVariableNode(nodeId, nodeMap, nodeOutputMapResult)
	if err != nil {
		return
	}

	nodeOutputMapResult = variable.RunVariable(ctx, nodeOutputMapResult)
	fmt.Printf("%+v", nodeOutputMapResult)
	return
}

func runLLMNode(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, err error) {
	nodeOutputMapResult = nodeOutputMap
	llm, err := NewLLMNode(nodeId, nodeMap)
	if err != nil {
		return
	}
	nodeOutputMapResult, err = llm.RunLLM(ctx, nodeMap, nodeOutputMapResult)
	if err != nil {
		return
	}
	fmt.Printf("%+v", nodeOutputMapResult)
	return
}

func runConditionNode(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (isSuccess bool, err error) {
	condition, err := NewConditionNode(nodeId, nodeMap, nodeOutputMap)
	if err != nil {
		return
	}
	isSuccess, _, err = condition.RunCondition(ctx, nodeOutputMap, nodeMap)
	if err != nil {
		return
	}
	fmt.Println(isSuccess)
	fmt.Printf("%+v", nodeOutputMap)
	return
}

func runCodeNode(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, err error) {
	nodeOutputMapResult = nodeOutputMap
	code, err := NewCodeNode(nodeId, nodeMap, nodeOutputMapResult)
	if err != nil {
		return
	}
	fmt.Printf("code:%+v", code)

	nodeOutputMapResult, err = code.RunCode(ctx, nodeOutputMapResult)
	return
}

func runEndNode(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (outputVariable, answerContent string, err error) {
	end, err := NewEndNode(nodeId, nodeMap, nodeOutputMap)
	outputVariable, answerContent, err = end.RunEnd(ctx)
	return
}

func runPluginsNode(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, err error) {
	nodeOutputMapResult = nodeOutputMap
	pluginsNode, err := NewPluginsNode(nodeId, nodeMap, nodeOutputMap)
	if err != nil {
		return
	}

	nodeOutputMap, err = pluginsNode.RunPlugins(ctx, nodeOutputMap, nodeMap)
	if err != nil {
		return
	}

	return
}

func runWorkflow(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, err error) {
	node, err := NewWorkflowNode(nodeId, nodeMap)
	if err != nil {
		return nodeOutputMapResult, err
	}

	nodeOutputMapResult, err = node.RunWorkflow(ctx, nodeMap, nodeOutputMap)
	if err != nil {
		return nodeOutputMapResult, err
	}
	return
}
