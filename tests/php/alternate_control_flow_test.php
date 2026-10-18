<?php
$trace = [];
for ($i = 0; $i < 3; $i++):
    if ($i === 0):
        $trace[] = 'zero';
    elseif ($i === 1):
        $trace[] = 'one';
    else:
        $trace[] = 'two';
    endif;
endfor;
$j = 0;
while ($j < 2):
    foreach (['a' => 3, 'b' => 4] as $key => $value):
        switch ($key):
            case 'a': $trace[] = $value + $j; break;
            default: $trace[] = $key; break;
        endswitch;
    endforeach;
    $j++;
endwhile;
if ($trace !== ['zero', 'one', 'two', 3, 'b', 4, 'b']) {
    throw new Exception('nested alternate control syntax must preserve order and branches');
}
echo "alternate control flow OK\n";
