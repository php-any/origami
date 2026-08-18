<?php

namespace App\Livewire\Admin\Orders;

use App\Models\Order;
use Livewire\Component;

class Index extends Component
{
    public $search = '';
    public $statusFilter = '';

    public function updateStatus($id, $status)
    {
        $order = Order::find($id);
        if ($order) {
            $order->status = $status;
            $order->save();
        }
    }

    public function render()
    {
        $query = Order::query();

        if ($this->search) {
            $query->where(function ($q) {
                $q->where('order_no', 'like', '%' . $this->search . '%')
                    ->orWhere('shipping_name', 'like', '%' . $this->search . '%')
                    ->orWhereHas('user', function ($q2) {
                        $q2->where('name', 'like', '%' . $this->search . '%')
                            ->orWhere('email', 'like', '%' . $this->search . '%');
                    });
            });
        }

        if ($this->statusFilter) {
            $query->where('status', $this->statusFilter);
        }

        $orders = $query->with('user', 'items')->latest()->paginate(10);

        return view('livewire.admin.orders.index', [
            'orders' => $orders,
        ]);
    }
}
