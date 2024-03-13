package kpworkflow

import (
	"context"
	"errors"
	"github.com/tmc/langchaingo/kaopuai/kpworkflow/basenode"
)

// KpWorkflow 是llm代理与工作流交互的工具
type KpWorkflowConcept interface {
	initWorkflow() error
	Call(ctx context.Context, input string) (string, error)
}

type KpWorkflow struct {
	SchemaJson     string
	NodeMap        map[string]basenode.Node
	EdgeMap        map[string]basenode.Edge
	StartNodeId    string
	WorkflowParams map[string]basenode.SchemaOutputs //启动参数
}

func NewKpWorkflow(jsonStr string, workflowParams map[string]basenode.SchemaOutputs) (*KpWorkflow, error) {
	k := &KpWorkflow{
		SchemaJson:     jsonStr,
		WorkflowParams: workflowParams,
	}
	err := k.initWorkflow()
	if err != nil {
		return nil, err
	}
	return k, nil
}

func (r *KpWorkflow) initWorkflow() (err error) {
	schema, nodeMap, err := basenode.NewSchema(r.SchemaJson)
	if err != nil {
		return
	}

	edgeList := schema.Edges
	if len(edgeList) == 0 {
		return errors.New("无路径节点")
	}

	//获取开始，结束点
	startNodeId, _, err := basenode.FindFirstAndEndNode(schema.Nodes)
	if err != nil {
		return err
	}
	r.StartNodeId = startNodeId

	//获取节点路径map
	edgeSourceNodeIdMap := basenode.GetEdgeMap(edgeList)
	r.NodeMap = nodeMap
	r.EdgeMap = edgeSourceNodeIdMap

	return nil
}

func (r *KpWorkflow) Call(ctx context.Context) (result string, err error) {
	// zhao kaishi jiedian
	// kaishi de xia yige jiedian
	// chuli wan canshu , zhixing zhaodaode jieidan

	result, err = basenode.DealWorkFlow(ctx, r.StartNodeId, r.EdgeMap, r.NodeMap, r.WorkflowParams)
	if err != nil {
		return
	}
	return
}
