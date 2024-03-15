package workflownode

import (
	"context"
	"fmt"
	"testing"
)

var schemaJson string = "{\"nodes\":[{\"id\":\"100001\",\"type\":\"1\",\"meta\":{\"position\":{\"x\":11.008421052631578,\"y\":-385.29473684210524}},\"data\":{\"outputs\":[{\"type\":\"string\",\"name\":\"s1\",\"required\":true,\"description\":\"1\"},{\"type\":\"integer\",\"name\":\"s2\",\"required\":true,\"description\":\"2\"},{\"type\":\"boolean\",\"name\":\"s3\",\"required\":true,\"description\":\"3\"},{\"type\":\"float\",\"name\":\"s4\",\"required\":true,\"description\":\"4\"}],\"nodeMeta\":{\"title\":\"Start\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Start.png\",\"description\":\"The starting node of the workflow, used to set the information needed to initiate the workflow.\",\"subTitle\":\"\"}}},{\"id\":\"900001\",\"type\":\"2\",\"meta\":{\"position\":{\"x\":1185.4144360902258,\"y\":-220.16842105263157}},\"data\":{\"nodeMeta\":{\"title\":\"End\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-End.png\",\"description\":\"The final node of the workflow, used to return the result information after the workflow runs.\",\"subTitle\":\"\"},\"inputs\":{\"terminatePlan\":\"returnVariables\",\"inputParameters\":[{\"name\":\"output\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"125658\",\"name\":\"output\"}}}}]}}},{\"id\":\"125658\",\"type\":\"9\",\"meta\":{\"position\":{\"x\":544.3542857142858,\"y\":-259.15187969924807}},\"data\":{\"nodeMeta\":{\"title\":\"test2\",\"description\":\"test2\",\"icon\":\"https://lf16-alice-tos-sign.oceanapi-i18n.com/obj/ocean-cloud-tos-sg/plugin_icon/workflow.png?lk3s=cd508e2b&x-expires=1710163728&x-signature=2T2bI3R%2Fx0gNVu%2FH%2Fgh46hGtveU%3D\"},\"inputs\":{\"workflowId\":\"7337993373015605266\",\"spaceId\":\"7332418282986913799\",\"inputDefs\":[{\"name\":\"a\",\"type\":\"string\",\"description\":\"参数a\",\"required\":true},{\"name\":\"b\",\"type\":\"integer\",\"description\":\"参数b\",\"required\":true},{\"name\":\"z\",\"type\":\"boolean\",\"description\":\"zzz\",\"required\":false},{\"name\":\"x\",\"type\":\"float\",\"description\":\"xxx\",\"required\":true},{\"name\":\"y1\",\"type\":\"string\",\"required\":true},{\"name\":\"y2\",\"type\":\"string\",\"required\":true}],\"type\":1,\"inputParameters\":[{\"name\":\"a\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"s1\"}}}},{\"name\":\"b\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"s2\"}}}},{\"name\":\"z\",\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"s3\"}}}},{\"name\":\"x\",\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"s4\"}}}},{\"name\":\"y1\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"s1\"}}}},{\"name\":\"y2\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"s1\"}}}}]},\"outputs\":[{\"type\":\"string\",\"name\":\"output\",\"required\":false}]}}],\"edges\":[{\"sourceNodeID\":\"100001\",\"targetNodeID\":\"125658\"},{\"sourceNodeID\":\"125658\",\"targetNodeID\":\"900001\"}]}"
var params = make(map[string]SchemaOutputs)

func TestEdges(t *testing.T) {
	params["s1"] = SchemaOutputs{
		Name:        "s1",
		Type:        "string",
		Required:    true,
		Description: "1",
		Value:       1,
	}
	params["s2"] = SchemaOutputs{
		Name:        "s2",
		Type:        "integer",
		Required:    true,
		Description: "2",
		Value:       2,
	}
	params["s3"] = SchemaOutputs{
		Name:        "s3",
		Type:        "boolean",
		Required:    true,
		Description: "3",
		Value:       true,
	}
	params["s4"] = SchemaOutputs{
		Name:        "s4",
		Type:        "float",
		Required:    true,
		Description: "4",
		Value:       1.5,
	}
	ctx := context.Background()

	edgesResult, err := RunEdges(ctx, schemaJson, params)
	if err != nil {
		return
	}

	fmt.Printf("%+v", edgesResult)
	return
}

//func TestEdgesV2(t *testing.T) {
//	params["s1"] = SchemaOutputs{
//		Name:        "s1",
//		Type:        "string",
//		Required:    true,
//		Description: "1",
//		Value:       1,
//	}
//	params["s2"] = SchemaOutputs{
//		Name:        "s2",
//		Type:        "integer",
//		Required:    true,
//		Description: "2",
//		Value:       2,
//	}
//	params["s3"] = SchemaOutputs{
//		Name:        "s3",
//		Type:        "boolean",
//		Required:    true,
//		Description: "3",
//		Value:       true,
//	}
//	params["s4"] = SchemaOutputs{
//		Name:        "s4",
//		Type:        "float",
//		Required:    true,
//		Description: "4",
//		Value:       1.5,
//	}
//	KpWorkflow, err := kpworkflow.NewKpWorkflow(schemaJson, params)
//	if err != nil {
//		return
//	}
//	ctx := context.Background()
//	nodeList, err := KpWorkflow.Call(ctx)
//	if err != nil {
//		return
//	}
//	fmt.Printf("%+v", nodeList)
//	return
//}
