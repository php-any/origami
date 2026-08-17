<?php

namespace App\Livewire\Admin\Roles;

use App\Models\Role;
use Livewire\Component;

class Index extends Component
{
    public function delete($id)
    {
        $role = Role::find($id);
        if ($role && $role->name !== 'super-admin') {
            $role->delete();
        }
    }

    public function render()
    {
        $roles = Role::withCount('admins', 'permissions')->latest()->get();

        return view('livewire.admin.roles.index', [
            'roles' => $roles,
        ]);
    }
}
