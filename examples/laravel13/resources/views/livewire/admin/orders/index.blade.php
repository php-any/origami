<div>
    <div class="mb-4 flex items-center justify-between">
        <div class="flex space-x-2">
            <input type="text" wire:model.live.debounce.300ms="search" placeholder="搜索订单号/客户..." class="px-3 py-2 border border-gray-300 rounded-md text-sm w-64">
            <select wire:model.live="statusFilter" class="px-3 py-2 border border-gray-300 rounded-md text-sm">
                <option value="">全部状态</option>
                <option value="pending">待付款</option>
                <option value="paid">已付款</option>
                <option value="shipped">已发货</option>
                <option value="completed">已完成</option>
                <option value="cancelled">已取消</option>
            </select>
        </div>
    </div>

    <div class="bg-white rounded-lg shadow overflow-hidden">
        <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
                <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">订单号</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">客户</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">金额</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">状态</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">创建时间</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">操作</th>
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                @forelse($orders as $order)
                <tr>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-mono font-medium text-blue-600">
                        <a href="{{ route('admin.orders.detail', $order->id) }}">{{ $order->order_no }}</a>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ $order->user?->name ?? $order->shipping_name }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">¥{{ number_format($order->total_amount, 2) }}</td>
                    <td class="px-6 py-4 whitespace-nowrap">
                        <select wire:change="updateStatus({{ $order->id }}, $event.target.value)"
                            class="text-xs px-2 py-1 rounded-md border border-gray-300">
                            <option value="pending" {{ $order->status === 'pending' ? 'selected' : '' }}>待付款</option>
                            <option value="paid" {{ $order->status === 'paid' ? 'selected' : '' }}>已付款</option>
                            <option value="shipped" {{ $order->status === 'shipped' ? 'selected' : '' }}>已发货</option>
                            <option value="completed" {{ $order->status === 'completed' ? 'selected' : '' }}>已完成</option>
                            <option value="cancelled" {{ $order->status === 'cancelled' ? 'selected' : '' }}>已取消</option>
                        </select>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ $order->created_at?->format('Y-m-d H:i') }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm">
                        <a href="{{ route('admin.orders.detail', $order->id) }}" class="text-blue-600 hover:text-blue-800">查看详情</a>
                    </td>
                </tr>
                @empty
                <tr>
                    <td colspan="6" class="px-6 py-4 text-center text-gray-500">暂无数据</td>
                </tr>
                @endforelse
            </tbody>
        </table>
        <div class="px-6 py-4">
            {{ $orders->links() }}
        </div>
    </div>
</div>
