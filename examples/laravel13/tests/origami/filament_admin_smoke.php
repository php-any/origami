<?php

/**
 * Filament v5 管理后台冒烟测试（轻量：类存在 + 表/种子数据，避免触发已知 Collection 递归缺口）。
 */

require __DIR__.'/../../vendor/autoload.php';

$app = require __DIR__.'/../../bootstrap/app.php';
$app->make(Illuminate\Contracts\Console\Kernel::class)->bootstrap();

$checks = [];

$checks['AdminPanelProvider'] = class_exists(App\Providers\Filament\AdminPanelProvider::class);
$checks['AdminResource'] = class_exists(App\Filament\Resources\Admins\AdminResource::class);
$checks['UserResource'] = class_exists(App\Filament\Resources\Users\UserResource::class);
$checks['RoleResource'] = class_exists(App\Filament\Resources\Roles\RoleResource::class);
$checks['PermissionResource'] = class_exists(App\Filament\Resources\Permissions\PermissionResource::class);
$checks['ProductResource'] = class_exists(App\Filament\Resources\Products\ProductResource::class);
$checks['OrderResource'] = class_exists(App\Filament\Resources\Orders\OrderResource::class);
$checks['CategoryResource'] = class_exists(App\Filament\Resources\Categories\CategoryResource::class);
$checks['MediaResource'] = class_exists(App\Filament\Resources\Media\MediaResource::class);
$checks['ActivityResource'] = class_exists(App\Filament\Resources\Activities\ActivityResource::class);
$checks['ManageSettingsPage'] = class_exists(App\Filament\Pages\ManageSettings::class);
$checks['StatsOverviewWidget'] = class_exists(App\Filament\Widgets\StatsOverviewWidget::class);
$checks['LatestOrdersWidget'] = class_exists(App\Filament\Widgets\LatestOrdersWidget::class);
$checks['LowStockNotification'] = class_exists(App\Notifications\LowStockNotification::class);
$checks['NewOrderNotification'] = class_exists(App\Notifications\NewOrderNotification::class);
$checks['AdminPolicy'] = class_exists(App\Policies\AdminPolicy::class);
$checks['SettingPolicy'] = class_exists(App\Policies\SettingPolicy::class);

$admin = Illuminate\Support\Facades\DB::table('admins')->where('email', 'admin@example.com')->first();
$checks['admin_seeded'] = $admin !== null;
$checks['admin_is_active'] = $admin !== null && (int) $admin->is_active === 1;

$manager = Illuminate\Support\Facades\DB::table('admins')->where('email', 'manager@example.com')->first();
$checks['manager_seeded'] = $manager !== null;

$checks['categories_seeded'] = Illuminate\Support\Facades\DB::table('categories')->count() >= 1;
$checks['products_seeded'] = Illuminate\Support\Facades\DB::table('products')->count() >= 1;
$checks['orders_seeded'] = Illuminate\Support\Facades\DB::table('orders')->count() >= 1;
$checks['settings_seeded'] = Illuminate\Support\Facades\DB::table('settings')->where('key', 'site_name')->exists();
$checks['permissions_seeded'] = Illuminate\Support\Facades\DB::table('permissions')->where('name', 'settings.edit')->exists();
$checks['categories_perm_seeded'] = Illuminate\Support\Facades\DB::table('permissions')->where('name', 'categories.view')->exists();
$checks['media_perm_seeded'] = Illuminate\Support\Facades\DB::table('permissions')->where('name', 'media.view')->exists();
$checks['activity_perm_seeded'] = Illuminate\Support\Facades\DB::table('permissions')->where('name', 'activity.view')->exists();

$pending = Illuminate\Support\Facades\DB::table('orders')->where('status', 'pending')->first();
$checks['pending_order_exists'] = $pending !== null;

$orderModel = new App\Models\Order();
$orderModel->status = 'pending';
$checks['order_transition_pending_to_paid'] = $orderModel->canTransitionTo('paid');
$checks['order_transition_pending_to_cancelled'] = $orderModel->canTransitionTo('cancelled');
$checks['order_transition_pending_not_completed'] = ! $orderModel->canTransitionTo('completed');

$schema = Illuminate\Support\Facades\Schema::class;
$checks['categories_table'] = $schema::hasTable('categories');
$checks['media_table'] = $schema::hasTable('media');
$checks['settings_table'] = $schema::hasTable('settings');
$checks['notifications_table'] = $schema::hasTable('notifications');
$checks['activity_log_table'] = $schema::hasTable('activity_log');
$checks['products_have_category_id'] = $schema::hasColumn('products', 'category_id');
$checks['products_have_image'] = $schema::hasColumn('products', 'image');

$failed = array_filter($checks, fn (bool $ok): bool => ! $ok);

foreach ($checks as $name => $ok) {
    echo ($ok ? '[OK] ' : '[FAIL] ').$name.PHP_EOL;
}

if ($failed !== []) {
    echo PHP_EOL.'FAILED: '.count($failed).' check(s)'.PHP_EOL;
    exit(1);
}

echo PHP_EOL.'All Filament admin smoke checks passed.'.PHP_EOL;
