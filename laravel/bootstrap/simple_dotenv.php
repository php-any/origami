<?php
/**
 * Simple .env loader for Origami - bypasses the vlucas/phpdotenv library
 * which relies on nested closures that Origami doesn't fully support.
 */
function origami_load_dotenv(string $path): void
{
    if (!file_exists($path)) {
        return;
    }

    $content = file_get_contents($path);
    if ($content === false) {
        return;
    }

    $lines = preg_split("/\r\n|\n|\r/", $content);
    if (!is_array($lines)) {
        return;
    }

    foreach ($lines as $line) {
        $line = trim($line);

        // Skip empty lines and comments
        if ($line === '' || $line[0] === '#') {
            continue;
        }

        // Skip lines without '='
        $eqPos = strpos($line, '=');
        if ($eqPos === false) {
            continue;
        }

        $name = trim(substr($line, 0, $eqPos));
        $value = trim(substr($line, $eqPos + 1));

        // Skip export prefix
        if (substr($name, 0, 7) === 'export ') {
            $name = trim(substr($name, 7));
        }

        // Remove surrounding quotes from value
        $len = strlen($value);
        if ($len >= 2) {
            $first = $value[0];
            $last = $value[$len - 1];
            if (($first === '"' && $last === '"') || ($first === "'" && $last === "'")) {
                $value = substr($value, 1, -1);
            }
        }

        // Handle inline comments for unquoted values
        if ($len === 0 || ($value[0] !== '"' && $value[0] !== "'")) {
            $commentPos = strpos($value, ' #');
            if ($commentPos !== false) {
                $value = substr($value, 0, $commentPos);
            }
        }

        // Set via putenv and $_ENV
        putenv($name . '=' . $value);
        $_ENV[$name] = $value;
        $_SERVER[$name] = $value;
    }
}
