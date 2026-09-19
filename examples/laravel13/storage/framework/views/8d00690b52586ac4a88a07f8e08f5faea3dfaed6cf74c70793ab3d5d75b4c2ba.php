<?php $attributes ??= new \Illuminate\View\ComponentAttributeBag;

$__newAttributes = [];
$__propNames = \Illuminate\View\ComponentAttributeBag::extractPropNames(([
    'columns' => null,
    'hasReorderableColumns',
    'hasToggleableColumns',
    'reorderAnimationDuration' => 300,
]));

foreach ($attributes->all() as $__key => $__value) {
    if (in_array($__key, $__propNames)) {
        $$__key = $$__key ?? $__value;
    } else {
        $__newAttributes[$__key] = $__value;
    }
}

$attributes = new \Illuminate\View\ComponentAttributeBag($__newAttributes);

unset($__propNames);
unset($__newAttributes);

foreach (array_filter(([
    'columns' => null,
    'hasReorderableColumns',
    'hasToggleableColumns',
    'reorderAnimationDuration' => 300,
]), 'is_string', ARRAY_FILTER_USE_KEY) as $__key => $__value) {
    $$__key = $$__key ?? $__value;
}

$__defined_vars = get_defined_vars();

foreach ($attributes->all() as $__key => $__value) {
    if (array_key_exists($__key, $__defined_vars)) unset($$__key);
}

unset($__defined_vars, $__key, $__value); ?>

<?php
    use Filament\Support\Enums\GridDirection;
    use Filament\Support\Icons\Heroicon;
    use Filament\Support\View\ComponentAttributeBag as FilamentComponentAttributeBag;
    use Filament\Tables\View\TablesIconAlias;
?>

<div
    <?php if($hasToggleableColumns): ?>
        
        x-id="['fi-ta-col-manager-group-checkbox', 'fi-ta-col-manager-column-checkbox']"
    <?php endif; ?>
    <?php if($hasReorderableColumns): ?>
        x-sortable
        x-on:end.stop="reorderColumns($event.target.sortable.toArray())"
        data-sortable-animation-duration="<?php echo e($reorderAnimationDuration); ?>"
    <?php endif; ?>
    <?php echo e((new FilamentComponentAttributeBag)
            ->grid($columns, GridDirection::Column)
            ->class(['fi-ta-col-manager-items'])); ?>

