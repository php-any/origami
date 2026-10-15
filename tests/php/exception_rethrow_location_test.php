<?php
function rethrow_location($error) { throw $error; }

$explicit = new ErrorException('warning', 0, E_WARNING, 'original-warning.php', 123);
try {
    rethrow_location($explicit);
} catch (Throwable $error) {
    if ($error !== $explicit || $error->getFile() !== 'original-warning.php' || $error->getLine() !== 123) {
        throw new Exception('throw must preserve explicitly supplied exception location');
    }
    try { rethrow_location($error); } catch (Throwable $again) {
        if ($again !== $explicit || $again->getFile() !== 'original-warning.php' || $again->getLine() !== 123) {
            throw new Exception('rethrow must preserve the original exception location');
        }
    }
}

try {
    $firstLine = __LINE__ + 1;
    throw new RuntimeException('first throw');
} catch (Throwable $original) {
    $file = $original->getFile();
    $line = $original->getLine();
    $trace = $original->getTrace();
    if ($file !== __FILE__ || $line !== $firstLine) { throw new Exception('first throw location: ' . json_encode([$file, __FILE__, $line, $firstLine])); }
    try { rethrow_location($original); } catch (Throwable $again) {
        if ($again !== $original || $again->getFile() !== $file || $again->getLine() !== $line || $again->getTrace() !== $trace) {
            throw new Exception('rethrow must retain file, line, trace and object identity');
        }
    }
}
echo "exception rethrow location: PASS\n";
