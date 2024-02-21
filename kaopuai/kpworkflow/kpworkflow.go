package kpworkflow

import "context"

// KpWorkflow 是llm代理与工作流交互的工具
type KpWorkflow interface {
	Name() string
	Description() string
	Call(ctx context.Context, input string) (string, error)
}
