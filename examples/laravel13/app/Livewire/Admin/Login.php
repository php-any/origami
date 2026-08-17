<?php

namespace App\Livewire\Admin;

use Livewire\Component;

class Login extends Component
{
    public string $email = '';
    public string $password = '';
    public bool $remember = false;

    public function login()
    {
        $this->validate([
            'email' => 'required|email',
            'password' => 'required',
        ]);

        if (auth('admin')->attempt(['email' => $this->email, 'password' => $this->password], $this->remember)) {
            /** @var \App\Models\Admin $admin */
            $admin = auth('admin')->user();
            if (!$admin->is_active) {
                auth('admin')->logout();
                session()->flash('error', '该账号已被禁用');
                return;
            }
            session()->regenerate();
            return redirect()->route('admin.dashboard');
        }

        session()->flash('error', '邮箱或密码错误');
    }

    public function render()
    {
        return view('livewire.admin.login');
    }
}
