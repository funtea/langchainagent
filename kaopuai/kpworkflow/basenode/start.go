package basenode

import (
	"errors"
	"fmt"
)

/*
startJson 	json
params      前端启动参数
*/
func NewStartNode(nodeMap map[string]Node, params map[string]SchemaOutputs) (start *Node, err error) {
	//err = json.Unmarshal([]byte(startJsonExample), &start)
	//if err != nil {
	//	return
	//}

	var startNode = Node{}
	for _, node := range nodeMap {
		if node.Type == TypeStartNode {
			startNode = node
		}
	}

	err = (&startNode).ParseStartInputs(params)
	if err != nil {
		return
	}

	return &startNode, nil
}

func (start *Node) ParseStartInputs(params map[string]SchemaOutputs) (err error) {
	if len(start.Data.Outputs) == 0 {
		return
	}

	for key, output := range start.Data.Outputs {
		tmpOutput, exist := params[output.Name]
		//check require
		if output.Required == true && exist == false {
			return errors.New(fmt.Sprintf("开始节点，参数%s为必填", output.Name))
		}

		//check 类型
		if output.Type != tmpOutput.Type {
			return errors.New(fmt.Sprintf("开始节点，参数%s类型不一致", output.Name))
		}

		//deal 赋值
		start.Data.Outputs[key].Value = tmpOutput.Value
	}

	return
}

// Run方法
func (start *Node) RunStart() map[string]map[string]SchemaOutputs {
	//逻辑处理，因为start逻辑简单，只需要把value提取出来，作为返回参数即可。  其他节点需要参与不同计算得出结果
	var outputs = make(map[string]map[string]SchemaOutputs)
	if len(start.Data.Outputs) == 0 {
		return outputs
	}

	for _, output := range start.Data.Outputs {
		// 检查outputs[start.Id]是否为nil，如果是，则初始化它
		if _, ok := outputs[start.Id]; !ok {
			outputs[start.Id] = make(map[string]SchemaOutputs)
		}

		// 现在outputs[start.Id]已经被初始化，可以安全地添加键值对
		outputs[start.Id][output.Name] = output
	}

	return outputs
}

var startJsonExample string = "{\"id\":\"100001\",\"type\":\"1\",\"meta\":{\"position\":{\"x\":-631.8697963254239,\"y\":-364.15624440949136}},\"data\":{\"outputs\":[{\"type\":\"string\",\"name\":\"a\",\"required\":true,\"description\":\"参数a\"},{\"type\":\"integer\",\"name\":\"b\",\"required\":true,\"description\":\"参数b\"}],\"nodeMeta\":{\"title\":\"Start\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Start.png\",\"description\":\"The starting node of the workflow, used to set the information needed to initiate the workflow.\",\"subTitle\":\"\"}}}"
