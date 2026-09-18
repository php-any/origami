<?php

namespace App\Policies;

use App\Models\Admin;
use App\Models\Role;

class RolePolicy extends AdminResourcePolicy
{
    public function viewAny(?Admin $admin): bool
    {
        return $this->allows($admin, 'roles.view');
    }

    public function view(?Admin $admin, Role $role): bool
    {
        return $this->allows($admin, 'roles.view');
    }

    public function create(?Admin $admin): bool
    {
        return $this->allows($admin, 'roles.create');
    }

    public function update(?Admin $admin, Role $role): bool
    {
        return $this->allows($admin, 'roles.edit');
    }

    public function delete(?Admin $admin, Role $role): bool
    {
        if ($role->isSuperAdmin()) {
            return false;
        }

        return $this->allows($admin, 'roles.delete');
    }
}
