<?php

namespace App\Policies;

use App\Models\Admin;

class AdminPolicy extends AdminResourcePolicy
{
    public function viewAny(?Admin $admin): bool
    {
        return $this->allows($admin, 'admins.view');
    }

    public function view(?Admin $admin, Admin $model): bool
    {
        return $this->allows($admin, 'admins.view');
    }

    public function create(?Admin $admin): bool
    {
        return $this->allows($admin, 'admins.create');
    }

    public function update(?Admin $admin, Admin $model): bool
    {
        return $this->allows($admin, 'admins.edit');
    }

    public function delete(?Admin $admin, Admin $model): bool
    {
        $actor = $this->admin($admin);

        if (! $actor || ! $this->allows($admin, 'admins.delete')) {
            return false;
        }

        if ($actor->id === $model->id) {
            return false;
        }

        return true;
    }
}
