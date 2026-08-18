<div>
    <div class="mb-4 flex items-center justify-between">
        <div class="flex space-x-2">
            <input type="text" wire:model.live.debounce.300ms="search" placeholder="搜索产品..." class="px-3 py-2 border border-gray-300 rounded-md text-sm w-64">
            <select wire:model.live="statusFilter" class="px-3 py-2 border border-gray-300 rounded-md text-sm">
                <option value="">全部状态</option>
                <option value="active">在售</option>
                <option value="inactive">下架</option>
            </select>
        </div>
        <a href="{{ route('admin.products.create') }}" class="bg-blue-600 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-700">+ 新建产品</a>
    </div>

    <div class="bg-white rounded-lg shadow overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
                <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">名称</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">SKU</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">价格</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">库存</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">状态</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">操作</th>
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                @forelse($products as $product)
                <tr>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ $product->id }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ $product->name }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-500">{{ $product->sku }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">¥{{ number_format($product->price, 2) }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ $product->stock }}</td>
                    <td class="px-6 py-4 whitespace-nowrap">
                        <button wire:click="toggleStatus({{ $product->id }})"
                            class="px-2 py-1 text-xs rounded-full {{ $product->status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800' }}">
                            {{ $product->status === 'active' ? '在售' : '下架' }}
                        </button>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm">
                        <a href="{{ route('admin.products.edit', $product->id) }}" class="text-blue-600 hover:text-blue-800 mr-3">编辑</a>
                        <button wire:click="delete({{ $product->id }})" wire:confirm="确定删除该产品？" class="text-red-600 hover:text-red-800">删除</button>
                    </td>
                </tr>
                @empty
                <tr>
                    <td colspan="7" class="px-6 py-4 text-center text-gray-500">暂无数据</td>
                </tr>
                @endforelse
            </tbody>
        </table>
        <div class="px-6 py-4">
            {{ $products->links() }}
        </div>
    </div>
</div>
