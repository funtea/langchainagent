package basenode

import (
	"encoding/json"
	"errors"
)

const (
	TypeStartNode     = "1"
	TypeEndNode       = "2"
	TypeCodeNode      = "5"
	TypeConditionNode = "8"
	TypeVariableNode  = "11"
)

type Schema struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type Node struct {
	Id   string     `json:"id"`
	Type string     `json:"type"`
	Data SchemaData `json:"data"`
	//Meta any        `json:"meta"`
}

type SchemaData struct {
	Outputs  []SchemaOutputs `json:"outputs"`
	Inputs   SchemaInputs    `json:"inputs"`
	NodeMeta SchemaNodeMeta  `json:"nodeMeta"`
}

type SchemaOutputs struct {
	Type         string            `json:"type"`
	Name         string            `json:"name"`
	Required     bool              `json:"required"`
	Description  string            `json:"description"`
	Value        any               `json:"value"` //自定义添加的
	ListSchema   OneOutputSchema   `json:"listSchema"`
	ObjectSchema []OneOutputSchema `json:"objectSchema"`
	pathKey      string
}

type OneOutputSchema struct {
	Type         string            `json:"type"`
	Value        any               `json:"value"` //自定义添加的
	Name         string            `json:"name"`  //object 才有该字段
	ListSchema   TwoOutputSchema   `json:"listSchema"`
	ObjectSchema []TwoOutputSchema `json:"objectSchema"`
}

type TwoOutputSchema struct {
	Type         string              `json:"type"`
	Value        any                 `json:"value"` //自定义添加的
	Name         string              `json:"name"`  //object 才有该字段
	ListSchema   ThreeOutputSchema   `json:"listSchema"`
	ObjectSchema []ThreeOutputSchema `json:"objectSchema"`
}

type ThreeOutputSchema struct {
	Type         string             `json:"type"`
	Value        any                `json:"value"` //自定义添加的
	Name         string             `json:"name"`  //object 才有该字段
	ListSchema   FourOutputSchema   `json:"listSchema"`
	ObjectSchema []FourOutputSchema `json:"objectSchema"`
}

type FourOutputSchema struct {
	Type       string           `json:"type"`
	Value      any              `json:"value"` //自定义添加的
	Name       string           `json:"name"`  //object 才有该字段
	ListSchema FiveOutputSchema `json:"listSchema"`
}

type FiveOutputSchema struct {
	Type  string `json:"type"`
	Value any    `json:"value"` //自定义添加的
	Name  string `json:"name"`  //object 才有该字段
}

