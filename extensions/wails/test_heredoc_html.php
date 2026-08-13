<?php
use Wails\Options\App;

$html = <<<HTML
<title>FROM_HEREDOC_PROBE</title>
<body>heredoc body</body>
HTML;

$app = new App([
    'HTML'  => $html,
    'Title' => 'Probe',
]);

file_put_contents(__DIR__ . '/.html_probe.txt', (string)($app->HTML ?? ''));
