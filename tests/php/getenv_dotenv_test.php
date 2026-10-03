<?php
namespace tests\php;
// Superglobals and the request environment are independent PHP stores.
putenv('ORIGAMI_DOTENV_PROBE');
$_ENV['ORIGAMI_DOTENV_PROBE'] = 'superglobal';
$_SERVER['ORIGAMI_DOTENV_PROBE'] = 'server';
if (getenv('ORIGAMI_DOTENV_PROBE') !== false) { Log::fatal('superglobal leaked into getenv'); }
$path = getenv('PATH');
$_ENV['PATH'] = '/custom/path';
if (getenv('PATH') !== $path) { Log::fatal('PATH overwritten by superglobal'); }
putenv('ORIGAMI_DOTENV_PROBE=environment');
if (getenv('ORIGAMI_DOTENV_PROBE') !== 'environment' || $_ENV['ORIGAMI_DOTENV_PROBE'] !== 'superglobal') { Log::fatal('putenv/getenv separation'); }
putenv('ORIGAMI_DOTENV_PROBE');
Log::info('getenv and superglobal isolation passed');
