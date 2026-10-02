<?php
function backtraceFrame(): array { return debug_backtrace(DEBUG_BACKTRACE_IGNORE_ARGS, 1)[0]; }
$frame = backtraceFrame();
if (!is_array($frame) || !is_string($frame['function'])) { throw new Exception('backtrace frame is not an array'); }
echo "backtrace frame array: PASS\n";
