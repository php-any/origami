<?php

namespace App\Livewire\Admin\Users;

use App\Models\User;
use Livewire\Component;

class Index extends Component
{
    public $search = '';

    public function delete($id)
    {
        $user = User::find($id);
        if ($user) {
            $user->delete();
        }
    }

    public function render()
    {
        $query = User::query();

        if ($this->search) {
            $query->where(function ($q) {
                $q->where('name', 'like', '%' . $this->search . '%')
                    ->orWhere('email', 'like', '%' . $this->search . '%');
            });
        }

        $users = $query->latest()->paginate(10);

        return view('livewire.admin.users.index', [
            'users' => $users,
        ]);
    }
}
