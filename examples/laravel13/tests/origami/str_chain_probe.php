<?php
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use Illuminate\Support\Str;

$expression = '$content, $logo, $isDarkMode = false';
echo "contains=".(Str::contains($expression, ',')?'yes':'no')."\n";

$s = Str::of($expression);
echo "of_class=".get_class($s)."\n";
echo "of_str=".(string)$s."\n";

$t = $s->trim();
echo "trim_str=".(string)$t."\n";

$e = $t->explode(',', 2);
echo "explode_class=".get_class($e)."\n";
echo "explode_all=".json_encode($e->all())."\n";
echo "explode_count=".$e->count()."\n";

$m = $e->map(fn ($part) => trim($part));
echo "map_class=".get_class($m)."\n";
echo "map_all=".json_encode($m->all())."\n";
echo "map_count=".$m->count()."\n";

$a = $m->toArray();
echo "toArray=".json_encode($a)." gettype=".gettype($a)." is_array=".(is_array($a)?'yes':'no')."\n";
echo "count=".count($a)."\n";
foreach ($a as $k=>$v) {
    echo "a[$k]=".json_encode($v)."\n";
}

[$name, $args] = $a;
echo "destruct name=".json_encode($name)." args=".json_encode($args)."\n";

// Direct list from array literal
[$n2, $a2] = ['$content', '$logo, $isDarkMode = false'];
echo "literal destruct=".json_encode($n2).",".json_encode($a2)."\n";
