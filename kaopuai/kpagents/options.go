package kpagents

import (
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/callbacks"
	"github.com/tmc/langchaingo/kaopuai/kpdatabase"
	"github.com/tmc/langchaingo/kaopuai/kpknowledge"
	"github.com/tmc/langchaingo/kaopuai/kpvariable"
	"github.com/tmc/langchaingo/kaopuai/kpworkflow"
	"github.com/tmc/langchaingo/kaopuai/longtermmemory"
	"github.com/tmc/langchaingo/prompts"
	"github.com/tmc/langchaingo/tools"
)

type KpCreationOptions struct {
	prompt         prompts.PromptTemplate
	tools          []tools.Tool
	workflows      []kpworkflow.KpWorkflow
	knowledge      kpknowledge.Knowledge
	variable       kpvariable.Variable
	database       kpdatabase.Database
	longTermMemory longtermmemory.LongTermMemory

	callbacksHandler callbacks.Handler
	errorHandler     *agents.ParserErrorHandler

	maxIterations           int
	returnIntermediateSteps bool

	// openai
	systemMessage string
	extraMessages []prompts.MessageFormatter
}

type KpCreationOption func(*KpCreationOptions)

func WithPrompt(prompt prompts.PromptTemplate) KpCreationOption {
	return func(co *KpCreationOptions) {
		co.prompt = prompt
	}
}

func WithTools(tools []tools.Tool) KpCreationOption {
	return func(co *KpCreationOptions) {
		co.tools = tools
	}
}

func WithWorkflows(workflows []kpworkflow.KpWorkflow) KpCreationOption {
	return func(co *KpCreationOptions) {
		co.workflows = workflows
	}
}

func WithKnowledge(knowledge kpknowledge.Knowledge) KpCreationOption {
	return func(co *KpCreationOptions) {
		co.knowledge = knowledge
	}
}

func WithVariable(variable kpvariable.Variable) KpCreationOption {
	return func(co *KpCreationOptions) {
		co.variable = variable
	}
}

func WithDatabase(database kpdatabase.Database) KpCreationOption {
	return func(co *KpCreationOptions) {
		co.database = database
	}
}

func WithLongTermMemory(longTermMemory longtermmemory.LongTermMemory) KpCreationOption {
	return func(co *KpCreationOptions) {
		co.longTermMemory = longTermMemory
	}
}

func WithCallbacksHandler(handler callbacks.Handler) KpCreationOption {
	return func(co *KpCreationOptions) {
		co.callbacksHandler = handler
	}
}

func WithParserErrorHandler(errorHandler *agents.ParserErrorHandler) KpCreationOption {
	return func(co *KpCreationOptions) {
		co.errorHandler = errorHandler
	}
}

func WithMaxIterations(iterations int) KpCreationOption {
	return func(o *KpCreationOptions) {
		o.maxIterations = iterations
	}

}

func WithReturnIntermediateSteps(returnIntermediateSteps bool) KpCreationOption {
	return func(o *KpCreationOptions) {
		o.returnIntermediateSteps = returnIntermediateSteps
	}
}
