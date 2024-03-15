package workflownode

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	TypeStartNode     = "1"
	TypeEndNode       = "2"
	TypeLLMNode       = "3"
	TypePluginsNode   = "4"
	TypeCodeNode      = "5"
	TypeKnowledgeNode = "6"
	TypeConditionNode = "8"
	TypeWorkflowNode  = "9"
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
	TestResult TestResult `json:"testResult"`
}

type TestResult struct {
	ResultJson     string `json:"resultJson"`
	OutputVariable string `json:"outputVariable"`
	AnswerContent  string `json:"answerContent"`
}

type SchemaData struct {
	Outputs  []SchemaOutputs `json:"outputs"`
	Inputs   SchemaInputs    `json:"inputs"`
	NodeMeta SchemaNodeMeta  `json:"nodeMeta"`
	Version  string          `json:"version"` //llm节点
}

type SchemaOutputs struct {
	Type         string            `json:"type"`
	Name         string            `json:"name"`
	Required     bool              `json:"required"`
	Description  string            `json:"description"`
	Value        any               `json:"value"` //自定义添加的
	ListSchema   OneOutputSchema   `json:"listSchema"`
	ObjectSchema []OneOutputSchema `json:"objectSchema"`
	ResultJson   string            `json:"resultJson"`
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
	Branches        []SchemaBranches        `json:"branches"`      //condition节点
	LlmParam        []SchemaInputParameters `json:"llmParam"`      //llm节点
	DatasetParam    []DatasetParam          `json:"datasetParam"`  //knowledge
	ApiParam        []SchemaInputParameters `json:"apiParam"`      //plugins
	WorkflowId      string                  `json:"workflowId"`    //workflow
	SpaceId         string                  `json:"spaceId"`       //workflow
	InputDefs       []SchemaOutputs         `json:"inputDefs"`     //workflow
	Type            int64                   `json:"type"`          //workflow
}

type DatasetParam struct {
	Name  string          `json:"name"`
	Input SchemaInputData `json:"input"`
}

type SchemaBranches struct {
	Condition SchemaBranchConditions `json:"condition"`
}

type SchemaBranchConditions struct {
	Logic      int64              `json:"logic"` //2and  1or
	Conditions []SchemaConditions `json:"conditions"`
}

type SchemaConditions struct {
	Operator int64                 `json:"operator"` //1：input类型equal。2notequal。9is empty
	Left     SchemaInputParameters `json:"left"`
	Right    SchemaInputParameters `json:"right"`
}

type SchemaInputParameters struct {
	Name  string          `json:"name"`
	Input SchemaInputData `json:"input"`
}

type SchemaInputData struct {
	Type       string                `json:"type"`
	Value      SchemaValueData       `json:"value"`
	ListSchema SchemaKnowledgeSchema `json:"listSchema"`
}

type SchemaKnowledgeSchema struct {
	Type string `json:"type"`
}

type SchemaValueData struct {
	Type               string            `json:"type"`
	Content            SchemaContentData `json:"content"`
	LiteralContent     string            `json:"literalContent,omitempty"`
	StringArrayContent []string          `json:"stringArrayContent,omitempty"`
}

type SchemaContentData struct {
	Source  string `json:"source"`
	BlockID string `json:"blockID"`
	Name    string `json:"name"`
	Value   any    `json:"value"`
}

type Edge struct {
	SourceNodeID string `json:"sourceNodeID"`           //起点
	TargetNodeID string `json:"targetNodeID"`           //下一个
	SourcePortID string `json:"sourcePortID,omitempty"` //“true” “false”
}

