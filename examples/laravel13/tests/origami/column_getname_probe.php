<?php
/**
 * 复现 Filament tables index.blade.php：$column->getName() 接到 IntValue 3。
 */
$_SERVER['HTTP_HOST'] = '127.0.0.1';
$_SERVER['SERVER_NAME'] = '127.0.0.1';
$_SERVER['REQUEST_URI'] = '/admin/roles';
$_SERVER['REQUEST_METHOD'] = 'GET';

require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use App\Filament\Resources\Roles\Pages\ListRoles;
use Filament\Facades\Filament;
use Illuminate\Support\Facades\Auth;

$dbg = storage_path('origami-debug');
if (!is_dir($dbg)) {
    mkdir($dbg, 0755, true);
}

function dumpv($label, $v)
{
    $t = gettype($v);
    $extra = '';
    if (is_object($v)) {
        $extra = ' class='.get_class($v);
        if (method_exists($v, 'getName')) {
            try {
                $extra .= ' getName='.var_export($v->getName(), true);
            } catch (Throwable $e) {
                $extra .= ' getName_ex='.$e->getMessage();
            }
        }
    } elseif (is_array($v)) {
        $extra = ' count='.count($v).' keys='.json_encode(array_keys($v));
    } else {
        $extra = ' val='.var_export($v, true);
    }
    echo $label.' type='.$t.$extra."\n";
}

function dumpEx(Throwable $e, $depth = 0)
{
    echo str_repeat('  ', $depth).'EX '.$e::class.' msg='.$e->getMessage()."\n";
    echo str_repeat('  ', $depth).'FILE '.$e->getFile().':'.$e->getLine()."\n";
    if ($e->getPrevious()) {
        dumpEx($e->getPrevious(), $depth + 1);
    }
}

Auth::guard('admin')->login(Admin::where('email', 'admin@example.com')->first());
Filament::setCurrentPanel(Filament::getPanel('admin'));

echo "globals_column_isset=".(array_key_exists('column', $GLOBALS) ? 'yes' : 'no')."\n";
if (array_key_exists('column', $GLOBALS)) {
    dumpv('GLOBALS[column]', $GLOBALS['column']);
}

echo "== table columns via ListRoles ==\n";
try {
    $lw = app(ListRoles::class);
    echo "lw_class=".get_class($lw)."\n";
    if (method_exists($lw, 'bootedInteractsWithTable')) {
        $lw->bootedInteractsWithTable();
    }
    if (method_exists($lw, 'mount')) {
        $lw->mount();
    }
    $table = $lw->getTable();
    dumpv('table', $table);
    $cols = $table->getColumns();
    dumpv('getColumns', $cols);
    if (is_array($cols) || $cols instanceof Traversable) {
        $i = 0;
        foreach ($cols as $k => $c) {
            echo "col[$i] key=".var_export($k, true)." ";
            dumpv('c', $c);
            $i++;
            if ($i > 20) {
                break;
            }
        }
    }
    $vis = $table->getVisibleColumns();
    dumpv('getVisibleColumns', $vis);
    $sortable = array_filter(
        $vis,
        fn (\Filament\Tables\Columns\Column $column): bool => $column->isSortable(),
    );
    dumpv('sortable', $sortable);
    foreach ($sortable as $column) {
        echo "sortable_name=".$column->getName()."\n";
    }
} catch (Throwable $e) {
    echo "table_setup_failed\n";
    dumpEx($e);
}

echo "== HTTP /admin/roles ==\n";
$kernel = $app->make(Illuminate\Contracts\Http\Kernel::class);
$paths = ['/admin/roles', '/admin/permissions', '/admin/users', '/admin/admins', '/admin/products', '/admin/categories'];
foreach ($paths as $path) {
    $req = Illuminate\Http\Request::create('http://127.0.0.1:18086'.$path, 'GET');
    $req->setLaravelSession($app['session']->driver());
    try {
        $resp = $kernel->handle($req);
        $body = (string) $resp->getContent();
        $hit = str_contains($body, 'getName') || str_contains($body, '不支持调用函数') || str_contains($body, 'Internal Server Error');
        echo $path.' status='.$resp->getStatusCode().' len='.strlen($body).' err='.($hit ? 'yes' : 'no')."\n";
        if ($hit) {
            if (preg_match('/当前值\([^)]*\)不支持调用函数[^<]{0,80}/u', $body, $m)) {
                echo "  snippet=".$m[0]."\n";
            }
            file_put_contents($dbg.'/column-getname-'.trim($path, '/').'.html', $body);
        }
        $kernel->terminate($req, $resp);
    } catch (Throwable $e) {
        echo $path." threw\n";
        dumpEx($e);
    }
}

echo "DONE\n";
