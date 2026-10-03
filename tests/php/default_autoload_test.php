<?php
function default_loader_assert($ok, $message) { if (!$ok) { throw new Exception($message); } }
$oldPath = ini_get('include_path');
ini_set('include_path', __DIR__ . '/fixtures');
default_loader_assert(spl_autoload_extensions() === '.inc,.php', 'default extensions');
default_loader_assert(spl_autoload_extensions('.php') === '.php', 'extension setter returns current');
default_loader_assert(spl_autoload_register() && spl_autoload_register(null), 'default registration');
default_loader_assert(spl_autoload_functions() === ['spl_autoload'], 'default identity and deduplication');
default_loader_assert(DefaultAutoload\CaseDefaultLoader::value() === 42, 'default lowercase namespace include path');
default_loader_assert(spl_autoload_call('MissingDefaultLoader') === null, 'autoload call returns void');
default_loader_assert(spl_autoload_unregister('spl_autoload'), 'default unregister');
function same_loader_declaration($name) {}
$first = same_loader_declaration(...); $second = same_loader_declaration(...);
spl_autoload_register($first); spl_autoload_register($second); spl_autoload_register($first);
default_loader_assert(count(spl_autoload_functions()) === 2, 'distinct first-class closure identities');
spl_autoload_unregister($first); spl_autoload_unregister($second);
spl_autoload_extensions('.inc,.php');
if ($oldPath !== false) { ini_set('include_path', $oldPath); }
echo "default autoload ok\n";
