<?php

namespace App\Livewire\Admin\Users;

use App\Models\User;
use Livewire\Component;

class Form extends Component
{
    public $userId;
    public $name = '';
    public $email = '';
    public $password = '';

    public function mount($userId = null)
    {
        $this->userId = $userId;

        if ($userId) {
            $user = User::findOrFail($userId);
            $this->name = $user->name;
            $this->email = $user->email;
        }
    }

    public function save()
    {
        $rules = [
            'name' => 'required|string|max:255',
            'email' => 'required|email|max:255|unique:users,email,' . ($this->userId ?? 'NULL'),
        ];

        if (!$this->userId) {
            $rules['password'] = 'required|string|min:8';
        }

        $this->validate($rules);

        $data = [
            'name' => $this->name,
            'email' => $this->email,
        ];

        if ($this->password) {
            $data['password'] = $this->password;
        }

        if ($this->userId) {
            $user = User::findOrFail($this->userId);
            $user->update($data);
        } else {
            User::create($data);
        }

        session()->flash('message', '保存成功');
        return redirect()->route('admin.users.index');
    }

    public function render()
    {
        return view('livewire.admin.users.form');
    }
}
