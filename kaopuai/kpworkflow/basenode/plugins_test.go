package basenode

import (
	"context"
	"testing"
)

func TestPlugins(t *testing.T) {
	var params = map[string]SchemaOutputs{
		"p1": {
			Type:  "string",
			Name:  "p1",
			Value: "345",
		},
		"p2": {
			Type:  "integer",
			Name:  "p2",
			Value: 789,
		},
		"p3": {
			Type:  "boolean",
			Name:  "p3",
			Value: true,
		},
		"p4": {
			Type:  "float",
			Name:  "p4",
			Value: 8.5,
		},
	}
	ctx := context.Background()
	var schemaJson = "{\"nodes\":[{\"id\":\"100001\",\"type\":\"1\",\"meta\":{\"position\":{\"x\":192,\"y\":70.5}},\"data\":{\"outputs\":[{\"type\":\"string\",\"name\":\"p1\",\"required\":true,\"description\":\"1\"},{\"type\":\"integer\",\"name\":\"p2\",\"required\":true,\"description\":\"2\"},{\"type\":\"boolean\",\"name\":\"p3\",\"required\":true,\"description\":\"p3\"},{\"type\":\"float\",\"name\":\"p4\",\"required\":true,\"description\":\"p4\"}],\"nodeMeta\":{\"title\":\"Start\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Start.png\",\"description\":\"The starting node of the workflow, used to set the information needed to initiate the workflow.\",\"subTitle\":\"\"}}},{\"id\":\"900001\",\"type\":\"2\",\"meta\":{\"position\":{\"x\":1894,\"y\":179}},\"data\":{\"nodeMeta\":{\"title\":\"End\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-End.png\",\"description\":\"The final node of the workflow, used to return the result information after the workflow runs.\",\"subTitle\":\"\"},\"inputs\":{\"terminatePlan\":\"returnVariables\",\"inputParameters\":[{\"name\":\"output\",\"input\":{\"type\":\"list\",\"listSchema\":{\"type\":\"string\"},\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"115892\",\"name\":\"content.result\"}}}}]}}},{\"id\":\"180778\",\"type\":\"4\",\"meta\":{\"position\":{\"x\":726,\"y\":0}},\"data\":{\"nodeMeta\":{\"title\":\"exchange\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon_Api.png\",\"subtitle\":\"shanda_exchange:exchange\",\"description\":\"This is an API used to query real-time exchange rates\"},\"inputs\":{\"apiParam\":[{\"name\":\"apiID\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"7343257522645680130\"}}},{\"name\":\"apiName\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"exchange\"}}},{\"name\":\"pluginID\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"7343252407138140166\"}}},{\"name\":\"pluginName\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"shanda_exchange\"}}},{\"name\":\"pluginVersion\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"\"}}},{\"name\":\"tips\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"\"}}},{\"name\":\"outDocLink\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"\"}}}],\"inputParameters\":[{\"name\":\"user_input_history\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"p1\"}}}},{\"name\":\"from\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"p2\"}}}},{\"name\":\"to\",\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"p3\"}}}},{\"name\":\"amount\",\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"p4\"}}}}]},\"outputs\":[{\"type\":\"string\",\"name\":\"error\",\"required\":false},{\"type\":\"string\",\"name\":\"request_id\",\"required\":false},{\"type\":\"object\",\"name\":\"result\",\"objectSchema\":[{\"type\":\"object\",\"name\":\"data\",\"objectSchema\":[{\"type\":\"string\",\"name\":\"base_currency_name\",\"required\":false},{\"type\":\"object\",\"name\":\"rates\",\"objectSchema\":[{\"type\":\"object\",\"name\":\"CNY\",\"objectSchema\":[{\"type\":\"string\",\"name\":\"currency_name\",\"required\":false},{\"type\":\"string\",\"name\":\"rate\",\"required\":false},{\"type\":\"string\",\"name\":\"rate_for_amount\",\"required\":false}],\"required\":false}],\"required\":false},{\"type\":\"string\",\"name\":\"status\",\"required\":false},{\"type\":\"string\",\"name\":\"updated_date\",\"required\":false},{\"type\":\"string\",\"name\":\"amount\",\"required\":false},{\"type\":\"string\",\"name\":\"base_currency_code\",\"required\":false}],\"required\":false}],\"required\":false},{\"type\":\"float\",\"name\":\"status\",\"required\":false}]}},{\"id\":\"115892\",\"type\":\"4\",\"meta\":{\"position\":{\"x\":1281.6465972169847,\"y\":76.60963333584736}},\"data\":{\"nodeMeta\":{\"title\":\"PushMarkdown\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon_Api.png\",\"subtitle\":\"PushDeer:PushMarkdown\",\"description\":\"Push markdown content as a notification to user's device\"},\"inputs\":{\"apiParam\":[{\"name\":\"apiID\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"7343179496817950725\"}}},{\"name\":\"apiName\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"PushMarkdown\"}}},{\"name\":\"pluginID\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"7343179375640166406\"}}},{\"name\":\"pluginName\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"PushDeer\"}}},{\"name\":\"pluginVersion\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"\"}}},{\"name\":\"tips\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"\"}}},{\"name\":\"outDocLink\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"\"}}}],\"inputParameters\":[{\"name\":\"pushkey\",\"input\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"object\",\"name\":\"data\",\"objectSchema\":[{\"type\":\"string\",\"name\":\"base_currency_name\",\"required\":false},{\"type\":\"object\",\"name\":\"rates\",\"objectSchema\":[{\"type\":\"object\",\"name\":\"CNY\",\"objectSchema\":[{\"type\":\"string\",\"name\":\"currency_name\",\"required\":false},{\"type\":\"string\",\"name\":\"rate\",\"required\":false},{\"type\":\"string\",\"name\":\"rate_for_amount\",\"required\":false}],\"required\":false}],\"required\":false},{\"type\":\"string\",\"name\":\"status\",\"required\":false},{\"type\":\"string\",\"name\":\"updated_date\",\"required\":false},{\"type\":\"string\",\"name\":\"amount\",\"required\":false},{\"type\":\"string\",\"name\":\"base_currency_code\",\"required\":false}],\"required\":false}],\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"180778\",\"name\":\"result\"}}}},{\"name\":\"text\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"180778\",\"name\":\"result.data.base_currency_name\"}}}},{\"name\":\"desp\",\"input\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"string\",\"name\":\"currency_name\",\"required\":false},{\"type\":\"string\",\"name\":\"rate\",\"required\":false},{\"type\":\"string\",\"name\":\"rate_for_amount\",\"required\":false}],\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"180778\",\"name\":\"result.data.rates.CNY\"}}}},{\"name\":\"type\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"180778\",\"name\":\"result.data.rates.CNY.currency_name\"}}}}]},\"outputs\":[{\"type\":\"object\",\"name\":\"content\",\"objectSchema\":[{\"type\":\"list\",\"name\":\"result\",\"listSchema\":{\"type\":\"string\"},\"required\":false}],\"required\":false},{\"type\":\"float\",\"name\":\"code\",\"required\":false,\"description\":\"Return code from push service, 0 means success\"}]}}],\"edges\":[{\"sourceNodeID\":\"100001\",\"targetNodeID\":\"180778\"},{\"sourceNodeID\":\"180778\",\"targetNodeID\":\"115892\"},{\"sourceNodeID\":\"115892\",\"targetNodeID\":\"900001\"}]}"
	schemaStruct, nodeMap, err := NewSchema(schemaJson)
	if err != nil {
		return
	}

	startNode, err := NewStartNode(nodeMap, params)
	if err != nil {
		return
	}
	nodeOutputMap := startNode.RunStart(ctx)

	node1 := schemaStruct.Nodes[2]
	plugins1, err := NewPluginsNode(node1.Id, nodeMap, nodeOutputMap)
	if err != nil {
		return
	}

	nodeOutputMap, err = plugins1.RunPlugins(ctx, nodeOutputMap, nodeMap)
	if err != nil {
		return
	}

	node2 := schemaStruct.Nodes[3]
	plugins2, err := NewPluginsNode(node2.Id, nodeMap, nodeOutputMap)
	if err != nil {
		return
	}
	nodeOutputMap, err = plugins2.RunPlugins(ctx, nodeOutputMap, nodeMap)
	if err != nil {
		return
	}

	return
}
