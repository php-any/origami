<div>
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4 mb-8">
        <div class="bg-white rounded-lg shadow p-5">
            <div class="text-sm text-gray-500">用户总数</div>
            <div class="text-2xl font-bold mt-1">{{ number_format($stats['total_users'] ?? 0) }}</div>
        </div>
        <div class="bg-white rounded-lg shadow p-5">
            <div class="text-sm text-gray-500">订单总数</div>
            <div class="text-2xl font-bold mt-1">{{ number_format($stats['total_orders'] ?? 0) }}</div>
        </div>
        <div class="bg-white rounded-lg shadow p-5">
            <div class="text-sm text-gray-500">产品总数</div>
            <div class="text-2xl font-bold mt-1">{{ number_format($stats['total_products'] ?? 0) }}</div>
        </div>
        <div class="bg-white rounded-lg shadow p-5">
            <div class="text-sm text-gray-500">总营收</div>
            <div class="text-2xl font-bold mt-1 text-green-600">¥{{ number_format($stats['total_revenue'] ?? 0, 2) }}</div>
        </div>
        <div class="bg-white rounded-lg shadow p-5">
            <div class="text-sm text-gray-500">待处理订单</div>
            <div class="text-2xl font-bold mt-1 text-yellow-600">{{ number_format($stats['pending_orders'] ?? 0) }}</div>
        </div>
    </div>

    <div class="bg-white rounded-lg shadow overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200">
            <h2 class="text-lg font-semibold text-gray-800">最近订单</h2>
        </div>
        <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
                <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">订单号</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">客户</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">金额</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">下单时间</th>
                </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
                @forelse($recentOrders as $order)
                <tr>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-blue-600">
                        <a href="{{ route('admin.orders.detail', $order->id) }}">{{ $order->order_no }}</a>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ $order->user?->name ?? $order->shipping_name }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">¥{{ number_format($order->total_amount, 2) }}</td>
                    <td class="px-6 py-4 whitespace-nowrap">
                        <span class="px-2 py-1 text-xs rounded-full
                            @if($order->status === 'pending') bg-yellow-100 text-yellow-800
                            @elseif($order->status === 'paid') bg-blue-100 text-blue-800
                            @elseif($order->status === 'shipped') bg-purple-100 text-purple-800
                            @elseif($order->status === 'completed') bg-green-100 text-green-800
                            @else bg-gray-100 text-gray-800 @endif">
                            {{ $order->status_display }}
                        </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ $order->created_at?->format('Y-m-d H:i') }}</td>
                </tr>
                @empty
                <tr>
                    <td colspan="5" class="px-6 py-4 text-center text-gray-500">暂无订单</td>
                </tr>
                @endforelse
            </tbody>
        </table>
    </div>
</div>
