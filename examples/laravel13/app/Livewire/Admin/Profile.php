<?php

namespace App\Livewire\Admin;

use Livewire\Component;

class Profile extends Component
{
    public $name = '';
    public $email = '';
    public $current_password = '';
    public $new_password = '';
    public $confirm_password = '';

    public function mount()
    {
        $admin = auth('admin')->user();
        $this->name = $admin->name;
        $this->email = $admin->email;
    }

    public function updateProfile()
    {
        $this->validate([
            'name' => 'required|string|max:255',
            'email' => 'required|email|max:255|unique:admins,email,' . auth('admin')->id(),
        ]);

        $admin = auth('admin')->user();
        $admin->name = $this->name;
        $admin->email = $this->email;
        $admin->save();

        session()->flash('message', '个人资料已更新');
    }

    public function changePassword()
    {
        $this->validate([
            'current_password' => 'required|string',
            'new_password' => 'required|string|min:8',
            'confirm_password' => 'required|string|same:new_password',
        ]);

        $admin = auth('admin')->user();
        if (!password_verify($this->current_password, $admin->password)) {
            session()->flash('error', '当前密码不正确');
            return;
        }

        $admin->password = $this->new_password;
        $admin->save();

        $this->current_password = '';
        $this->new_password = '';
        $this->confirm_password = '';

        session()->flash('message', '密码已更新');
    }

    public function render()
    {
        return view('livewire.admin.profile');
    }
}
