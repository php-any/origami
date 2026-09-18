<?php
namespace tests\php;
`$fn = function() { return 1; };
`$a = `$fn ?? "fallback";
\Log::info("coalesce_type=".gettype(`$a)." is_callable=".(is_callable(`$a)?"yes":"no"));
`$b = null ?? `$fn;
\Log::info("null_coalesce_fn_type=".gettype(`$b));
