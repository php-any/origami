package openai

import "github.com/php-any/origami/data"

// Load 将 OpenAI 扩展注册到 VM 中
func Load(vm data.VM) {
	// 主客户端
	vm.AddClass(NewClientClass())

	// 返回结果类型
	vm.AddClass(NewChatCompletionClass())
	vm.AddClass(NewEmbeddingResultClass())
	vm.AddClass(NewImageResultClass())
	vm.AddClass(NewTranscriptionResultClass())

	// 消息类型（对应 openai-go 的 SystemMessage/UserMessage/... 构造函数）
	vm.AddClass(NewSystemMessageClass())
	vm.AddClass(NewUserMessageClass())
	vm.AddClass(NewAssistantMessageClass())
	vm.AddClass(NewDeveloperMessageClass())
	vm.AddClass(NewToolMessageClass())

	// Chat 选项键枚举
	vm.AddClass(NewChatOptionClass())
}
