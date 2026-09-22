<?php
namespace tests\origami;
$c = collect([1,2,3]);
if ($c->count() !== 3) { \Log::fatal("count"); }
if ($c->first() != 1) { \Log::fatal("first"); }
\Log::info("collection ok");
