<?php
function andi($i, $j)
{
    for ($k=$i ; $k<=$j ; $k++) {
        if ($k >5) continue;
        echo "$k\n";
    }
}
andi (3,10);
