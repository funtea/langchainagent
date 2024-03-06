package basenode

import (
	"testing"
)

func TestLLM(t *testing.T) {
	a := map[string]any{
		"name": "小明",
	}
	_, err := runLLM("现在你的名字叫{{.name}}，请问你的名字叫什么？", "", 0.7, []string{"name"}, a)
	if err != nil {
		return
	}
	return
}
