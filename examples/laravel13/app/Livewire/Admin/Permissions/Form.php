<?php

namespace App\Livewire\Admin\Permissions;

use App\Models\Permission;
use Livewire\Component;

class Form extends Component
{
    public $permissionId;
    public $name = '';
    public $display_name = '';
    public $group = 'general';
    public $description = '';

    public function mount($permissionId = null)
    {
        $this->permissionId = $permissionId;

        if ($permissionId) {
            $perm = Permission::findOrFail($permissionId);
            $this->name = $perm->name;
            $this->display_name = $perm->display_name;
            $this->group = $perm->group;
            $this->description = $perm->description;
        }
    }

    public function save()
    {
        $this->validate([
            'name' => 'required|string|max:255|unique:permissions,name,' . ($this->permissionId ?? 'NULL'),
            'display_name' => 'required|string|max:255',
            'group' => 'required|string|max:50',
            'description' => 'nullable|string',
        ]);

        $data = [
            'name' => $this->name,
            'display_name' => $this->display_name,
            'group' => $this->group,
            'description' => $this->description,
        ];

        if ($this->permissionId) {
            Permission::findOrFail($this->permissionId)->update($data);
        } else {
            Permission::create($data);
        }

        session()->flash('message', '保存成功');
        return redirect()->route('admin.permissions.index');
    }

    public function render()
    {
        return view('livewire.admin.permissions.form');
    }
}
