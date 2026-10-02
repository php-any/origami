<?php

$app = require __DIR__.'/runtime_bootstrap.php';
$app->instance('request', Illuminate\Http\Request::create('/admin'));

use App\Filament\Pages\Dashboard;
use App\Models\Admin;
use Filament\Facades\Filament;

class WorkbenchAccessAdmin extends Admin
{
    public array $allowedPermissions = [];

    public function hasPermission(string $permission): bool
    {
        return in_array($permission, $this->allowedPermissions, true);
    }
}

class WorkbenchAccessPage extends Dashboard
{
    public function viewDataForTest(): array
    {
        return $this->getViewData();
    }
}

Filament::setCurrentPanel(Filament::getPanel('admin'));
$admin = new WorkbenchAccessAdmin;
auth('admin')->setUser($admin);
$page = new WorkbenchAccessPage;

$data = $page->viewDataForTest();
if ($data['metrics'] !== [] || $data['shortcuts'] !== [] || $data['orderStages'] !== [] || $data['lowStockProducts'] !== [] || $data['latestOrders'] !== []) {
    throw new RuntimeException('Workbench leaked business data without permissions');
}

$admin->allowedPermissions = ['orders.view'];
$data = $page->viewDataForTest();
if (count($data['metrics']) !== 2 || count($data['shortcuts']) !== 1 || $data['canViewProducts'] || count($data['orderStages']) !== 5) {
    throw new RuntimeException('Order access exposed other business modules');
}
foreach ($data['orderStages'] as $stage) {
    if (! str_contains($stage['url'], 'tab='.$stage['status'])) {
        throw new RuntimeException('Order stage does not use the native tab URL alias');
    }
}

$admin->allowedPermissions = ['products.view'];
$data = $page->viewDataForTest();
if (count($data['metrics']) !== 1 || count($data['shortcuts']) !== 1 || $data['canViewOrders'] || $data['orderStages'] !== [] || $data['latestOrders'] !== []) {
    throw new RuntimeException('Product access exposed order or customer data');
}
if (str_contains($data['shortcuts'][0]['url'], '/create')) {
    throw new RuntimeException('Read-only product access received a create shortcut');
}

echo "workbench access: PASS\n";
return true;
