package workflownode

import (
	"fmt"
	"strings"
)

func getResultJson(nodeType string, nodeOutputList []SchemaOutputs) (resultJson string) {
	resultJson = "{"
	for _, outputParam := range nodeOutputList {
		resultJson += `"` + outputParam.Name + `":`
		if outputParam.Type == "string" {
			resultJson += `"` + fmt.Sprintf("%s", outputParam.Value) + `",`
		} else if outputParam.Type == "integer" {
			resultJson += fmt.Sprintf("%d", outputParam.Value) + `,`
		} else if outputParam.Type == "boolean" {
			resultJson += fmt.Sprintf("%t", outputParam.Value) + `,`
		} else if outputParam.Type == "float" {
			resultJson += fmt.Sprintf("%f", outputParam.Value) + `,`
		} else if outputParam.Type == "list" {
			resultJson += fmt.Sprintf("%s", outputParam.Value) + `,`
		} else if outputParam.Type == "object" {
			resultJson += fmt.Sprintf("%s", outputParam.Value) + `,`
		}
	}
	resultJson = strings.TrimRight(resultJson, ",")
	resultJson += "}"
	return
}
