package workflownode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"os/exec"
	v8 "rogchap.com/v8go"
	"strings"
)

/**
 *nodeMap        key节点id   value  节点
 *nodeOutputMap  key节点id   value  节点输出的变量值
 */
func NewCodeNode(nodeId string, nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (code *Node, err error) {
	node := nodeMap[nodeId]
	if node.Type != TypeCodeNode {
		return nil, errors.New("code节点类型错误")
	}

	err = node.ParseCodeInput(nodeMap, nodeOutputMap)
	if err != nil {
		return nil, err
	}
	return &node, nil
}

func (node *Node) ParseCodeInput(nodeMap map[string]Node, nodeOutputMap map[string]map[string]SchemaOutputs) (err error) {
	if len(node.Data.Inputs.InputParameters) == 0 {
		//参数可以为空
		return nil
	}

	for key, inputParameter := range node.Data.Inputs.InputParameters {
		if inputParameter.Input.Value.Type == "ref" {
			//引用类型
			nodeId := node.Data.Inputs.InputParameters[key].Input.Value.Content.BlockID
			varName := node.Data.Inputs.InputParameters[key].Input.Value.Content.Name
			refNode := nodeOutputMap[nodeId]

			if nodeMap[nodeId].Type == TypeCodeNode ||
				nodeMap[nodeId].Type == TypeLLMNode ||
				nodeMap[nodeId].Type == TypeKnowledgeNode {

				//if inputParameter.Input.Type == "object" || inputParameter.Input.Type == "list" {
				//	return errors.New(fmt.Sprintf("参数%s 不能为object或者list", varName))
				//}

				//如果是code节点,  只取第一层数据
				scriptJsonAny := refNode["outputList"].Value
				scriptJson := scriptJsonAny.(string)
				codeGjson := gjson.Parse(scriptJson)
				tmpValue := codeGjson.Get(varName).String()
				node.Data.Inputs.InputParameters[key].Input.Value.Content.Value = tmpValue
			} else {
				//todo 检查变量名是否存在  varName
				node.Data.Inputs.InputParameters[key].Input.Value.Content.Value = refNode[varName].Value
			}
		} else if inputParameter.Input.Value.Type == "literal" {
			//直接使用节点设置的类型值,  自身就是值
			//node.Data.Inputs.InputParameters[0].Input.Value.LiteralContent
		}
	}

	return
}

/**
 * nodeOutputMap 其他节点输出的值
 */
func (code *Node) RunCode(ctx context.Context, nodeOutputMap map[string]map[string]SchemaOutputs) (nodeOutputMapResult map[string]map[string]SchemaOutputs, resultJson string, err error) {
	fmt.Printf("%+v", code)
	nodeOutputMapResult = nodeOutputMap
	todoCode := code.Data.Inputs.Code
	var scriptResult string
	if code.Data.Inputs.Language == 5 {
		//javascript
		javaScriptResult, err := runJavaScript(todoCode, code.Data.Inputs.InputParameters)
		if err != nil {
			return nodeOutputMapResult, resultJson, err
		}

		scriptResult = javaScriptResult
		//nodeOutputMap[code.Id]["runJson"] = scriptResult
		//outputList := code.Data.Outputs
		//nodeOutputMap = code.dealJavascriptResult(scriptResult, nodeOutputMap, outputList)
	} else if code.Data.Inputs.Language == 3 {
		//python
		pythonResult, err := runPython3(todoCode, code.Data.Inputs.InputParameters)
		if err != nil {
			return nodeOutputMap, resultJson, err
		}

		scriptResult = pythonResult
	}

	var schemaOutputs = SchemaOutputs{
		Name:  "outputList",
		Value: scriptResult,
	}

	if _, ok := nodeOutputMap[code.Id]; !ok {
		nodeOutputMap[code.Id] = make(map[string]SchemaOutputs)
	}
	nodeOutputMap[code.Id]["outputList"] = schemaOutputs
	return nodeOutputMap, resultJson, nil
}

func runJavaScript(todoCode string, inputParamList []SchemaInputParameters) (string, error) {
	var tmpCodeVariable = "var pa = {"
	for _, inputParam := range inputParamList {
		tmpCodeVariable += `"` + inputParam.Name + `":`
		if inputParam.Input.Type == "string" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += `"` + fmt.Sprintf("%s", inputParam.Input.Value.Content.Value) + `",`
			} else {
				tmpCodeVariable += `"` + fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `",`
			}
		} else if inputParam.Input.Type == "integer" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += fmt.Sprintf("%d", inputParam.Input.Value.Content.Value) + `,`
			} else {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `,`
			}
		} else if inputParam.Input.Type == "boolean" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += fmt.Sprintf("%t", inputParam.Input.Value.Content.Value) + `,`
			} else {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `,`
			}
		} else if inputParam.Input.Type == "float" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += fmt.Sprintf("%f", inputParam.Input.Value.Content.Value) + `,`
			} else {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `,`
			}
		}
	}
	tmpCodeVariable = strings.TrimRight(tmpCodeVariable, ",")
	tmpCodeVariable += "};"

	todoCode = tmpCodeVariable + todoCode + "main(pa);"

	ctx := v8.NewContext()
	val, err := ctx.RunScript(todoCode, "main.js") // global object will have the property set within the JS VM
	if err != nil {
		e := err.(*v8.JSError)    // JavaScript errors will be returned as the JSError struct
		fmt.Println(e.Message)    // the message of the exception thrown
		fmt.Println(e.Location)   // the filename, line number and the column where the error occured
		fmt.Println(e.StackTrace) // the full stack trace of the error, if available

		fmt.Printf("javascript error: %v", e)        // will format the standard error message
		fmt.Printf("javascript stack trace: %+v", e) // will format the full error stack trace
		return "", err
	}

	marshal, err := json.Marshal(val)
	if err != nil {
		return "", err
	}
	fmt.Printf("result: %+v", string(marshal))

	return string(marshal), nil
}