>
    <template
        x-for="(column, index) in columns.filter((column) => ! column.isHidden && column.label)"
        x-bind:key="(column.type === 'group' ? 'group::' : 'column::') + column.name + '_' + index"
    >
        <div
            <?php if($hasReorderableColumns): ?>
                x-bind:x-sortable-item="column.type === 'group' ? 'group::' + column.name : 'column::' + column.name"
            <?php endif; ?>
        >
            <template x-if="column.type === 'group'">
                <div class="fi-ta-col-manager-group">
                    <div class="fi-ta-col-manager-item">
                        <label
                            <?php if($hasToggleableColumns): ?> x-bind:for="$id('fi-ta-col-manager-group-checkbox', column.name)" <?php endif; ?>
                            class="fi-ta-col-manager-label"
                        >
                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasToggleableColumns): ?>
                                <input
                                    type="checkbox"
                                    class="fi-checkbox-input fi-valid"
                                    x-bind:id="$id('fi-ta-col-manager-group-checkbox', column.name)"
                                    x-bind:checked="(groupedColumns[column.name] || {}).checked || false"
                                    x-bind:disabled="(groupedColumns[column.name] || {}).disabled || false"
                                    x-effect="$el.indeterminate = (groupedColumns[column.name] || {}).indeterminate || false"
                                    x-on:change="toggleGroup(column.name)"
                                />
                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                            <span x-html="column.label"></span>
                        </label>

                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasReorderableColumns): ?>
                            <button
                                x-sortable-handle
                                x-bind:aria-label="<?php echo \Illuminate\Support\Js::from(__('filament-tables::table.column_manager.actions.reorder.label'))->toHtml() ?> + (column.name ? ' ' + column.name : '')"
                                x-on:click.stop
                                class="fi-ta-col-manager-reorder-handle fi-icon-btn"
                                type="button"
                            >
                                <?php echo e(\Filament\Support\generate_icon_html(Heroicon::Bars2, alias: TablesIconAlias::REORDER_HANDLE)); ?>

                            </button>
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                    </div>
                    <div
                        <?php if($hasReorderableColumns): ?>
                            x-sortable
                            x-on:end.stop="reorderGroupColumns($event.target.sortable.toArray(), column.name)"
                            data-sortable-animation-duration="<?php echo e($reorderAnimationDuration); ?>"
                        <?php endif; ?>
                        class="fi-ta-col-manager-group-items"
                    >
                        <template
                            x-for="
                                (groupColumn, index) in
                                    column.columns.filter((column) => ! column.isHidden && column.label)
                            "
                            x-bind:key="'column::' + groupColumn.name + '_' + index"
                        >
                            <div
                                <?php if($hasReorderableColumns): ?>
                                    x-bind:x-sortable-item="'column::' + groupColumn.name"
                                <?php endif; ?>
                            >
                                <div class="fi-ta-col-manager-item">
                                    <label
                                        <?php if($hasToggleableColumns): ?> x-bind:for="$id('fi-ta-col-manager-column-checkbox', groupColumn.name)" <?php endif; ?>
                                        class="fi-ta-col-manager-label"
                                    >
                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasToggleableColumns): ?>
                                            <input
                                                type="checkbox"
                                                class="fi-checkbox-input fi-valid"
                                                x-bind:id="$id('fi-ta-col-manager-column-checkbox', groupColumn.name)"
                                                x-bind:checked="(getColumn(groupColumn.name, column.name) || {}).isToggled || false"
                                                x-bind:disabled="(getColumn(groupColumn.name, column.name) || {}).isToggleable === false"
                                                x-on:change="toggleColumn(groupColumn.name, column.name)"
                                            />
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                        <span
                                            x-html="groupColumn.label"
                                        ></span>
                                    </label>

                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasReorderableColumns): ?>
                                        <button
                                            x-sortable-handle
                                            x-bind:aria-label="<?php echo \Illuminate\Support\Js::from(__('filament-tables::table.column_manager.actions.reorder.label'))->toHtml() ?> + (groupColumn.name ? ' ' + groupColumn.name : '')"
                                            x-on:click.stop
                                            class="fi-ta-col-manager-reorder-handle fi-icon-btn"
                                            type="button"
                                        >
                                            <?php echo e(\Filament\Support\generate_icon_html(Heroicon::Bars2, alias: TablesIconAlias::REORDER_HANDLE)); ?>

                                        </button>
                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                </div>
                            </div>
                        </template>
                    </div>
                </div>
            </template>
            <template x-if="column.type !== 'group'">
                <div class="fi-ta-col-manager-item">
                    <label
                        <?php if($hasToggleableColumns): ?> x-bind:for="$id('fi-ta-col-manager-column-checkbox', column.name)" <?php endif; ?>
                        class="fi-ta-col-manager-label"
                    >
                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasToggleableColumns): ?>
                            <input
                                type="checkbox"
                                class="fi-checkbox-input fi-valid"
                                x-bind:id="$id('fi-ta-col-manager-column-checkbox', column.name)"
                                x-bind:checked="(getColumn(column.name, null) || {}).isToggled || false"
                                x-bind:disabled="(getColumn(column.name, null) || {}).isToggleable === false"
                                x-on:change="toggleColumn(column.name)"
                            />
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                        <span x-html="column.label"></span>
                    </label>

                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasReorderableColumns): ?>
                        <button
                            x-sortable-handle
                            x-bind:aria-label="<?php echo \Illuminate\Support\Js::from(__('filament-tables::table.column_manager.actions.reorder.label'))->toHtml() ?> + (column.name ? ' ' + column.name : '')"
                            x-on:click.stop
                            class="fi-ta-col-manager-reorder-handle fi-icon-btn"
                            type="button"
                        >
                            <?php echo e(\Filament\Support\generate_icon_html(Heroicon::Bars2, alias: TablesIconAlias::REORDER_HANDLE)); ?>

                        </button>
                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                </div>
            </template>
        </div>
    </template>
</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\tables\resources\views/components/column-manager/content.blade.php ENDPATH**/ ?>