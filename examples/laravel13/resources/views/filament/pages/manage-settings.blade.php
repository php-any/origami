<x-filament::page class="origami-settings">
    <p class="origami-settings-intro">管理站点信息、访问状态和订单编号规则。</p>
    <div class="origami-settings-panel">
        <div class="origami-settings-panel-heading">
            <div><x-filament::icon icon="heroicon-o-adjustments-horizontal" class="icon" /><h2>常规配置</h2></div>
            <span>站点与订单</span>
        </div>
        {{ $this->form }}
    </div>
</x-filament::page>
