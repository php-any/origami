<?php

namespace App\Livewire\Admin\Admins;

use App\Models\Admin;
use App\Models\Role;
use Livewire\Component;

class Form extends Component
{
    public $adminId;
    public $name = '';
    public $email = '';
    public $password = '';
    public $is_active = true;
    public $role_ids = [];

    public function mount($adminId = null)
    {
        $this->adminId = $adminId;

        if ($adminId) {
            $admin = Admin::findOrFail($adminId);
            $this->name = $admin->name;
            $this->email = $admin->email;
            $this->is_active = $admin->is_active;
            $this->role_ids = $admin->roles()->pluck('roles.id')->all();
        }
    }

    public function save()
    {
        $rules = [
            'name' => 'required|string|max:255',
            'email' => 'required|email|max:255|unique:admins,email,' . ($this->adminId ?? 'NULL'),
        ];

        if (!$this->adminId) {
            $rules['password'] = 'required|string|min:8';
        }

        $this->validate($rules);

        $data = [
            'name' => $this->name,
            'email' => $this->email,
            'is_active' => $this->is_active,
        ];

        if ($this->password) {
            $data['password'] = $this->password;
        }

        if ($this->adminId) {
            $admin = Admin::findOrFail($this->adminId);
            $admin->update($data);
        } else {
            $admin = Admin::create($data);
        }

        $admin->roles()->sync($this->role_ids ?? []);

        session()->flash('message', '保存成功');

        return redirect()->route('admin.admins.index');
    }

    public function render()
    {
        $roles = Role::all();
        return view('livewire.admin.admins.form', [
            'roles' => $roles,
        ]);
    }
}
