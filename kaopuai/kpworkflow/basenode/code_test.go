package basenode

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCode(t *testing.T) {
	//python()
	//testgjson()
	nodeJs()
	return
	//return
	var todoCode = "function main(params) {\n    const ret = {\n        \"key0\": params.input1,\n        \"key1\": \"hi\",\n        \"key2\": [\"hello\", \"world\"],\n        \"key3\": {\n            \"key31\": \"hi\"\n        },\n        \"key4\": [{\n            \"key41\": true,\n            \"key42\": 1,\n            \"key43\": 12.88,\n            \"key44\": [\"hello\"],\n            \"key45\": {\n                \"key451\": \"hello\"\n            }\n        }]\n    };\n\n    return ret;\n} main(pa)"

	var inputParamList []SchemaInputParameters
	err := json.Unmarshal([]byte(paramJson), &inputParamList)
	if err != nil {
		return
	}

	var tmpCodeVariable = "var pa = {"
	for _, inputParam := range inputParamList {
		tmpCodeVariable += `"` + inputParam.Name + `":`
		if inputParam.Input.Type == "string" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += `"` + fmt.Sprintf("%s", inputParam.Input.Value.Content.Value) + `",`
			} else {
				tmpCodeVariable += `"` + fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `",`
			}
		} else if inputParam.Input.Type == "inter" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.Content.Value) + `,`
			} else {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `,`
			}
		} else if inputParam.Input.Type == "boolean" {
			if inputParam.Input.Value.Type == "ref" {
				tmpCodeVariable += fmt.Sprintf("%t", inputParam.Input.Value.Content.Value) + `,`
			} else {
				tmpCodeVariable += fmt.Sprintf("%s", inputParam.Input.Value.LiteralContent) + `,`
			}
		}
	}
	tmpCodeVariable = strings.TrimRight(tmpCodeVariable, ",")
	tmpCodeVariable += "};"

	todoCode = tmpCodeVariable + todoCode
	//V8gotest(todoCode, tmpCodeVariable)

	v8test2()
}

var paramJson = "[{\"name\":\"input1\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\",\"value\":\"123\"}}}},{\"name\":\"input2\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"3232\"}}},{\"name\":\"input3\",\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\",\"value\":true}}}}]"

//var paramJson = "{\"params\":{\"input1\":\"naqxugqevc\",\"input2\":\"chjqwf\",\"input3\":false,\"abc\":5753}}"

func v8test2() {
	//var nodeOutputMap = make(map[string]SchemaOutputs)
	jsonResult := "{\"key0\":\"fp1\",\"key1\":\"fp2\",\"key2\":[\"hello\",\"world\"],\"key3\":{\"key31\":\"hi\"},\"key4\":[{\"key41\":true,\"key42\":1,\"key43\":12.88,\"key44\":[\"hello\"],\"key45\":{\"key451\":{\"key4511\":[1,2,3]}}},{\"key51\":true,\"key52\":1,\"key53\":12.88,\"key54\":[\"hello\"],\"key55\":{\"key551\":{\"key5511\":[1,2,3]}}}]}"
	schema := parseSchema()
	schemaOutputList := schema.Nodes[3].Data.Outputs
	fmt.Printf("%+v", schemaOutputList)

	if len(schemaOutputList) == 0 {
		return
	}
	gjsonObject := gjson.Parse(jsonResult)

	//var outputParamType = "string,integer,float,boolean,"
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

	return
}

