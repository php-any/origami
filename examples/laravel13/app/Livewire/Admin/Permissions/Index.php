<?php

namespace App\Livewire\Admin\Permissions;

use App\Models\Permission;
use Livewire\Component;

class Index extends Component
{
    public $search = '';

    public function delete($id)
    {
        $perm = Permission::find($id);
        if ($perm) {
            $perm->delete();
        }
    }

    public function render()
    {
        $query = Permission::query();

        if ($this->search) {
            $query->where(function ($q) {
                $q->where('name', 'like', '%' . $this->search . '%')
                    ->orWhere('display_name', 'like', '%' . $this->search . '%');
            });
        }

        $permissions = $query->orderBy('group')->orderBy('name')->get();

        return view('livewire.admin.permissions.index', [
            'permissions' => $permissions,
        ]);
    }
}
