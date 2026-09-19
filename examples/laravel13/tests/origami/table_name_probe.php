<?php
/**
 * 打印 Eloquent 表名与管理员哈希，定位 db:seed 的 `user)`。
 */
$_SERVER['HTTP_HOST'] = '127.0.0.1';
$_SERVER['SERVER_NAME'] = '127.0.0.1';
$_SERVER['REQUEST_URI'] = '/';
$_SERVER['REQUEST_METHOD'] = 'GET';
require __DIR__.'/../../vendor/autoload.php';
$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

use App\Models\Admin;
use App\Models\User;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Str;

$u = new User();
$rp = new ReflectionProperty(User::class, 'table');
$rp->setAccessible(true);
echo "raw_table=".var_export($rp->getValue($u), true)."\n";
echo "raw_init=".($rp->isInitialized($u) ? 'yes' : 'no')."\n";
echo "get_class=".var_export(get_class($u), true)."\n";
echo "basename_this=".var_export(class_basename($u), true)."\n";
echo "basename_fqcn=".var_export(class_basename(User::class), true)."\n";
echo "computed=".var_export(Str::snake(Str::pluralStudly(class_basename($u))), true)."\n";
echo "user_table=".var_export($u->getTable(), true)."\n";
echo "admin_table=".var_export((new Admin)->getTable(), true)."\n";
echo "pluralStudly_User=".var_export(Str::pluralStudly('User'), true)."\n";
echo "snake_Users=".var_export(Str::snake('Users'), true)."\n";
echo "plural_User=".var_export(Str::plural('User'), true)."\n";

$admin = Admin::where('email', 'admin@example.com')->first();
echo "admin_found=".($admin ? 'yes' : 'no')."\n";
if ($admin) {
    echo "hash_ok=".(Hash::check('password', $admin->password) ? 'yes' : 'no')."\n";
    echo "is_active=".var_export($admin->is_active, true)."\n";
}

class TableName_UserProbe extends User
{
    public function debugTable(): void
    {
        echo "inside_basename=".var_export(class_basename($this), true)."\n";
        echo "inside_computed=".var_export(Str::snake(Str::pluralStudly(class_basename($this))), true)."\n";
        echo "inside_table_prop=".var_export($this->table, true)."\n";
        echo "inside_coalesce=".var_export($this->table ?? Str::snake(Str::pluralStudly(class_basename($this))), true)."\n";
    }
}

$probe = new TableName_UserProbe();
$probe->debugTable();
echo "probe_getTable=".var_export($probe->getTable(), true)."\n";
