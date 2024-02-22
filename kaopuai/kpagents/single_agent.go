package kpagents

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

// SingleAgent 实现了agents.Agent接口
type SingleAgent struct {
}

func (s *SingleAgent) Plan(ctx context.Context, intermediateSteps []schema.AgentStep, inputs map[string]string) ([]schema.AgentAction, *schema.AgentFinish, error) {
	return nil, nil, nil
}

func (s *SingleAgent) GetInputKeys() []string {
	return nil
}

func (s *SingleAgent) GetOutputKeys() []string {
	return nil
}

func NewSingleAgent(llm llms.Model, opts ...KpCreationOption) *SingleAgent {
	return &SingleAgent{}
}
