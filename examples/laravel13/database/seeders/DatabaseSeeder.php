<?php

namespace Database\Seeders;

use App\Models\Admin;
use App\Models\Permission;
use App\Models\Product;
use App\Models\Role;
use App\Models\User;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\DB;

class DatabaseSeeder extends Seeder
{
    public function run(): void
    {
        // Create default roles
        $superAdminRoleId = DB::table('roles')->insertGetId([
            'name' => 'super-admin',
            'display_name' => '超级管理员',
            'description' => '拥有所有权限的超级管理员',
            'created_at' => now(),
            'updated_at' => now(),
        ]);

        $adminRoleId = DB::table('roles')->insertGetId([
            'name' => 'admin',
            'display_name' => '管理员',
            'description' => '管理用户、订单和产品',
            'created_at' => now(),
            'updated_at' => now(),
        ]);

        $operatorRoleId = DB::table('roles')->insertGetId([
            'name' => 'operator',
            'display_name' => '运营人员',
            'description' => '负责订单和产品运营',
            'created_at' => now(),
            'updated_at' => now(),
        ]);

        // Create permissions
        $permissions = [
            // Dashboard
            ['name' => 'dashboard.view', 'display_name' => '查看仪表盘', 'group' => 'dashboard', 'description' => '查看后台仪表盘'],

            // Admin management
            ['name' => 'admins.view', 'display_name' => '查看管理员', 'group' => 'admin', 'description' => '查看管理员列表'],
            ['name' => 'admins.create', 'display_name' => '创建管理员', 'group' => 'admin', 'description' => '创建新的管理员账号'],
            ['name' => 'admins.edit', 'display_name' => '编辑管理员', 'group' => 'admin', 'description' => '编辑管理员信息'],
            ['name' => 'admins.delete', 'display_name' => '删除管理员', 'group' => 'admin', 'description' => '删除管理员账号'],

            // User management
            ['name' => 'users.view', 'display_name' => '查看用户', 'group' => 'user', 'description' => '查看用户列表'],
            ['name' => 'users.edit', 'display_name' => '编辑用户', 'group' => 'user', 'description' => '编辑用户信息'],
            ['name' => 'users.delete', 'display_name' => '删除用户', 'group' => 'user', 'description' => '删除用户账号'],

            // Role management
            ['name' => 'roles.view', 'display_name' => '查看角色', 'group' => 'role', 'description' => '查看角色列表'],
            ['name' => 'roles.create', 'display_name' => '创建角色', 'group' => 'role', 'description' => '创建新的角色'],
            ['name' => 'roles.edit', 'display_name' => '编辑角色', 'group' => 'role', 'description' => '编辑角色信息'],
            ['name' => 'roles.delete', 'display_name' => '删除角色', 'group' => 'role', 'description' => '删除角色'],

            // Permission management
            ['name' => 'permissions.view', 'display_name' => '查看权限', 'group' => 'permission', 'description' => '查看权限列表'],
            ['name' => 'permissions.create', 'display_name' => '创建权限', 'group' => 'permission', 'description' => '创建新的权限'],
            ['name' => 'permissions.edit', 'display_name' => '编辑权限', 'group' => 'permission', 'description' => '编辑权限信息'],
            ['name' => 'permissions.delete', 'display_name' => '删除权限', 'group' => 'permission', 'description' => '删除权限'],

            // Product management
            ['name' => 'products.view', 'display_name' => '查看产品', 'group' => 'product', 'description' => '查看产品列表'],
            ['name' => 'products.create', 'display_name' => '创建产品', 'group' => 'product', 'description' => '创建新产品'],
            ['name' => 'products.edit', 'display_name' => '编辑产品', 'group' => 'product', 'description' => '编辑产品信息'],
            ['name' => 'products.delete', 'display_name' => '删除产品', 'group' => 'product', 'description' => '删除产品'],

            // Order management
            ['name' => 'orders.view', 'display_name' => '查看订单', 'group' => 'order', 'description' => '查看订单列表'],
            ['name' => 'orders.edit', 'display_name' => '编辑订单', 'group' => 'order', 'description' => '编辑订单状态'],

            // Profile
            ['name' => 'profile.edit', 'display_name' => '编辑个人信息', 'group' => 'profile', 'description' => '修改个人资料'],
        ];

        $permissionIds = [];
        foreach ($permissions as $perm) {
            $permissionIds[] = DB::table('permissions')->insertGetId([
                'name' => $perm['name'],
                'display_name' => $perm['display_name'],
                'group' => $perm['group'],
                'description' => $perm['description'],
                'created_at' => now(),
                'updated_at' => now(),
            ]);
        }

        // Assign all permissions to super-admin role
        foreach ($permissionIds as $permId) {
            DB::table('role_permission')->insert([
                'role_id' => $superAdminRoleId,
                'permission_id' => $permId,
            ]);
        }

        // Assign selected permissions to admin role
        $adminPerms = DB::table('permissions')->whereIn('name', [
            'dashboard.view',
            'admins.view', 'admins.create', 'admins.edit',
            'users.view', 'users.edit',
            'roles.view',
            'permissions.view',
            'products.view', 'products.create', 'products.edit',
            'orders.view', 'orders.edit',
            'profile.edit',
        ])->pluck('id');
        foreach ($adminPerms as $permId) {
            DB::table('role_permission')->insert([
                'role_id' => $adminRoleId,
                'permission_id' => $permId,
            ]);
        }

        // Assign order and product permissions to operator role
        $operatorPerms = DB::table('permissions')->whereIn('name', [
            'dashboard.view',
            'orders.view', 'orders.edit',
            'products.view',
            'profile.edit',
        ])->pluck('id');
        foreach ($operatorPerms as $permId) {
            DB::table('role_permission')->insert([
                'role_id' => $operatorRoleId,
                'permission_id' => $permId,
            ]);
        }

        // Create super admin
        $superAdminId = DB::table('admins')->insertGetId([
            'name' => '超级管理员',
            'email' => 'admin@example.com',
            'password' => password_hash('password', PASSWORD_DEFAULT),
            'is_active' => 1,
            'created_at' => now(),
            'updated_at' => now(),
        ]);
        DB::table('admin_role')->insert([
            'admin_id' => $superAdminId,
            'role_id' => $superAdminRoleId,
        ]);

        // Create a regular admin
        $regularAdminId = DB::table('admins')->insertGetId([
            'name' => '普通管理员',
            'email' => 'manager@example.com',
            'password' => password_hash('password', PASSWORD_DEFAULT),
            'is_active' => 1,
            'created_at' => now(),
            'updated_at' => now(),
        ]);
        DB::table('admin_role')->insert([
            'admin_id' => $regularAdminId,
            'role_id' => $adminRoleId,
        ]);

        // Create test user
        DB::table('users')->insert([
            'name' => '测试用户',
            'email' => 'user@example.com',
            'password' => password_hash('password', PASSWORD_DEFAULT),
            'created_at' => now(),
            'updated_at' => now(),
        ]);

        for ($i = 0; $i < 5; $i++) {
            DB::table('users')->insert([
                'name' => '用户' . ($i + 1),
                'email' => 'user' . ($i + 1) . '@example.com',
                'password' => password_hash('password', PASSWORD_DEFAULT),
                'created_at' => now(),
                'updated_at' => now(),
            ]);
        }

        // Create sample products
        $products = [
            ['name' => '苹果 iPhone 15 Pro', 'sku' => 'APL-IP15P', 'description' => '最新款苹果手机', 'price' => 8999.00, 'stock' => 50],
            ['name' => '华为 Mate 60 Pro', 'sku' => 'HUA-M60P', 'description' => '华为旗舰手机', 'price' => 6999.00, 'stock' => 30],
            ['name' => '小米 14 Ultra', 'sku' => 'XIA-M14U', 'description' => '小米旗舰手机', 'price' => 5999.00, 'stock' => 80],
            ['name' => '联想拯救者笔记本', 'sku' => 'LEN-LJZ', 'description' => '游戏笔记本电脑', 'price' => 7999.00, 'stock' => 20],
            ['name' => '戴尔显示器 27寸', 'sku' => 'DEL-27D', 'description' => '4K 高清显示器', 'price' => 1999.00, 'stock' => 45],
        ];
        foreach ($products as $product) {
            DB::table('products')->insert([
                'name' => $product['name'],
                'sku' => $product['sku'],
                'description' => $product['description'],
                'price' => $product['price'],
                'stock' => $product['stock'],
                'status' => 'active',
                'created_at' => now(),
                'updated_at' => now(),
            ]);
        }

        // Create sample orders
        $userIds = DB::table('users')->orderBy('id')->take(3)->pluck('id');
        $productList = DB::table('products')->orderBy('id')->take(3)->get();

        $statuses = ['pending', 'paid', 'shipped'];
        $i = 0;
        foreach ($userIds as $userId) {
            $product = $productList[$i % count($productList)];
            $quantity = $i + 1;
            $subtotal = $product->price * $quantity;

            $orderId = DB::table('orders')->insertGetId([
                'order_no' => 'ORD' . date('YmdHis') . str_pad((string) ($i + 1), 4, '0', STR_PAD_LEFT),
                'user_id' => $userId,
                'status' => $statuses[$i],
                'total_amount' => $subtotal,
                'shipping_name' => DB::table('users')->where('id', $userId)->value('name'),
                'shipping_phone' => '1380000' . str_pad((string) ($i + 1), 4, '0', STR_PAD_LEFT),
                'shipping_address' => '北京市朝阳区示例路' . ($i + 1) . '号',
                'created_at' => now(),
                'updated_at' => now(),
            ]);

            DB::table('order_items')->insert([
                'order_id' => $orderId,
                'product_id' => $product->id,
                'product_name' => $product->name,
                'price' => $product->price,
                'quantity' => $quantity,
                'subtotal' => $subtotal,
                'created_at' => now(),
                'updated_at' => now(),
            ]);

            $i++;
        }
    }
}
