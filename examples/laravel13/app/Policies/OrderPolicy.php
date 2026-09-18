<?php

namespace App\Policies;

use App\Models\Admin;
use App\Models\Order;

class OrderPolicy extends AdminResourcePolicy
{
    public function viewAny(?Admin $admin): bool
    {
        return $this->allows($admin, 'orders.view');
    }

    public function view(?Admin $admin, Order $order): bool
    {
        return $this->allows($admin, 'orders.view');
    }

    public function create(?Admin $admin): bool
    {
        return false;
    }

    public function update(?Admin $admin, Order $order): bool
    {
        return $this->allows($admin, 'orders.edit');
    }

    public function delete(?Admin $admin, Order $order): bool
    {
        return false;
    }
}
