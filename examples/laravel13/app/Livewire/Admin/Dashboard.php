<?php

namespace App\Livewire\Admin;

use App\Models\Order;
use App\Models\Product;
use App\Models\User;
use Livewire\Component;

class Dashboard extends Component
{
    public $stats = [];

    public function mount()
    {
        $this->stats = [
            'total_users' => User::count(),
            'total_orders' => Order::count(),
            'total_products' => Product::count(),
            'total_revenue' => Order::where('status', '!=', 'cancelled')->sum('total_amount'),
            'pending_orders' => Order::where('status', 'pending')->count(),
        ];
    }

    public function render()
    {
        $recentOrders = Order::with('user')->latest()->take(5)->get();

        return view('livewire.admin.dashboard', [
            'recentOrders' => $recentOrders,
        ]);
    }
}
