<?php

namespace App\Livewire\Admin\Orders;

use App\Models\Order;
use Livewire\Component;

class Detail extends Component
{
    public $order;

    public function mount($orderId)
    {
        $this->order = Order::with('user', 'items.product')->findOrFail($orderId);
    }

    public function updateStatus($status)
    {
        $validStatuses = ['pending', 'paid', 'shipped', 'completed', 'cancelled'];
        if (in_array($status, $validStatuses)) {
            $this->order->status = $status;
            $this->order->save();
            session()->flash('message', '订单状态已更新');
        }
    }

    public function render()
    {
        return view('livewire.admin.orders.detail');
    }
}
