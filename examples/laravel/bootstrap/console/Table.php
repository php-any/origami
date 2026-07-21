<?php

namespace Bootstrap\Console;

/**
 * 终端表格输出
 */
class Table
{
    private array $headers = [];

    /** @var list<list<string>> */
    private array $rows = [];

    public function __construct(
        private readonly Output $output,
    ) {}

    public function setHeaders(array $headers): self
    {
        $this->headers = array_map('strval', $headers);

        return $this;
    }

    public function addRow(array $row): self
    {
        $this->rows[] = array_map('strval', $row);

        return $this;
    }

    public function render(): void
    {
        if ($this->headers === [] && $this->rows === []) {
            return;
        }

        $columns = count($this->headers);
        if ($columns === 0 && count($this->rows) > 0) {
            $columns = count($this->rows[0]);
        }

        $widths = [];
        for ($i = 0; $i < $columns; $i++) {
            $widths[$i] = 0;
        }

        foreach ($this->headers as $index => $header) {
            $widths[$index] = max($widths[$index], strlen($header));
        }

        foreach ($this->rows as $row) {
            foreach ($row as $index => $cell) {
                if ($index >= $columns) {
                    continue;
                }
                $widths[$index] = max($widths[$index], strlen($cell));
            }
        }

        if ($this->headers !== []) {
            $this->output->writeln('  ' . $this->formatRow($this->headers, $widths));
            $this->output->writeln('  ' . $this->separator($widths));
        }

        foreach ($this->rows as $row) {
            $this->output->writeln('  ' . $this->formatRow($row, $widths));
        }
    }

    private function formatRow(array $row, array $widths): string
    {
        $parts = [];
        foreach ($widths as $index => $width) {
            $value = $row[$index] ?? '';
            $parts[] = str_pad($value, $width);
        }

        return implode('  ', $parts);
    }

    private function separator(array $widths): string
    {
        $parts = [];
        foreach ($widths as $width) {
            $parts[] = str_repeat('-', $width);
        }

        return implode('  ', $parts);
    }
}
