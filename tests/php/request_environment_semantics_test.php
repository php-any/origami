<?php
$key = 'ORIGAMI_ENVIRONMENT_SEMANTICS';
putenv($key);
$_ENV[$key] = 'array';
$_SERVER[$key] = 'server';
if (getenv($key) !== false) { throw new Exception('getenv reads array'); }
if (!putenv($key . '=value=with=equals') || getenv($key, true) !== 'value=with=equals') { throw new Exception('putenv value'); }
if ($_ENV[$key] !== 'array' || $_SERVER[$key] !== 'server') { throw new Exception('putenv changes arrays'); }
$all = getenv();
if (!is_array($all) || $all[$key] !== 'value=with=equals') { throw new Exception('getenv all'); }
if (!putenv($key) || getenv($key) !== false || isset(getenv()[$key])) { throw new Exception('putenv removal'); }
echo "request environment OK\n";