func runPython3(todoCode string, inputParamList []SchemaInputParameters) (string, error) {
	//构造python 参数
	var tmpCodeVariable = "\nparams = {"
	for _, inputParam := range inputParamList {
		tmpCodeVariable += "'" + inputParam.Name + "':"
		if inputParam.Input.Type == "string" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += "'" + fmt.Sprintf("%s", inputParam.Input.Value.Content.Value) + "',"
			} else {
				tmpCodeVariable += "'" + fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + "',"
			}
		} else if inputParam.Input.Type == "integer" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += fmt.Sprintf("%d", inputParam.Input.Value.Content.Value) + `,`
			} else {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `,`
			}
		} else if inputParam.Input.Type == "boolean" {
			if inputParam.Input.Value.Type == "ref" {
				if (inputParam.Input.Value.Content.Value).(bool) {
					tmpCodeVariable += "True,"
				} else {
					tmpCodeVariable += "False,"
				}
				//tmpCodeVariable += fmt.Sprintf("%t", inputParam.Input.Value.Content.Value) + `,`
			} else {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `,`
			}
		} else if inputParam.Input.Type == "float" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += fmt.Sprintf("%f", inputParam.Input.Value.Content.Value) + `,`
			} else {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `,`
			}
		}
	}
	tmpCodeVariable = strings.TrimRight(tmpCodeVariable, ",")
	tmpCodeVariable += "}\n"

	//python 程序
	pythonCode := "import json\n"
	pythonCode += "import asyncio\n"
	pythonCode += "from typing import Dict, Any, List\n"
	pythonCode += "class Args:\n"
	pythonCode += "    def __init__(self, params: Dict[str, Any]):\n"
	pythonCode += "        self.params = params\n"
	pythonCode += "class Output:\n"
	pythonCode += "    def __init__(self, data: Dict[str, Any]):\n"
	pythonCode += "        self.data = data\n"

	// 定义一个Python代码字符串
	//pythonCode += "def main(args: Args) -> Output:\n    params = args.params\n    ret: Output = {\n        \"key0\": params['input1'] + params['input2'],\n        \"key1\": \"hi\",\n        \"key2\": [\"hello\", \"world\"],\n        \"key3\": {\n            \"key31\": \"hi\"\n        },\n        \"key4\": [{\n            \"key41\": True,\n            \"key42\": 1,\n            \"key43\": 12.88,\n            \"key44\": [\"hello\"],\n            \"key45\": {\n                \"key451\": {\n                    \"key4511\":[1,2,3]\n                }\n            }\n        },{\n            \"key41\": True,\n            \"key42\": 1,\n            \"key43\": 12.88,\n            \"key44\": [\"hello\"],\n            \"key45\": {\n                \"key451\": {\n                    \"key4511\":[1,2,3]\n                }\n            }\n        }]\n    }\n    return ret\n"
	pythonCode += todoCode

	//pythonCode += "params = {'input1': 'Hello from ','input2': 'Go!'}\n"
	pythonCode += tmpCodeVariable
	pythonCode += "args = Args(params)\n"
	pythonCode += "loop = asyncio.get_event_loop()\n"
	pythonCode += "result = loop.run_until_complete(main(args))\n"
	pythonCode += "json_data = json.dumps(result)\n"
	pythonCode += "print(json_data)\n"
	cmd := exec.Command("python3", "-c", pythonCode) // 确保python3命令可用，或者根据你的环境修改为python或python2等
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("\nError executing command:", err)
		return "", err
	}
	fmt.Println("\nPython script output:", string(output))
	return string(output), nil
}

