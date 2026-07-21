<?php

namespace Bootstrap\Console;

/**
 * CLI 格式化输出（类似 Symfony Console OutputStyle）
 */
class Output
{
    public function write(string $message): void
    {
        echo $message;
    }

    public function writeln(string $message = ''): void
    {
        echo $message . "\n";
    }

    public function newLine(int $count = 1): void
    {
        echo str_repeat("\n", max(0, $count));
    }

    public function title(string $message): void
    {
        $this->newLine();
        $this->writeln('  ' . Style::bold($message));
        $this->writeln('  ' . str_repeat('=', min(strlen($message), 60)));
    }

    public function section(string $message): void
    {
        $this->newLine();
        $this->writeln('  ' . Style::bold($message));
        $this->writeln('  ' . str_repeat('-', min(strlen($message), 60)));
    }

    public function info(string $message): void
    {
        $this->writeln('  ' . Style::green('[info]') . ' ' . $message);
    }

    public function success(string $message): void
    {
        $this->writeln('  ' . Style::green('[OK]') . ' ' . $message);
    }

    public function warning(string $message): void
    {
        $this->writeln('  ' . Style::yellow('[warn]') . ' ' . $message);
    }

    public function error(string $message): void
    {
        fwrite(STDERR, '  ' . Style::red('[error]') . ' ' . $message . "\n");
    }

    public function comment(string $message): void
    {
        $this->writeln('  ' . Style::dim($message));
    }

    /**
     * Symfony Console 风格双列详情（标签 .... 值）
     */
    public function twoColumnDetail(string $first, ?string $second = null): void
    {
        if ($second === null) {
            $this->writeln('  ' . $first);
            return;
        }

        $width = 44;
        $plainLen = strlen($this->stripAnsi($first));
        if ($plainLen >= $width) {
            $this->writeln('  ' . $first . ' ' . $second);
            return;
        }

        $dots = max(1, $width - $plainLen);
        $this->writeln('  ' . $first . ' ' . str_repeat('.', $dots) . ' ' . $second);
    }

    /**
     * 命令列表行：命令名 + 描述（Symfony list 风格）
     */
    public function describeCommand(string $name, string $description, int $nameWidth, int $indent = 2): void
    {
        $prefix = str_repeat(' ', $indent);
        $padding = max(1, $nameWidth - strlen($name) + 2);
        $this->writeln($prefix . Style::green($name) . str_repeat(' ', $padding) . $description);
    }

    /**
     * 命名空间标题行（Symfony list 风格）
     */
    public function describeNamespace(string $namespace): void
    {
        $this->writeln(' ' . Style::yellow($namespace));
    }

    private function stripAnsi(string $text): string
    {
        return preg_replace('/\033\[[0-9;]*m/', '', $text) ?? $text;
    }

    public function listing(array $items): void
    {
        foreach ($items as $item) {
            $this->writeln('  * ' . $item);
        }
    }
}
