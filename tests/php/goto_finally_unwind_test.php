<?php
function gotoFinallyTrace(): array {
    $trace = [];
    try {
        try {
            $trace[] = 'body';
            goto finished;
        } finally {
            $trace[] = 'inner';
        }
    } finally {
        $trace[] = 'outer';
    }
    $trace[] = 'unreachable';
    finished:
    $trace[] = 'target';
    return $trace;
}
if (gotoFinallyTrace() !== ['body', 'inner', 'outer', 'target']) {
    throw new Exception('goto must unwind finally blocks before reaching target');
}
echo "goto finally unwind OK\n";
