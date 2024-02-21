package kpagents

import (
	"context"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/schema"
)

type MultiAgent struct {
}

func (m *MultiAgent) Plan(ctx context.Context, intermediateSteps []schema.AgentStep, inputs map[string]string) ([]schema.AgentAction, *schema.AgentFinish, error) {
	return nil, nil, nil
}

func (m *MultiAgent) GetInputKeys() []string {
	return nil
}

func (m *MultiAgent) GetOutputKeys() []string {
	return nil
}

func NewMultiAgent(llm llms.Model, opts ...KpCreationOption) *MultiAgent {
	return &MultiAgent{}
}