/**
 *v8Result v8运行结果json
 *nodeOutputMap  各个节点输出的汇总结果，将本次结果放到该map中
 *outputList   code节点outputparamers格式
 */
func (code *Node) dealJavascriptResult(jsonResult string, nodeOutputMap map[string]map[string]SchemaOutputs, schemaOutputList []SchemaOutputs) map[string]map[string]SchemaOutputs {
	fmt.Printf("%+v", schemaOutputList)
	if len(schemaOutputList) == 0 {
		return nodeOutputMap
	}
	gjsonObject := gjson.Parse(jsonResult)
	for oneKey, oneSchemaOutput := range schemaOutputList {
		if oneSchemaOutput.Type == "string" {
			//直接赋值
			tmpValue := gjsonObject.Get(oneSchemaOutput.Name)
			schemaOutputList[oneKey].Value = tmpValue.String()
		} else if oneSchemaOutput.Type == "integer" {
			//直接赋值
			tmpValue := gjsonObject.Get(oneSchemaOutput.Name)
			schemaOutputList[oneKey].Value = tmpValue.Int()
		} else if oneSchemaOutput.Type == "float" {
			//直接赋值
			tmpValue := gjsonObject.Get(oneSchemaOutput.Name)
			schemaOutputList[oneKey].Value = tmpValue.Float()
		} else if oneSchemaOutput.Type == "boolean" {
			//直接赋值
			tmpValue := gjsonObject.Get(oneSchemaOutput.Name)
			schemaOutputList[oneKey].Value = tmpValue.Bool()
		} else if oneSchemaOutput.Type == "list" {
			tmpValue := gjsonObject.Get(oneSchemaOutput.Name)
			if oneSchemaOutput.ListSchema.Type != "object" {
				tmpValueList := tmpValue.Array()
				schemaOutputList[oneKey].ListSchema.Value = tmpValueList

			} else {
				//还需要处理object
				fmt.Println()

				TwoResult := DealThreeCodeOutputObject(tmpValue.Raw, oneSchemaOutput.ListSchema.ObjectSchema, true)
				schemaOutputList[oneKey].ListSchema.ObjectSchema = TwoResult
				fmt.Printf("%+v", TwoResult)
			}
		} else if oneSchemaOutput.Type == "object" {
			tmpValue := gjsonObject.Get(oneSchemaOutput.Name)

			TwoResult := DealTwoCodeOutputObject(tmpValue.Raw, oneSchemaOutput.ObjectSchema, false)
			schemaOutputList[oneKey].ObjectSchema = TwoResult
			fmt.Printf("%+v", TwoResult)
		}
	}

	//返回
	for _, tmpSchemaOutput := range schemaOutputList {
		if _, ok := nodeOutputMap[code.Id]; !ok {
			nodeOutputMap[code.Id] = make(map[string]SchemaOutputs)
		}
		nodeOutputMap[code.Id][tmpSchemaOutput.Name] = tmpSchemaOutput
	}

	return nodeOutputMap
}

