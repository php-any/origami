<?php

namespace Database\Seeders;

use App\Models\Admin;
use App\Models\Category;
use App\Models\Order;
use App\Models\Permission;
use App\Models\Product;
use App\Models\Role;
use App\Models\Setting;
use App\Models\User;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Hash;

class DatabaseSeeder extends Seeder
{
    public function run(): void
    {
        $superAdminRole = Role::query()->create([
            'name' => 'super-admin',
            'display_name' => '超级管理员',
            'description' => '拥有所有权限的超级管理员',
        ]);

        $adminRole = Role::query()->create([
            'name' => 'admin',
            'display_name' => '管理员',
            'description' => '管理用户、订单和产品',
        ]);

        $operatorRole = Role::query()->create([
            'name' => 'operator',
            'display_name' => '运营人员',
            'description' => '负责订单和产品运营',
        ]);

        $permissions = [
            ['name' => 'dashboard.view', 'display_name' => '查看仪表盘', 'group' => 'dashboard'],
            ['name' => 'admins.view', 'display_name' => '查看管理员', 'group' => 'admin'],
            ['name' => 'admins.create', 'display_name' => '创建管理员', 'group' => 'admin'],
            ['name' => 'admins.edit', 'display_name' => '编辑管理员', 'group' => 'admin'],
            ['name' => 'admins.delete', 'display_name' => '删除管理员', 'group' => 'admin'],
            ['name' => 'users.view', 'display_name' => '查看用户', 'group' => 'user'],
            ['name' => 'users.create', 'display_name' => '创建用户', 'group' => 'user'],
            ['name' => 'users.edit', 'display_name' => '编辑用户', 'group' => 'user'],
            ['name' => 'users.delete', 'display_name' => '删除用户', 'group' => 'user'],
            ['name' => 'roles.view', 'display_name' => '查看角色', 'group' => 'role'],
            ['name' => 'roles.create', 'display_name' => '创建角色', 'group' => 'role'],
            ['name' => 'roles.edit', 'display_name' => '编辑角色', 'group' => 'role'],
            ['name' => 'roles.delete', 'display_name' => '删除角色', 'group' => 'role'],
            ['name' => 'permissions.view', 'display_name' => '查看权限', 'group' => 'permission'],
            ['name' => 'permissions.create', 'display_name' => '创建权限', 'group' => 'permission'],
            ['name' => 'permissions.edit', 'display_name' => '编辑权限', 'group' => 'permission'],
            ['name' => 'permissions.delete', 'display_name' => '删除权限', 'group' => 'permission'],
            ['name' => 'products.view', 'display_name' => '查看产品', 'group' => 'product'],
            ['name' => 'products.create', 'display_name' => '创建产品', 'group' => 'product'],
            ['name' => 'products.edit', 'display_name' => '编辑产品', 'group' => 'product'],
            ['name' => 'products.delete', 'display_name' => '删除产品', 'group' => 'product'],
            ['name' => 'orders.view', 'display_name' => '查看订单', 'group' => 'order'],
            ['name' => 'orders.edit', 'display_name' => '编辑订单', 'group' => 'order'],
            ['name' => 'categories.view', 'display_name' => '查看分类', 'group' => 'category'],
            ['name' => 'categories.create', 'display_name' => '创建分类', 'group' => 'category'],
            ['name' => 'categories.edit', 'display_name' => '编辑分类', 'group' => 'category'],
            ['name' => 'categories.delete', 'display_name' => '删除分类', 'group' => 'category'],
            ['name' => 'media.view', 'display_name' => '查看媒体', 'group' => 'media'],
            ['name' => 'media.create', 'display_name' => '上传媒体', 'group' => 'media'],
            ['name' => 'media.edit', 'display_name' => '编辑媒体', 'group' => 'media'],
            ['name' => 'media.delete', 'display_name' => '删除媒体', 'group' => 'media'],
            ['name' => 'activity.view', 'display_name' => '查看操作日志', 'group' => 'activity'],
            ['name' => 'settings.edit', 'display_name' => '系统设置', 'group' => 'settings'],
            ['name' => 'notifications.view', 'display_name' => '查看通知', 'group' => 'notifications'],
            ['name' => 'profile.edit', 'display_name' => '编辑个人信息', 'group' => 'profile'],
        ];

        foreach ($permissions as $perm) {
            Permission::query()->create([
                'name' => $perm['name'],
                'display_name' => $perm['display_name'],
                'group' => $perm['group'],
                'description' => $perm['display_name'],
            ]);
        }

        $allPermissionIds = Permission::query()->pluck('id');
        $superAdminRole->permissions()->sync($allPermissionIds);

        $adminPermNames = [
            'dashboard.view',
            'admins.view', 'admins.create', 'admins.edit',
            'users.view', 'users.create', 'users.edit',
            'roles.view',
            'permissions.view',
            'products.view', 'products.create', 'products.edit',
            'orders.view', 'orders.edit',
            'categories.view', 'categories.create', 'categories.edit',
            'media.view', 'media.create',
            'activity.view',
            'notifications.view',
            'profile.edit',
        ];
        $adminRole->permissions()->sync(
            Permission::query()->whereIn('name', $adminPermNames)->pluck('id')
        );

        $operatorPermNames = [
            'dashboard.view',
            'orders.view', 'orders.edit',
            'products.view',
            'categories.view',
            'media.view',
            'profile.edit',
            'notifications.view',
        ];
        $operatorRole->permissions()->sync(
            Permission::query()->whereIn('name', $operatorPermNames)->pluck('id')
        );

        $superAdmin = Admin::query()->create([
            'name' => '超级管理员',
            'email' => 'admin@example.com',
            'password' => Hash::make('password'),
            'is_active' => true,
        ]);
        $superAdmin->roles()->attach($superAdminRole);

        $manager = Admin::query()->create([
            'name' => '普通管理员',
            'email' => 'manager@example.com',
            'password' => Hash::make('password'),
            'is_active' => true,
        ]);
        $manager->roles()->attach($adminRole);

        User::query()->create([
            'name' => '测试用户',
            'email' => 'user@example.com',
            'password' => Hash::make('password'),
        ]);

        for ($i = 1; $i <= 5; $i++) {
            User::query()->create([
                'name' => '用户'.$i,
                'email' => "user{$i}@example.com",
                'password' => Hash::make('password'),
            ]);
        }

        $phones = Category::query()->create([
            'name' => '手机',
            'slug' => 'phones',
            'sort' => 1,
            'is_active' => true,
        ]);
        $computers = Category::query()->create([
            'name' => '电脑配件',
            'slug' => 'computers',
            'sort' => 2,
            'is_active' => true,
        ]);

        $products = [
            ['name' => '苹果 iPhone 15 Pro', 'sku' => 'APL-IP15P', 'description' => '最新款苹果手机', 'price' => 8999.00, 'stock' => 50, 'category_id' => $phones->id],
            ['name' => '华为 Mate 60 Pro', 'sku' => 'HUA-M60P', 'description' => '华为旗舰手机', 'price' => 6999.00, 'stock' => 30, 'category_id' => $phones->id],
            ['name' => '小米 14 Ultra', 'sku' => 'XIA-M14U', 'description' => '小米旗舰手机', 'price' => 5999.00, 'stock' => 80, 'category_id' => $phones->id],
            ['name' => '联想拯救者笔记本', 'sku' => 'LEN-LJZ', 'description' => '游戏笔记本电脑', 'price' => 7999.00, 'stock' => 20, 'category_id' => $computers->id],
            ['name' => '戴尔显示器 27寸', 'sku' => 'DEL-27D', 'description' => '4K 高清显示器', 'price' => 1999.00, 'stock' => 45, 'category_id' => $computers->id],
        ];

        foreach ($products as $product) {
            Product::query()->create([
                ...$product,
                'status' => 'active',
            ]);
        }

        $userIds = User::query()->orderBy('id')->take(3)->pluck('id');
        $productList = Product::query()->orderBy('id')->take(3)->get();
        $statuses = ['pending', 'paid', 'shipped'];

        foreach ($userIds as $i => $userId) {
            $product = $productList[$i % $productList->count()];
            $quantity = $i + 1;
            $subtotal = $product->price * $quantity;

            $order = Order::query()->create([
                'order_no' => 'ORD'.date('YmdHis').str_pad((string) ($i + 1), 4, '0', STR_PAD_LEFT),
                'user_id' => $userId,
                'status' => $statuses[$i],
                'total_amount' => $subtotal,
                'shipping_name' => User::query()->where('id', $userId)->value('name'),
                'shipping_phone' => '1380000'.str_pad((string) ($i + 1), 4, '0', STR_PAD_LEFT),
                'shipping_address' => '北京市朝阳区示例路'.($i + 1).'号',
            ]);

            DB::table('order_items')->insert([
                'order_id' => $order->id,
                'product_id' => $product->id,
                'product_name' => $product->name,
                'price' => $product->price,
                'quantity' => $quantity,
                'subtotal' => $subtotal,
                'created_at' => now(),
                'updated_at' => now(),
            ]);
        }

        Setting::setValue('site_name', 'Origami Admin', 'general', '站点名称');
        Setting::setValue('maintenance_mode', false, 'general', '维护模式');
        Setting::setValue('order_prefix', 'ORD', 'orders', '订单号前缀');
    }
}
