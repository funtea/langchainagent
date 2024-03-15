package workflownode

import (
	"context"
	"errors"
	"fmt"
)

func RunEdges(ctx context.Context, schemaJson string, params map[string]SchemaOutputs) (dealWorkFlowResult DealWorkFlowResult, err error) {
	schema, nodeMap, err := NewSchema(schemaJson)
	if err != nil {
		return
	}

	edgeList := schema.Edges
	if len(edgeList) == 0 {
		return dealWorkFlowResult, errors.New("无路径节点")
	}

	//获取开始，结束点
	startNodeId, _, err := FindFirstAndEndNode(schema.Nodes)
	if err != nil {
		return dealWorkFlowResult, err
	}

	//获取节点路径map
	edgeSourceNodeIdMap := GetEdgeMap(edgeList)

	nodeList := schema.Nodes
	fmt.Printf("%+v", nodeList)

	dealWorkFlowResult, err = DealWorkFlow(ctx, startNodeId, edgeSourceNodeIdMap, nodeMap, params)

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

type DealWorkFlowResult struct {
	StartNode     *Node
	EndNode       *Node
	CodeNode      *Node
	ConditionNode *Node
	KnowledgeNode *Node
	LlmNode       *Node
	PluginsNode   *Node
	VariableNode  *Node
	WorkflowNode  *Node
}

// 处理工作流主函数
func DealWorkFlow(ctx context.Context, startNodeId string, edgeSourceNodeIdMap map[string]Edge, nodeMap map[string]Node, params map[string]SchemaOutputs) (dealWorkFlowResult DealWorkFlowResult, err error) {
	var i int64 = 0
	//初始化节点id，默认开始节点
	var thisNodeId = startNodeId
	var nodeOutputMap = make(map[string]map[string]SchemaOutputs)

	var startNode, endNode, codeNode, knowledgeNode, llmNode, pluginsNode, variableNode, workflowNode *Node
	for {
		if nodeMap[thisNodeId].Type == TypeStartNode {
			//开始节点
			nodeOutputMap, startNode, err = runStart(ctx, nodeMap, params)
			if err != nil {
				return dealWorkFlowResult, err
			}
			dealWorkFlowResult.StartNode = startNode
		}

		if nodeMap[thisNodeId].Type == TypeKnowledgeNode {
			//知识点节点
			nodeOutputMap, knowledgeNode, err = runKnowledge(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return dealWorkFlowResult, err
			}
			dealWorkFlowResult.KnowledgeNode = knowledgeNode
		}

		if nodeMap[thisNodeId].Type == TypeWorkflowNode {
			//工作流节点
			nodeOutputMap, workflowNode, err = runWorkflow(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return dealWorkFlowResult, err
			}
			dealWorkFlowResult.WorkflowNode = workflowNode
		}

		if nodeMap[thisNodeId].Type == TypeVariableNode {
			//变量节点
			nodeOutputMap, variableNode, err = runVariable(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return dealWorkFlowResult, err
			}
			dealWorkFlowResult.VariableNode = variableNode
		}

		if nodeMap[thisNodeId].Type == TypeLLMNode {
			//llm节点
			nodeOutputMap, llmNode, err = runLLMNode(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return dealWorkFlowResult, err
			}
			dealWorkFlowResult.LlmNode = llmNode
		}

		if nodeMap[thisNodeId].Type == TypeConditionNode {
			//condition节点
			isSuccess, conditionNode, err := runConditionNode(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return dealWorkFlowResult, err
			}

			//todo nextNodeId
			if isSuccess == true {
				thisNodeId = edgeSourceNodeIdMap[thisNodeId+"true"].TargetNodeID
			} else {
				thisNodeId = edgeSourceNodeIdMap[thisNodeId+"false"].TargetNodeID
			}
			dealWorkFlowResult.ConditionNode = conditionNode
			continue
		}

		if nodeMap[thisNodeId].Type == TypeCodeNode {
			//code节点
			nodeOutputMap, codeNode, err = runCodeNode(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return dealWorkFlowResult, err
			}
			dealWorkFlowResult.CodeNode = codeNode
		}

		if nodeMap[thisNodeId].Type == TypePluginsNode {
			//变量节点
			nodeOutputMap, pluginsNode, err = runPluginsNode(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return dealWorkFlowResult, err
			}
			dealWorkFlowResult.PluginsNode = pluginsNode
		}

		//赋值下一轮节点id
		if nodeMap[thisNodeId].Type == TypeEndNode {
			endNode, err = runEndNode(ctx, thisNodeId, nodeMap, nodeOutputMap)
			if err != nil {
				return dealWorkFlowResult, err
			}
			dealWorkFlowResult.EndNode = endNode
			break
		}

		thisNodeId = edgeSourceNodeIdMap[thisNodeId].TargetNodeID
		i++
	}

	return
}

func runStart(ctx context.Context, nodeMap map[string]Node, params map[string]SchemaOutputs) (nodeOutputMap map[string]map[string]SchemaOutputs, node *Node, err error) {
	startNode, err := NewStartNode(nodeMap, params)
	if err != nil {
		return
	}
	nodeOutputMap, resultJson := startNode.RunStart(ctx)
	fmt.Printf("startNode 结构体:%+v \r\n", startNode)
	fmt.Printf("startResult 结构体:%+v \r\n", nodeOutputMap)
	startNode.TestResult.ResultJson = resultJson
	return nodeOutputMap, startNode, err
}

func runKnowledge(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, node *Node, err error) {
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
	return nodeOutputMapResult, knowledge, err
}

func runVariable(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, node *Node, err error) {
	nodeOutputMapResult = nodeOutputMap
	variable, err := NewVariableNode(nodeId, nodeMap, nodeOutputMapResult)
	if err != nil {
		return
	}

	nodeOutputMapResult, resultJson, err := variable.RunVariable(ctx, nodeOutputMapResult)
	fmt.Printf("%+v", nodeOutputMapResult)
	variable.TestResult.ResultJson = resultJson
	return nodeOutputMapResult, variable, err
}

func runLLMNode(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, node *Node, err error) {
	nodeOutputMapResult = nodeOutputMap
	llm, err := NewLLMNode(nodeId, nodeMap)
	if err != nil {
		return
	}
	nodeOutputMapResult, resultJson, err := llm.RunLLM(ctx, nodeMap, nodeOutputMapResult)
	if err != nil {
		return
	}
	fmt.Printf("%+v", nodeOutputMapResult)
	llm.TestResult.ResultJson = resultJson
	return nodeOutputMapResult, llm, err
}

func runConditionNode(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (isSuccess bool, node *Node, err error) {
	condition, err := NewConditionNode(nodeId, nodeMap, nodeOutputMap)
	if err != nil {
		return
	}
	isSuccess, resultJson, err := condition.RunCondition(ctx, nodeOutputMap, nodeMap)
	if err != nil {
		return
	}
	fmt.Println(isSuccess)
	fmt.Printf("%+v", nodeOutputMap)
	condition.TestResult.ResultJson = resultJson
	return isSuccess, condition, err
}

func runCodeNode(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, node *Node, err error) {
	nodeOutputMapResult = nodeOutputMap
	code, err := NewCodeNode(nodeId, nodeMap, nodeOutputMapResult)
	if err != nil {
		return
	}
	fmt.Printf("code:%+v", code)

	nodeOutputMapResult, resultJson, err := code.RunCode(ctx, nodeOutputMapResult)
	code.TestResult.ResultJson = resultJson
	return nodeOutputMapResult, code, err
}

func runEndNode(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (node *Node, err error) {
	end, err := NewEndNode(nodeId, nodeMap, nodeOutputMap)
	outputVariable, answerContent, err := end.RunEnd(ctx)
	end.TestResult.OutputVariable = outputVariable
	end.TestResult.AnswerContent = answerContent
	return end, err
}

func runPluginsNode(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, node *Node, err error) {
	nodeOutputMapResult = nodeOutputMap
	pluginsNode, err := NewPluginsNode(nodeId, nodeMap, nodeOutputMap)
	if err != nil {
		return
	}

	nodeOutputMap, resultJson, err := pluginsNode.RunPlugins(ctx, nodeOutputMap, nodeMap)
	if err != nil {
		return
	}
	pluginsNode.TestResult.ResultJson = resultJson
	return nodeOutputMapResult, pluginsNode, err
}

func runWorkflow(ctx context.Context, nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, node *Node, err error) {
	workflowNode, err := NewWorkflowNode(nodeId, nodeMap)
	if err != nil {
		return nodeOutputMapResult, workflowNode, err
	}

	nodeOutputMapResult, resultJson, err := workflowNode.RunWorkflow(ctx, nodeMap, nodeOutputMap)
	if err != nil {
		return nodeOutputMapResult, workflowNode, err
	}

	workflowNode.TestResult.ResultJson = resultJson
	return nodeOutputMapResult, workflowNode, err
}
