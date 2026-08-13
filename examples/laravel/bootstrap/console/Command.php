<?php

namespace Bootstrap\Console;

use Illuminate\Console\Command as IlluminateCommand;
use Illuminate\Console\OutputStyle;
use Symfony\Component\Console\Input\ArgvInput;
use Symfony\Component\Console\Input\InputArgument;
use Symfony\Component\Console\Input\InputOption;
use Symfony\Component\Console\Output\ConsoleOutput;

/**
 * Artisan 命令基类：继承 illuminate/console，保留 Origami #[CliApplication] 无参 execute() 入口。
 *
 * 双轨说明：
 * - Origami：#[Command] 扫描 + ExecuteCommand → 本类 execute()（无参）→ handle()
 * - Illuminate/Symfony：可走 parent::execute(Input, Output)（需已 setLaravel）
 */
abstract class Command extends IlluminateCommand implements CommandInterface
{
    public function __construct()
    {
        if (!isset($this->signature) && empty($this->name)) {
            $this->name = 'bootstrap-command';
        }

        parent::__construct();
        $this->mergeDefineInput();
    }

    /**
     * Origami CLI 入口；传入 Symfony Input/Output 时走父类。
     *
     * @param  \Symfony\Component\Console\Input\InputInterface|null  $input
     * @param  \Symfony\Component\Console\Output\OutputInterface|null  $output
     * @return int
     */
    public function execute($input = null, $output = null)
    {
        if ($input === null) {
            $this->runViaOrigamiCli();

            return self::SUCCESS;
        }

        return parent::execute($input, $output);
    }

    /**
     * 将 $_SERVER['argv'] 绑到本命令定义，注入 OutputStyle 后调用 handle()。
     * 不走 Termwind components / Laravel\Prompts（避免上游 blocker）。
     */
    protected function runViaOrigamiCli(): void
    {
        $this->setLaravel(illuminate_container());

        $input = new ArgvInput(null, $this->getDefinition());
        $style = new OutputStyle($input, new ConsoleOutput());

        $this->input = $input;
        $this->output = $style;

        $this->handle();
    }

    /**
     * 兼容旧命令的 InputDefinition；新命令优先用 $signature。
     */
    protected function defineInput(): InputDefinition
    {
        return new InputDefinition();
    }

    private function mergeDefineInput(): void
    {
        $definition = $this->defineInput();

        foreach ($definition->getArguments() as $name => $meta) {
            if ($this->getDefinition()->hasArgument($name)) {
                continue;
            }

            // 勿用 empty($required)：Origami 下 empty(false) 可能非 PHP 语义
            $mode = ($meta['required'] === true) ? InputArgument::REQUIRED : InputArgument::OPTIONAL;
            $this->getDefinition()->addArgument(new InputArgument(
                $name,
                $mode,
                (string) ($meta['description'] ?? ''),
                $meta['default'] ?? null
            ));
        }

        foreach ($definition->getOptions() as $name => $meta) {
            if ($this->getDefinition()->hasOption($name)) {
                continue;
            }

            $acceptValue = ($meta['accept_value'] === true);
            $mode = $acceptValue ? InputOption::VALUE_OPTIONAL : InputOption::VALUE_NONE;
            // VALUE_NONE 不可传 default（含 false）；Symfony 内部默认为 false
            $default = $acceptValue ? ($meta['default'] ?? null) : null;

            $this->getDefinition()->addOption(new InputOption(
                $name,
                $meta['shortcut'] ?? null,
                $mode,
                (string) ($meta['description'] ?? ''),
                $default
            ));
        }
    }

    /**
     * 标志位是否已启用（对应旧 ArgvInput::hasOption 语义，而非 Symfony hasOption「是否已定义」）。
     */
    protected function optionEnabled(string $name): bool
    {
        $value = $this->option($name);

        return $value !== null && $value !== false && $value !== '';
    }

    /** Symfony 风格双列详情（不依赖 Termwind） */
    protected function twoColumnDetail(string $first, ?string $second = null): void
    {
        if ($second === null) {
            $this->line('  ' . $first);
            return;
        }

        $width = 44;
        $plainLen = strlen($this->stripAnsi($first));
        if ($plainLen >= $width) {
            $this->line('  ' . $first . ' ' . $second);
            return;
        }

        $dots = max(1, $width - $plainLen);
        $this->line('  ' . $first . ' ' . str_repeat('.', $dots) . ' ' . $second);
    }

    /** 命令列表行 */
    protected function describeCommand(string $name, string $description, int $nameWidth, int $indent = 2): void
    {
        $prefix = str_repeat(' ', $indent);
        $padding = max(1, $nameWidth - strlen($name) + 2);
        $this->line($prefix . Style::green($name) . str_repeat(' ', $padding) . $description);
    }

    /** 命名空间标题行 */
    protected function describeNamespace(string $namespace): void
    {
        $this->line(' ' . Style::yellow($namespace));
    }

    private function stripAnsi(string $text): string
    {
        return preg_replace('/\033\[[0-9;]*m/', '', $text) ?? $text;
    }
}
