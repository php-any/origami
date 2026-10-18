<?php
declare(ticks=1);
goto done;
throw new Exception('forward goto did not skip statement');
done:
$count = 0;
again:
$count++;
if ($count < 3) { goto again; }
function tickGotoFunction(): int {
    goto end;
    return -1;
    end:
    return 7;
}
if ($count !== 3 || tickGotoFunction() !== 7) { throw new Exception('tick-wrapped goto target'); }
echo "declare ticks goto OK\n";
