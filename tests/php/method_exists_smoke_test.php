<?php
namespace tests\php;
class ME_Host { public function foo() {} }
`$o = new ME_Host();
\Log::info("foo=".(method_exists(`$o,"foo")?"y":"n"));
\Log::info("bar=".(method_exists(`$o,"bar")?"y":"n"));
\Log::info("rendered=".(method_exists(`$o,"rendered")?"y":"n"));
