<?php

namespace App\Policies;

use App\Models\Admin;

abstract class AdminResourcePolicy
{
    protected function admin(?Admin $admin): ?Admin
    {
        return $admin ?? auth('admin')->user();
    }

    protected function allows(?Admin $admin, string $permission): bool
    {
        $admin = $this->admin($admin);

        return $admin instanceof Admin && $admin->hasPermission($permission);
    }
}
