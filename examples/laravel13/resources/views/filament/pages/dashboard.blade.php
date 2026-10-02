<x-filament-panels::page>
    <div class="tabler-ui" data-bs-theme="light">
        <div class="page-header mb-4">
            <div class="row g-2 align-items-center">
                <div class="col">
                    <div class="page-pretitle">管理控制台</div>
                    <h1 class="page-title">工作台</h1>
                </div>
                <div class="col-auto ms-auto d-print-none">
                    <div class="btn-list">
                        @if (\App\Filament\Resources\Orders\OrderResource::canViewAny())
                            <a href="{{ \App\Filament\Resources\Orders\OrderResource::getUrl() }}" class="btn"><x-filament::icon icon="heroicon-o-shopping-bag" class="icon" />订单管理</a>
                        @endif
                        @if (\App\Filament\Resources\Products\ProductResource::canCreate())
                            <a href="{{ \App\Filament\Resources\Products\ProductResource::getUrl('create') }}" class="btn btn-primary"><x-filament::icon icon="heroicon-o-plus" class="icon" />新增商品</a>
                        @endif
                    </div>
                </div>
            </div>
        </div>

        <div class="row row-cards">
            @foreach ($metrics as $metric)
                <div class="col-sm-6 col-xl-3">
                    <div class="card">
                        <div class="card-body">
                            <div class="row align-items-center">
                                <div class="col-auto"><span class="avatar bg-{{ $metric['tone'] }}-lt"><x-filament::icon :icon="$metric['icon']" class="icon" /></span></div>
                                <div class="col"><div class="subheader">{{ $metric['label'] }}</div><div class="h1 mb-0 mt-1">{{ $metric['value'] }}</div></div>
                            </div>
                            <div class="text-secondary mt-3 small">{{ $metric['note'] }}</div>
                        </div>
                    </div>
                </div>
            @endforeach

            @if ($canViewOrders)
                <div class="col-lg-8">
                    <div class="card">
                        <div class="card-header">
                            <div><h2 class="card-title">订单概览</h2><div class="card-subtitle">按状态查看和处理订单</div></div>
                            <div class="card-actions"><a href="{{ \App\Filament\Resources\Orders\OrderResource::getUrl() }}" class="btn btn-sm">查看全部<x-filament::icon icon="heroicon-o-arrow-right" class="icon ms-1" /></a></div>
                        </div>
                        <div class="table-responsive">
                            <table class="table table-vcenter card-table">
                                <thead><tr><th>订单状态</th><th class="text-end">数量</th><th class="w-50">占比</th><th class="w-1"></th></tr></thead>
                                <tbody>
                                    @foreach ($orderStages as $stage)
                                        <tr>
                                            <td>{{ $stage['label'] }}</td><td class="text-end fw-bold">{{ $stage['count'] }}</td>
                                            <td><div class="d-flex align-items-center gap-3"><div class="progress flex-fill"><div class="progress-bar" role="progressbar" style="width: {{ $stage['percent'] }}%" aria-valuenow="{{ $stage['percent'] }}" aria-valuemin="0" aria-valuemax="100" aria-label="{{ $stage['label'] }}占比"></div></div><span class="text-secondary small">{{ $stage['percent'] }}%</span></div></td>
                                            <td><a href="{{ $stage['url'] }}" class="btn btn-sm btn-ghost-primary" aria-label="查看{{ $stage['label'] }}订单">查看</a></td>
                                        </tr>
                                    @endforeach
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
            @endif

            <div class="{{ $canViewOrders ? 'col-lg-4' : 'col-12' }}">
                <div class="card">
                    <div class="card-header"><h2 class="card-title">常用入口</h2></div>
                    <div class="list-group list-group-flush">
                        @foreach ($shortcuts as $shortcut)
                            <a class="list-group-item list-group-item-action" href="{{ $shortcut['url'] }}">
                                <div class="row align-items-center">
                                    <div class="col-auto"><span class="avatar avatar-sm"><x-filament::icon :icon="$shortcut['icon']" class="icon" /></span></div>
                                    <div class="col text-truncate"><div class="fw-medium">{{ $shortcut['label'] }}</div><div class="text-secondary small">{{ $shortcut['note'] }}</div></div>
                                    <div class="col-auto text-secondary"><x-filament::icon icon="heroicon-o-chevron-right" class="icon" /></div>
                                </div>
                            </a>
                        @endforeach
                        @if (! count($shortcuts))<div class="empty"><p class="empty-title">暂无可用模块</p><p class="empty-subtitle text-secondary">当前账号未分配业务访问权限。</p></div>@endif
                    </div>
                </div>
            </div>

            @if ($canViewOrders)
                <div class="col-lg-8">
                    <div class="card">
                        <div class="card-header"><h2 class="card-title">最新订单</h2><div class="card-actions"><span class="badge bg-secondary-lt">最近 5 笔</span></div></div>
                        @if (count($latestOrders))
                            <div class="table-responsive">
                                <table class="table table-vcenter card-table">
                                    <thead><tr><th>订单编号</th><th>客户</th><th>状态</th><th class="text-end">金额</th><th>创建时间</th></tr></thead>
                                    <tbody>
                                        @foreach ($latestOrders as $order)
                                            <tr><td><a href="{{ \App\Filament\Resources\Orders\OrderResource::getUrl('view', ['record' => $order]) }}">{{ $order->order_no }}</a></td><td>{{ $order->user?->name ?? '—' }}</td><td><span class="badge bg-secondary-lt">{{ $order->status_display }}</span></td><td class="text-end">¥ {{ number_format((float) $order->total_amount, 2) }}</td><td class="text-secondary">{{ $order->created_at->format('m-d H:i') }}</td></tr>
                                        @endforeach
                                    </tbody>
                                </table>
                            </div>
                        @else
                            <div class="empty py-5"><div class="empty-icon"><x-filament::icon icon="heroicon-o-shopping-bag" class="icon" /></div><p class="empty-title">暂无订单</p><p class="empty-subtitle text-secondary">客户提交订单后，最新交易会显示在这里。</p></div>
                        @endif
                    </div>
                </div>
            @endif

            @if ($canViewProducts)
                <div class="{{ $canViewOrders ? 'col-lg-4' : 'col-12' }}">
                    <div class="card">
                        <div class="card-header"><h2 class="card-title">库存预警</h2><div class="card-actions"><span class="badge bg-yellow-lt">{{ $lowStockCount }} 件</span></div></div>
                        @if ($lowStockCount)
                            <div class="list-group list-group-flush">
                                @foreach ($lowStockProducts as $product)
                                    <div class="list-group-item"><div class="row align-items-center"><div class="col"><div class="fw-medium">{{ $product->name }}</div><div class="text-secondary small">{{ $product->sku }}</div></div><div class="col-auto"><span class="badge bg-{{ $product->stock === 0 ? 'red' : 'yellow' }}-lt">剩余 {{ $product->stock }} 件</span></div></div></div>
                                @endforeach
                            </div>
                        @else
                            <div class="empty py-5"><div class="empty-icon"><x-filament::icon icon="heroicon-o-check-circle" class="icon" /></div><p class="empty-title">暂无低库存商品</p><p class="empty-subtitle text-secondary">在售商品库存不超过 10 件时触发预警。</p></div>
                        @endif
                        <div class="card-footer"><a href="{{ \App\Filament\Resources\Products\ProductResource::getUrl('index', ['tab' => 'low_stock']) }}" class="btn btn-outline-primary w-100">查看预警商品</a></div>
                    </div>
                </div>
            @endif
        </div>
        <div class="text-secondary small mt-3">统计范围：全部历史数据 · 金额单位：人民币</div>
    </div>
</x-filament-panels::page>
