<?php

namespace App\Policies;

use App\Models\Admin;
use App\Models\Product;

class ProductPolicy extends AdminResourcePolicy
{
    public function viewAny(?Admin $admin): bool
    {
        return $this->allows($admin, 'products.view');
    }

    public function view(?Admin $admin, Product $product): bool
    {
        return $this->allows($admin, 'products.view');
    }

    public function create(?Admin $admin): bool
    {
        return $this->allows($admin, 'products.create');
    }

    public function update(?Admin $admin, Product $product): bool
    {
        return $this->allows($admin, 'products.edit');
    }

    public function delete(?Admin $admin, Product $product): bool
    {
        return $this->allows($admin, 'products.delete');
    }
}