// 处理第二层
func DealTwoCodeOutputObject(toParseJson string, schemaOutputList []OneOutputSchema, isArray bool) []OneOutputSchema {
	if isArray {
		gjsonObject := gjson.Parse(toParseJson)
		gjsonObject.ForEach(func(tmpKey, value gjson.Result) bool {
			//gjsonkey := int64(tmpKey.Index)
			//tmpValue := value.Get(schemaOutputList[gjsonkey].Name)
			//if tmpValue.Raw == "" {
			//	return false
			//}

			if tmpKey.Index > 0 {
				//coze 截止2月28日，只有一个索引。     不要再次赋值
				fmt.Println("coze 增加了索引")
				return false
			}

			for key, TwoSchemaOutput := range schemaOutputList {
				tmpValue := value.Get(schemaOutputList[key].Name)
				if TwoSchemaOutput.Type == "string" {
					schemaOutputList[key].Value = tmpValue.String()
				} else if TwoSchemaOutput.Type == "integer" {
					//直接赋值
					schemaOutputList[key].Value = tmpValue.Int()
				} else if TwoSchemaOutput.Type == "float" {
					//直接赋值
					schemaOutputList[key].Value = tmpValue.Float()
				} else if TwoSchemaOutput.Type == "boolean" {
					//直接赋值
					schemaOutputList[key].Value = tmpValue.Bool()
				} else if TwoSchemaOutput.Type == "list" {
					if TwoSchemaOutput.ListSchema.Type != "object" {
						schemaOutputList[key].ListSchema.Value = tmpValue.Array()
					} else {
						//还需要处理object
						TwoResult := DealFourCodeOutputObject(tmpValue.Raw, TwoSchemaOutput.ListSchema.ObjectSchema, true)
						schemaOutputList[key].ListSchema.ObjectSchema = TwoResult
						fmt.Printf("%+v", TwoResult)
					}
				} else if TwoSchemaOutput.Type == "object" {
					fmt.Println()
					//递归调用
					TwoResult := DealThreeCodeOutputObject(tmpValue.Raw, TwoSchemaOutput.ObjectSchema, false)
					schemaOutputList[key].ObjectSchema = TwoResult
				}

				fmt.Printf("%+v", TwoSchemaOutput)
			}

			return true
		})
	} else {
		gjsonObject := gjson.Parse(toParseJson)
		for key, TwoSchemaOutput := range schemaOutputList {
			tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
			if TwoSchemaOutput.Type == "string" {
				schemaOutputList[key].Value = tmpValue.String()
			} else if TwoSchemaOutput.Type == "integer" {
				//直接赋值
				schemaOutputList[key].Value = tmpValue.Int()
			} else if TwoSchemaOutput.Type == "float" {
				//直接赋值
				schemaOutputList[key].Value = tmpValue.Float()
			} else if TwoSchemaOutput.Type == "boolean" {
				//直接赋值
				schemaOutputList[key].Value = tmpValue.Bool()
			} else if TwoSchemaOutput.Type == "list" {
				if TwoSchemaOutput.ListSchema.Type != "object" {
					schemaOutputList[key].ListSchema.Value = tmpValue.Array()
				} else {
					//还需要处理object
					TwoResult := DealFourCodeOutputObject(tmpValue.Raw, TwoSchemaOutput.ListSchema.ObjectSchema, true)
					schemaOutputList[key].ListSchema.ObjectSchema = TwoResult
					fmt.Printf("%+v", TwoResult)
				}
			} else if TwoSchemaOutput.Type == "object" {
				fmt.Println()
				//递归调用
				TwoResult := DealThreeCodeOutputObject(tmpValue.Raw, TwoSchemaOutput.ObjectSchema, false)
				schemaOutputList[key].ObjectSchema = TwoResult
			}

			fmt.Printf("%+v", TwoSchemaOutput)
		}
	}

	return schemaOutputList
}

