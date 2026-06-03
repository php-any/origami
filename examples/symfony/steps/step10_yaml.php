<?php

use Symfony\Component\Yaml\Yaml;

// step10: symfony/yaml

$parsed = Yaml::parse("name: origami\nversion: 1\nfeatures:\n  - routing\n  - console\n");
step_check('Yaml::parse name', ($parsed['name'] ?? '') === 'origami');
step_check('Yaml::parse features count', count($parsed['features'] ?? []) === 2);

$dumped = Yaml::dump(['foo' => 'bar', 'num' => 1]);
step_check('Yaml::dump', str_contains($dumped, 'foo: bar'));

step_check('step10_yaml', true, 'symfony/yaml');
