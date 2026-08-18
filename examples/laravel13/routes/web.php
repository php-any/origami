<?php

use App\Livewire\Admin\Admins\Form as AdminForm;
use App\Livewire\Admin\Admins\Index as AdminIndex;
use App\Livewire\Admin\Dashboard;
use App\Livewire\Admin\Login;
use App\Livewire\Admin\Orders\Detail as OrderDetail;
use App\Livewire\Admin\Orders\Index as OrderIndex;
use App\Livewire\Admin\Permissions\Form as PermissionForm;
use App\Livewire\Admin\Permissions\Index as PermissionIndex;
use App\Livewire\Admin\Products\Form as ProductForm;
use App\Livewire\Admin\Products\Index as ProductIndex;
use App\Livewire\Admin\Profile;
use App\Livewire\Admin\Roles\Form as RoleForm;
use App\Livewire\Admin\Roles\Index as RoleIndex;
use App\Livewire\Admin\Users\Form as UserForm;
use App\Livewire\Admin\Users\Index as UserIndex;
use Illuminate\Support\Facades\Route;

Route::get('/', function () {
    return view('welcome');
});

Route::get('/origami-health', fn () => 'OK');

// Admin login
Route::get('/login', Login::class)->name('login');
Route::post('/login', Login::class)->name('login.post');

// Admin routes (all protected)
Route::middleware(['admin.auth'])->prefix('admin')->name('admin.')->group(function () {
    // Dashboard
    Route::get('/', Dashboard::class)->name('dashboard');

    // Admin management
    Route::get('/admins', AdminIndex::class)->name('admins.index');
    Route::get('/admins/create', AdminForm::class)->name('admins.create');
    Route::get('/admins/{adminId}/edit', AdminForm::class)->name('admins.edit');

    // User management
    Route::get('/users', UserIndex::class)->name('users.index');
    Route::get('/users/create', UserForm::class)->name('users.create');
    Route::get('/users/{userId}/edit', UserForm::class)->name('users.edit');

    // Role management
    Route::get('/roles', RoleIndex::class)->name('roles.index');
    Route::get('/roles/create', RoleForm::class)->name('roles.create');
    Route::get('/roles/{roleId}/edit', RoleForm::class)->name('roles.edit');

    // Permission management
    Route::get('/permissions', PermissionIndex::class)->name('permissions.index');
    Route::get('/permissions/create', PermissionForm::class)->name('permissions.create');
    Route::get('/permissions/{permissionId}/edit', PermissionForm::class)->name('permissions.edit');

    // Product management
    Route::get('/products', ProductIndex::class)->name('products.index');
    Route::get('/products/create', ProductForm::class)->name('products.create');
    Route::get('/products/{productId}/edit', ProductForm::class)->name('products.edit');

    // Order management
    Route::get('/orders', OrderIndex::class)->name('orders.index');
    Route::get('/orders/{orderId}', OrderDetail::class)->name('orders.detail');

    // Profile
    Route::get('/profile', Profile::class)->name('profile');

    // Logout
    Route::post('/logout', function () {
        auth('admin')->logout();
        request()->session()->invalidate();
        request()->session()->regenerateToken();
        return redirect('/login');
    })->name('logout');
});

// TEMP DIAG
Route::get('/diag', function () {
    $shared = app('view')->getShared();
    $groups = app('router')->getMiddlewareGroups();
    $web = isset($groups['web']) ? implode(',', $groups['web']) : 'NO WEB GROUP';
    $errors = array_key_exists('errors', $shared) ? gettype($shared['errors']) : 'NOT SET';
    $pipeline = class_exists('Illuminate\\View\\Middleware\\ShareErrorsFromSession') ? 'class-exists' : 'no-class';
    return "web=[$web]\nerrors=$errors\nsharePipeline=$pipeline";
});

// TEMP DIAG2
Route::get('/diag2', function () {
    $route = app('router')->getCurrentRoute() ?? request()->route();
    $mw = $route ? implode(',', $route->gatherMiddleware()) : 'NO ROUTE';
    $routeMw = $route ? implode(',', $route->middleware() ?? []) : 'NONE';
    return "routeMw=[$routeMw]\ngatherMw=[$mw]";
});

// TEMP DIAG3
Route::get('/diag3', function () {
    $hasSession = method_exists(request(), 'session') ? (request()->session() ? 'session-obj' : 'session-null') : 'no-session-method';
    $sessErrors = session('errors') === null ? 'null' : 'has-value';
    $mc = app('router')->getMiddlewareGroups();
    $prevent = isset($mc['web']) ? implode(',', $mc['web']) : 'NONE';
    return "session=$hasSession\nsessErrors=$sessErrors\nweb=$prevent";
});