type SchemaNodeMeta struct {
	Title       string `json:"title"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
	SubTitle    string `json:"subTitle"`
}

type SchemaInputs struct {
	Mode            string                  `json:"mode"` //variable节点 set get
	InputParameters []SchemaInputParameters `json:"inputParameters"`
	Code            string                  `json:"code"`
	Language        int64                   `json:"language"`
	TerminatePlan   string                  `json:"terminatePlan"` //end节点 returnVariables：不需要拼接内容  "useAnswerContent" 需要拼接内容
	Content         SchemaInputData         `json:"content"`       //end节点 拼接内容用
}

type SchemaInputParameters struct {
	Name  string          `json:"name"`
	Input SchemaInputData `json:"input"`
}

type SchemaInputData struct {
	Type  string          `json:"type"`
	Value SchemaValueData `json:"value"`
}

type SchemaValueData struct {
	Type           string            `json:"type"`
	Content        SchemaContentData `json:"content"`
	LiteralContent string            `json:"literalContent,omitempty"`
}

type SchemaContentData struct {
	Source  string `json:"source"`
	BlockID string `json:"blockID"`
	Name    string `json:"name"`
	Value   any    `json:"value"`
}

type Edge struct {
	SourceNodeID string `json:"sourceNodeID"`
	TargetNodeID string `json:"targetNodeID"`
}

func NewSchema(schemaJson string) (schema *Schema, nodeMap map[string]Node, err error) {
	err = json.Unmarshal([]byte(schemaJsonExample), &schema)
	if err != nil {
		return
	}

	nodeMap, err = schema.ParseInputs()
	if err != nil {
		return
	}

	return
}

func (schema *Schema) ParseInputs() (map[string]Node, error) {
	var nodeMap = make(map[string]Node)
	if len(schema.Nodes) < 0 {
		return nodeMap, errors.New("无节点")
	}

	for _, node := range schema.Nodes {
		nodeMap[node.Id] = node
	}
	return nodeMap, nil
}

var schemaJsonExample string = "{\"nodes\":[{\"id\":\"100001\",\"type\":\"1\",\"meta\":{\"position\":{\"x\":192,\"y\":110.5}},\"data\":{\"outputs\":[{\"type\":\"string\",\"name\":\"a\",\"required\":true,\"description\":\"参数a\"},{\"type\":\"integer\",\"name\":\"b\",\"required\":true,\"description\":\"参数b\"},{\"type\":\"boolean\",\"name\":\"z\",\"required\":false,\"description\":\"zzz\"},{\"type\":\"float\",\"name\":\"x\",\"required\":true,\"description\":\"xxx\"}],\"nodeMeta\":{\"title\":\"Start\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Start.png\",\"description\":\"The starting node of the workflow, used to set the information needed to initiate the workflow.\",\"subTitle\":\"\"}}},{\"id\":\"900001\",\"type\":\"2\",\"meta\":{\"position\":{\"x\":2957.001897533207,\"y\":578.2362428842505}},\"data\":{\"nodeMeta\":{\"title\":\"End\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-End.png\",\"description\":\"The final node of the workflow, used to return the result information after the workflow runs.\",\"subTitle\":\"\"},\"inputs\":{\"terminatePlan\":\"useAnswerContent\",\"inputParameters\":[{\"name\":\"c\",\"input\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"list\",\"name\":\"key4511\",\"listSchema\":{\"type\":\"integer\"}}],\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"120710\",\"name\":\"key4.key45.key451\"}}}}],\"content\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"我的content\"}}}}},{\"id\":\"117411\",\"type\":\"11\",\"meta\":{\"position\":{\"x\":778,\"y\":687.5}},\"data\":{\"nodeMeta\":{\"title\":\"Variable\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Variable.png\",\"description\":\"Used for reading and writing variables in your bot. The variable name must match the variable name in Bot.\",\"subTitle\":\"Variable\"},\"inputs\":{\"mode\":\"set\",\"inputParameters\":[{\"name\":\"botVariable\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"bbb\"}}}]},\"outputs\":[{\"type\":\"boolean\",\"name\":\"isSuccess\"}]}},{\"id\":\"120710\",\"type\":\"5\",\"meta\":{\"position\":{\"x\":2377.684210526316,\"y\":-451.57894736842104}},\"data\":{\"nodeMeta\":{\"title\":\"Code\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Code.png\",\"description\":\"Write code to process input variables to generate return values.\",\"subTitle\":\"Code\"},\"inputs\":{\"inputParameters\":[{\"name\":\"input1\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}},{\"name\":\"input2\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"3232\"}}},{\"name\":\"input3\",\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}},{\"name\":\"input4\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}],\"code\":\"function main( params ){\\n    const ret = {\\n        \\\"key0\\\": params.input + params.input,\\n        \\\"key1\\\": \\\"hi\\\",\\n        \\\"key2\\\": [\\\"hello\\\", \\\"world\\\"],\\n        \\\"key3\\\": {\\n            \\\"key31\\\": \\\"hi\\\"\\n        },\\n        \\\"key4\\\": [{\\n            \\\"key41\\\": true,\\n            \\\"key42\\\": 1,\\n            \\\"key43\\\": 12.88,\\n            \\\"key44\\\": [\\\"hello\\\"],\\n            \\\"key45\\\": {\\n                \\\"key451\\\": {\\n                    \\\"key4511\\\":[9,4,2]\\n                }\\n            }\\n        }]\\n    };\\n\\n    return ret;\\n}\",\"language\":5},\"outputs\":[{\"type\":\"string\",\"name\":\"key0\"},{\"type\":\"string\",\"name\":\"key1\"},{\"type\":\"list\",\"name\":\"key2\",\"listSchema\":{\"type\":\"string\"}},{\"type\":\"object\",\"name\":\"key3\",\"objectSchema\":[{\"type\":\"string\",\"name\":\"key31\"}]},{\"type\":\"list\",\"name\":\"key4\",\"listSchema\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"boolean\",\"name\":\"key41\"},{\"type\":\"integer\",\"name\":\"key42\"},{\"type\":\"float\",\"name\":\"key43\"},{\"type\":\"list\",\"name\":\"key44\",\"listSchema\":{\"type\":\"string\"}},{\"type\":\"object\",\"name\":\"key45\",\"objectSchema\":[{\"type\":\"object\",\"name\":\"key451\",\"objectSchema\":[{\"type\":\"list\",\"name\":\"key4511\",\"listSchema\":{\"type\":\"integer\"}}]}]}]}}]}},{\"id\":\"171236\",\"type\":\"8\",\"meta\":{\"position\":{\"x\":1603,\"y\":530.5}},\"data\":{\"nodeMeta\":{\"title\":\"Condition\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Condition.png\",\"description\":\"Connect two downstream branches. If the set conditions are met, run only the ‘if’ branch; otherwise, run only the ‘else’ branch.\",\"subTitle\":\"Condition\"},\"inputs\":{\"branches\":[{\"condition\":{\"logic\":2,\"conditions\":[{\"operator\":1,\"left\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}},\"right\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"literal\",\"literalContent\":\"true\"}}}},{\"operator\":2,\"left\":{\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}},\"right\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}}},{\"operator\":9,\"left\":{\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}},{\"operator\":10,\"left\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"z\"}}}}},{\"operator\":13,\"left\":{\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"x\"}}}},\"right\":{\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}}},{\"operator\":11,\"left\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}}},{\"operator\":12,\"left\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}}},{\"operator\":14,\"left\":{\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"x\"}}}},\"right\":{\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}},{\"operator\":15,\"left\":{\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"x\"}}}},\"right\":{\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}},{\"operator\":16,\"left\":{\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"x\"}}}},\"right\":{\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}}]}}]}}}],\"edges\":[{\"sourceNodeID\":\"100001\",\"targetNodeID\":\"117411\"},{\"sourceNodeID\":\"120710\",\"targetNodeID\":\"900001\"},{\"sourceNodeID\":\"117411\",\"targetNodeID\":\"171236\"},{\"sourceNodeID\":\"171236\",\"targetNodeID\":\"120710\",\"sourcePortID\":\"true\"},{\"sourceNodeID\":\"171236\",\"targetNodeID\":\"900001\",\"sourcePortID\":\"false\"}]}"

//var schemaJsonExample string = "{\"nodes\":[{\"id\":\"100001\",\"type\":\"1\",\"meta\":{\"position\":{\"x\":-819.3697963254239,\"y\":-411.65624440949136}},\"data\":{\"outputs\":[{\"type\":\"string\",\"name\":\"a\",\"required\":true,\"description\":\"参数a\"},{\"type\":\"integer\",\"name\":\"b\",\"required\":true,\"description\":\"参数b\"},{\"type\":\"boolean\",\"name\":\"z\",\"required\":false,\"description\":\"zzz\"},{\"type\":\"float\",\"name\":\"x\",\"required\":true,\"description\":\"xxx\"}],\"nodeMeta\":{\"title\":\"Start\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Start.png\",\"description\":\"The starting node of the workflow, used to set the information needed to initiate the workflow.\",\"subTitle\":\"\"}}},{\"id\":\"900001\",\"type\":\"2\",\"meta\":{\"position\":{\"x\":1042.8874766175206,\"y\":-342.7518870012242}},\"data\":{\"nodeMeta\":{\"title\":\"End\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-End.png\",\"description\":\"The final node of the workflow, used to return the result information after the workflow runs.\",\"subTitle\":\"\"},\"inputs\":{\"terminatePlan\":\"useAnswerContent\",\"inputParameters\":[{\"name\":\"c\",\"input\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"list\",\"name\":\"key4511\",\"objectSchema\":{\"type\":\"integer\"}}],\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"120710\",\"name\":\"key4.0.key45.key451\"}}}}],\"content\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"我的content\"}}}}},{\"id\":\"117411\",\"type\":\"11\",\"meta\":{\"position\":{\"x\":-237.04042005996448,\"y\":-307.84672281764773}},\"data\":{\"nodeMeta\":{\"title\":\"Variable\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Variable.png\",\"description\":\"Used for reading and writing variables in your bot. The variable name must match the variable name in Bot.\",\"subTitle\":\"Variable\"},\"inputs\":{\"mode\":\"set\",\"inputParameters\":[{\"name\":\"botVariable\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"bbb\"}}}]},\"outputs\":[{\"type\":\"boolean\",\"name\":\"isSuccess\"}]}},{\"id\":\"120710\",\"type\":\"5\",\"meta\":{\"position\":{\"x\":414.02340143315394,\"y\":-712.5092562838116}},\"data\":{\"nodeMeta\":{\"title\":\"Code\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Code.png\",\"description\":\"Write code to process input variables to generate return values.\",\"subTitle\":\"Code\"},\"inputs\":{\"inputParameters\":[{\"name\":\"input1\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}},{\"name\":\"input2\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"3232\"}}},{\"name\":\"input3\",\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}},{\"name\":\"input4\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}],\"code\":\"async def main(args: Args) -> Output:\\n    params = args.params\\n    ret: Output = {\\n        \\\"key0\\\": params['input1'] + params['input2'],\\n        \\\"key1\\\": \\\"hi\\\",\\n        \\\"key2\\\": [\\\"hello\\\", \\\"world\\\"],\\n        \\\"key3\\\": {\\n            \\\"key31\\\": \\\"hi\\\"\\n        },\\n        \\\"key4\\\": [{\\n            \\\"key41\\\": True,\\n            \\\"key42\\\": 1,\\n            \\\"key43\\\": 12.88,\\n            \\\"key44\\\": [\\\"hello\\\"],\\n            \\\"key45\\\": {\\n                \\\"key451\\\": {\\n                    \\\"key4511\\\":[4,5,6]\\n                }\\n            }\\n        }]\\n    }\\n    return ret\",\"language\":3},\"outputs\":[{\"type\":\"string\",\"name\":\"key0\"},{\"type\":\"string\",\"name\":\"key1\"},{\"type\":\"list\",\"name\":\"key2\",\"listSchema\":{\"type\":\"string\"}},{\"type\":\"object\",\"name\":\"key3\",\"objectSchema\":[{\"type\":\"string\",\"name\":\"key31\"}]},{\"type\":\"list\",\"name\":\"key4\",\"listSchema\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"boolean\",\"name\":\"key41\"},{\"type\":\"integer\",\"name\":\"key42\"},{\"type\":\"float\",\"name\":\"key43\"},{\"type\":\"list\",\"name\":\"key44\",\"listSchema\":{\"type\":\"string\"}},{\"type\":\"object\",\"name\":\"key45\",\"objectSchema\":[{\"type\":\"object\",\"name\":\"key451\",\"objectSchema\":[{\"type\":\"list\",\"name\":\"key4511\",\"listSchema\":{\"type\":\"integer\"}}]}]}]}}]}}],\"edges\":[{\"sourceNodeID\":\"100001\",\"targetNodeID\":\"117411\"},{\"sourceNodeID\":\"117411\",\"targetNodeID\":\"120710\"},{\"sourceNodeID\":\"120710\",\"targetNodeID\":\"900001\"}]}"
