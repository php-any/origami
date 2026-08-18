<div>
    @if(session('message'))
    <div class="bg-green-100 border-l-4 border-green-500 text-green-700 p-4 mb-4">
        {{ session('message') }}
    </div>
    @endif

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div class="lg:col-span-2 bg-white rounded-lg shadow overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
                <h2 class="text-lg font-semibold text-gray-800">订单 #{{ $order->order_no }}</h2>
                <span class="px-3 py-1 text-sm rounded-full
                    @if($order->status === 'pending') bg-yellow-100 text-yellow-800
                    @elseif($order->status === 'paid') bg-blue-100 text-blue-800
                    @elseif($order->status === 'shipped') bg-purple-100 text-purple-800
                    @elseif($order->status === 'completed') bg-green-100 text-green-800
                    @else bg-gray-100 text-gray-800 @endif">
                    {{ $order->status_display }}
                </span>
            </div>

            <div class="px-6 py-4">
                <h3 class="text-sm font-semibold text-gray-700 mb-3">订单商品</h3>
                <table class="min-w-full">
                    <thead>
                        <tr class="border-b">
                            <th class="py-2 text-left text-xs font-medium text-gray-500">商品</th>
                            <th class="py-2 text-left text-xs font-medium text-gray-500">单价</th>
                            <th class="py-2 text-left text-xs font-medium text-gray-500">数量</th>
                            <th class="py-2 text-left text-xs font-medium text-gray-500">小计</th>
                        </tr>
                    </thead>
                    <tbody>
                        @foreach($order->items as $item)
                        <tr class="border-b">
                            <td class="py-2 text-sm text-gray-900">{{ $item->product_name }}</td>
                            <td class="py-2 text-sm text-gray-500">¥{{ number_format($item->price, 2) }}</td>
                            <td class="py-2 text-sm text-gray-500">{{ $item->quantity }}</td>
                            <td class="py-2 text-sm text-gray-900">¥{{ number_format($item->subtotal, 2) }}</td>
                        </tr>
                        @endforeach
                    </tbody>
                    <tfoot>
                        <tr>
                            <td colspan="3" class="py-2 text-right text-sm font-semibold text-gray-700">合计：</td>
                            <td class="py-2 text-sm font-bold text-red-600">¥{{ number_format($order->total_amount, 2) }}</td>
                        </tr>
                    </tfoot>
                </table>
            </div>
        </div>

        <div class="space-y-6">
            <div class="bg-white rounded-lg shadow p-6">
                <h3 class="text-sm font-semibold text-gray-700 mb-3">收货信息</h3>
                <div class="text-sm space-y-2 text-gray-600">
                    <p>收货人：{{ $order->shipping_name }}</p>
                    <p>电话：{{ $order->shipping_phone }}</p>
                    <p>地址：{{ $order->shipping_address }}</p>
                    @if($order->remark)
                    <p>备注：{{ $order->remark }}</p>
                    @endif
                </div>
            </div>

            <div class="bg-white rounded-lg shadow p-6">
                <h3 class="text-sm font-semibold text-gray-700 mb-3">订单操作</h3>
                <div class="space-y-2">
                    <button wire:click="updateStatus('paid')" class="w-full bg-blue-600 text-white px-4 py-2 rounded-md text-sm hover:bg-blue-700">标记已付款</button>
                    <button wire:click="updateStatus('shipped')" class="w-full bg-purple-600 text-white px-4 py-2 rounded-md text-sm hover:bg-purple-700">标记已发货</button>
                    <button wire:click="updateStatus('completed')" class="w-full bg-green-600 text-white px-4 py-2 rounded-md text-sm hover:bg-green-700">标记已完成</button>
                    <button wire:click="updateStatus('cancelled')" class="w-full bg-red-600 text-white px-4 py-2 rounded-md text-sm hover:bg-red-700">取消订单</button>
                </div>
            </div>
        </div>
    </div>
</div>
