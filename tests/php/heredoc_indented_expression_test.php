<?php
$name = 'Ada';
$values = [
    <<<TEXT
    Hello $name
      indented
    TEXT,
    <<<'RAW'
    Hello $name
      indented
    RAW,
];
if ($values !== ["Hello Ada\n  indented", 'Hello $name' . "\n  indented"]) {
    throw new Exception('indented heredoc and nowdoc must strip closing indentation');
}
$joined = (<<<TEXT
    left
    TEXT) . ':right';
if ($joined !== 'left:right') {
    throw new Exception('heredoc closing marker must permit following expression tokens');
}
echo "indented heredoc expressions OK\n";
