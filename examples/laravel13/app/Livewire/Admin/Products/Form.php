<?php

namespace App\Livewire\Admin\Products;

use App\Models\Product;
use Livewire\Component;

class Form extends Component
{
    public $productId;
    public $name = '';
    public $sku = '';
    public $description = '';
    public $price = 0;
    public $stock = 0;
    public $status = 'active';

    public function mount($productId = null)
    {
        $this->productId = $productId;

        if ($productId) {
            $product = Product::findOrFail($productId);
            $this->name = $product->name;
            $this->sku = $product->sku;
            $this->description = $product->description;
            $this->price = $product->price;
            $this->stock = $product->stock;
            $this->status = $product->status;
        }
    }

    public function save()
    {
        $this->validate([
            'name' => 'required|string|max:255',
            'sku' => 'required|string|max:50|unique:products,sku,' . ($this->productId ?? 'NULL'),
            'description' => 'nullable|string',
            'price' => 'required|numeric|min:0',
            'stock' => 'required|integer|min:0',
            'status' => 'required|in:active,inactive',
        ]);

        $data = [
            'name' => $this->name,
            'sku' => $this->sku,
            'description' => $this->description,
            'price' => $this->price,
            'stock' => $this->stock,
            'status' => $this->status,
        ];

        if ($this->productId) {
            Product::findOrFail($this->productId)->update($data);
        } else {
            Product::create($data);
        }

        session()->flash('message', '保存成功');
        return redirect()->route('admin.products.index');
    }

    public function render()
    {
        return view('livewire.admin.products.form');
    }
}