// 处理第三层
func DealThreeCodeOutputObject(toParseJson string, schemaOutputList []TwoOutputSchema, isArray bool) []TwoOutputSchema {
	if isArray {
		gjsonObject := gjson.Parse(toParseJson)
		gjsonObject.ForEach(func(tmpKey, value gjson.Result) bool {
			//gjsonkey := int64(tmpKey.Index)

			//if tmpValue.Raw == "" {
			//	return false
			//}

			if tmpKey.Index > 0 {
				//coze 截止2月28日，只有一个索引。     不要再次赋值
				fmt.Println("coze 增加了索引")
				return false
			}

			for key, TwoSchemaOutput := range schemaOutputList {
				tmpValue := value.Get(schemaOutputList[key].Name)
				if TwoSchemaOutput.Type == "string" {
					schemaOutputList[key].Value = tmpValue.String()
				} else if TwoSchemaOutput.Type == "integer" {
					//直接赋值
					schemaOutputList[key].Value = tmpValue.Int()
				} else if TwoSchemaOutput.Type == "float" {
					//直接赋值
					schemaOutputList[key].Value = tmpValue.Float()
				} else if TwoSchemaOutput.Type == "boolean" {
					//直接赋值
					schemaOutputList[key].Value = tmpValue.Bool()
				} else if TwoSchemaOutput.Type == "list" {
					if TwoSchemaOutput.ListSchema.Type != "object" {
						schemaOutputList[key].ListSchema.Value = tmpValue.Array()
					} else {
						//还需要处理object
						TwoResult := DealFiveCodeOutputObject(tmpValue.Raw, TwoSchemaOutput.ListSchema.ObjectSchema, true)
						schemaOutputList[key].ListSchema.ObjectSchema = TwoResult
						fmt.Printf("%+v", TwoResult)
					}
				} else if TwoSchemaOutput.Type == "object" {
					fmt.Println()
					//递归调用
					TwoResult := DealFourCodeOutputObject(tmpValue.Raw, TwoSchemaOutput.ObjectSchema, false)
					schemaOutputList[key].ObjectSchema = TwoResult
				}

				fmt.Printf("%+v", TwoSchemaOutput)
			}

			return true
		})
	} else {
		gjsonObject := gjson.Parse(toParseJson)
		for key, TwoSchemaOutput := range schemaOutputList {
			tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
			if TwoSchemaOutput.Type == "string" {
				schemaOutputList[key].Value = tmpValue.String()
			} else if TwoSchemaOutput.Type == "integer" {
				//直接赋值
				schemaOutputList[key].Value = tmpValue.Int()
			} else if TwoSchemaOutput.Type == "float" {
				//直接赋值
				schemaOutputList[key].Value = tmpValue.Float()
			} else if TwoSchemaOutput.Type == "boolean" {
				//直接赋值
				schemaOutputList[key].Value = tmpValue.Bool()
			} else if TwoSchemaOutput.Type == "list" {
				if TwoSchemaOutput.ListSchema.Type != "object" {
					schemaOutputList[key].ListSchema.Value = tmpValue.Array()
				} else {
					//还需要处理object
					TwoResult := DealFiveCodeOutputObject(tmpValue.Raw, TwoSchemaOutput.ListSchema.ObjectSchema, true)
					schemaOutputList[key].ListSchema.ObjectSchema = TwoResult
					fmt.Printf("%+v", TwoResult)
				}
			} else if TwoSchemaOutput.Type == "object" {
				fmt.Println()
				//递归调用
				TwoResult := DealFourCodeOutputObject(tmpValue.Raw, TwoSchemaOutput.ObjectSchema, false)
				schemaOutputList[key].ObjectSchema = TwoResult
			}

			fmt.Printf("%+v", TwoSchemaOutput)
		}
	}

	return schemaOutputList
}

//func getValue(keysetting, valueMap, key) {
//
//}

// 处理第三层
func DealFourCodeOutputObject(toParseJson string, schemaOutputList []ThreeOutputSchema, isArray bool) []ThreeOutputSchema {
	if isArray {
		gjsonObject := gjson.Parse(toParseJson)
		gjsonObject.ForEach(func(tmpKey, value gjson.Result) bool {
			//gjsonkey := int64(tmpKey.Index)
			//tmpValue := value.Get(schemaOutputList[gjsonkey].Name)
			//if tmpValue.Raw == "" {
			//	return false
			//}

			if tmpKey.Index > 0 {
				//coze 截止2月28日，只有一个索引。     不要再次赋值
				fmt.Println("coze 增加了索引")
				return false
			}

			for key, TwoSchemaOutput := range schemaOutputList {
				tmpValue := value.Get(schemaOutputList[key].Name)
				if TwoSchemaOutput.Type == "string" {
					schemaOutputList[key].Value = tmpValue.String()
				} else if TwoSchemaOutput.Type == "integer" {
					//直接赋值
					schemaOutputList[key].Value = tmpValue.Int()
				} else if TwoSchemaOutput.Type == "float" {
					//直接赋值
					schemaOutputList[key].Value = tmpValue.Float()
				} else if TwoSchemaOutput.Type == "boolean" {
					//直接赋值
					schemaOutputList[key].Value = tmpValue.Bool()
				} else if TwoSchemaOutput.Type == "list" {
					if TwoSchemaOutput.ListSchema.Type != "object" {
						schemaOutputList[key].ListSchema.Value = tmpValue.Array()
					} else {
						//还需要处理object
						fmt.Println("不应该进到这里！")
					}
				} else if TwoSchemaOutput.Type == "object" {
					fmt.Println()
					//递归调用
					TwoResult := DealFiveCodeOutputObject(tmpValue.Raw, TwoSchemaOutput.ObjectSchema, false)
					schemaOutputList[key].ObjectSchema = TwoResult
				}

				fmt.Printf("%+v", TwoSchemaOutput)
			}

			return true
		})
	} else {
		gjsonObject := gjson.Parse(toParseJson)
		for key, TwoSchemaOutput := range schemaOutputList {
			if TwoSchemaOutput.Type == "string" {
				tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
				schemaOutputList[key].Value = tmpValue.String()
			} else if TwoSchemaOutput.Type == "integer" {
				//直接赋值
				tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
				schemaOutputList[key].Value = tmpValue.Int()
			} else if TwoSchemaOutput.Type == "float" {
				//直接赋值
				tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
				schemaOutputList[key].Value = tmpValue.Float()
			} else if TwoSchemaOutput.Type == "boolean" {
				//直接赋值
				tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
				schemaOutputList[key].Value = tmpValue.Bool()
			} else if TwoSchemaOutput.Type == "list" {
				tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
				if TwoSchemaOutput.ListSchema.Type != "object" {
					schemaOutputList[key].ListSchema.Value = tmpValue.Array()
				} else {
					//还需要处理object
					fmt.Println("不应该进到这里！")
				}
			} else if TwoSchemaOutput.Type == "object" {
				fmt.Println()
				//递归调用
				tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)

				TwoResult := DealFiveCodeOutputObject(tmpValue.Raw, TwoSchemaOutput.ObjectSchema, false)
				schemaOutputList[key].ObjectSchema = TwoResult
			}

			fmt.Printf("%+v", TwoSchemaOutput)
		}
	}

	return schemaOutputList
}

