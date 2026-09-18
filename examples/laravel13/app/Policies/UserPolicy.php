<?php

namespace App\Policies;

use App\Models\Admin;
use App\Models\User;

class UserPolicy extends AdminResourcePolicy
{
    public function viewAny(?Admin $admin): bool
    {
        return $this->allows($admin, 'users.view');
    }

    public function view(?Admin $admin, User $user): bool
    {
        return $this->allows($admin, 'users.view');
    }

    public function create(?Admin $admin): bool
    {
        return $this->allows($admin, 'users.create');
    }

    public function update(?Admin $admin, User $user): bool
    {
        return $this->allows($admin, 'users.edit');
    }

    public function delete(?Admin $admin, User $user): bool
    {
        return $this->allows($admin, 'users.delete');
    }
}
