<?php

namespace App\Filament\Pages;

use App\Filament\Resources\Categories\CategoryResource;
use App\Filament\Resources\Media\MediaResource;
use App\Filament\Resources\Orders\OrderResource;
use App\Filament\Resources\Products\ProductResource;
use App\Filament\Resources\Users\UserResource;
use App\Models\Order;
use App\Models\Product;
use App\Models\User;
use Filament\Pages\Dashboard as BaseDashboard;

class Dashboard extends BaseDashboard
{
    protected static string|\BackedEnum|null $navigationIcon = 'heroicon-o-squares-2x2';

    protected static ?string $navigationLabel = '工作台';

    protected static ?string $title = '工作台';

    protected string $view = 'filament.pages.dashboard';

    public function getHeading(): ?string
    {
        return null;
    }

    protected function getViewData(): array
    {
        $canViewOrders = OrderResource::canViewAny();
        $canViewProducts = ProductResource::canViewAny();
        $metrics = [];
        $orderStages = [];
        $shortcuts = [];

        if ($canViewOrders) {
            $orderCount = Order::count();
            $metrics[] = ['label' => '成交金额', 'value' => '¥ ' . number_format((float) Order::whereIn('status', ['paid', 'shipped', 'completed'])->sum('total_amount'), 2), 'note' => '已付款、已发货及已完成订单', 'icon' => 'heroicon-o-banknotes', 'tone' => 'blue'];
            $metrics[] = ['label' => '订单总量', 'value' => number_format($orderCount), 'note' => '全部订单，包含已取消订单', 'icon' => 'heroicon-o-shopping-bag', 'tone' => 'teal'];
            foreach (['pending' => '待付款', 'paid' => '待发货', 'shipped' => '待收货', 'completed' => '已完成', 'cancelled' => '已取消'] as $status => $label) {
                $count = Order::where('status', $status)->count();
                $orderStages[] = ['label' => $label, 'count' => $count, 'status' => $status, 'percent' => $orderCount > 0 ? round($count / $orderCount * 100) : 0, 'url' => OrderResource::getUrl('index', ['tab' => $status])];
            }
            $shortcuts[] = ['label' => '订单管理', 'note' => '查询订单与处理发货', 'icon' => 'heroicon-o-shopping-bag', 'url' => OrderResource::getUrl()];
        }

        if ($canViewProducts) {
            $metrics[] = ['label' => '在售商品', 'value' => number_format(Product::where('status', 'active')->count()), 'note' => '当前已上架商品', 'icon' => 'heroicon-o-cube', 'tone' => 'purple'];
            $shortcuts[] = ['label' => ProductResource::canCreate() ? '新增商品' : '商品管理', 'note' => '维护商品、价格与库存', 'icon' => 'heroicon-o-cube', 'url' => ProductResource::getUrl(ProductResource::canCreate() ? 'create' : 'index')];
        }
        if (UserResource::canViewAny()) {
            $metrics[] = ['label' => '客户总数', 'value' => number_format(User::count()), 'note' => '已注册的前台用户', 'icon' => 'heroicon-o-users', 'tone' => 'yellow'];
            $shortcuts[] = ['label' => '客户管理', 'note' => '查看客户资料与订单', 'icon' => 'heroicon-o-users', 'url' => UserResource::getUrl()];
        }
        if (CategoryResource::canViewAny()) {
            $shortcuts[] = ['label' => '商品分类', 'note' => '整理商品分类结构', 'icon' => 'heroicon-o-folder', 'url' => CategoryResource::getUrl()];
        }
        if (MediaResource::canViewAny()) {
            $shortcuts[] = ['label' => '媒体资源', 'note' => '管理图片与文件', 'icon' => 'heroicon-o-photo', 'url' => MediaResource::getUrl()];
        }
        if (ManageSettings::canAccess()) {
            $shortcuts[] = ['label' => '系统设置', 'note' => '维护站点与订单配置', 'icon' => 'heroicon-o-cog-6-tooth', 'url' => ManageSettings::getUrl()];
        }

        return [
            'metrics' => $metrics,
            'orderStages' => $orderStages,
            'shortcuts' => $shortcuts,
            'canViewOrders' => $canViewOrders,
            'canViewProducts' => $canViewProducts,
            'latestOrders' => $canViewOrders ? Order::with('user')->latest()->limit(5)->get() : [],
            'lowStockCount' => $canViewProducts ? Product::where('status', 'active')->where('stock', '<=', 10)->count() : 0,
            'lowStockProducts' => $canViewProducts ? Product::where('status', 'active')->where('stock', '<=', 10)->orderBy('stock')->limit(5)->get() : [],
        ];
    }
}
