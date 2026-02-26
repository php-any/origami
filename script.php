<?php
// 临时测试：Symfony/Console 风格标签输出颜色
// 例如 "Test Console Application <info>1.0.0</info>" 中 1.0.0 应以绿色输出到控制台

function formatConsoleOutput($message) {
    $esc = "\033";  // 双引号解析八进制转义，与 chr(27) 等价
    $reset = $esc . '[0m';
    $styles = [
        'info'     => $esc . '[32m',
        'comment'  => $esc . '[33m',
        'error'    => $esc . '[31m',
        'question' => $esc . '[36m',
    ];
    $out = '';
    $i = 0;
    $len = strlen($message);
    while ($i < $len) {
        $open = strpos($message, '<', $i);
        if ($open === false) {
            $out .= substr($message, $i);
            break;
        }
        $out .= substr($message, $i, $open - $i);
        $close = strpos($message, '>', $open);
        if ($close === false) {
            $out .= substr($message, $open);
            break;
        }
        $tag = substr($message, $open + 1, $close - $open - 1);
        $endTag = '</' . $tag . '>';
        $end = strpos($message, $endTag, $close + 1);
        if ($end !== false && isset($styles[$tag])) {
            $text = substr($message, $close + 1, $end - $close - 1);
            $out .= $styles[$tag] . $text . $reset;
            $i = $end + strlen($endTag);
        } else {
            $out .= substr($message, $open, $close - $open + 1);
            $i = $close + 1;
        }
    }
    return $out;
}

$line = 'Test Console Application <info>1.0.0</info>';
echo formatConsoleOutput($line);
echo "\n";
echo "[INFO] Console 颜色输出测试（终端中 1.0.0 应为绿色）\n";