var schemaJsonExample1 string = "{\"nodes\":[{\"id\":\"100001\",\"type\":\"1\",\"meta\":{\"position\":{\"x\":-819.3697963254239,\"y\":-411.65624440949136}},\"data\":{\"outputs\":[{\"type\":\"string\",\"name\":\"a\",\"required\":true,\"description\":\"参数a\"},{\"type\":\"integer\",\"name\":\"b\",\"required\":true,\"description\":\"参数b\"}],\"nodeMeta\":{\"title\":\"Start\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Start.png\",\"description\":\"The starting node of the workflow, used to set the information needed to initiate the workflow.\",\"subTitle\":\"\"}}},{\"id\":\"900001\",\"type\":\"2\",\"meta\":{\"position\":{\"x\":1085.4634805630683,\"y\":-301.1839790394442}},\"data\":{\"nodeMeta\":{\"title\":\"End\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-End.png\",\"description\":\"The final node of the workflow, used to return the result information after the workflow runs.\",\"subTitle\":\"\"},\"inputs\":{\"terminatePlan\":\"returnVariables\",\"inputParameters\":[{\"name\":\"c\",\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}}],\"content\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"我的contetn\"}}}}},{\"id\":\"117411\",\"type\":\"11\",\"meta\":{\"position\":{\"x\":-240.5690339275338,\"y\":-303.14190432755527}},\"data\":{\"nodeMeta\":{\"title\":\"Variable\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Variable.png\",\"description\":\"Used for reading and writing variables in your bot. The variable name must match the variable name in Bot.\",\"subTitle\":\"Variable\"},\"inputs\":{\"mode\":\"set\",\"inputParameters\":[{\"name\":\"botVariable\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"bbb\"}}}]},\"outputs\":[{\"type\":\"boolean\",\"name\":\"isSuccess\"}]}},{\"id\":\"120710\",\"type\":\"5\",\"meta\":{\"position\":{\"x\":414.02340143315394,\"y\":-712.5092562838116}},\"data\":{\"nodeMeta\":{\"title\":\"Code\",\"icon\":\"https://sf16-va.tiktokcdn.com/obj/eden-va2/dvsmryvd_avi_dvsm/ljhwZthlaukjlkulzlp/icon/icon-Code.png\",\"description\":\"Write code to process input variables to generate return values.\",\"subTitle\":\"Code\"},\"inputs\":{\"inputParameters\":[{\"name\":\"input1\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"a\"}}}},{\"name\":\"input2\",\"input\":{\"type\":\"string\",\"value\":{\"type\":\"literal\",\"literalContent\":\"3232\"}}},{\"name\":\"input3\",\"input\":{\"type\":\"boolean\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"117411\",\"name\":\"isSuccess\"}}}},{\"name\":\"abc\",\"input\":{\"type\":\"integer\",\"value\":{\"type\":\"ref\",\"content\":{\"source\":\"block-output\",\"blockID\":\"100001\",\"name\":\"b\"}}}}],\"code\":\"async function main({ params }: Args): Promise {\\n    const ret = {\\n        \\\"key0\\\": params.input1 + params.input2,\\n        \\\"key1\\\": \\\"\\\",\\n        \\\"key2\\\": [\\\"hello\\\", \\\"world\\\"],\\n        \\\"key3\\\": {\\n            \\\"key31\\\": \\\"hi\\\"\\n        },\\n        \\\"key4\\\": [{\\n            \\\"key41\\\": true,\\n            \\\"key42\\\": 1,\\n            \\\"key43\\\": 12.88,\\n            \\\"key44\\\": [\\\"hello\\\"],\\n            \\\"key45\\\": {\\n                \\\"key451\\\": {\\n                    \\\"key4511\\\":[1,2,3]\\n                }\\n            }\\n        }]\\n    };\\n\\n    return ret;\\n}\",\"language\":5},\"outputs\":[{\"type\":\"string\",\"name\":\"key0\"},{\"type\":\"string\",\"name\":\"key1\"},{\"type\":\"list\",\"name\":\"key2\",\"listSchema\":{\"type\":\"string\"}},{\"type\":\"object\",\"name\":\"key3\",\"objectSchema\":[{\"type\":\"string\",\"name\":\"key31\"}]},{\"type\":\"list\",\"name\":\"key4\",\"listSchema\":{\"type\":\"object\",\"objectSchema\":[{\"type\":\"boolean\",\"name\":\"key41\"},{\"type\":\"integer\",\"name\":\"key42\"},{\"type\":\"float\",\"name\":\"key43\"},{\"type\":\"list\",\"name\":\"key44\",\"listSchema\":{\"type\":\"string\"}},{\"type\":\"object\",\"name\":\"key45\",\"objectSchema\":[{\"type\":\"object\",\"name\":\"key451\",\"objectSchema\":[{\"type\":\"list\",\"name\":\"key4511\",\"listSchema\":{\"type\":\"integer\"}}]}]}]}}]}}],\"edges\":[{\"sourceNodeID\":\"100001\",\"targetNodeID\":\"117411\"},{\"sourceNodeID\":\"117411\",\"targetNodeID\":\"120710\"},{\"sourceNodeID\":\"120710\",\"targetNodeID\":\"900001\"}]}"

func parseSchema() (schema *Schema) {
	err := json.Unmarshal([]byte(schemaJsonExample1), &schema)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	return
}

func testgjson() {
	var a = "{\n  \"programmers\": [\n    {\n      \"firstName\": \"Janet\", \n      \"lastName\": \"McLaughlin\", \n    }, {\n      \"firstName\": \"Elliotte\", \n      \"lastName\": \"Hunter\", \n    }, {\n      \"firstName\": \"Jason\", \n      \"lastName\": \"Harold\", \n    }\n  ]\n}"

	//result := gjson.Get(a, "programmers.#.lastName")
	//for _, name := range result.Array() {
	//	println(name.String())
	//}

	name := gjson.Parse(a).Get(`programmers.0.firstName`)
	println(name.String()) // prints "Elliotte"
	return
}

