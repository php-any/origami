<?php
    use Filament\Support\Enums\IconPosition;
    use Filament\Support\Facades\FilamentAsset;
    use Filament\Support\View\ComponentAttributeBag as FilamentComponentAttributeBag;
    use Filament\Widgets\View\Components\StatsOverviewWidgetComponent\StatComponent\DescriptionComponent;
    use Filament\Widgets\View\Components\StatsOverviewWidgetComponent\StatComponent\StatsOverviewWidgetStatChartComponent;

    $chartColor = $getChartColor() ?? 'gray';
    $descriptionColor = $getDescriptionColor() ?? 'gray';
    $descriptionIcon = $getDescriptionIcon();
    $descriptionIconPosition = $getDescriptionIconPosition();
    $url = $getUrl();
    $tag = $url ? 'a' : 'div';
?>

<<?php echo $tag; ?>

    <?php if($url): ?>
        <?php echo e(\Filament\Support\generate_href_html($url, $shouldOpenUrlInNewTab())); ?>

    <?php endif; ?>
    <?php echo e($getExtraAttributeBag()
            ->class([
                'fi-wi-stats-overview-stat',
            ])); ?>

>
    <div class="fi-wi-stats-overview-stat-content">
        <div class="fi-wi-stats-overview-stat-label-ctn">
            <?php echo e(\Filament\Support\generate_icon_html($getIcon())); ?>


            <span class="fi-wi-stats-overview-stat-label">
                <?php echo e($getLabel()); ?>

            </span>
        </div>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($value = $getValue())): ?>
            <div class="fi-wi-stats-overview-stat-value">
                <?php echo e($value); ?>

            </div>
        <?php else: ?>
            <div class="fi-wi-stats-overview-stat-placeholder">
                <?php echo e($getPlaceholder()); ?>

            </div>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($description = $getDescription()): ?>
            <div
                <?php echo e((new FilamentComponentAttributeBag)->color(DescriptionComponent::class, $descriptionColor)->class(['fi-wi-stats-overview-stat-description'])); ?>

            >
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($descriptionIcon && in_array($descriptionIconPosition, [IconPosition::Before, 'before'])): ?>
                    <?php echo e(\Filament\Support\generate_icon_html($descriptionIcon, attributes: (new FilamentComponentAttributeBag))); ?>

                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                <span>
                    <?php echo e($description); ?>

                </span>

                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($descriptionIcon && in_array($descriptionIconPosition, [IconPosition::After, 'after'])): ?>
                    <?php echo e(\Filament\Support\generate_icon_html($descriptionIcon, attributes: (new FilamentComponentAttributeBag))); ?>

                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
            </div>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    </div>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($chart = $getChart()): ?>
        
        <div x-data="{ statsOverviewStatChart() {} }">
            <div
                x-load
                x-load-src="<?php echo e(FilamentAsset::getAlpineComponentSrc('stats-overview/stat/chart', 'filament/widgets')); ?>"
                wire:ignore
                x-data="statsOverviewStatChart({
                            key: <?php echo \Illuminate\Support\Js::from($getKey(false))->toHtml() ?>,
                            labels: <?php echo \Illuminate\Support\Js::from(array_keys($chart))->toHtml() ?>,
                            values: <?php echo \Illuminate\Support\Js::from(array_values($chart))->toHtml() ?>,
                        })"
                <?php echo e((new FilamentComponentAttributeBag)->color(StatsOverviewWidgetStatChartComponent::class, $chartColor)->class(['fi-wi-stats-overview-stat-chart'])); ?>

            >
                
                <canvas x-ref="canvas" aria-hidden="true"></canvas>

                
                <span
                    aria-hidden="true"
                    x-ref="backgroundColorElement"
                    class="fi-wi-stats-overview-stat-chart-bg-color"
                ></span>

                <span
                    aria-hidden="true"
                    x-ref="borderColorElement"
                    class="fi-wi-stats-overview-stat-chart-border-color"
                ></span>
            </div>
        </div>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
</<?php echo $tag; ?>>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\widgets\resources\views/stats-overview-widget/stat.blade.php ENDPATH**/ ?>