// 处理第5层
func DealFiveCodeOutputObject(toParseJson string, schemaOutputList []FourOutputSchema, isArray bool) []FourOutputSchema {
	if isArray {
		gjsonObject := gjson.Parse(toParseJson)
		gjsonObject.ForEach(func(tmpKey, value gjson.Result) bool {
			//gjsonkey := int64(tmpKey.Index)
			//tmpValue := value.Get(schemaOutputList[gjsonkey].Name)
			//if tmpValue.Raw == "" {
			//	return false
			//}

			if tmpKey.Index > 0 {
				//coze 截止2月28日，只有一个索引。     不要再次赋值
				fmt.Println("coze 增加了索引")
				return false
			}

			for key, TwoSchemaOutput := range schemaOutputList {
				tmpValue := value.Get(schemaOutputList[key].Name)
				if TwoSchemaOutput.Type == "string" {
					schemaOutputList[key].Value = tmpValue.String()
				} else if TwoSchemaOutput.Type == "integer" {
					//直接赋值
					schemaOutputList[key].Value = tmpValue.Int()
				} else if TwoSchemaOutput.Type == "float" {
					//直接赋值
					schemaOutputList[key].Value = tmpValue.Float()
				} else if TwoSchemaOutput.Type == "boolean" {
					//直接赋值
					schemaOutputList[key].Value = tmpValue.Bool()
				} else if TwoSchemaOutput.Type == "list" {
					schemaOutputList[key].ListSchema.Value = tmpValue.Array()
				} else if TwoSchemaOutput.Type == "object" {
					//递归调用
					fmt.Printf("不应该存在第五层。tmpValue：%+v", tmpValue)
				}

				fmt.Printf("%+v", TwoSchemaOutput)
			}

			return true
		})
	} else {
		gjsonObject := gjson.Parse(toParseJson)
		for key, TwoSchemaOutput := range schemaOutputList {
			if TwoSchemaOutput.Type == "string" {
				tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
				schemaOutputList[key].Value = tmpValue.String()
			} else if TwoSchemaOutput.Type == "integer" {
				//直接赋值
				tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
				schemaOutputList[key].Value = tmpValue.Int()
			} else if TwoSchemaOutput.Type == "float" {
				//直接赋值
				tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
				schemaOutputList[key].Value = tmpValue.Float()
			} else if TwoSchemaOutput.Type == "boolean" {
				//直接赋值
				tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
				schemaOutputList[key].Value = tmpValue.Bool()
			} else if TwoSchemaOutput.Type == "list" {
				tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
				schemaOutputList[key].ListSchema.Value = tmpValue.Array()
			} else if TwoSchemaOutput.Type == "object" {
				//递归调用
				tmpValue := gjsonObject.Get(TwoSchemaOutput.Name)
				fmt.Printf("不应该存在第五层。tmpValue：%+v", tmpValue)
			}

			fmt.Printf("%+v", TwoSchemaOutput)
		}
	}

	return schemaOutputList
}
