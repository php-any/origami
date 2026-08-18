<?php

namespace App\Livewire\Admin\Roles;

use App\Models\Permission;
use App\Models\Role;
use Livewire\Component;

class Form extends Component
{
    public $roleId;
    public $name = '';
    public $display_name = '';
    public $description = '';
    public $permission_ids = [];

    public function mount($roleId = null)
    {
        $this->roleId = $roleId;

        if ($roleId) {
            $role = Role::findOrFail($roleId);
            $this->name = $role->name;
            $this->display_name = $role->display_name;
            $this->description = $role->description;
            $this->permission_ids = $role->permissions()->pluck('permissions.id')->all();
        }
    }

    public function save()
    {
        $this->validate([
            'name' => 'required|string|max:255|unique:roles,name,' . ($this->roleId ?? 'NULL'),
            'display_name' => 'required|string|max:255',
            'description' => 'nullable|string',
        ]);

        $data = [
            'name' => $this->name,
            'display_name' => $this->display_name,
            'description' => $this->description,
        ];

        if ($this->roleId) {
            $role = Role::findOrFail($this->roleId);
            if ($role->name === 'super-admin') {
                $data['name'] = 'super-admin'; // 不允许修改 super-admin 的角色名
            }
            $role->update($data);
        } else {
            $role = Role::create($data);
        }

        $role->permissions()->sync($this->permission_ids ?? []);

        session()->flash('message', '保存成功');
        return redirect()->route('admin.roles.index');
    }

    public function render()
    {
        $permissionGroups = Permission::all()->groupBy('group');

        return view('livewire.admin.roles.form', [
            'permissionGroups' => $permissionGroups,
        ]);
    }
}