func NewSchema(schemaJson string) (schema *Schema, nodeMap map[string]Node, err error) {
	if len(schemaJson) > 0 {
		err = json.Unmarshal([]byte(schemaJson), &schema)
	} else {
		err = json.Unmarshal([]byte(schemaJsonExample), &schema)
	}

	if err != nil {
		fmt.Println(err.Error())
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

var schemaJsonExample string = "{\"nodes\":[{\"id\":\"100001\",\"type\":\"1\",\"meta\":{\"position\":{\"x\":192,\"y\":0}},\"data\":{\"outputs\":[{\"type\":\"string\",\"name\":\"a\",\"required\":true,\"description\":\"参数a\"},{\"type\":\"integer\",\"name\":\"b\",\"required\":true,\"description\":\"参数b\"},{\"type\":\"boolean\",\"name\":\"z\",\"required\":false,\"description\":\"zzz\"},{\"type\":\"float\",\"name\":\"x\",\"required\":true,\"description\":\"xxx\"},{\"type\":\"string\",\"name\":\"y1\",\"required\":true},{\"type\":\"string\",\"name\":\"y2\",\"required\":true}],\"nodeMeta\":{\"title\":\"Start\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Start.png\",\"description\":\"The starting node of the workflow, used to set the information needed to initiate the workflow.\",\"subTitle\":\"\"}}},{\"id\":\"900001\",\"type\":\"2\",\"meta\":{\"position\":{\"x\":4343,\"y\":878}},\"data\":{\"nodeMeta\":{\"title\":\"End\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-End.png\",\"description\":\"The final node of the workflow, used to return the result information after the workflow runs.\",\"subTitle\":\"\"},\"inputs\":{\"terminatePlan\":\"useAnswerContent\",\"inputParameters\":[{\"name\":\"c\",\"input\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"list\",\"name\":\"key4511\",\"listSchema\":{\"type\":\"integer\"}}],\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"120710\",\"name\":\"key4.0.key45.key451\"}}}}],\"content\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"我的content\"}}}}},{\"id\":\"117411\",\"type\":\"11\",\"meta\":{\"position\":{\"x\":1362,\"y\":922}},\"data\":{\"nodeMeta\":{\"title\":\"Variable\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Variable.png\",\"description\":\"Used for reading and writing variables in your bot. The variable name must match the variable name in Bot.\",\"subTitle\":\"Variable\"},\"inputs\":{\"mode\":\"set\",\"inputParameters\":[{\"name\":\"botVariable\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"bbb\"}}}]},\"outputs\":[{\"type\":\"boolean\",\"name\":\"isSuccess\"}]}},{\"id\":\"120710\",\"type\":\"5\",\"meta\":{\"position\":{\"x\":3719.884210526316,\"y\":-78.4263157894737}},\"data\":{\"nodeMeta\":{\"title\":\"Code\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Code.png\",\"description\":\"Write code to process input variables to generate return values.\",\"subTitle\":\"Code\"},\"inputs\":{\"inputParameters\":[{\"name\":\"input1\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}},{\"name\":\"input2\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"3232\"}}},{\"name\":\"input3\",\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}},{\"name\":\"input4\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}],\"code\":\"function main( params ){\\n    const ret = {\\n        \\\"key0\\\": params.input1 + params.input2,\\n        \\\"key1\\\": \\\"hi\\\",\\n        \\\"key2\\\": [\\\"hello\\\", \\\"world\\\"],\\n        \\\"key3\\\": {\\n            \\\"key31\\\": \\\"hi\\\"\\n        },\\n        \\\"key4\\\": [{\\n            \\\"key41\\\": true,\\n            \\\"key42\\\": 1,\\n            \\\"key43\\\": 12.88,\\n            \\\"key44\\\": [\\\"hello\\\"],\\n            \\\"key45\\\": {\\n                \\\"key451\\\": {\\n                    \\\"key4511\\\":[9,4,2]\\n                }\\n            }\\n        }]\\n    };\\n\\n    return ret;\\n}\",\"language\":5},\"outputs\":[{\"type\":\"string\",\"name\":\"key0\"},{\"type\":\"string\",\"name\":\"key1\"},{\"type\":\"list\",\"name\":\"key2\",\"listSchema\":{\"type\":\"string\"}},{\"type\":\"object\",\"name\":\"key3\",\"objectSchema\":[{\"type\":\"string\",\"name\":\"key31\"}]},{\"type\":\"list\",\"name\":\"key4\",\"listSchema\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"boolean\",\"name\":\"key41\"},{\"type\":\"integer\",\"name\":\"key42\"},{\"type\":\"float\",\"name\":\"key43\"},{\"type\":\"list\",\"name\":\"key44\",\"listSchema\":{\"type\":\"string\"}},{\"type\":\"object\",\"name\":\"key45\",\"objectSchema\":[{\"type\":\"object\",\"name\":\"key451\",\"objectSchema\":[{\"type\":\"list\",\"name\":\"key4511\",\"listSchema\":{\"type\":\"integer\"}}]}]}]}}]}},{\"id\":\"171236\",\"type\":\"8\",\"meta\":{\"position\":{\"x\":2992,\"y\":751}},\"data\":{\"nodeMeta\":{\"title\":\"Condition\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Condition.png\",\"description\":\"Connect two downstream branches. If the set conditions are met, run only the ‘if’ branch; otherwise, run only the ‘else’ branch.\",\"subTitle\":\"Condition\"},\"inputs\":{\"branches\":[{\"condition\":{\"logic\":2,\"conditions\":[{\"operator\":1,\"left\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}},\"right\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"literal\",\"literalContent\":\"false\"}}}},{\"operator\":2,\"left\":{\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}},\"right\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}}},{\"operator\":9,\"left\":{\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}},{\"operator\":10,\"left\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"z\"}}}}},{\"operator\":13,\"left\":{\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"x\"}}}},\"right\":{\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}}},{\"operator\":11,\"left\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}}},{\"operator\":12,\"left\":{\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}}},{\"operator\":14,\"left\":{\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"x\"}}}},\"right\":{\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}},{\"operator\":15,\"left\":{\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"x\"}}}},\"right\":{\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}},{\"operator\":16,\"left\":{\"input\":{\"type\":\"float\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"x\"}}}},\"right\":{\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}},{\"operator\":10,\"left\":{\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"181217\",\"name\":\"outputList.output\"}}}}}]}}]}}},{\"id\":\"196739\",\"type\":\"3\",\"meta\":{\"position\":{\"x\":2111.5,\"y\":624.011385199241}},\"data\":{\"nodeMeta\":{\"title\":\"LLM\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-LLM.png\",\"description\":\"Invoke the large language model, generate responses using variables and prompt words.\",\"subTitle\":\"LLM\"},\"inputs\":{\"inputParameters\":[{\"name\":\"input1\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"a\"}}},{\"name\":\"input2\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"1\"}}},{\"name\":\"input3\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"false\"}}},{\"name\":\"input4\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"1.5\"}}}],\"llmParam\":[{\"name\":\"modleName\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"GPT-3.5 (16K)\"}}},{\"name\":\"modelType\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"literal\",\"literalContent\":\"113\"}}},{\"name\":\"prompt\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"input1：{{.input1}}\\ninput2：{{.input2}}\\ninput3：{{.input3}}\\ninput4：{{.input4}}\"}}},{\"name\":\"temperature\",\"input\":{\"type\":\"float\",\"value\":{\"type\":\"literal\",\"literalContent\":\"0.7\"}}}]},\"outputs\":[{\"type\":\"string\",\"name\":\"output1\",\"description\":\"是input1的值\"},{\"type\":\"integer\",\"name\":\"output2\",\"description\":\"是input2的值\"},{\"type\":\"boolean\",\"name\":\"output3\",\"description\":\"是input3的值\"},{\"type\":\"float\",\"name\":\"output4\",\"description\":\"是input4的值\"},{\"type\":\"list\",\"name\":\"output5\",\"listSchema\":{\"type\":\"string\"},\"description\":\"是input1的集合\"},{\"type\":\"list\",\"name\":\"output6\",\"listSchema\":{\"type\":\"integer\"},\"description\":\"是input2的集合\"},{\"type\":\"list\",\"name\":\"output7\",\"listSchema\":{\"type\":\"boolean\"},\"description\":\"是input3的集合\"},{\"type\":\"list\",\"name\":\"output8\",\"listSchema\":{\"type\":\"float\"},\"description\":\"是input4的集合\"}],\"version\":\"2\"}},{\"id\":\"181217\",\"type\":\"6\",\"meta\":{\"position\":{\"x\":726,\"y\":755.1052631578947}},\"data\":{\"nodeMeta\":{\"title\":\"Knowledge\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Knowledge.png\",\"description\":\"In the selected knowledge, the best matching information is recalled based on the input variable and returned as an Array.\",\"subTitle\":\"Knowledge\"},\"outputs\":[{\"type\":\"list\",\"name\":\"outputList\",\"listSchema\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"string\",\"name\":\"output\"}]}}],\"inputs\":{\"inputParameters\":[{\"name\":\"Query\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}}],\"datasetParam\":[{\"name\":\"datasetList\",\"input\":{\"type\":\"list\",\"listSchema\":{\"type\":\"string\"},\"value\":{\"type\":\"literal\",\"stringArrayContent\":[\"7338252484214980609\"]}}},{\"name\":\"topK\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"literal\",\"literalContent\":\"3\"}}},{\"name\":\"minScore\",\"input\":{\"type\":\"number\",\"value\":{\"type\":\"literal\",\"literalContent\":\"0.5\"}}},{\"name\":\"strategy\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"literal\",\"literalContent\":\"1\"}}}]}}}],\"edges\":[{\"sourceNodeID\":\"120710\",\"targetNodeID\":\"900001\"},{\"sourceNodeID\":\"171236\",\"targetNodeID\":\"120710\",\"sourcePortID\":\"true\"},{\"sourceNodeID\":\"171236\",\"targetNodeID\":\"900001\",\"sourcePortID\":\"false\"},{\"sourceNodeID\":\"196739\",\"targetNodeID\":\"171236\"},{\"sourceNodeID\":\"100001\",\"targetNodeID\":\"181217\"},{\"sourceNodeID\":\"181217\",\"targetNodeID\":\"117411\"},{\"sourceNodeID\":\"117411\",\"targetNodeID\":\"196739\"}]}"
