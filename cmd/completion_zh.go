package cmd

import "fmt"

// localizeCompletionCmd 将 Cobra 自动注册的 completion 子命令文案改为中文。
func localizeCompletionCmd() {
	rootCmd.InitDefaultCompletionCmd()

	completionCmd, _, err := rootCmd.Find([]string{"completion"})
	if err != nil {
		return
	}

	bin := rootCmd.Name()

	completionCmd.Short = "生成 shell 命令行自动补全脚本"
	completionCmd.Long = fmt.Sprintf(`为 %[1]s 生成 bash/zsh/fish/PowerShell 的 Tab 补全脚本。

配置后，在终端输入 %[1]s 再按 Tab，可自动提示子命令和参数，无需手打完整命令。
详见各子命令帮助中的安装说明。`, bin)

	for _, sub := range completionCmd.Commands() {
		switch sub.Name() {
		case "bash":
			sub.Short = "生成 bash 自动补全脚本"
			sub.Long = fmt.Sprintf(`生成 bash shell 的自动补全脚本。

依赖 bash-completion 包，若未安装可通过系统包管理器安装。

在当前 shell 会话中加载补全：

	source <(%[1]s completion bash)

为每个新会话加载补全（仅需执行一次）：

#### Linux:

	%[1]s completion bash > /etc/bash_completion.d/%[1]s

#### macOS:

	%[1]s completion bash > $(brew --prefix)/etc/bash_completion.d/%[1]s

配置生效需要重新打开 shell。`, bin)
		case "zsh":
			sub.Short = "生成 zsh 自动补全脚本"
			sub.Long = fmt.Sprintf(`生成 zsh shell 的自动补全脚本。

若环境中尚未启用 shell 补全，可先执行：

	echo "autoload -U compinit; compinit" >> ~/.zshrc

在当前 shell 会话中加载补全：

	source <(%[1]s completion zsh)

为每个新会话加载补全（仅需执行一次）：

#### Linux:

	%[1]s completion zsh > "${fpath[1]}/_%[1]s"

#### macOS:

	%[1]s completion zsh > $(brew --prefix)/share/zsh/site-functions/_%[1]s

配置生效需要重新打开 shell。`, bin)
		case "fish":
			sub.Short = "生成 fish 自动补全脚本"
			sub.Long = fmt.Sprintf(`生成 fish shell 的自动补全脚本。

在当前 shell 会话中加载补全：

	%[1]s completion fish | source

为每个新会话加载补全（仅需执行一次）：

	%[1]s completion fish > ~/.config/fish/completions/%[1]s.fish

配置生效需要重新打开 shell。`, bin)
		case "powershell":
			sub.Short = "生成 PowerShell 自动补全脚本"
			sub.Long = fmt.Sprintf(`生成 PowerShell 的自动补全脚本。

在当前 shell 会话中加载补全：

	%[1]s completion powershell | Out-String | Invoke-Expression

为每个新会话加载补全，将上述命令的输出添加到 PowerShell 配置文件中。`, bin)
		}

		if f := sub.Flags().Lookup("no-descriptions"); f != nil {
			f.Usage = "禁用补全说明"
		}
	}
}
