package kpagents

import (
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/kaopuai/kpdatabase"
	"github.com/tmc/langchaingo/kaopuai/kpknowledge"
	"github.com/tmc/langchaingo/kaopuai/kpvariable"
	"github.com/tmc/langchaingo/kaopuai/kpworkflow"
	"github.com/tmc/langchaingo/tools"
)

// 当agent执行器

type AgentExecutor struct {
	Agent     agents.Agent
	Tools     []tools.Tool
	Workflows []kpworkflow.KpWorkflow
	Knowledge kpknowledge.Knowledge
	Variable  kpvariable.Variable
	Database  kpdatabase.Database

	CallbacksHandler callbacks.Handler
	ErrorHandler     *agents.ParserErrorHandler

	MaxIterations           int
	ReturnIntermediateSteps bool
}

func NewAgentExecutor(agent agents.Agent, opts ...KpCreationOption) AgentExecutor {
	return AgentExecutor{}
}
