package kpagents

import (
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/tools"
)

// AgentType agent类型
type AgentType string

const (
	SingleAgentType AgentType = "single_agent_type"
	MultiAgentType  AgentType = "multi_agent_type"
)

// Initialize is a function that creates a new executor with the specified LLM
// model, tools, agent type, and options. It returns an Executor or an error
// if there is any issues during the creation process.
func Initialize(
	llm llms.Model,
	tools []tools.Tool,
	agentType AgentType,
	opts ...KpCreationOption,
) (AgentExecutor, error) {
	var agent agents.Agent
	switch agentType {
	case SingleAgentType:
		agent = NewSingleAgent(llm, opts...)
	case MultiAgentType:
		agent = NewMultiAgent(llm, opts...)
	default:
		return AgentExecutor{}, agents.ErrUnknownAgentType
	}
	return NewAgentExecutor(agent, opts...), nil
}