func python() {
	pythonCode := "from typing import Dict, Any, List\n"
	pythonCode += "class Args:\n"
	pythonCode += "    def __init__(self, params: Dict[str, Any]):\n"
	pythonCode += "        self.params = params\n"
	pythonCode += "class Output:\n"
	pythonCode += "    def __init__(self, data: Dict[str, Any]):\n"
	pythonCode += "        self.data = data\n"

	// 定义一个Python代码字符串
	pythonCode += "def main(args: Args) -> Output:\n    params = args.params\n    ret: Output = {\n        \"key0\": params['input1'] + params['input2'],\n        \"key1\": \"hi\",\n        \"key2\": [\"hello\", \"world\"],\n        \"key3\": {\n            \"key31\": \"hi\"\n        },\n        \"key4\": [{\n            \"key41\": True,\n            \"key42\": 1,\n            \"key43\": 12.88,\n            \"key44\": [\"hello\"],\n            \"key45\": {\n                \"key451\": {\n                    \"key4511\":[1,2,3]\n                }\n            }\n        },{\n            \"key41\": True,\n            \"key42\": 1,\n            \"key43\": 12.88,\n            \"key44\": [\"hello\"],\n            \"key45\": {\n                \"key451\": {\n                    \"key4511\":[1,2,3]\n                }\n            }\n        }]\n    }\n    return ret\n"

	pythonCode += "params = {'input1': 'Hello from ','input2': 'Go!'}\n"
	pythonCode += "args = Args(params)\n"
	pythonCode += "result = main(args)\n"
	pythonCode += "print(result)\n"
	cmd := exec.Command("python3", "-c", pythonCode) // 确保python3命令可用，或者根据你的环境修改为python或python2等
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}
	fmt.Println("Python script output:", string(output))
	return
}

func nodeJs() {
	// Define the input data for the TypeScript function.
	//params := Args{
	//	Input1: "Hello",
	//	Input2: 42,
	//}
	params := "{\"input1\": \"Hello\",\"input2\": 42}"

	// Convert the input data to JSON format.
	jsonParams, err := json.Marshal(params)
	if err != nil {
		log.Fatalf("Failed to marshal params: %v", err)
	}

	// Construct the TypeScript code with the JSON data injected.
	tsCode := `  
async function main({ params }: Args): Promise<any> {  
	const ret = {  
		"key0": params.input1 + params.input2,  
		"key1": "hi",  
		"key2": ["hello", "world"],  
		"key3": {  
			"key31": "hi"  
		},  
		"key4": [{  
			"key41": true,  
			"key42": 1,  
			"key43": 12.88,  
			"key44": ["hello"],  
			"key45": {  
				"key451": {  
					"key4511":[7,8,9]  
				}  
			}  
		}]  
	};  
  
	return ret;  
}  
  
// A placeholder type for the function parameters, used for TypeScript type checking.  
interface Args {  
	input1: string;  
	input2: number;  
}  
  
// Call the main function with the provided parameters and log the result.  
main({ params: JSON.parse(%s) })  
	.then(result => console.log(result))  
	.catch(error => console.error(error));  
	`

	// Inject the JSON params into the TypeScript code template.
	tsCodeWithParams := fmt.Sprintf(tsCode, jsonParams)

	// Write the TypeScript code to a temporary file.
	tmpFile, err := ioutil.TempFile("", "typescript-program.ts")
	if err != nil {
		log.Fatalf("Failed to create temporary TypeScript file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up the temporary file when done.

	if _, err := tmpFile.WriteString(tsCodeWithParams); err != nil {
		log.Fatalf("Failed to write TypeScript code to file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		log.Fatalf("Failed to close temporary TypeScript file: %v", err)
	}

	// Compile the TypeScript code to JavaScript.
	tscCmd := exec.Command("tsc", tmpFile.Name())
	tscCmd.Stdout = os.Stdout
	tscCmd.Stderr = os.Stderr
	if err := tscCmd.Run(); err != nil {
		log.Fatalf("Failed to compile TypeScript: %v", err)
		return
	}

	// Execute the generated JavaScript code.
	jsFile := tmpFile.Name()[:len(tmpFile.Name())-3] + "js" // Replace ".ts" with ".js".
	nodeCmd := exec.Command("node", jsFile)
	nodeCmd.Stdout = os.Stdout
	nodeCmd.Stderr = os.Stderr
	//if err := nodeCmd.Output()
	fmt.Printf("%+v", nodeCmd)
	return
}
