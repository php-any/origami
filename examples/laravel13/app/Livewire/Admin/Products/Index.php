<?php

namespace App\Livewire\Admin\Products;

use App\Models\Product;
use Livewire\Component;

class Index extends Component
{
    public $search = '';
    public $statusFilter = '';

    public function toggleStatus($id)
    {
        $product = Product::find($id);
        if ($product) {
            $product->status = $product->status === 'active' ? 'inactive' : 'active';
            $product->save();
        }
    }

    public function delete($id)
    {
        $product = Product::find($id);
        if ($product) {
            $product->delete();
        }
    }

    public function render()
    {
        $query = Product::query();

        if ($this->search) {
            $query->where(function ($q) {
                $q->where('name', 'like', '%' . $this->search . '%')
                    ->orWhere('sku', 'like', '%' . $this->search . '%');
            });
        }

        if ($this->statusFilter) {
            $query->where('status', $this->statusFilter);
        }

        $products = $query->latest()->paginate(10);

        return view('livewire.admin.products.index', [
            'products' => $products,
        ]);
    }
}
