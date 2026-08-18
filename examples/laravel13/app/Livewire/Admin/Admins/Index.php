<?php

namespace App\Livewire\Admin\Admins;

use App\Models\Admin;
use Livewire\Component;

class Index extends Component
{
    public $search = '';
    public $statusFilter = '';

    public function toggleActive($id)
    {
        $admin = Admin::find($id);
        if ($admin && $admin->id !== auth('admin')->id()) {
            $admin->is_active = !$admin->is_active;
            $admin->save();
        }
    }

    public function delete($id)
    {
        $admin = Admin::find($id);
        if ($admin && $admin->id !== auth('admin')->id()) {
            $admin->delete();
        }
    }

    public function render()
    {
        $query = Admin::query();

        if ($this->search) {
            $query->where(function ($q) {
                $q->where('name', 'like', '%' . $this->search . '%')
                    ->orWhere('email', 'like', '%' . $this->search . '%');
            });
        }

        if ($this->statusFilter !== '') {
            $query->where('is_active', $this->statusFilter === 'active');
        }

        $admins = $query->with('roles')->latest()->paginate(10);

        return view('livewire.admin.admins.index', [
            'admins' => $admins,
        ]);
    }
}
