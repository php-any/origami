<?php

namespace App\Policies;

use App\Models\Admin;
use App\Models\Category;

class CategoryPolicy extends AdminResourcePolicy
{
    public function viewAny(?Admin $admin): bool
    {
        return $this->allows($admin, 'categories.view');
    }

    public function view(?Admin $admin, Category $category): bool
    {
        return $this->allows($admin, 'categories.view');
    }

    public function create(?Admin $admin): bool
    {
        return $this->allows($admin, 'categories.create');
    }

    public function update(?Admin $admin, Category $category): bool
    {
        return $this->allows($admin, 'categories.edit');
    }

    public function delete(?Admin $admin, Category $category): bool
    {
        return $this->allows($admin, 'categories.delete');
    }
}
