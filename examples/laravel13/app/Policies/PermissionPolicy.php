<?php

namespace App\Policies;

use App\Models\Admin;
use App\Models\Permission;

class PermissionPolicy extends AdminResourcePolicy
{
    public function viewAny(?Admin $admin): bool
    {
        return $this->allows($admin, 'permissions.view');
    }

    public function view(?Admin $admin, Permission $permission): bool
    {
        return $this->allows($admin, 'permissions.view');
    }

    public function create(?Admin $admin): bool
    {
        return $this->allows($admin, 'permissions.create');
    }

    public function update(?Admin $admin, Permission $permission): bool
    {
        return $this->allows($admin, 'permissions.edit');
    }

    public function delete(?Admin $admin, Permission $permission): bool
    {
        if ($permission->roles()->exists()) {
            return false;
        }

        return $this->allows($admin, 'permissions.delete');
    }
}
