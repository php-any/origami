<?php
    use Filament\Actions\Action;
    use Filament\Actions\ActionGroup;
    use Filament\Actions\BulkAction;
    use Filament\Actions\BulkActionGroup;
    use Filament\Support\Enums\Alignment;
    use Filament\Support\Enums\IconSize;
    use Filament\Support\Enums\VerticalAlignment;
    use Filament\Support\Enums\Width;
    use Filament\Support\Facades\FilamentView;
    use Filament\Support\Icons\Heroicon;
    use Filament\Support\View\ComponentAttributeBag as FilamentComponentAttributeBag;
    use Filament\Tables\Actions\HeaderActionsPosition;
    use Filament\Tables\Columns\Column;
    use Filament\Tables\Columns\ColumnGroup;
    use Filament\Tables\Enums\ColumnManagerLayout;
    use Filament\Tables\Enums\ColumnManagerResetActionPosition;
    use Filament\Tables\Enums\FiltersLayout;
    use Filament\Tables\Enums\FiltersResetActionPosition;
    use Filament\Tables\Enums\RecordActionsPosition;
    use Filament\Tables\Enums\RecordCheckboxPosition;
    use Filament\Tables\Filters\Indicator;
    use Filament\Tables\Table;
    use Filament\Tables\View\TablesIconAlias;
    use Filament\Tables\View\TablesRenderHook;
    use Illuminate\Contracts\Pagination\CursorPaginator;
    use Illuminate\Contracts\Pagination\LengthAwarePaginator;
    use Illuminate\Contracts\Pagination\Paginator;
    use Illuminate\Contracts\Support\Htmlable;
    use Illuminate\Database\Query\Builder;
    use Illuminate\Support\Number;
    use Illuminate\Support\Str;
    use Illuminate\View\ComponentAttributeBag;

    $defaultRecordActions = $getRecordActions();
    $flatRecordActionsCount = count($getFlatRecordActions());
    $recordActionsAlignment = $getRecordActionsAlignment();
    $recordActionsPosition = $getRecordActionsPosition();
    $recordActionsColumnLabel = $getRecordActionsColumnLabel();

    if (! $recordActionsAlignment instanceof Alignment) {
        $recordActionsAlignment = filled($recordActionsAlignment) ? (Alignment::tryFrom($recordActionsAlignment) ?? $recordActionsAlignment) : null;
    }

    $activeFiltersCount = $getActiveFiltersCount();
    $isSelectionDisabled = $isSelectionDisabled();
    $maxSelectableRecords = $getMaxSelectableRecords();
    $columns = $getVisibleColumns();
    $collapsibleColumnsLayout = $getCollapsibleColumnsLayout();
    $columnsLayout = $getColumnsLayout();
    $content = $getContent();
    $contentGrid = $getContentGrid();
    $contentFooter = $getContentFooter();
    $filterIndicators = $getFilterIndicators();
    $filtersApplyAction = $getFiltersApplyAction();
    $filtersForm = $getFiltersForm();
    $filtersFormWidth = $getFiltersFormWidth();
    $filtersResetActionPosition = $getFiltersResetActionPosition();
    $columnManagerResetActionPosition = $getColumnManagerResetActionPosition();
    $hasColumnGroups = $hasColumnGroups();
    $hasColumnsLayout = $hasColumnsLayout();
    $hasPageSummary = $hasPageSummary();
    $hasAllTableSummary = $hasAllTableSummary();
    $hasSummary = $hasSummary($this->getAllTableSummaryQuery());
    $hasTopLevelSummary = $hasSummary && ($hasPageSummary || $hasAllTableSummary);
    $header = $getHeader();
    $headerActions = array_filter(
        $getHeaderActions(),
        fn (\Filament\Actions\Action | ActionGroup $action): bool => $action->isVisible(),
    );
    $headerActionsPosition = $getHeaderActionsPosition();
    $heading = $getHeading();
    $group = $getGrouping();
    $toolbarActions = array_filter(
        $getToolbarActions(),
        fn (\Filament\Actions\Action | ActionGroup $action): bool => $action->isVisible(),
    );

    $hasNonBulkToolbarAction = false;

    foreach ($toolbarActions as $toolbarAction) {
        if ($toolbarAction instanceof BulkActionGroup) {
            continue;
        }

        if ($toolbarAction instanceof ActionGroup) {
            if ($toolbarAction->hasNonBulkAction()) {
                $hasNonBulkToolbarAction = true;

                break;
            }

            continue;
        }

        if (! $toolbarAction->isBulk()) {
            $hasNonBulkToolbarAction = true;

            break;
        }
    }

    $groups = $getGroups();
    $description = $getDescription();
    $isGroupsOnly = $isGroupsOnly() && $group;
    $isReorderable = $isReorderable();
    $isReordering = $isReordering();
    $areGroupingSettingsVisible = (! $isReordering) && count($groups) && (! $areGroupingSettingsHidden());
    $isGroupingDirectionSettingHidden = $isGroupingDirectionSettingHidden();
    $areGroupsCollapsedByDefault = $areGroupsCollapsedByDefault();
    $areGroupingSettingsInDropdownOnDesktop = $areGroupingSettingsInDropdownOnDesktop();
    $isColumnSearchVisible = $isSearchableByColumn();
    $isGlobalSearchVisible = $isSearchable();
    $isSearchOnBlur = $isSearchOnBlur();
    $isSelectionEnabled = $isSelectionEnabled() && (! $isGroupsOnly);
    $selectsCurrentPageOnly = $selectsCurrentPageOnly();
    $selectsGroupsOnly = $selectsGroupsOnly();
    $recordCheckboxPosition = $getRecordCheckboxPosition();
    $isStriped = $isStriped();
    $isStackedOnMobile = $isStackedOnMobile();
    $isLoaded = $isLoaded();
    $hasFilters = $isFilterable();
    $filtersLayout = $getFiltersLayout();
    $filtersTriggerAction = $getFiltersTriggerAction();
    $hasFiltersDialog = $hasFilters && in_array($filtersLayout, [FiltersLayout::Dropdown, FiltersLayout::Modal]);
    $hasFiltersAboveContent = $hasFilters && in_array($filtersLayout, [FiltersLayout::AboveContent, FiltersLayout::AboveContentCollapsible]);
    $hasFiltersBelowContent = $hasFilters && ($filtersLayout === FiltersLayout::BelowContent);
    $hasFiltersBeforeContent = $hasFilters && in_array($filtersLayout, [FiltersLayout::BeforeContent, FiltersLayout::BeforeContentCollapsible]);
    $hasFiltersAfterContent = $hasFilters && in_array($filtersLayout, [FiltersLayout::AfterContent, FiltersLayout::AfterContentCollapsible]);
    $hasCollapsibleFilters = $hasFilters && in_array($filtersLayout, [FiltersLayout::AboveContentCollapsible, FiltersLayout::BeforeContentCollapsible, FiltersLayout::AfterContentCollapsible]);
    $hasFiltersTrigger = $hasFilters && ($hasFiltersDialog || $hasFiltersBeforeContent || $hasFiltersAfterContent);
    $filtersFormMaxHeight = $getFiltersFormMaxHeight();
    $hasColumnManager = $hasColumnManager();
    $columnManagerLayout = $getColumnManagerLayout();
    $hasReorderableColumns = $hasReorderableColumns();
    $hasToggleableColumns = $hasToggleableColumns();
    $columnManagerApplyAction = $getColumnManagerApplyAction();
    $columnManagerTriggerAction = $getColumnManagerTriggerAction();
    $hasHeader = $header || $heading || $description || ($headerActions && (! $isReordering)) || $isReorderable || $areGroupingSettingsVisible || $isGlobalSearchVisible || $hasFilters || count($filterIndicators) || $hasColumnManager;
    $hasHeaderToolbar = $isReorderable || $areGroupingSettingsVisible || $isGlobalSearchVisible || $hasFiltersTrigger || $hasColumnManager;

    // https://github.com/filamentphp/filament/pull/19787
    $headerVisibilityMode = ($hasHeader || $hasNonBulkToolbarAction)
        ? 'visible'
        : (count($toolbarActions) ? 'selection' : 'hidden');
    $headerToolbarVisibilityMode = ($hasHeaderToolbar || $hasNonBulkToolbarAction)
        ? 'visible'
        : (count($toolbarActions) ? 'selection' : 'hidden');
    $headingTag = $getHeadingTag();
    $secondLevelHeadingTag = $heading ? $getHeadingTag(1) : $headingTag;
    $pluralModelLabel = $getPluralModelLabel();
    $records = $isLoaded ? $getRecords() : null;
    $hasPagination = (($records instanceof Paginator) || ($records instanceof CursorPaginator)) && (($records instanceof LengthAwarePaginator) ? $records->total() : $records->isNotEmpty());
    $hasEmptyState = ($records !== null) && ! count($records);
    $hasContentLayout = $content || $hasColumnsLayout;
    $searchDebounce = $getSearchDebounce();
    $allSelectableRecordsCount = ($isSelectionEnabled && $isLoaded) ? $getAllSelectableRecordsCount() : null;
    $columnsCount = count($columns);
    $reorderRecordsTriggerAction = $getReorderRecordsTriggerAction($isReordering);
    $page = $this->getTablePage();
    $defaultSortOptionLabel = $getDefaultSortOptionLabel();
    $sortDirection = $getSortDirection();

    $reduceVisibleRecordActions = function ($record) use ($defaultRecordActions): array {
        return array_reduce(
            $defaultRecordActions,
            function (array $carry, $action) use ($record): array {
                $action = $action->getClone();

                if (! $action instanceof BulkAction) {
                    $action->record($record);
                }

                if ($action->isHidden()) {
                    return $carry;
                }

                $carry[] = $action;

                return $carry;
            },
            initial: [],
        );
    };

    $recordActionsByRecordKey = [];
    $hasRecordActionsForAnyRecord = false;

    // Determine whether the record actions column should be rendered by scanning
    // records until one exposes a visible action. Every record that gets checked
    // is cached so the row loop below never re-evaluates its visibility.
    if (($records !== null) && (! $isReordering)) {
        foreach ($records as $record) {
            $recordActionsByRecordKey[$getRecordKey($record)] = $currentRecordActions = $reduceVisibleRecordActions($record);

            if ($currentRecordActions !== []) {
                $hasRecordActionsForAnyRecord = true;

                break;
            }
        }
    }

    if ($hasRecordActionsForAnyRecord && (! $isReordering)) {
        $columnsCount++;
    }

    if ($isSelectionEnabled || $isReordering) {
        $columnsCount++;
    }

    if ($group) {
        $groupedSummarySelectedState = $this->getTableSummarySelectedState($this->getAllTableSummaryQuery(), modifyQueryUsing: fn (Builder $query) => $group->groupQuery($query, model: $getQuery()->getModel()));
    }

    if (is_string($filtersFormWidth)) {
        $filtersFormWidth = Width::tryFrom($filtersFormWidth) ?? $filtersFormWidth;
    }

    $loadingTargetsWireTarget = implode(',', Table::LOADING_TARGETS);
?>

<div
    <?php if(! $isLoaded): ?>
        wire:init="loadTable"
    <?php endif; ?>
    x-data="filamentTable({
                areGroupsCollapsedByDefault: <?php echo \Illuminate\Support\Js::from($areGroupsCollapsedByDefault)->toHtml() ?>,
                canTrackDeselectedRecords: <?php echo \Illuminate\Support\Js::from($canTrackDeselectedRecords())->toHtml() ?>,
                currentSelectionLivewireProperty: <?php echo \Illuminate\Support\Js::from($getCurrentSelectionLivewireProperty())->toHtml() ?>,
                maxSelectableRecords: <?php echo \Illuminate\Support\Js::from($maxSelectableRecords)->toHtml() ?>,
                selectsCurrentPageOnly: <?php echo \Illuminate\Support\Js::from($selectsCurrentPageOnly)->toHtml() ?>,
                $wire,
            })"
    <?php echo e($getExtraAttributeBag()->class([
            'fi-ta',
            'fi-loading' => $records === null,
        ])); ?>

>
    <input
        type="hidden"
        value="<?php echo e($allSelectableRecordsCount); ?>"
        x-ref="allSelectableRecordsCount"
    />

    <div
        class="<?php echo \Illuminate\Support\Arr::toCssClasses([
            'fi-ta-ctn',
            'fi-ta-ctn-with-content-layout' => $hasContentLayout,
            'fi-ta-ctn-with-footer' => $hasPagination || $hasEmptyState || $hasFiltersBelowContent,
            'fi-ta-ctn-with-header' => $hasHeader,
        ]); ?>"
    >
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasFiltersBeforeContent): ?>
            <div
                wire:ignore.self
                x-ref="filtersContentContainer"
                x-transition:enter-start="fi-opacity-0"
                x-transition:leave-end="fi-opacity-0"
                x-bind:class="{ 'fi-open': areFiltersOpen }"
                class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                    'fi-ta-filters-before-content-ctn',
                    'lg:fi-open' => ! $hasCollapsibleFilters,
                    (($filtersFormWidth ??= Width::ExtraSmall) instanceof Width) ? "fi-width-{$filtersFormWidth->value}" : (is_string($filtersFormWidth) ? $filtersFormWidth : null),
                ]); ?>"
            >
                <?php if (isset($component)) { $__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.filters','data' => ['applyAction' => $filtersApplyAction,'form' => $filtersForm,'headingTag' => $secondLevelHeadingTag,'class' => 'fi-ta-filters-before-content','resetActionPosition' => $filtersResetActionPosition]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::filters'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['apply-action' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersApplyAction),'form' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersForm),'heading-tag' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($secondLevelHeadingTag),'class' => 'fi-ta-filters-before-content','reset-action-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersResetActionPosition)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b)): ?>
<?php $attributes = $__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b; ?>
<?php unset($__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b)): ?>
<?php $component = $__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b; ?>
<?php unset($__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b); ?>
<?php endif; ?>
            </div>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <div class="fi-ta-main">
            <div
                <?php if(! $hasHeader): ?> x-cloak <?php endif; ?>
                x-show="<?php echo \Illuminate\Support\Js::from($hasHeader)->toHtml() ?> || <?php echo \Illuminate\Support\Js::from($hasNonBulkToolbarAction)->toHtml() ?> || (getSelectedRecordsCount() && <?php echo \Illuminate\Support\Js::from(count($toolbarActions))->toHtml() ?>)"
                <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.header.'.e($headerVisibilityMode).''; ?>wire:key="<?php echo e($this->getId()); ?>.table.header.<?php echo e($headerVisibilityMode); ?>"
                class="fi-ta-header-ctn"
            >
                <?php echo e(FilamentView::renderHook(TablesRenderHook::HEADER_BEFORE, scopes: static::class)); ?>


                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($header): ?>
                    <?php echo e($header); ?>

                <?php elseif($heading || $description || ($headerActions && (! $isReordering))): ?>
                    <div
                        class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                            'fi-ta-header',
                            'fi-ta-header-adaptive-actions-position' => $headerActions && (! $isReordering) && ($headerActionsPosition === HeaderActionsPosition::Adaptive),
                        ]); ?>"
                    >
                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($heading || $description): ?>
                            <div>
                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($heading): ?>
                                    <<?php echo e($headingTag); ?>

                                        class="fi-ta-header-heading"
                                    >
                                        <?php echo e($heading); ?>

                                    </<?php echo e($headingTag); ?>>
                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($description): ?>
                                    <p class="fi-ta-header-description">
                                        <?php echo e($description); ?>

                                    </p>
                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                            </div>
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if((! $isReordering) && $headerActions): ?>
                            <div
                                class="fi-ta-actions fi-align-start fi-wrapped"
                            >
                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $headerActions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $action): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                    <?php echo e($action); ?>

                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                            </div>
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                    </div>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                <?php echo e(FilamentView::renderHook(TablesRenderHook::HEADER_AFTER, scopes: static::class)); ?>


                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasFiltersAboveContent): ?>
                    <div
                        <?php if($hasCollapsibleFilters): ?>
                            x-bind:class="{ 'fi-open': areFiltersOpen }"
                        <?php endif; ?>
                        class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                            'fi-ta-filters-above-content-ctn',
                        ]); ?>"
                    >
                        <?php if (isset($component)) { $__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.filters','data' => ['applyAction' => $filtersApplyAction,'form' => $filtersForm,'headingTag' => $secondLevelHeadingTag,'xCloak' => true,'xShow' => $hasCollapsibleFilters ? 'areFiltersOpen' : null,'resetActionPosition' => $filtersResetActionPosition]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::filters'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['apply-action' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersApplyAction),'form' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersForm),'heading-tag' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($secondLevelHeadingTag),'x-cloak' => true,'x-show' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasCollapsibleFilters ? 'areFiltersOpen' : null),'reset-action-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersResetActionPosition)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b)): ?>
<?php $attributes = $__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b; ?>
<?php unset($__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b)): ?>
<?php $component = $__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b; ?>
<?php unset($__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b); ?>
<?php endif; ?>

                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasCollapsibleFilters): ?>
                            <span
                                x-on:click="areFiltersOpen = ! areFiltersOpen"
                                class="fi-ta-filters-trigger-action-ctn"
                            >
                                <?php echo e($filtersTriggerAction->badge($activeFiltersCount)); ?>

                            </span>
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                    </div>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                <?php echo e(FilamentView::renderHook(TablesRenderHook::TOOLBAR_BEFORE, scopes: static::class)); ?>


                <div
                    <?php if(! $hasHeaderToolbar): ?> x-cloak <?php endif; ?>
                    x-show="<?php echo \Illuminate\Support\Js::from($hasHeaderToolbar)->toHtml() ?> || <?php echo \Illuminate\Support\Js::from($hasNonBulkToolbarAction)->toHtml() ?> || (getSelectedRecordsCount() && <?php echo \Illuminate\Support\Js::from(count($toolbarActions))->toHtml() ?>)"
                    <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.header-toolbar.'.e($headerToolbarVisibilityMode).''; ?>wire:key="<?php echo e($this->getId()); ?>.table.header-toolbar.<?php echo e($headerToolbarVisibilityMode); ?>"
                    class="fi-ta-header-toolbar"
                >
                    <?php echo e(FilamentView::renderHook(TablesRenderHook::TOOLBAR_START, scopes: static::class)); ?>


                    <div class="fi-ta-actions fi-align-start fi-wrapped">
                        <?php echo e(FilamentView::renderHook(TablesRenderHook::TOOLBAR_REORDER_TRIGGER_BEFORE, scopes: static::class)); ?>


                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isReorderable): ?>
                            <?php echo e($reorderRecordsTriggerAction); ?>

                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                        <?php echo e(FilamentView::renderHook(TablesRenderHook::TOOLBAR_REORDER_TRIGGER_AFTER, scopes: static::class)); ?>


                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if((! $isReordering) && count($toolbarActions)): ?>
                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $toolbarActions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $action): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                <?php echo e($action); ?>

                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                        <?php echo e(FilamentView::renderHook(TablesRenderHook::TOOLBAR_GROUPING_SELECTOR_BEFORE, scopes: static::class)); ?>


                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($areGroupingSettingsVisible): ?>
                            <div
                                x-data="{
                                    grouping: $wire.$entangle('tableGrouping', true),
                                    group: null,
                                    direction: null,
                                }"
                                x-init="
                                    if (grouping) {
                                        ;[group, direction] = grouping.split(':')
                                        direction ??= 'asc'
                                    }

                                    $watch('grouping', function () {
                                        if (! grouping) {
                                            group = null
                                            direction = null

                                            return
                                        }

                                        ;[group, direction] = grouping.split(':')
                                        direction ??= 'asc'
                                    })

                                    $watch('direction', function () {
                                        grouping = group ? `${group}:${direction}` : null
                                    })

                                    $watch('group', function (newGroup, oldGroup) {
                                        if (! newGroup) {
                                            direction = null
                                            grouping = group ? `${group}:${direction}` : null

                                            return
                                        }

                                        if (oldGroup) {
                                            grouping = group ? `${group}:${direction}` : null

                                            return
                                        }

                                        direction ??= 'asc'
                                        grouping = group ? `${group}:${direction}` : null
                                    })
                                "
                                class="fi-ta-grouping-settings"
                            >
                                <?php if (isset($component)) { $__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.dropdown.index','data' => ['placement' => 'bottom-start','shift' => true,'width' => 'xs','wire:key' => $this->getId() . '.table.grouping','class' => \Illuminate\Support\Arr::toCssClasses([
                                        'sm:fi-hidden' => ! $areGroupingSettingsInDropdownOnDesktop,
                                    ])]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::dropdown'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['placement' => 'bottom-start','shift' => true,'width' => 'xs','wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . '.table.grouping'),'class' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(\Illuminate\Support\Arr::toCssClasses([
                                        'sm:fi-hidden' => ! $areGroupingSettingsInDropdownOnDesktop,
                                    ]))]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                     <?php $__env->slot('trigger', null, []); ?> 
                                        <?php echo e($getGroupRecordsTriggerAction()); ?>

                                     <?php $__env->endSlot(); ?>

                                    <div class="fi-ta-grouping-settings-fields">
                                        <label>
                                            <span>
                                                <?php echo e(__('filament-tables::table.grouping.fields.group.label')); ?>

                                            </span>

                                            <?php if (isset($component)) { $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.wrapper','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.wrapper'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                <?php if (isset($component)) { $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.select','data' => ['xModel' => 'group','xOn:change' => 'resetCollapsedGroups()']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.select'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['x-model' => 'group','x-on:change' => 'resetCollapsedGroups()']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                    <option value="">-</option>

                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $groups; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $groupOption): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                        <option
                                                            value="<?php echo e($groupOption->getId()); ?>"
                                                        >
                                                            <?php echo e($groupOption->getLabel()); ?>

                                                        </option>
                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $attributes = $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $component = $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
                                             <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $attributes = $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $component = $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
                                        </label>

                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $isGroupingDirectionSettingHidden): ?>
                                            <label x-cloak x-show="group">
                                                <span>
                                                    <?php echo e(__('filament-tables::table.grouping.fields.direction.label')); ?>

                                                </span>

                                                <?php if (isset($component)) { $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.wrapper','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.wrapper'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                    <?php if (isset($component)) { $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.select','data' => ['xModel' => 'direction']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.select'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['x-model' => 'direction']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                        <option value="asc">
                                                            <?php echo e(__('filament-tables::table.grouping.fields.direction.options.asc')); ?>

                                                        </option>

                                                        <option value="desc">
                                                            <?php echo e(__('filament-tables::table.grouping.fields.direction.options.desc')); ?>

                                                        </option>
                                                     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $attributes = $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $component = $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
                                                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $attributes = $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $component = $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
                                            </label>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                    </div>
                                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b)): ?>
<?php $attributes = $__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b; ?>
<?php unset($__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b)): ?>
<?php $component = $__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b; ?>
<?php unset($__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b); ?>
<?php endif; ?>

                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $areGroupingSettingsInDropdownOnDesktop): ?>
                                    <div class="fi-ta-grouping-settings-fields">
                                        <label>
                                            <?php if (isset($component)) { $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.wrapper','data' => ['prefix' => __('filament-tables::table.grouping.fields.group.label')]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.wrapper'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['prefix' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament-tables::table.grouping.fields.group.label'))]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                <?php if (isset($component)) { $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.select','data' => ['xModel' => 'group','xOn:change' => 'resetCollapsedGroups()']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.select'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['x-model' => 'group','x-on:change' => 'resetCollapsedGroups()']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                    <option value="">-</option>

                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $groups; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $groupOption): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                        <option
                                                            value="<?php echo e($groupOption->getId()); ?>"
                                                        >
                                                            <?php echo e($groupOption->getLabel()); ?>

                                                        </option>
                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $attributes = $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $component = $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
                                             <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $attributes = $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $component = $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
                                        </label>

                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $isGroupingDirectionSettingHidden): ?>
                                            <label x-cloak x-show="group">
                                                <span class="fi-sr-only">
                                                    <?php echo e(__('filament-tables::table.grouping.fields.direction.label')); ?>

                                                </span>

                                                <?php if (isset($component)) { $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.wrapper','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.wrapper'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                    <?php if (isset($component)) { $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.select','data' => ['xModel' => 'direction']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.select'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['x-model' => 'direction']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                        <option value="asc">
                                                            <?php echo e(__('filament-tables::table.grouping.fields.direction.options.asc')); ?>

                                                        </option>

                                                        <option value="desc">
                                                            <?php echo e(__('filament-tables::table.grouping.fields.direction.options.desc')); ?>

                                                        </option>
                                                     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $attributes = $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $component = $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
                                                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $attributes = $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $component = $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
                                            </label>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                    </div>
                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                            </div>
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                        <?php echo e(FilamentView::renderHook(TablesRenderHook::TOOLBAR_GROUPING_SELECTOR_AFTER, scopes: static::class)); ?>

                    </div>

                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isGlobalSearchVisible || $hasFiltersTrigger || $hasColumnManager): ?>
                        <div>
                            <?php echo e(FilamentView::renderHook(TablesRenderHook::TOOLBAR_SEARCH_BEFORE, scopes: static::class)); ?>


                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isGlobalSearchVisible): ?>
                                <?php
                                    $searchPlaceholder = $getSearchPlaceholder();
                                ?>

                                <?php if (isset($component)) { $__componentOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.search-field','data' => ['debounce' => $searchDebounce,'onBlur' => $isSearchOnBlur,'placeholder' => $searchPlaceholder]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::search-field'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['debounce' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($searchDebounce),'on-blur' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isSearchOnBlur),'placeholder' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($searchPlaceholder)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc)): ?>
<?php $attributes = $__attributesOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc; ?>
<?php unset($__attributesOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc)): ?>
<?php $component = $__componentOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc; ?>
<?php unset($__componentOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc); ?>
<?php endif; ?>
                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                            <?php echo e(FilamentView::renderHook(TablesRenderHook::TOOLBAR_SEARCH_AFTER, scopes: static::class)); ?>


                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasFiltersTrigger || $hasColumnManager): ?>
                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasFiltersDialog): ?>
                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(($filtersLayout === FiltersLayout::Modal) || $filtersTriggerAction->isModalSlideOver()): ?>
                                        <?php
                                            $filtersTriggerActionModalAlignment = $filtersTriggerAction->getModalAlignment();
                                            $filtersTriggerActionIsModalAutofocused = $filtersTriggerAction->isModalAutofocused();
                                            $filtersTriggerActionHasModalCloseButton = $filtersTriggerAction->hasModalCloseButton();
                                            $filtersTriggerActionIsModalClosedByClickingAway = $filtersTriggerAction->isModalClosedByClickingAway();
                                            $filtersTriggerActionIsModalClosedByEscaping = $filtersTriggerAction->isModalClosedByEscaping();
                                            $filtersTriggerActionModalDescription = $filtersTriggerAction->getModalDescription();
                                            $filtersTriggerActionExtraModalWindowAttributeBag = $filtersTriggerAction->getExtraModalWindowAttributeBag();
                                            $filtersTriggerActionExtraModalOverlayAttributeBag = $filtersTriggerAction->getExtraModalOverlayAttributeBag();
                                            $filtersTriggerActionVisibleModalFooterActions = $filtersTriggerAction->getVisibleModalFooterActions();
                                            $filtersTriggerActionModalFooterActionsAlignment = $filtersTriggerAction->getModalFooterActionsAlignment();
                                            $filtersTriggerActionModalHeading = $filtersTriggerAction->getCustomModalHeading() ?? __('filament-tables::table.filters.heading');
                                            $filtersTriggerActionModalIcon = $filtersTriggerAction->getModalIcon();
                                            $filtersTriggerActionModalIconColor = $filtersTriggerAction->getModalIconColor();
                                            $filtersTriggerActionIsModalSlideOver = $filtersTriggerAction->isModalSlideOver();
                                            $filtersTriggerActionModalSlideOverPosition = $filtersTriggerAction->getModalSlideOverPosition();
                                            $filtersTriggerActionIsModalFooterSticky = $filtersTriggerAction->isModalFooterSticky();
                                            $filtersTriggerActionIsModalHeaderSticky = $filtersTriggerAction->isModalHeaderSticky();
                                        ?>

                                        <?php if (isset($component)) { $__componentOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.modal.index','data' => ['alignment' => $filtersTriggerActionModalAlignment,'autofocus' => $filtersTriggerActionIsModalAutofocused,'closeButton' => $filtersTriggerActionHasModalCloseButton,'closeByClickingAway' => $filtersTriggerActionIsModalClosedByClickingAway,'closeByEscaping' => $filtersTriggerActionIsModalClosedByEscaping,'description' => $filtersTriggerActionModalDescription,'extraModalWindowAttributeBag' => $filtersTriggerActionExtraModalWindowAttributeBag,'extraModalOverlayAttributeBag' => $filtersTriggerActionExtraModalOverlayAttributeBag,'footerActions' => $filtersTriggerActionVisibleModalFooterActions,'footerActionsAlignment' => $filtersTriggerActionModalFooterActionsAlignment,'heading' => $filtersTriggerActionModalHeading,'icon' => $filtersTriggerActionModalIcon,'iconColor' => $filtersTriggerActionModalIconColor,'slideOver' => $filtersTriggerActionIsModalSlideOver,'slideOverPosition' => $filtersTriggerActionModalSlideOverPosition,'stickyFooter' => $filtersTriggerActionIsModalFooterSticky,'stickyHeader' => $filtersTriggerActionIsModalHeaderSticky,'width' => $filtersFormWidth,'wire:key' => $this->getId() . '.table.filters','class' => 'fi-ta-filters-modal']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::modal'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['alignment' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionModalAlignment),'autofocus' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionIsModalAutofocused),'close-button' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionHasModalCloseButton),'close-by-clicking-away' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionIsModalClosedByClickingAway),'close-by-escaping' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionIsModalClosedByEscaping),'description' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionModalDescription),'extra-modal-window-attribute-bag' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionExtraModalWindowAttributeBag),'extra-modal-overlay-attribute-bag' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionExtraModalOverlayAttributeBag),'footer-actions' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionVisibleModalFooterActions),'footer-actions-alignment' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionModalFooterActionsAlignment),'heading' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionModalHeading),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionModalIcon),'icon-color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionModalIconColor),'slide-over' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionIsModalSlideOver),'slide-over-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionModalSlideOverPosition),'sticky-footer' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionIsModalFooterSticky),'sticky-header' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersTriggerActionIsModalHeaderSticky),'width' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersFormWidth),'wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . '.table.filters'),'class' => 'fi-ta-filters-modal']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                             <?php $__env->slot('trigger', null, []); ?> 
                                                <?php echo e($filtersTriggerAction->badge($activeFiltersCount)); ?>

                                             <?php $__env->endSlot(); ?>

                                            <?php echo e($filtersTriggerAction->getModalContent()); ?>


                                            <?php echo e($filtersForm); ?>


                                            <?php echo e($filtersTriggerAction->getModalContentFooter()); ?>

                                         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c)): ?>
<?php $attributes = $__attributesOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c; ?>
<?php unset($__attributesOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c)): ?>
<?php $component = $__componentOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c; ?>
<?php unset($__componentOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c); ?>
<?php endif; ?>
                                    <?php else: ?>
                                        <?php if (isset($component)) { $__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.dropdown.index','data' => ['maxHeight' => $filtersFormMaxHeight,'placement' => 'bottom-end','shift' => true,'flip' => false,'width' => $filtersFormWidth ?? Width::ExtraSmall,'wire:key' => $this->getId() . '.table.filters','class' => 'fi-ta-filters-dropdown']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::dropdown'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['max-height' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersFormMaxHeight),'placement' => 'bottom-end','shift' => true,'flip' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(false),'width' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersFormWidth ?? Width::ExtraSmall),'wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . '.table.filters'),'class' => 'fi-ta-filters-dropdown']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                             <?php $__env->slot('trigger', null, []); ?> 
                                                <?php echo e($filtersTriggerAction->badge($activeFiltersCount)); ?>

                                             <?php $__env->endSlot(); ?>

                                            <?php if (isset($component)) { $__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.filters','data' => ['applyAction' => $filtersApplyAction,'form' => $filtersForm,'headingTag' => $secondLevelHeadingTag,'resetActionPosition' => $filtersResetActionPosition]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::filters'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['apply-action' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersApplyAction),'form' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersForm),'heading-tag' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($secondLevelHeadingTag),'reset-action-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersResetActionPosition)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b)): ?>
<?php $attributes = $__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b; ?>
<?php unset($__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b)): ?>
<?php $component = $__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b; ?>
<?php unset($__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b); ?>
<?php endif; ?>
                                         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b)): ?>
<?php $attributes = $__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b; ?>
<?php unset($__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b)): ?>
<?php $component = $__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b; ?>
<?php unset($__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b); ?>
<?php endif; ?>
                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                <?php elseif($hasFiltersBeforeContent || $hasFiltersAfterContent): ?>
                                    <span
                                        x-ref="filtersTriggerActionContainer"
                                        x-on:click="toggleFiltersDropdown"
                                        class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                            'fi-ta-filters-trigger-action-ctn',
                                            'lg:fi-hidden' => ! $hasCollapsibleFilters,
                                        ]); ?>"
                                    >
                                        <?php echo e($filtersTriggerAction->badge($activeFiltersCount)); ?>

                                    </span>
                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                <?php echo e(FilamentView::renderHook(TablesRenderHook::TOOLBAR_COLUMN_MANAGER_TRIGGER_BEFORE, scopes: static::class)); ?>


                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasColumnManager): ?>
                                    <?php
                                        $columnManagerMaxHeight = $getColumnManagerMaxHeight();
                                        $columnManagerWidth = $getColumnManagerWidth();
                                        $columnManagerColumns = $getColumnManagerColumns();
                                    ?>

                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(($columnManagerLayout === ColumnManagerLayout::Modal) || $columnManagerTriggerAction->isModalSlideOver()): ?>
                                        <?php
                                            $columnManagerTriggerActionModalAlignment = $columnManagerTriggerAction->getModalAlignment();
                                            $columnManagerTriggerActionIsModalAutofocused = $columnManagerTriggerAction->isModalAutofocused();
                                            $columnManagerTriggerActionHasModalCloseButton = $columnManagerTriggerAction->hasModalCloseButton();
                                            $columnManagerTriggerActionIsModalClosedByClickingAway = $columnManagerTriggerAction->isModalClosedByClickingAway();
                                            $columnManagerTriggerActionIsModalClosedByEscaping = $columnManagerTriggerAction->isModalClosedByEscaping();
                                            $columnManagerTriggerActionModalDescription = $columnManagerTriggerAction->getModalDescription();
                                            $columnManagerTriggerActionExtraModalWindowAttributeBag = $columnManagerTriggerAction->getExtraModalWindowAttributeBag();
                                            $columnManagerTriggerActionExtraModalOverlayAttributeBag = $columnManagerTriggerAction->getExtraModalOverlayAttributeBag();
                                            $columnManagerTriggerActionVisibleModalFooterActions = $columnManagerTriggerAction->getVisibleModalFooterActions();
                                            $columnManagerTriggerActionModalFooterActionsAlignment = $columnManagerTriggerAction->getModalFooterActionsAlignment();
                                            $columnManagerTriggerActionModalHeading = $columnManagerTriggerAction->getCustomModalHeading() ?? __('filament-tables::table.column_manager.heading');
                                            $columnManagerTriggerActionModalIcon = $columnManagerTriggerAction->getModalIcon();
                                            $columnManagerTriggerActionModalIconColor = $columnManagerTriggerAction->getModalIconColor();
                                            $columnManagerTriggerActionIsModalSlideOver = $columnManagerTriggerAction->isModalSlideOver();
                                            $columnManagerTriggerActionModalSlideOverPosition = $columnManagerTriggerAction->getModalSlideOverPosition();
                                            $columnManagerTriggerActionIsModalFooterSticky = $columnManagerTriggerAction->isModalFooterSticky();
                                            $columnManagerTriggerActionIsModalHeaderSticky = $columnManagerTriggerAction->isModalHeaderSticky();
                                        ?>

                                        <?php if (isset($component)) { $__componentOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.modal.index','data' => ['alignment' => $columnManagerTriggerActionModalAlignment,'autofocus' => $columnManagerTriggerActionIsModalAutofocused,'closeButton' => $columnManagerTriggerActionHasModalCloseButton,'closeByClickingAway' => $columnManagerTriggerActionIsModalClosedByClickingAway,'closeByEscaping' => $columnManagerTriggerActionIsModalClosedByEscaping,'description' => $columnManagerTriggerActionModalDescription,'extraModalWindowAttributeBag' => $columnManagerTriggerActionExtraModalWindowAttributeBag,'extraModalOverlayAttributeBag' => $columnManagerTriggerActionExtraModalOverlayAttributeBag,'footerActions' => $columnManagerTriggerActionVisibleModalFooterActions,'footerActionsAlignment' => $columnManagerTriggerActionModalFooterActionsAlignment,'heading' => $columnManagerTriggerActionModalHeading,'icon' => $columnManagerTriggerActionModalIcon,'iconColor' => $columnManagerTriggerActionModalIconColor,'slideOver' => $columnManagerTriggerActionIsModalSlideOver,'slideOverPosition' => $columnManagerTriggerActionModalSlideOverPosition,'stickyFooter' => $columnManagerTriggerActionIsModalFooterSticky,'stickyHeader' => $columnManagerTriggerActionIsModalHeaderSticky,'width' => $columnManagerWidth,'wire:key' => $this->getId() . '.table.column-manager','class' => 'fi-ta-col-manager-modal']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::modal'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['alignment' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionModalAlignment),'autofocus' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionIsModalAutofocused),'close-button' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionHasModalCloseButton),'close-by-clicking-away' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionIsModalClosedByClickingAway),'close-by-escaping' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionIsModalClosedByEscaping),'description' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionModalDescription),'extra-modal-window-attribute-bag' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionExtraModalWindowAttributeBag),'extra-modal-overlay-attribute-bag' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionExtraModalOverlayAttributeBag),'footer-actions' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionVisibleModalFooterActions),'footer-actions-alignment' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionModalFooterActionsAlignment),'heading' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionModalHeading),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionModalIcon),'icon-color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionModalIconColor),'slide-over' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionIsModalSlideOver),'slide-over-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionModalSlideOverPosition),'sticky-footer' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionIsModalFooterSticky),'sticky-header' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerTriggerActionIsModalHeaderSticky),'width' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerWidth),'wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . '.table.column-manager'),'class' => 'fi-ta-col-manager-modal']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                             <?php $__env->slot('trigger', null, []); ?> 
                                                <?php echo e($columnManagerTriggerAction); ?>

                                             <?php $__env->endSlot(); ?>

                                            <?php echo e($columnManagerTriggerAction->getModalContent()); ?>


                                            <div
                                                x-data="filamentTableColumnManager({
                                                            columns: $wire.entangle('tableColumns'),
                                                            isLive: <?php echo e($columnManagerApplyAction->isVisible() ? 'false' : 'true'); ?>,
                                                        })"
                                                x-on:apply-table-column-manager.window="applyTableColumnManager()"
                                                x-on:reset-table-column-manager.window="resetDeferredColumns()"
                                                class="fi-ta-col-manager"
                                            >
                                                <?php if (isset($component)) { $__componentOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.column-manager.content','data' => ['columns' => $columnManagerColumns,'hasReorderableColumns' => $hasReorderableColumns,'hasToggleableColumns' => $hasToggleableColumns,'reorderAnimationDuration' => $getReorderAnimationDuration()]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::column-manager.content'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerColumns),'has-reorderable-columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasReorderableColumns),'has-toggleable-columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasToggleableColumns),'reorder-animation-duration' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($getReorderAnimationDuration())]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d)): ?>
<?php $attributes = $__attributesOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d; ?>
<?php unset($__attributesOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d); ?>
<?php endif; ?>
<?php if (isset($__componentOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d)): ?>
<?php $component = $__componentOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d; ?>
<?php unset($__componentOriginaldec847254e0650d09ee33cdc9b002db8a50a2582739ebd8387aecbe508e6d91d); ?>
<?php endif; ?>
                                            </div>

                                            <?php echo e($columnManagerTriggerAction->getModalContentFooter()); ?>

                                         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c)): ?>
<?php $attributes = $__attributesOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c; ?>
<?php unset($__attributesOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c)): ?>
<?php $component = $__componentOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c; ?>
<?php unset($__componentOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c); ?>
<?php endif; ?>
                                    <?php else: ?>
                                        <?php if (isset($component)) { $__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.dropdown.index','data' => ['maxHeight' => $columnManagerMaxHeight,'placement' => 'bottom-end','shift' => true,'flip' => false,'width' => $columnManagerWidth,'wire:key' => $this->getId() . '.table.column-manager','class' => 'fi-ta-col-manager-dropdown']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::dropdown'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['max-height' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerMaxHeight),'placement' => 'bottom-end','shift' => true,'flip' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(false),'width' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerWidth),'wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . '.table.column-manager'),'class' => 'fi-ta-col-manager-dropdown']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                             <?php $__env->slot('trigger', null, []); ?> 
                                                <?php echo e($columnManagerTriggerAction); ?>

                                             <?php $__env->endSlot(); ?>

                                            <?php if (isset($component)) { $__componentOriginal4ee5d255e9e8a2f7aa184577323605201469ab4a2d834b3cd552c033c1a4d2e4 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal4ee5d255e9e8a2f7aa184577323605201469ab4a2d834b3cd552c033c1a4d2e4 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.column-manager.index','data' => ['applyAction' => $columnManagerApplyAction,'columns' => $columnManagerColumns,'resetActionPosition' => $columnManagerResetActionPosition,'hasReorderableColumns' => $hasReorderableColumns,'hasToggleableColumns' => $hasToggleableColumns,'headingTag' => $secondLevelHeadingTag,'reorderAnimationDuration' => $getReorderAnimationDuration()]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::column-manager'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['apply-action' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerApplyAction),'columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerColumns),'reset-action-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columnManagerResetActionPosition),'has-reorderable-columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasReorderableColumns),'has-toggleable-columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasToggleableColumns),'heading-tag' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($secondLevelHeadingTag),'reorder-animation-duration' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($getReorderAnimationDuration())]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal4ee5d255e9e8a2f7aa184577323605201469ab4a2d834b3cd552c033c1a4d2e4)): ?>
<?php $attributes = $__attributesOriginal4ee5d255e9e8a2f7aa184577323605201469ab4a2d834b3cd552c033c1a4d2e4; ?>
<?php unset($__attributesOriginal4ee5d255e9e8a2f7aa184577323605201469ab4a2d834b3cd552c033c1a4d2e4); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal4ee5d255e9e8a2f7aa184577323605201469ab4a2d834b3cd552c033c1a4d2e4)): ?>
<?php $component = $__componentOriginal4ee5d255e9e8a2f7aa184577323605201469ab4a2d834b3cd552c033c1a4d2e4; ?>
<?php unset($__componentOriginal4ee5d255e9e8a2f7aa184577323605201469ab4a2d834b3cd552c033c1a4d2e4); ?>
<?php endif; ?>
                                         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b)): ?>
<?php $attributes = $__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b; ?>
<?php unset($__attributesOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b)): ?>
<?php $component = $__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b; ?>
<?php unset($__componentOriginala2d829441b5c0d22b79e927e8c82a2eeedeac44ab591e1f59e226d7013ec6f7b); ?>
<?php endif; ?>
                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                <?php echo e(FilamentView::renderHook(TablesRenderHook::TOOLBAR_COLUMN_MANAGER_TRIGGER_AFTER, scopes: static::class)); ?>

                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                        </div>
                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                    <?php echo e(FilamentView::renderHook(TablesRenderHook::TOOLBAR_END)); ?>

                </div>

                <?php echo e(FilamentView::renderHook(TablesRenderHook::TOOLBAR_AFTER)); ?>

            </div>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isReordering): ?>
                <div
                    x-cloak
                    role="status"
                    aria-live="polite"
                    <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.reorder.indicator'; ?>wire:key="<?php echo e($this->getId()); ?>.table.reorder.indicator"
                    class="fi-ta-reorder-indicator"
                >
                    <?php echo e(\Filament\Support\generate_loading_indicator_html(new Filament\Support\View\ComponentAttributeBag([
                            'wire:loading.delay.' . config('filament.livewire_loading_delay', 'default') => '',
                            'wire:target' => 'reorderTable',
                        ]))); ?>


                    <?php echo e(__('filament-tables::table.reorder_indicator')); ?>

                </div>
            <?php elseif($isSelectionEnabled && ($maxSelectableRecords !== 1) && $isLoaded): ?>
                <div
                    x-cloak
                    role="status"
                    aria-live="polite"
                    aria-atomic="true"
                    x-bind:hidden="! getSelectedRecordsCount()"
                    x-show="getSelectedRecordsCount()"
                    <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.selection.indicator'; ?>wire:key="<?php echo e($this->getId()); ?>.table.selection.indicator"
                    class="fi-ta-selection-indicator"
                >
                    <div>
                        <?php echo e(\Filament\Support\generate_loading_indicator_html(new Filament\Support\View\ComponentAttributeBag([
                                'x-show' => 'isLoading',
                            ]))); ?>


                        <span
                            x-text="
                                window.pluralize(<?php echo \Illuminate\Support\Js::from(__('filament-tables::table.selection_indicator.selected_count'))->toHtml() ?>, getSelectedRecordsCount(), {
                                    count: new Intl.NumberFormat(<?php echo \Illuminate\Support\Js::from(str_replace('_', '-', app()->getLocale()))->toHtml() ?>).format(getSelectedRecordsCount()),
                                })
                            "
                        ></span>
                    </div>

                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $isSelectionDisabled): ?>
                        <div>
                            <?php echo e(FilamentView::renderHook(TablesRenderHook::SELECTION_INDICATOR_ACTIONS_BEFORE, scopes: static::class)); ?>


                            <div class="fi-ta-selection-indicator-actions-ctn">
                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $selectsGroupsOnly): ?>
                                    <?php if (isset($component)) { $__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.link','data' => ['color' => 'primary','tag' => 'button','xOn:click' => 'selectAllRecords','xShow' => 'canSelectAllRecords()','wire:key' => $this->getId() . 'table.selection.indicator.actions.select-all.' . $allSelectableRecordsCount . '.' . $page]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::link'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'primary','tag' => 'button','x-on:click' => 'selectAllRecords','x-show' => 'canSelectAllRecords()','wire:key' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($this->getId() . 'table.selection.indicator.actions.select-all.' . $allSelectableRecordsCount . '.' . $page)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                        <?php echo e(trans_choice('filament-tables::table.selection_indicator.actions.select_all.label', $allSelectableRecordsCount, ['count' => Number::format($allSelectableRecordsCount, locale: app()->getLocale())])); ?>

                                     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe)): ?>
<?php $attributes = $__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe; ?>
<?php unset($__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe); ?>
<?php endif; ?>
<?php if (isset($__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe)): ?>
<?php $component = $__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe; ?>
<?php unset($__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe); ?>
<?php endif; ?>
                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                <?php if (isset($component)) { $__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.link','data' => ['color' => 'danger','tag' => 'button','xOn:click' => 'deselectAllRecords']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::link'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => 'danger','tag' => 'button','x-on:click' => 'deselectAllRecords']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                    <?php echo e(__('filament-tables::table.selection_indicator.actions.deselect_all.label')); ?>

                                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe)): ?>
<?php $attributes = $__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe; ?>
<?php unset($__attributesOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe); ?>
<?php endif; ?>
<?php if (isset($__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe)): ?>
<?php $component = $__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe; ?>
<?php unset($__componentOriginale1f44def72a41d4a45b5ae8992164b530404c0c18bd6d60b9c027f9699c5fcbe); ?>
<?php endif; ?>
                            </div>

                            <?php echo e(FilamentView::renderHook(TablesRenderHook::SELECTION_INDICATOR_ACTIONS_AFTER, scopes: static::class)); ?>

                        </div>
                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                </div>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($filterIndicators): ?>
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($filterIndicatorsView = FilamentView::renderHook(TablesRenderHook::FILTER_INDICATORS, scopes: static::class, data: ['filterIndicators' => $filterIndicators]))): ?>
                    <?php echo e($filterIndicatorsView); ?>

                <?php else: ?>
                    <div class="fi-ta-filter-indicators">
                        <div>
                            <span class="fi-ta-filter-indicators-label">
                                <?php echo e(__('filament-tables::table.filters.indicator')); ?>

                            </span>

                            <div
                                class="fi-ta-filter-indicators-badges-ctn"
                                role="list"
                            >
                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $filterIndicators; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $indicator): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                    <?php
                                        $indicatorColor = $indicator->getColor();
                                    ?>

                                    <?php if (isset($component)) { $__componentOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.badge','data' => ['color' => $indicatorColor,'role' => 'listitem']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::badge'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($indicatorColor),'role' => 'listitem']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                        <?php echo e($indicator->getLabel()); ?>


                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($indicator->isRemovable()): ?>
                                            <?php
                                                $indicatorRemoveLivewireClickHandler = $indicator->getRemoveLivewireClickHandler();
                                            ?>

                                             <?php $__env->slot('deleteButton', null, ['label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament-tables::table.filters.actions.remove.label')),'wire:click' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($indicatorRemoveLivewireClickHandler),'wire:loading.attr' => 'disabled','wire:target' => 'removeTableFilter']); ?>  <?php $__env->endSlot(); ?>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f)): ?>
<?php $attributes = $__attributesOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f; ?>
<?php unset($__attributesOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f)): ?>
<?php $component = $__componentOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f; ?>
<?php unset($__componentOriginal8ae421462de53d87fcf23baf86e2223b3f2cea2f7c9077e9da301cabdd5c562f); ?>
<?php endif; ?>
                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                            </div>
                        </div>

                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(collect($filterIndicators)->contains(fn (Indicator $indicator): bool => $indicator->isRemovable())): ?>
                            <?php echo e($getFiltersRemoveAllAction()); ?>

                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                    </div>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(((! $content) && (! $hasColumnsLayout)) || ($records === null) || count($records)): ?>
                <div
                    <?php if((! $isReordering) && ($pollingInterval = $getPollingInterval())): ?>
                        wire:poll.<?php echo e($pollingInterval); ?>

                    <?php endif; ?>
                    class="fi-ta-content-ctn fi-fixed-positioning-context"
                >
                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($records !== null): ?>
                        <?php
                            // The total across all pages, not the current page's count — and since the total is
                            // stable across pages, the live region only announces when the result set really
                            // changes, not on every pagination click. Non-length-aware paginators fall back to
                            // the page count.
                            $resultCount = ($records instanceof LengthAwarePaginator)
                                ? $records->total()
                                : count($records);
                      ?>hp

                        <div
                            role="status"
                            aria-live="polite"
                            aria-atomic="true"
                            class="fi-sr-only"
                        >
                            <?php echo e(trans_choice('filament-tables::table.result_count', $resultCount, ['count' => $resultCount])); ?>

                        </div>
                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasContentLayout && ($records !== null) && count($records)): ?>
                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $isReordering): ?>
                          <?phphp
                                $sortableColumns = array_filter(
                                    $columns,
                                    fn (Column $column): bool => $column->isSortable(),
                                );
                          ?>hp

                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(($isSelectionEnabled && ($maxSelectableRecords !== 1) && (! $isReordering) && (! $selectsGroupsOnly)) || count($sortableColumns)): ?>
                                <div class="fi-ta-content-header">
                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && ($maxSelectableRecords !== 1) && (! $isReordering) && (! $selectsGroupsOnly)): ?>
                                        <input
                                            aria-label="<?php echo e(__('filament-tables::table.fields.bulk_select_page.label')); ?>"
                                            type="checkbox"
                                            <?php if($isSelectionDisabled): ?>
                                                disabled
                                            <?php elseif($maxSelectableRecords): ?>
                                                x-bind:disabled="
                                                    const recordsOnPage = getRecordsOnPage()

                                                    return recordsOnPage.length && ! areRecordsToggleable(recordsOnPage)
                                                "
                                            <?php endif; ?>
                                            x-bind:checked="
                                                const recordsOnPage = getRecordsOnPage()

                                                if (recordsOnPage.length && areRecordsSelected(recordsOnPage)) {
                                                    $el.checked = true
                                                    $el.indeterminate = false

                                                    return 'checked'
                                                }

                                                $el.checked = false
                                                $el.indeterminate =
                                                    recordsOnPage.length && areRecordsPartiallySelected(recordsOnPage)

                                                return null
                                            "
                                            x-on:click="toggleSelectRecordsOnPage"
                                            
                                            <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk-select-page.checkbox.'.e(Str::random()).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk-select-page.checkbox.'.e(Str::random()).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk-select-page.checkbox.'.e(Str::random()).''; ?>wire:key="<?php echo e($this->getId()); ?>.table.bulk-select-page.checkbox.<?php echo e(Str::random()); ?>"
                                            wire:loading.attr="disabled"
                                            wire:target="<?php echo e($loadingTargetsWireTarget); ?>"
                                            class="fi-ta-page-checkbox fi-checkbox-input"
                                        />
                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(count($sortableColumns)): ?>
                                        <div
                                            x-data="{
                                                sort: $wire.$entangle('tableSort', true),
                                                column: null,
                                                direction: null,
                                            }"
                                            x-init="
                                                if (sort) {
                                                    ;[column, direction] = sort.split(':')
                                                    direction ??= 'asc'
                                                }

                                                $watch('sort', function () {
                                                    if (! sort) {
                                                        return
                                                    }

                                                    ;[column, direction] = sort.split(':')
                                                    direction ??= 'asc'
                                                })

                                                $watch('direction', function () {
                                                    sort = column ? `${column}:${direction}` : null
                                                })

                                                $watch('column', function (newColumn, oldColumn) {
                                                    if (! newColumn) {
                                                        direction = null
                                                        sort = column ? `${column}:${direction}` : null

                                                        return
                                                    }

                                                    if (oldColumn) {
                                                        sort = column ? `${column}:${direction}` : null

                                                        return
                                                    }

                                                    direction = 'asc'
                                                    sort = column ? `${column}:${direction}` : null
                                                })
                                            "
                                            class="fi-ta-sorting-settings"
                                        >
                                            <label>
                                                <?php if (isset($component)) { $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.wrapper','data' => ['prefix' => __('filament-tables::table.sorting.fields.column.label')]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.wrapper'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['prefix' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament-tables::table.sorting.fields.column.label'))]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                    <?php if (isset($component)) { $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.select','data' => ['xModel' => 'column']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.select'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['x-model' => 'column']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                        <option value="">
                                                            <?php echo e($defaultSortOptionLabel); ?>

                                                        </option>

                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $sortableColumns; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $column): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                            <option
                                                                value="<?php echo e($column->getName()); ?>"
                                                            >
                                                                <?php echo e($column->getLabel()); ?>

                                                            </option>
                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $attributes = $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $component = $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
                                                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $attributes = $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $component = $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
                                            </label>

                                            <label x-cloak x-show="column">
                                                <span class="fi-sr-only">
                                                    <?php echo e(__('filament-tables::table.sorting.fields.direction.label')); ?>

                                                </span>

                                                <?php if (isset($component)) { $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.wrapper','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.wrapper'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                    <?php if (isset($component)) { $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.select','data' => ['xModel' => 'direction']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.select'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['x-model' => 'direction']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                        <option value="asc">
                                                            <?php echo e(__('filament-tables::table.sorting.fields.direction.options.asc')); ?>

                                                        </option>

                                                        <option value="desc">
                                                            <?php echo e(__('filament-tables::table.sorting.fields.direction.options.desc')); ?>

                                                        </option>
                                                     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $attributes = $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $component = $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
                                                 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $attributes = $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $component = $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
                                            </label>
                                        </div>
                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                </div>
                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($content): ?>
                            <?php echo e($content->with(['records' => $records])); ?>

                        <?php else: ?>
                            <div
                                <?php if($isReorderable): ?>
                                    x-on:end.stop="
                                        $wire.reorderTable(
                                            $event.target.sortable.toArray(),
                                            $event.item.getAttribute('x-sortable-item'),
                                        )
                                    "
                                    x-sortable
                                    data-sortable-animation-duration="<?php echo e($getReorderAnimationDuration()); ?>"
                                <?php endif; ?>
                                aria-label="<?php echo e($pluralModelLabel); ?>"
                                role="list"
                                <?php echo e((new FilamentComponentAttributeBag)
                                        ->when($contentGrid, fn (ComponentAttributeBag $attributes) => $attributes->grid($contentGrid))
                                        ->class([
                                            'fi-ta-content',
                                            'fi-ta-content-grid' => $contentGrid,
                                            'fi-ta-content-grouped' => $this->getTableGrouping(),
                                        ])); ?>

                            >
                              <?phphp
                                    $previousRecord = null;
                                    $previousRecordGroupKey = null;
                                    $previousRecordGroupTitle = null;
                              ?>hp

                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $records; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $record): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                  <?phphp
                                        $recordAction = $getRecordAction($record);
                                        $recordKey = $getRecordKey($record);
                                        $recordUrl = $getRecordUrl($record);
                                        $openRecordUrlInNewTab = $shouldOpenRecordUrlInNewTab($record);
                                        $recordGroupKey = $group?->getStringKey($record);
                                        $recordGroupTitle = $group?->getTitle($record, $recordGroupKey);
                                        $isRecordGroupCollapsible = $group?->isCollapsible();
                                        $recordIsSelectable = $isSelectionEnabled && $isRecordSelectable($record);

                                        $collapsibleColumnsLayout?->record($record)->recordKey($recordKey);
                                        $hasCollapsibleColumnsLayout = (bool) $collapsibleColumnsLayout?->isVisible();

                                        $recordActions = $recordActionsByRecordKey[$recordKey]
                                            ?? ($hasRecordActionsForAnyRecord ? $reduceVisibleRecordActions($record) : []);
                                  ?>hp

                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if((string) $recordGroupTitle !== (string) $previousRecordGroupTitle): ?>
                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasSummary && (! $isReordering) && filled($previousRecordGroupTitle)): ?>
                                            <table
                                                class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                                    'fi-ta-table',
                                                    'fi-ta-table-reordering' => $isReordering,
                                                ]); ?>"
                                            >
                                                <tbody>
                                                  <?phphp
                                                        $groupScopedAllTableSummaryQuery = $group->scopeQuery($this->getAllTableSummaryQuery(), $previousRecord);
                                                  ?>hp

                                                    <?php if (isset($component)) { $__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.summary.row','data' => ['columns' => $columns,'extraHeadingColumn' => true,'heading' => 
                                                            __('filament-tables::table.summary.subheadings.group', [
                                                                'group' => $previousRecordGroupTitle,
                                                                'label' => $pluralModelLabel,
                                                            ])
                                                        ,'placeholderColumns' => false,'query' => $groupScopedAllTableSummaryQuery,'selectedState' => $groupedSummarySelectedState[$previousRecordGroupKey] ?? []]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::summary.row'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columns),'extra-heading-column' => true,'heading' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
                                                            __('filament-tables::table.summary.subheadings.group', [
                                                                'group' => $previousRecordGroupTitle,
                                                                'label' => $pluralModelLabel,
                                                            ])
                                                        ),'placeholder-columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(false),'query' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupScopedAllTableSummaryQuery),'selected-state' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupedSummarySelectedState[$previousRecordGroupKey] ?? [])]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da)): ?>
<?php $attributes = $__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da; ?>
<?php unset($__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da)): ?>
<?php $component = $__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da; ?>
<?php unset($__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da); ?>
<?php endif; ?>
                                                </tbody>
                                            </table>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                        <div
                                            <?php if($isRecordGroupCollapsible = $group->isCollapsible()): ?>
                                                x-on:click="toggleCollapseGroup(<?php echo \Illuminate\Support\Js::from($recordGroupTitle)->toHtml() ?>)"
                                                <?php if(! $hasSummary): ?>
                                                    x-bind:class="{ 'fi-collapsed': isGroupCollapsed(<?php echo \Illuminate\Support\Js::from($recordGroupTitle)->toHtml() ?>) }"
                                                <?php endif; ?>
                                            <?php endif; ?>
                                            class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                                'fi-ta-group-header',
                                                'fi-collapsible' => $isRecordGroupCollapsible,
                                            ]); ?>"
                                        >
                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && ($maxSelectableRecords !== 1)): ?>
                                                <input
                                                    aria-label="<?php echo e(__('filament-tables::table.fields.bulk_select_group.label', ['title' => $recordGroupTitle])); ?>"
                                                    type="checkbox"
                                                    data-group-selectable-record-keys="<?php echo e(json_encode($this->getGroupedSelectableTableRecordKeys($recordGroupKey))); ?>"
                                                    <?php if($isSelectionDisabled): ?>
                                                        disabled
                                                    <?php else: ?>
                                                        x-on:click="toggleSelectRecords(JSON.parse($el.dataset.groupSelectableRecordKeys))"
                                                        <?php if($maxSelectableRecords): ?>
                                                            x-bind:disabled="
                                                                const recordsInGroup = JSON.parse($el.dataset.groupSelectableRecordKeys)

                                                                return recordsInGroup.length && ! areRecordsToggleable(recordsInGroup)
                                                            "
                                                        <?php endif; ?>
                                                    <?php endif; ?>
                                                    x-bind:checked="
                                                        const recordsInGroup = JSON.parse($el.dataset.groupSelectableRecordKeys)

                                                        if (recordsInGroup.length && areRecordsSelected(recordsInGroup)) {
                                                            $el.checked = true
                                                            $el.indeterminate = false

                                                            return 'checked'
                                                        }

                                                        $el.checked = false
                                                        $el.indeterminate =
                                                            recordsInGroup.length && areRecordsPartiallySelected(recordsInGroup)

                                                        return null
                                                    "
                                                    <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk_select_group.checkbox.'.e($page).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk_select_group.checkbox.'.e($page).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk_select_group.checkbox.'.e($page).''; ?>wire:key="<?php echo e($this->getId()); ?>.table.bulk_select_group.checkbox.<?php echo e($page); ?>"
                                                    wire:loading.attr="disabled"
                                                    wire:target="<?php echo e($loadingTargetsWireTarget); ?>"
                                                    class="fi-ta-group-checkbox fi-checkbox-input"
                                                />
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                            <div>
                                                <<?php echo e($secondLevelHeadingTag); ?>

                                                    class="fi-ta-group-heading"
                                                >
                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($recordGroupLabel = ($group->isTitlePrefixedWithLabel() ? $group->getLabel() : null))): ?>
                                                            <?php echo e($recordGroupLabel); ?>:
                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                    <?php echo e($recordGroupTitle); ?>

                                                </<?php echo e($secondLevelHeadingTag); ?>>

                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($recordGroupDescription = $group->getDescription($record, $recordGroupTitle))): ?>
                                                    <p
                                                        class="fi-ta-group-description"
                                                    >
                                                        <?php echo e($recordGroupDescription); ?>

                                                    </p>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                            </div>

                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isRecordGroupCollapsible): ?>
                                                <button
                                                    aria-label="<?php echo e(filled($recordGroupLabel) ? ($recordGroupLabel . ': ' . $recordGroupTitle) : $recordGroupTitle); ?>"
                                                    x-bind:aria-expanded="! isGroupCollapsed(<?php echo \Illuminate\Support\Js::from($recordGroupTitle)->toHtml() ?>)"
                                                    type="button"
                                                    class="fi-icon-btn fi-size-sm"
                                                >
                                                    <?php echo e(\Filament\Support\generate_icon_html(Heroicon::ChevronUp, alias: TablesIconAlias::GROUPING_COLLAPSE_BUTTON, size: IconSize::Small)); ?>

                                                </button>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                        </div>
                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                    <div
                                        <?php if($hasCollapsibleColumnsLayout): ?>
                                            x-data="{ isCollapsed: <?php echo \Illuminate\Support\Js::from($collapsibleColumnsLayout->isCollapsed())->toHtml() ?> }"
                                            x-init="$dispatch('collapsible-table-row-initialized')"
                                            x-on:collapse-all-table-rows.window="isCollapsed = true"
                                            x-on:expand-all-table-rows.window="isCollapsed = false"
                                            x-bind:class="isCollapsed && 'fi-ta-record-collapsed'"
                                        <?php endif; ?>
                                        role="listitem"
                                        <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.records.'.e($recordKey).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.records.'.e($recordKey).''; ?>wire:key="<?php echo e($this->getId()); ?>.table.records.<?php echo e($recordKey); ?>"
                                        <?php if($isReordering): ?>
                                            x-sortable-item="<?php echo e($recordKey); ?>"
                                            x-sortable-handle
                                        <?php endif; ?>
                                        class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                            'fi-ta-record',
                                            'fi-clickable' => $recordUrl || $recordAction,
                                            'fi-ta-record-with-content-prefix' => $isReordering || $recordIsSelectable,
                                            'fi-ta-record-with-content-suffix' => $hasCollapsibleColumnsLayout && (! $isReordering),
                                            ...$getRecordClasses($record),
                                        ]); ?>"
                                        x-bind:class="{
                                            <?php echo e($group?->isCollapsible() ? '\'fi-collapsed\': isGroupCollapsed(' . \Illuminate\Support\Js::from($recordGroupTitle) . '),' : ''); ?>

                                            'fi-selected': <?php echo \Illuminate\Support\Js::from($recordIsSelectable)->toHtml() ?> && isRecordSelected(<?php echo \Illuminate\Support\Js::from($recordKey)->toHtml() ?>),
                                        }"
                                    >
                                      <?phphp
                                            $hasItemBeforeRecordContent = $isReordering || $recordIsSelectable;
                                            $hasItemAfterRecordContent = $hasCollapsibleColumnsLayout && (! $isReordering);
                                      ?>hp

                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isReordering): ?>
                                            <button
                                                aria-label="<?php echo e(__('filament-tables::table.actions.reorder_record.label', ['key' => $recordKey])); ?>"
                                                class="fi-ta-reorder-handle fi-icon-btn"
                                                type="button"
                                            >
                                                <?php echo e(\Filament\Support\generate_icon_html(Heroicon::Bars2, alias: TablesIconAlias::REORDER_HANDLE)); ?>

                                            </button>
                                        <?php elseif($recordIsSelectable): ?>
                                            <input
                                                aria-label="<?php echo e(__('filament-tables::table.fields.bulk_select_record.label', ['key' => $recordKey])); ?>"
                                                type="checkbox"
                                                <?php if($isSelectionDisabled): ?>
                                                    disabled
                                                <?php elseif($maxSelectableRecords && ($maxSelectableRecords !== 1)): ?>
                                                    x-bind:disabled="! areRecordsToggleable([<?php echo \Illuminate\Support\Js::from($recordKey)->toHtml() ?>])"
                                                <?php endif; ?>
                                                value="<?php echo e($recordKey); ?>"
                                                x-on:click="toggleSelectedRecord(<?php echo \Illuminate\Support\Js::from($recordKey)->toHtml() ?>)"
                                                x-bind:checked="isRecordSelected(<?php echo \Illuminate\Support\Js::from($recordKey)->toHtml() ?>) ? 'checked' : null"
                                                data-group="<?php echo e($recordGroupKey); ?>"
                                                wire:loading.attr="disabled"
                                                wire:target="<?php echo e($loadingTargetsWireTarget); ?>"
                                                class="fi-ta-record-checkbox fi-checkbox-input"
                                            />
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                        <div class="fi-ta-record-content-ctn">
                                            <div>
                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($recordUrl): ?>
                                                    <a
                                                        <?php echo e(\Filament\Support\generate_href_html($recordUrl, $openRecordUrlInNewTab, hasNestedClickEventHandler: true)); ?>

                                                        <?php echo e($getExtraRecordLinkAttributeBag($record)->class(['fi-ta-record-content'])); ?>

                                                    >
                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $columnsLayout; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $columnsLayoutComponent): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                            <?php echo e($columnsLayoutComponent
                                                                    ->record($record)
                                                                    ->recordKey($recordKey)
                                                                    ->rowLoop($loop)
                                                                    ->renderInLayout()); ?>

                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                    </a>
                                                <?php elseif($recordAction): ?>
                                                  <?phphp
                                                        $recordWireClickAction = $getRecordAction($record)
                                                            ? "mountTableAction('{$recordAction}', '{$recordKey}')"
                                                            : $recordWireClickAction = "{$recordAction}('{$recordKey}')";
                                                  ?>hp

                                                    <button
                                                        type="button"
                                                        wire:click="<?php echo e($recordWireClickAction); ?>"
                                                        wire:loading.attr="disabled"
                                                        wire:target="<?php echo e($recordWireClickAction); ?>"
                                                        class="fi-ta-record-content"
                                                    >
                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $columnsLayout; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $columnsLayoutComponent): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                            <?php echo e($columnsLayoutComponent
                                                                    ->record($record)
                                                                    ->recordKey($recordKey)
                                                                    ->rowLoop($loop)
                                                                    ->renderInLayout()); ?>

                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                    </button>
                                                <?php else: ?>
                                                    <div
                                                        class="fi-ta-record-content"
                                                    >
                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $columnsLayout; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $columnsLayoutComponent): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                            <?php echo e($columnsLayoutComponent
                                                                    ->record($record)
                                                                    ->recordKey($recordKey)
                                                                    ->rowLoop($loop)
                                                                    ->renderInLayout()); ?>

                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                    </div>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasCollapsibleColumnsLayout && (! $isReordering)): ?>
                                                    <div
                                                        x-collapse
                                                        x-show="! isCollapsed"
                                                        class="fi-ta-record-content fi-collapsible"
                                                    >
                                                        <?php echo e($collapsibleColumnsLayout); ?>

                                                    </div>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                            </div>

                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($recordActions && (! $isReordering)): ?>
                                                <div
                                                    class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                                        'fi-ta-actions fi-wrapped sm:fi-not-wrapped',
                                                        match ($recordActionsAlignment ?? Alignment::Start) {
                                                            Alignment::Start => 'fi-align-start',
                                                            Alignment::Center => 'fi-align-center',
                                                            Alignment::End => 'fi-align-end',
                                                        } => $contentGrid,
                                                        'fi-align-start md:fi-align-end' => ! $contentGrid,
                                                        'fi-ta-actions-before-columns-position' => $recordActionsPosition === RecordActionsPosition::BeforeColumns,
                                                    ]); ?>"
                                                >
                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $recordActions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $action): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                        <?php echo e($action); ?>

                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                </div>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                        </div>

                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasCollapsibleColumnsLayout && (! $isReordering)): ?>
                                            <button
                                                aria-label="<?php echo e(__('filament-tables::table.actions.toggle_record_content.label', ['key' => $recordKey])); ?>"
                                                x-bind:aria-expanded="! isCollapsed"
                                                type="button"
                                                x-on:click="isCollapsed = ! isCollapsed"
                                                class="fi-ta-record-collapse-btn fi-icon-btn"
                                            >
                                                <?php echo e(\Filament\Support\generate_icon_html(Heroicon::ChevronDown, alias: TablesIconAlias::COLUMNS_COLLAPSE_BUTTON)); ?>

                                            </button>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                    </div>

                                  <?phphp
                                        $previousRecordGroupKey = $recordGroupKey;
                                        $previousRecordGroupTitle = $recordGroupTitle;
                                        $previousRecord = $record;
                                  ?>hp
                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>

                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasSummary && (! $isReordering) && filled($previousRecordGroupTitle) && $this->shouldRenderTrailingGroupedTableSummary($previousRecord)): ?>
                                    <table class="fi-ta-table">
                                        <tbody>
                                          <?phphp
                                                $groupScopedAllTableSummaryQuery = $group->scopeQuery($this->getAllTableSummaryQuery(), $previousRecord);
                                          ?>hp

                                            <?php if (isset($component)) { $__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.summary.row','data' => ['columns' => $columns,'extraHeadingColumn' => true,'heading' => __('filament-tables::table.summary.subheadings.group', ['group' => $previousRecordGroupTitle, 'label' => $pluralModelLabel]),'placeholderColumns' => false,'query' => $groupScopedAllTableSummaryQuery,'selectedState' => $groupedSummarySelectedState[$previousRecordGroupKey] ?? []]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::summary.row'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columns),'extra-heading-column' => true,'heading' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament-tables::table.summary.subheadings.group', ['group' => $previousRecordGroupTitle, 'label' => $pluralModelLabel])),'placeholder-columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(false),'query' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupScopedAllTableSummaryQuery),'selected-state' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupedSummarySelectedState[$previousRecordGroupKey] ?? [])]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da)): ?>
<?php $attributes = $__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da; ?>
<?php unset($__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da)): ?>
<?php $component = $__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da; ?>
<?php unset($__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da); ?>
<?php endif; ?>
                                        </tbody>
                                    </table>
                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                            </div>
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(($content || $hasColumnsLayout) && $contentFooter): ?>
                            <?php echo e($contentFooter->with([
                                    'columns' => $columns,
                                    'records' => $records,
                                ])); ?>

                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasTopLevelSummary && (! $isReordering)): ?>
                            <table class="fi-ta-table">
                                <tbody>
                                    <?php if (isset($component)) { $__componentOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.summary.index','data' => ['allTableSummary' => $hasAllTableSummary,'columns' => $columns,'extraHeadingColumn' => true,'pageSummary' => $hasPageSummary,'placeholderColumns' => false,'pluralModelLabel' => $pluralModelLabel,'records' => $records]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::summary'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['all-table-summary' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasAllTableSummary),'columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columns),'extra-heading-column' => true,'page-summary' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasPageSummary),'placeholder-columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(false),'plural-model-label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($pluralModelLabel),'records' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($records)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed)): ?>
<?php $attributes = $__attributesOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed; ?>
<?php unset($__attributesOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed)): ?>
<?php $component = $__componentOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed; ?>
<?php unset($__componentOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed); ?>
<?php endif; ?>
                                </tbody>
                            </table>
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                    <?php elseif((! ($content || $hasColumnsLayout)) && ($records !== null)): ?>
                      <?phphp
                            $sortableColumns = $isStackedOnMobile ? array_filter(
                                $columns,
                                fn (Column $column): bool => $column->isSortable(),
                            ) : [];
                      ?>hp

                        <table
                            aria-label="<?php echo e(filled($tableAccessibleLabel = trim(strip_tags((string) $heading))) ? $tableAccessibleLabel : $pluralModelLabel); ?>"
                            class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                'fi-ta-table',
                                'fi-ta-table-stacked-on-mobile' => $isStackedOnMobile,
                            ]); ?>"
                        >
                            <thead>
                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isStackedOnMobile && (count($sortableColumns) || ($isSelectionEnabled && ($maxSelectableRecords !== 1) && (! $selectsGroupsOnly))) && (! $isReordering)): ?>
                                    <tr class="fi-ta-table-stacked-header-row">
                                        <td
                                            colspan="100%"
                                            class="fi-ta-table-stacked-header-cell"
                                        >
                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(count($sortableColumns)): ?>
                                                <div
                                                    x-data="{
                                                        sort: $wire.$entangle('tableSort', true),
                                                        column: null,
                                                        direction: null,
                                                    }"
                                                    x-init="
                                                        if (sort) {
                                                            ;[column, direction] = sort.split(':')
                                                            direction ??= 'asc'
                                                        }

                                                        $watch('sort', function () {
                                                            if (! sort) {
                                                                return
                                                            }

                                                            ;[column, direction] = sort.split(':')
                                                            direction ??= 'asc'
                                                        })

                                                        $watch('direction', function () {
                                                            sort = column ? `${column}:${direction}` : null
                                                        })

                                                        $watch('column', function (newColumn, oldColumn) {
                                                            if (! newColumn) {
                                                                direction = null
                                                                sort = column ? `${column}:${direction}` : null

                                                                return
                                                            }

                                                            if (oldColumn) {
                                                                sort = column ? `${column}:${direction}` : null

                                                                return
                                                            }

                                                            direction = 'asc'
                                                            sort = column ? `${column}:${direction}` : null
                                                        })
                                                    "
                                                    class="fi-ta-table-stacked-sorting"
                                                >
                                                    <label>
                                                        <?php if (isset($component)) { $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.wrapper','data' => ['prefix' => __('filament-tables::table.sorting.fields.column.label')]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.wrapper'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['prefix' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(__('filament-tables::table.sorting.fields.column.label'))]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                            <?php if (isset($component)) { $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.select','data' => ['xModel' => 'column']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.select'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['x-model' => 'column']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                                <option
                                                                    value=""
                                                                >
                                                                    <?php echo e($defaultSortOptionLabel); ?>

                                                                </option>

                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $sortableColumns; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $sortableColumn): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                                    <option
                                                                        value="<?php echo e($sortableColumn->getName()); ?>"
                                                                    >
                                                                        <?php echo e($sortableColumn->getLabel()); ?>

                                                                    </option>
                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                             <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $attributes = $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $component = $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
                                                         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $attributes = $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $component = $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
                                                    </label>

                                                    <label
                                                        x-cloak
                                                        x-show="column"
                                                    >
                                                        <span
                                                            class="fi-sr-only"
                                                        >
                                                            <?php echo e(__('filament-tables::table.sorting.fields.direction.label')); ?>

                                                        </span>

                                                        <?php if (isset($component)) { $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.wrapper','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.wrapper'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                            <?php if (isset($component)) { $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.select','data' => ['xModel' => 'direction']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.select'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['x-model' => 'direction']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                                                                <option
                                                                    value="asc"
                                                                >
                                                                    <?php echo e(__('filament-tables::table.sorting.fields.direction.options.asc')); ?>

                                                                </option>

                                                                <option
                                                                    value="desc"
                                                                >
                                                                    <?php echo e(__('filament-tables::table.sorting.fields.direction.options.desc')); ?>

                                                                </option>
                                                             <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $attributes = $__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__attributesOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b)): ?>
<?php $component = $__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b; ?>
<?php unset($__componentOriginal819ab5a8a54f49cec5eda6eb1911322c039e48909b5a13bde4d2fb526ea3253b); ?>
<?php endif; ?>
                                                         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $attributes = $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7)): ?>
<?php $component = $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7; ?>
<?php unset($__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7); ?>
<?php endif; ?>
                                                    </label>
                                                </div>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && ($maxSelectableRecords !== 1) && (! $selectsGroupsOnly)): ?>
                                                <input
                                                    aria-label="<?php echo e(__('filament-tables::table.fields.bulk_select_page.label')); ?>"
                                                    type="checkbox"
                                                    <?php if($isSelectionDisabled): ?>
                                                        disabled
                                                    <?php elseif($maxSelectableRecords): ?>
                                                        x-bind:disabled="
                                                            const recordsOnPage = getRecordsOnPage()

                                                            return recordsOnPage.length && ! areRecordsToggleable(recordsOnPage)
                                                        "
                                                    <?php endif; ?>
                                                    x-bind:checked="
                                                        const recordsOnPage = getRecordsOnPage()

                                                        if (recordsOnPage.length && areRecordsSelected(recordsOnPage)) {
                                                            $el.checked = true
                                                            $el.indeterminate = false

                                                            return 'checked'
                                                        }

                                                        $el.checked = false
                                                        $el.indeterminate =
                                                            recordsOnPage.length && areRecordsPartiallySelected(recordsOnPage)

                                                        return null
                                                    "
                                                    x-on:click="toggleSelectRecordsOnPage"
                                                    
                                                    <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk-select-page.checkbox.stacked.'.e(Str::random()).''; ?>wire:key="<?php echo e($this->getId()); ?>.table.bulk-select-page.checkbox.stacked.<?php echo e(Str::random()); ?>"
                                                    wire:loading.attr="disabled"
                                                    wire:target="<?php echo e($loadingTargetsWireTarget); ?>"
                                                    class="fi-ta-page-checkbox fi-checkbox-input"
                                                />
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                        </td>
                                    </tr>
                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasColumnGroups): ?>
                                    <tr class="fi-ta-table-head-groups-row">
                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(count($records)): ?>
                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isReordering): ?>
                                                <th></th>
                                            <?php else: ?>
                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && in_array($recordActionsPosition, [RecordActionsPosition::BeforeCells, RecordActionsPosition::BeforeColumns])): ?>
                                                    <th></th>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && $recordCheckboxPosition === RecordCheckboxPosition::BeforeCells): ?>
                                                    <th></th>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $columnsLayout; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $columnGroup): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($columnGroup instanceof Column): ?>
                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($columnGroup->isVisible() && (! $columnGroup->isToggledHidden())): ?>
                                                    <th></th>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                            <?php elseif($columnGroup instanceof ColumnGroup): ?>
                                              <?phphp
                                                    $columnGroupColumnsCount = count($columnGroup->getVisibleColumns());
                                              ?>hp

                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($columnGroupColumnsCount): ?>
                                                    <th
                                                        colspan="<?php echo e($columnGroupColumnsCount); ?>"
                                                        scope="colgroup"
                                                        <?php echo e($columnGroup->getExtraHeaderAttributeBag()->class([
                                                                'fi-ta-header-group-cell',
                                                                'fi-wrapped' => $columnGroup->canHeaderWrap(),
                                                                ((($columnGroupAlignment = $columnGroup->getAlignment()) instanceof Filament\Support\Enums\Alignment) ? "fi-align-{$columnGroupAlignment->value}" : (is_string($columnGroupAlignment) ? $columnGroupAlignment : '')),
                                                                (filled($columnGroupHiddenFrom = $columnGroup->getHiddenFrom()) ? "{$columnGroupHiddenFrom}:fi-hidden" : ''),
                                                                (filled($columnGroupVisibleFrom = $columnGroup->getVisibleFrom()) ? "{$columnGroupVisibleFrom}:fi-visible" : ''),
                                                            ])); ?>

                                                    >
                                                        <?php echo e($columnGroup->getLabel()); ?>

                                                    </th>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>

                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if((! $isReordering) && count($records)): ?>
                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && in_array($recordActionsPosition, [RecordActionsPosition::AfterColumns, RecordActionsPosition::AfterCells])): ?>
                                                <th></th>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && $recordCheckboxPosition === RecordCheckboxPosition::AfterCells): ?>
                                                <th></th>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                    </tr>
                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                <tr>
                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(count($records)): ?>
                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isReordering): ?>
                                            <th></th>
                                        <?php else: ?>
                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && $recordActionsPosition === RecordActionsPosition::BeforeCells): ?>
                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($recordActionsColumnLabel): ?>
                                                    <th
                                                        scope="col"
                                                        class="fi-ta-header-cell"
                                                    >
                                                        <?php echo e($recordActionsColumnLabel); ?>

                                                    </th>
                                                <?php else: ?>
                                                    <th
                                                        aria-label="<?php echo e(trans_choice('filament-tables::table.columns.actions.label', $flatRecordActionsCount)); ?>"
                                                        scope="col"
                                                        class="fi-ta-actions-header-cell fi-ta-empty-header-cell"
                                                    ></th>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && $recordCheckboxPosition === RecordCheckboxPosition::BeforeCells): ?>
                                                <th
                                                    scope="col"
                                                    class="fi-ta-cell fi-ta-selection-cell"
                                                >
                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(($maxSelectableRecords !== 1) && (! $selectsGroupsOnly)): ?>
                                                        <input
                                                            aria-label="<?php echo e(__('filament-tables::table.fields.bulk_select_page.label')); ?>"
                                                            type="checkbox"
                                                            <?php if($isSelectionDisabled): ?>
                                                                disabled
                                                            <?php elseif($maxSelectableRecords): ?>
                                                                x-bind:disabled="
                                                                    const recordsOnPage = getRecordsOnPage()

                                                                    return recordsOnPage.length && ! areRecordsToggleable(recordsOnPage)
                                                                "
                                                            <?php endif; ?>
                                                            x-bind:checked="
                                                                const recordsOnPage = getRecordsOnPage()

                                                                if (recordsOnPage.length && areRecordsSelected(recordsOnPage)) {
                                                                    $el.checked = true
                                                                    $el.indeterminate = false

                                                                    return 'checked'
                                                                }

                                                                $el.checked = false
                                                                $el.indeterminate =
                                                                    recordsOnPage.length && areRecordsPartiallySelected(recordsOnPage)

                                                                return null
                                                            "
                                                            x-on:click="toggleSelectRecordsOnPage"
                                                            
                                                            <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk-select-page.checkbox.'.e(Str::random()).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk-select-page.checkbox.'.e(Str::random()).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk-select-page.checkbox.'.e(Str::random()).''; ?>wire:key="<?php echo e($this->getId()); ?>.table.bulk-select-page.checkbox.<?php echo e(Str::random()); ?>"
                                                            wire:loading.attr="disabled"
                                                            wire:target="<?php echo e($loadingTargetsWireTarget); ?>"
                                                            class="fi-ta-page-checkbox fi-checkbox-input"
                                                        />
                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                                </th>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && $recordActionsPosition === RecordActionsPosition::BeforeColumns): ?>
                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($recordActionsColumnLabel): ?>
                                                    <th
                                                        scope="col"
                                                        class="fi-ta-header-cell"
                                                    >
                                                        <?php echo e($recordActionsColumnLabel); ?>

                                                    </th>
                                                <?php else: ?>
                                                    <th
                                                        aria-label="<?php echo e(trans_choice('filament-tables::table.columns.actions.label', $flatRecordActionsCount)); ?>"
                                                        scope="col"
                                                        class="fi-ta-actions-header-cell fi-ta-empty-header-cell"
                                                    ></th>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                  <?phphp
                                        $hasHeaderCellRenderHook = FilamentView::hasRenderHook(TablesRenderHook::HEADER_CELL, scopes: static::class);
                                  ?>hp

                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $columns; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $column): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasHeaderCellRenderHook && filled($headerCellView = FilamentView::renderHook(TablesRenderHook::HEADER_CELL, scopes: static::class, data: [
                                                 'column' => $column,
                                                 'isReordering' => $isReordering,
                                             ]))): ?>
                                            <?php echo e($headerCellView); ?>

                                        <?php else: ?>
                                          <?phphp
                                                $columnName = $column->getName();
                                                $columnLabel = $column->getLabel();
                                                $columnAlignment = $column->getAlignment();
                                                $columnWidth = $column->getWidth();
                                                $isColumnActivelySorted = $getSortColumn() === $column->getName();
                                                $isColumnSortable = $column->isSortable() && (! $isReordering);

                                                // A custom label may contain interactive elements, which are invalid inside a native `<button>`, so a `<span>` with button semantics is used instead.
                                                $columnSortControlTag = ($columnLabel instanceof Htmlable) ? 'span' : 'button';

                                                $columnHeaderTooltip = $column->getHeaderTooltip();
                                                $columnHeaderTooltipAttribute = ($columnHeaderTooltip instanceof Htmlable)
                                                    ? 'x-tooltip.html'
                                                    : 'x-tooltip';
                                          ?>hp

                                            <th
                                                <?php if($isColumnSortable): ?>
                                                    aria-sort="<?php echo e($isColumnActivelySorted ? ($sortDirection === 'asc' ? 'ascending' : 'descending') : 'none'); ?>"
                                                <?php endif; ?>
                                                scope="col"
                                                <?php echo e($column->getExtraHeaderAttributeBag()
                                                        ->class([
                                                            'fi-ta-header-cell',
                                                            'fi-ta-header-cell-' . str($columnName)->camel()->kebab(),
                                                            'fi-growable' => blank($columnWidth) && $column->canGrow(default: false),
                                                            'fi-grouped' => $column->getGroup(),
                                                            'fi-wrapped' => $column->canHeaderWrap(),
                                                            'fi-ta-header-cell-sorted' => $isColumnActivelySorted,
                                                            ((($columnAlignment = $column->getAlignment()) instanceof Filament\Support\Enums\Alignment) ? "fi-align-{$columnAlignment->value}" : (is_string($columnAlignment) ? $columnAlignment : '')),
                                                            (filled($columnHiddenFrom = $column->getHiddenFrom()) ? "{$columnHiddenFrom}:fi-hidden" : ''),
                                                            (filled($columnVisibleFrom = $column->getVisibleFrom()) ? "{$columnVisibleFrom}:fi-visible" : ''),
                                                        ])
                                                        ->style([
                                                            ('width: ' . e($columnWidth)) => filled($columnWidth),
                                                        ])); ?>

                                            >
                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isColumnSortable): ?>
                                                    <<?php echo e($columnSortControlTag); ?>

                                                        <?php if($columnSortControlTag === 'button'): ?>
                                                            type="button"
                                                        <?php else: ?>
                                                            role="button"
                                                            tabindex="0"
                                                            x-on:keydown.enter.prevent.stop="$wire.sortTable('<?php echo e($columnName); ?>')"
                                                            x-on:keydown.space.prevent.stop="$wire.sortTable('<?php echo e($columnName); ?>')"
                                                        <?php endif; ?>
                                                        wire:click="sortTable('<?php echo e($columnName); ?>')"
                                                        wire:loading.attr="disabled"
                                                        wire:target="sortTable('<?php echo e($columnName); ?>')"
                                                        class="fi-ta-header-cell-sort-btn"
                                                    >
                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($columnHeaderTooltip)): ?>
                                                            <span
                                                                <?php echo e($columnHeaderTooltipAttribute); ?>="{
                                                                    content: <?php echo \Illuminate\Support\Js::from($columnHeaderTooltip)->toHtml() ?>,
                                                                    theme: $store.theme,
                                                                }"
                                                                class="fi-ta-header-cell-tooltip"
                                                            >
                                                                <?php echo e($columnLabel); ?>

                                                            </span>
                                                        <?php else: ?>
                                                            <?php echo e($columnLabel); ?>

                                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                        <?php echo e(\Filament\Support\generate_icon_html(($isColumnActivelySorted && $sortDirection === 'asc') ? Heroicon::ChevronUp : Heroicon::ChevronDown, alias: match (true) {
                                                                $isColumnActivelySorted && ($sortDirection === 'asc') => TablesIconAlias::HEADER_CELL_SORT_ASC_BUTTON,
                                                                $isColumnActivelySorted && ($sortDirection === 'desc') => TablesIconAlias::HEADER_CELL_SORT_DESC_BUTTON,
                                                                default => TablesIconAlias::HEADER_CELL_SORT_BUTTON,
                                                            }, attributes: (new Filament\Support\View\ComponentAttributeBag([
                                                                'wire:loading.remove.delay.' . config('filament.livewire_loading_delay', 'default') => true,
                                                                'wire:target' => "sortTable('{$columnName}')",
                                                            ])))); ?>


                                                        <?php echo e(\Filament\Support\generate_loading_indicator_html(new Filament\Support\View\ComponentAttributeBag([
                                                                'wire:loading.delay.' . config('filament.livewire_loading_delay', 'default') => '',
                                                                'wire:target' => "sortTable('{$columnName}')",
                                                            ]))); ?>

                                                    </<?php echo e($columnSortControlTag); ?>>
                                                <?php else: ?>
                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($columnHeaderTooltip)): ?>
                                                        <span
                                                            <?php echo e($columnHeaderTooltipAttribute); ?>="{
                                                                content: <?php echo \Illuminate\Support\Js::from($columnHeaderTooltip)->toHtml() ?>,
                                                                theme: $store.theme,
                                                            }"
                                                            class="fi-ta-header-cell-tooltip"
                                                        >
                                                            <?php echo e($columnLabel); ?>

                                                        </span>
                                                    <?php else: ?>
                                                        <?php echo e($columnLabel); ?>

                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                            </th>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>

                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if((! $isReordering) && count($records)): ?>
                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && $recordActionsPosition === RecordActionsPosition::AfterColumns): ?>
                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($recordActionsColumnLabel): ?>
                                                <th
                                                    scope="col"
                                                    class="fi-ta-header-cell fi-align-end"
                                                >
                                                    <?php echo e($recordActionsColumnLabel); ?>

                                                </th>
                                            <?php else: ?>
                                                <th
                                                    aria-label="<?php echo e(trans_choice('filament-tables::table.columns.actions.label', $flatRecordActionsCount)); ?>"
                                                    scope="col"
                                                    class="fi-ta-actions-header-cell fi-ta-empty-header-cell"
                                                ></th>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && $recordCheckboxPosition === RecordCheckboxPosition::AfterCells): ?>
                                            <th
                                                scope="col"
                                                class="fi-ta-cell fi-ta-selection-cell"
                                            >
                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(($maxSelectableRecords !== 1) && (! $selectsGroupsOnly)): ?>
                                                    <input
                                                        aria-label="<?php echo e(__('filament-tables::table.fields.bulk_select_page.label')); ?>"
                                                        type="checkbox"
                                                        <?php if($isSelectionDisabled): ?>
                                                            disabled
                                                        <?php elseif($maxSelectableRecords): ?>
                                                            x-bind:disabled="
                                                                const recordsOnPage = getRecordsOnPage()

                                                                return recordsOnPage.length && ! areRecordsToggleable(recordsOnPage)
                                                            "
                                                        <?php endif; ?>
                                                        x-bind:checked="
                                                            const recordsOnPage = getRecordsOnPage()

                                                            if (recordsOnPage.length && areRecordsSelected(recordsOnPage)) {
                                                                $el.checked = true
                                                                $el.indeterminate = false

                                                                return 'checked'
                                                            }

                                                            $el.checked = false
                                                            $el.indeterminate =
                                                                recordsOnPage.length && areRecordsPartiallySelected(recordsOnPage)

                                                            return null
                                                        "
                                                        x-on:click="toggleSelectRecordsOnPage"
                                                        
                                                        <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk-select-page.checkbox.'.e(Str::random()).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk-select-page.checkbox.'.e(Str::random()).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk-select-page.checkbox.'.e(Str::random()).''; ?>wire:key="<?php echo e($this->getId()); ?>.table.bulk-select-page.checkbox.<?php echo e(Str::random()); ?>"
                                                        wire:loading.attr="disabled"
                                                        wire:target="<?php echo e($loadingTargetsWireTarget); ?>"
                                                        class="fi-ta-page-checkbox fi-checkbox-input"
                                                    />
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                            </th>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && $recordActionsPosition === RecordActionsPosition::AfterCells): ?>
                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($recordActionsColumnLabel): ?>
                                                <th
                                                    scope="col"
                                                    class="fi-ta-header-cell fi-align-end"
                                                >
                                                    <?php echo e($recordActionsColumnLabel); ?>

                                                </th>
                                            <?php else: ?>
                                                <th
                                                    aria-label="<?php echo e(trans_choice('filament-tables::table.columns.actions.label', $flatRecordActionsCount)); ?>"
                                                    scope="col"
                                                    class="fi-ta-actions-header-cell fi-ta-empty-header-cell"
                                                ></th>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                </tr>
                            </thead>

                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isColumnSearchVisible || count($records)): ?>
                                <tbody
                                    <?php if($isReorderable): ?>
                                        x-on:end.stop="
                                            $wire.reorderTable(
                                                $event.target.sortable.toArray(),
                                                $event.item.getAttribute('x-sortable-item'),
                                            )
                                        "
                                        x-sortable
                                        data-sortable-animation-duration="<?php echo e($getReorderAnimationDuration()); ?>"
                                    <?php endif; ?>
                                >
                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isColumnSearchVisible): ?>
                                        <tr
                                            class="fi-ta-row fi-ta-row-not-reorderable"
                                        >
                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(count($records)): ?>
                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isReordering): ?>
                                                    <td></td>
                                                <?php else: ?>
                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && in_array($recordActionsPosition, [RecordActionsPosition::BeforeCells, RecordActionsPosition::BeforeColumns])): ?>
                                                        <td></td>
                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && $recordCheckboxPosition === RecordCheckboxPosition::BeforeCells): ?>
                                                        <td></td>
                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $columns; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $column): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                              <?phphp
                                                    $columnName = $column->getName();
                                              ?>hp

                                                <td
                                                    class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                                        'fi-ta-cell',
                                                        'fi-ta-individual-search-cell' => $isIndividuallySearchable = $column->isIndividuallySearchable(),
                                                        'fi-ta-individual-search-cell-' . str($columnName)->camel()->kebab() => $isIndividuallySearchable,
                                                    ]); ?>"
                                                >
                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isIndividuallySearchable): ?>
                                                        <?php if (isset($component)) { $__componentOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.search-field','data' => ['debounce' => $searchDebounce,'onBlur' => $isSearchOnBlur,'wireModel' => 'tableColumnSearches.' . $columnName]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::search-field'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['debounce' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($searchDebounce),'on-blur' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isSearchOnBlur),'wire-model' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute('tableColumnSearches.' . $columnName)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc)): ?>
<?php $attributes = $__attributesOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc; ?>
<?php unset($__attributesOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc)): ?>
<?php $component = $__componentOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc; ?>
<?php unset($__componentOriginalb0392522c9cdb744cefaf7380e4bf6d08d7e15c7cad01fc4516bdc4540f978cc); ?>
<?php endif; ?>
                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                                </td>
                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>

                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if((! $isReordering) && count($records)): ?>
                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && in_array($recordActionsPosition, [RecordActionsPosition::AfterColumns, RecordActionsPosition::AfterCells])): ?>
                                                    <td></td>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && $recordCheckboxPosition === RecordCheckboxPosition::AfterCells): ?>
                                                    <td></td>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                        </tr>
                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(count($records)): ?>
                                      <?phphp
                                            $isRecordRowStriped = false;
                                            $previousRecord = null;
                                            $previousRecordGroupKey = null;
                                            $previousRecordGroupTitle = null;
                                      ?>hp

                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $records; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $record): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                          <?phphp
                                                $recordAction = $getRecordAction($record);
                                                $recordKey = $getRecordKey($record);
                                                $recordUrl = $getRecordUrl($record);
                                                $openRecordUrlInNewTab = $shouldOpenRecordUrlInNewTab($record);
                                                $recordGroupKey = $group?->getStringKey($record);
                                                $recordGroupTitle = $group?->getTitle($record, $recordGroupKey);
                                                $recordIsSelectable = $isSelectionEnabled && $isRecordSelectable($record);

                                                $recordActions = $recordActionsByRecordKey[$recordKey]
                                                    ?? ($hasRecordActionsForAnyRecord ? $reduceVisibleRecordActions($record) : []);
                                          ?>hp

                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if((string) $recordGroupTitle !== (string) $previousRecordGroupTitle): ?>
                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasSummary && (! $isReordering) && filled($previousRecordGroupTitle)): ?>
                                                  <?phphp
                                                        $groupColumn = $group->getColumn();
                                                        $groupScopedAllTableSummaryQuery = $group->scopeQuery($this->getAllTableSummaryQuery(), $previousRecord);
                                                  ?>hp

                                                    <?php if (isset($component)) { $__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.summary.row','data' => ['actions' => $hasRecordActionsForAnyRecord,'actionsPosition' => $recordActionsPosition,'columns' => $columns,'groupColumn' => $groupColumn,'groupsOnly' => $isGroupsOnly,'heading' => $isGroupsOnly ? $previousRecordGroupTitle : __('filament-tables::table.summary.subheadings.group', ['group' => $previousRecordGroupTitle, 'label' => $pluralModelLabel]),'query' => $groupScopedAllTableSummaryQuery,'recordCheckboxPosition' => $recordCheckboxPosition,'selectedState' => $groupedSummarySelectedState[$previousRecordGroupKey] ?? [],'selectionEnabled' => $isSelectionEnabled]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::summary.row'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['actions' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasRecordActionsForAnyRecord),'actions-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($recordActionsPosition),'columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columns),'group-column' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupColumn),'groups-only' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isGroupsOnly),'heading' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isGroupsOnly ? $previousRecordGroupTitle : __('filament-tables::table.summary.subheadings.group', ['group' => $previousRecordGroupTitle, 'label' => $pluralModelLabel])),'query' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupScopedAllTableSummaryQuery),'record-checkbox-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($recordCheckboxPosition),'selected-state' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupedSummarySelectedState[$previousRecordGroupKey] ?? []),'selection-enabled' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isSelectionEnabled)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da)): ?>
<?php $attributes = $__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da; ?>
<?php unset($__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da)): ?>
<?php $component = $__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da; ?>
<?php unset($__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da); ?>
<?php endif; ?>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $isGroupsOnly): ?>
                                                    <tr
                                                        class="fi-ta-row fi-ta-group-header-row"
                                                    >
                                                      <?phphp
                                                            $isRecordGroupCollapsible = $group?->isCollapsible();
                                                            $groupHeaderColspan = $columnsCount;

                                                            if ($isSelectionEnabled) {
                                                                $groupHeaderColspan--;

                                                                if (
                                                                    ($recordCheckboxPosition === RecordCheckboxPosition::BeforeCells) &&
                                                                    $hasRecordActionsForAnyRecord &&
                                                                    ($recordActionsPosition === RecordActionsPosition::BeforeCells)
                                                                ) {
                                                                    $groupHeaderColspan--;
                                                                }
                                                            }
                                                      ?>hp

                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && $recordCheckboxPosition === RecordCheckboxPosition::BeforeCells): ?>
                                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && $recordActionsPosition === RecordActionsPosition::BeforeCells): ?>
                                                                <td></td>
                                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                            <td
                                                                class="fi-ta-cell fi-ta-group-selection-cell"
                                                            >
                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($maxSelectableRecords !== 1): ?>
                                                                    <input
                                                                        aria-label="<?php echo e(__('filament-tables::table.fields.bulk_select_group.label', ['title' => $recordGroupTitle])); ?>"
                                                                        type="checkbox"
                                                                        data-group-selectable-record-keys="<?php echo e(json_encode($this->getGroupedSelectableTableRecordKeys($recordGroupKey))); ?>"
                                                                        <?php if($isSelectionDisabled): ?>
                                                                            disabled
                                                                        <?php else: ?>
                                                                            x-on:click="toggleSelectRecords(JSON.parse($el.dataset.groupSelectableRecordKeys))"
                                                                            <?php if($maxSelectableRecords): ?>
                                                                                x-bind:disabled="
                                                                                    const recordsInGroup = JSON.parse($el.dataset.groupSelectableRecordKeys)

                                                                                    return recordsInGroup.length && ! areRecordsToggleable(recordsInGroup)
                                                                                "
                                                                            <?php endif; ?>
                                                                        <?php endif; ?>
                                                                        x-bind:checked="
                                                                            const recordsInGroup = JSON.parse($el.dataset.groupSelectableRecordKeys)

                                                                            if (recordsInGroup.length && areRecordsSelected(recordsInGroup)) {
                                                                                $el.checked = true
                                                                                $el.indeterminate = false

                                                                                return 'checked'
                                                                            }

                                                                            $el.checked = false
                                                                            $el.indeterminate =
                                                                                recordsInGroup.length && areRecordsPartiallySelected(recordsInGroup)

                                                                            return null
                                                                        "
                                                                        <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk_select_group.checkbox.'.e($page).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk_select_group.checkbox.'.e($page).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk_select_group.checkbox.'.e($page).''; ?>wire:key="<?php echo e($this->getId()); ?>.table.bulk_select_group.checkbox.<?php echo e($page); ?>"
                                                                        wire:loading.attr="disabled"
                                                                        wire:target="<?php echo e($loadingTargetsWireTarget); ?>"
                                                                        class="fi-ta-group-checkbox fi-checkbox-input"
                                                                    />
                                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                                            </td>
                                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                        <td
                                                            colspan="<?php echo e($groupHeaderColspan); ?>"
                                                            class="fi-ta-group-header-cell"
                                                        >
                                                            <div
                                                                <?php if($isRecordGroupCollapsible): ?>
                                                                    x-on:click="toggleCollapseGroup(<?php echo \Illuminate\Support\Js::from($recordGroupTitle)->toHtml() ?>)"
                                                                    x-bind:class="isGroupCollapsed(<?php echo \Illuminate\Support\Js::from($recordGroupTitle)->toHtml() ?>) ? 'fi-collapsed' : null"
                                                                <?php endif; ?>
                                                                class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                                                    'fi-ta-group-header',
                                                                    'fi-collapsible' => $isRecordGroupCollapsible,
                                                                ]); ?>"
                                                            >
                                                                <div>
                                                                    <<?php echo e($secondLevelHeadingTag); ?>

                                                                        class="fi-ta-group-heading"
                                                                    >
                                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($recordGroupLabel = ($group->isTitlePrefixedWithLabel() ? $group->getLabel() : null))): ?>
                                                                                <?php echo e($recordGroupLabel); ?>:
                                                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                                        <?php echo e($recordGroupTitle); ?>

                                                                    </<?php echo e($secondLevelHeadingTag); ?>>

                                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($recordGroupDescription = $group->getDescription($record, $recordGroupTitle))): ?>
                                                                        <p
                                                                            class="fi-ta-group-description"
                                                                        >
                                                                            <?php echo e($recordGroupDescription); ?>

                                                                        </p>
                                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                                                </div>

                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isRecordGroupCollapsible): ?>
                                                                    <button
                                                                        aria-label="<?php echo e(filled($recordGroupLabel) ? ($recordGroupLabel . ': ' . $recordGroupTitle) : $recordGroupTitle); ?>"
                                                                        x-bind:aria-expanded="! isGroupCollapsed(<?php echo \Illuminate\Support\Js::from($recordGroupTitle)->toHtml() ?>)"
                                                                        type="button"
                                                                        class="fi-icon-btn fi-size-sm"
                                                                    >
                                                                        <?php echo e(\Filament\Support\generate_icon_html(Heroicon::ChevronUp, alias: TablesIconAlias::GROUPING_COLLAPSE_BUTTON, size: IconSize::Small)); ?>

                                                                    </button>
                                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                                            </div>
                                                        </td>

                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && $recordCheckboxPosition === RecordCheckboxPosition::AfterCells): ?>
                                                            <td
                                                                class="fi-ta-cell fi-ta-group-selection-cell"
                                                            >
                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($maxSelectableRecords !== 1): ?>
                                                                    <input
                                                                        aria-label="<?php echo e(__('filament-tables::table.fields.bulk_select_group.label', ['title' => $recordGroupTitle])); ?>"
                                                                        type="checkbox"
                                                                        data-group-selectable-record-keys="<?php echo e(json_encode($this->getGroupedSelectableTableRecordKeys($recordGroupKey))); ?>"
                                                                        <?php if($isSelectionDisabled): ?>
                                                                            disabled
                                                                        <?php else: ?>
                                                                            x-on:click="toggleSelectRecords(JSON.parse($el.dataset.groupSelectableRecordKeys))"
                                                                            <?php if($maxSelectableRecords): ?>
                                                                                x-bind:disabled="
                                                                                    const recordsInGroup = JSON.parse($el.dataset.groupSelectableRecordKeys)

                                                                                    return recordsInGroup.length && ! areRecordsToggleable(recordsInGroup)
                                                                                "
                                                                            <?php endif; ?>
                                                                        <?php endif; ?>
                                                                        x-bind:checked="
                                                                            const recordsInGroup = JSON.parse($el.dataset.groupSelectableRecordKeys)

                                                                            if (recordsInGroup.length && areRecordsSelected(recordsInGroup)) {
                                                                                $el.checked = true
                                                                                $el.indeterminate = false

                                                                                return 'checked'
                                                                            }

                                                                            $el.checked = false
                                                                            $el.indeterminate =
                                                                                recordsInGroup.length && areRecordsPartiallySelected(recordsInGroup)

                                                                            return null
                                                                        "
                                                                        <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk_select_group.checkbox.'.e($page).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk_select_group.checkbox.'.e($page).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.bulk_select_group.checkbox.'.e($page).''; ?>wire:key="<?php echo e($this->getId()); ?>.table.bulk_select_group.checkbox.<?php echo e($page); ?>"
                                                                        wire:loading.attr="disabled"
                                                                        wire:target="<?php echo e($loadingTargetsWireTarget); ?>"
                                                                        class="fi-ta-group-checkbox fi-checkbox-input"
                                                                    />
                                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                                            </td>
                                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                                    </tr>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                              <?phphp
                                                    $isRecordRowStriped = false;
                                              ?>hp
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! $isGroupsOnly): ?>
                                                <tr
                                                    <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.records.'.e($recordKey).''; ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.records.'.e($recordKey).''; ?>wire:key="<?php echo e($this->getId()); ?>.table.records.<?php echo e($recordKey); ?>"
                                                    <?php echo e($isReordering ? 'x-sortable-handle' : null); ?>

                                                    <?php echo $isReordering ? 'x-sortable-item="' . e($recordKey) . '"' : null; ?>

                                                    x-bind:class="{
                                                        <?php echo e($group?->isCollapsible() ? '\'fi-collapsed\': isGroupCollapsed(' . \Illuminate\Support\Js::from($recordGroupTitle) . '),' : ''); ?>

                                                        'fi-selected': <?php echo \Illuminate\Support\Js::from($recordIsSelectable)->toHtml() ?> && isRecordSelected(<?php echo \Illuminate\Support\Js::from($recordKey)->toHtml() ?>),
                                                    }"
                                                    class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                                        'fi-ta-row',
                                                        'fi-clickable' => $recordAction || $recordUrl,
                                                        'fi-striped' => $isStriped && $isRecordRowStriped,
                                                        ...$getRecordClasses($record),
                                                    ]); ?>"
                                                >
                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isReordering): ?>
                                                        <td class="fi-ta-cell">
                                                            <button
                                                                aria-label="<?php echo e(__('filament-tables::table.actions.reorder_record.label', ['key' => $recordKey])); ?>"
                                                                class="fi-ta-reorder-handle fi-icon-btn"
                                                                type="button"
                                                            >
                                                                <?php echo e(\Filament\Support\generate_icon_html(Heroicon::Bars2, alias: TablesIconAlias::REORDER_HANDLE)); ?>

                                                            </button>
                                                        </td>
                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && $recordActionsPosition === RecordActionsPosition::BeforeCells && (! $isReordering)): ?>
                                                        <td class="fi-ta-cell">
                                                            <div
                                                                class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                                                    'fi-ta-actions',
                                                                    match ($recordActionsAlignment) {
                                                                        Alignment::Center => 'fi-align-center',
                                                                        Alignment::Start, Alignment::Left => 'fi-align-start',
                                                                        Alignment::Between, Alignment::Justify => 'fi-align-between',
                                                                        Alignment::End, Alignment::Right => '',
                                                                        default => is_string($recordActionsAlignment) ? $recordActionsAlignment : '',
                                                                    },
                                                                ]); ?>"
                                                            >
                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $recordActions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $action): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                                    <?php echo e($action); ?>

                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                            </div>
                                                        </td>
                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && ($recordCheckboxPosition === RecordCheckboxPosition::BeforeCells) && (! $isReordering)): ?>
                                                        <td
                                                            class="fi-ta-cell fi-ta-selection-cell"
                                                        >
                                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($recordIsSelectable): ?>
                                                                <input
                                                                    aria-label="<?php echo e(__('filament-tables::table.fields.bulk_select_record.label', ['key' => $recordKey])); ?>"
                                                                    type="checkbox"
                                                                    <?php if($isSelectionDisabled): ?>
                                                                        disabled
                                                                    <?php elseif($maxSelectableRecords && ($maxSelectableRecords !== 1)): ?>
                                                                        x-bind:disabled="! areRecordsToggleable([<?php echo \Illuminate\Support\Js::from($recordKey)->toHtml() ?>])"
                                                                    <?php endif; ?>
                                                                    value="<?php echo e($recordKey); ?>"
                                                                    x-on:click="toggleSelectedRecord(<?php echo \Illuminate\Support\Js::from($recordKey)->toHtml() ?>)"
                                                                    x-bind:checked="isRecordSelected(<?php echo \Illuminate\Support\Js::from($recordKey)->toHtml() ?>) ? 'checked' : null"
                                                                    data-group="<?php echo e($recordGroupKey); ?>"
                                                                    wire:loading.attr="disabled"
                                                                    wire:target="<?php echo e($loadingTargetsWireTarget); ?>"
                                                                    class="fi-ta-record-checkbox fi-checkbox-input"
                                                                />
                                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                                        </td>
                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && $recordActionsPosition === RecordActionsPosition::BeforeColumns && (! $isReordering)): ?>
                                                        <td class="fi-ta-cell">
                                                            <div
                                                                class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                                                    'fi-ta-actions',
                                                                    match ($recordActionsAlignment) {
                                                                        Alignment::Center => 'fi-align-center',
                                                                        Alignment::Start, Alignment::Left => 'fi-align-start',
                                                                        Alignment::Between, Alignment::Justify => 'fi-align-between',
                                                                        Alignment::End, Alignment::Right => '',
                                                                        default => is_string($recordActionsAlignment) ? $recordActionsAlignment : '',
                                                                    },
                                                                ]); ?>"
                                                            >
                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $recordActions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $action): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                                    <?php echo e($action); ?>

                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                            </div>
                                                        </td>
                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $columns; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $column): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                      <?phphp
                                                            $column->record($record);
                                                            $column->rowLoop($loop->parent);
                                                            $column->recordKey($recordKey);

                                                            $columnAction = $column->getAction();
                                                            $columnUrl = $column->getUrl();
                                                            $columnHasStateBasedUrls = $column->hasStateBasedUrls();
                                                            $isColumnClickDisabled = $column->isClickDisabled() || $isReordering;

                                                            $columnWrapperTag = match (true) {
                                                                ($columnUrl || ($recordUrl && $columnAction === null)) && (! $columnHasStateBasedUrls) && (! $isColumnClickDisabled) => 'a',
                                                                ($columnAction || $recordAction) && (! $columnHasStateBasedUrls) && (! $isColumnClickDisabled) => 'button',
                                                                default => 'div',
                                                            };

                                                            if ($columnWrapperTag === 'button') {
                                                                if ($columnAction instanceof Action) {
                                                                    $columnWireClickAction = "mountTableAction('{$columnAction->getName()}', '{$recordKey}')";
                                                                } elseif ($columnAction) {
                                                                    $columnWireClickAction = "callTableColumnAction('{$column->getName()}', '{$recordKey}')";
                                                                } else {
                                                                    if ($this->getTable()->getAction($recordAction)) {
                                                                        $columnWireClickAction = "mountTableAction('{$recordAction}', '{$recordKey}')";
                                                                    } else {
                                                                        $columnWireClickAction = "{$recordAction}('{$recordKey}')";
                                                                    }
                                                                }
                                                            }
                                                      ?>hp

                                                        <td
                                                            <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($this->getId()).'.table.record.'.e($recordKey).'.column.'.e($column->getName()).''; ?>wire:key="<?php echo e($this->getId()); ?>.table.record.<?php echo e($recordKey); ?>.column.<?php echo e($column->getName()); ?>"
                                                            <?php echo $column->getCachedCellAttributeHtml(); ?>

                                                        >
                                                            <?php echo $isStackedOnMobile ? '<div class="fi-ta-cell-label">' . e($column->getLabel()) . '</div><div class="fi-ta-cell-content">' : ''; ?>

                                                            <<?php echo e($columnWrapperTag); ?>

                                                                <?php if($columnWrapperTag === 'a'): ?>
                                                                    <?php echo e(\Filament\Support\generate_href_html($columnUrl ?: $recordUrl, $columnUrl ? $column->shouldOpenUrlInNewTab() : $openRecordUrlInNewTab, hasNestedClickEventHandler: true)); ?>

                                                                    <?php if(blank($columnUrl) && filled($recordUrl)): ?>
                                                                        <?php echo e($getExtraRecordLinkAttributeBag($record)); ?>

                                                                    <?php endif; ?>
                                                                <?php elseif($columnWrapperTag === 'button'): ?>
                                                                    type="button"
                                                                    wire:click.prevent.stop="<?php echo e($columnWireClickAction); ?>"
                                                                    wire:loading.attr="disabled"
                                                                    wire:target="<?php echo e($columnWireClickAction); ?>"
                                                                <?php endif; ?>
                                                                class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                                                    'fi-ta-col',
                                                                    'fi-ta-col-has-column-url' => ($columnWrapperTag === 'a') && filled($columnUrl),
                                                                ]); ?>"
                                                            >
                                                                <?php echo e($column); ?>

                                                            </<?php echo e($columnWrapperTag); ?>>
                                                            <?php echo $isStackedOnMobile ? '</div>' : ''; ?>

                                                        </td>
                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>

                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && $recordActionsPosition === RecordActionsPosition::AfterColumns && (! $isReordering)): ?>
                                                        <td class="fi-ta-cell">
                                                            <div
                                                                class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                                                    'fi-ta-actions',
                                                                    match ($recordActionsAlignment) {
                                                                        Alignment::Center => 'fi-align-center',
                                                                        Alignment::Start, Alignment::Left => 'fi-align-start',
                                                                        Alignment::Between, Alignment::Justify => 'fi-align-between',
                                                                        Alignment::End, Alignment::Right => '',
                                                                        default => is_string($recordActionsAlignment) ? $recordActionsAlignment : '',
                                                                    },
                                                                ]); ?>"
                                                            >
                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $recordActions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $action): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                                    <?php echo e($action); ?>

                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                            </div>
                                                        </td>
                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isSelectionEnabled && $recordCheckboxPosition === RecordCheckboxPosition::AfterCells && (! $isReordering)): ?>
                                                        <td
                                                            class="fi-ta-cell fi-ta-selection-cell"
                                                        >
                                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($recordIsSelectable): ?>
                                                                <input
                                                                    aria-label="<?php echo e(__('filament-tables::table.fields.bulk_select_record.label', ['key' => $recordKey])); ?>"
                                                                    type="checkbox"
                                                                    <?php if($isSelectionDisabled): ?>
                                                                        disabled
                                                                    <?php elseif($maxSelectableRecords && ($maxSelectableRecords !== 1)): ?>
                                                                        x-bind:disabled="! areRecordsToggleable([<?php echo \Illuminate\Support\Js::from($recordKey)->toHtml() ?>])"
                                                                    <?php endif; ?>
                                                                    value="<?php echo e($recordKey); ?>"
                                                                    x-on:click="toggleSelectedRecord(<?php echo \Illuminate\Support\Js::from($recordKey)->toHtml() ?>)"
                                                                    x-bind:checked="isRecordSelected(<?php echo \Illuminate\Support\Js::from($recordKey)->toHtml() ?>) ? 'checked' : null"
                                                                    data-group="<?php echo e($recordGroupKey); ?>"
                                                                    wire:loading.attr="disabled"
                                                                    wire:target="<?php echo e($loadingTargetsWireTarget); ?>"
                                                                    class="fi-ta-record-checkbox fi-checkbox-input"
                                                                />
                                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                                        </td>
                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasRecordActionsForAnyRecord && $recordActionsPosition === RecordActionsPosition::AfterCells && (! $isReordering)): ?>
                                                        <td class="fi-ta-cell">
                                                            <div
                                                                class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                                                    'fi-ta-actions',
                                                                    match ($recordActionsAlignment) {
                                                                        Alignment::Center => 'fi-align-center',
                                                                        Alignment::Start, Alignment::Left => 'fi-align-start',
                                                                        Alignment::Between, Alignment::Justify => 'fi-align-between',
                                                                        Alignment::End, Alignment::Right => '',
                                                                        default => is_string($recordActionsAlignment) ? $recordActionsAlignment : '',
                                                                    },
                                                                ]); ?>"
                                                            >
                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $recordActions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $action): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                                    <?php echo e($action); ?>

                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                            </div>
                                                        </td>
                                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                                </tr>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                          <?phphp
                                                $isRecordRowStriped = ! $isRecordRowStriped;
                                                $previousRecord = $record;
                                                $previousRecordGroupKey = $recordGroupKey;
                                                $previousRecordGroupTitle = $recordGroupTitle;
                                          ?>hp
                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>

                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasSummary && (! $isReordering) && filled($previousRecordGroupTitle) && $this->shouldRenderTrailingGroupedTableSummary($previousRecord)): ?>
                                          <?phphp
                                                $groupColumn = $group->getColumn();
                                                $groupScopedAllTableSummaryQuery = $group->scopeQuery($this->getAllTableSummaryQuery(), $previousRecord);
                                          ?>hp

                                            <?php if (isset($component)) { $__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.summary.row','data' => ['actions' => $hasRecordActionsForAnyRecord,'actionsPosition' => $recordActionsPosition,'columns' => $columns,'groupColumn' => $groupColumn,'groupsOnly' => $isGroupsOnly,'heading' => $isGroupsOnly ? $previousRecordGroupTitle : __('filament-tables::table.summary.subheadings.group', ['group' => $previousRecordGroupTitle, 'label' => $pluralModelLabel]),'query' => $groupScopedAllTableSummaryQuery,'recordCheckboxPosition' => $recordCheckboxPosition,'selectedState' => $groupedSummarySelectedState[$previousRecordGroupKey] ?? [],'selectionEnabled' => $isSelectionEnabled]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::summary.row'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['actions' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasRecordActionsForAnyRecord),'actions-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($recordActionsPosition),'columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columns),'group-column' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupColumn),'groups-only' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isGroupsOnly),'heading' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isGroupsOnly ? $previousRecordGroupTitle : __('filament-tables::table.summary.subheadings.group', ['group' => $previousRecordGroupTitle, 'label' => $pluralModelLabel])),'query' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupScopedAllTableSummaryQuery),'record-checkbox-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($recordCheckboxPosition),'selected-state' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupedSummarySelectedState[$previousRecordGroupKey] ?? []),'selection-enabled' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isSelectionEnabled)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da)): ?>
<?php $attributes = $__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da; ?>
<?php unset($__attributesOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da)): ?>
<?php $component = $__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da; ?>
<?php unset($__componentOriginal4ff23d7e5feff2d83e72a158f14a878b319e764006eca7c75e2d453aa77c89da); ?>
<?php endif; ?>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasSummary && (! $isReordering)): ?>
                                          <?phphp
                                                $groupColumn = $group?->getColumn();
                                          ?>hp

                                            <?php if (isset($component)) { $__componentOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.summary.index','data' => ['actions' => $hasRecordActionsForAnyRecord,'actionsPosition' => $recordActionsPosition,'allTableSummary' => $hasAllTableSummary,'columns' => $columns,'groupColumn' => $groupColumn,'groupsOnly' => $isGroupsOnly,'pageSummary' => $hasPageSummary,'pluralModelLabel' => $pluralModelLabel,'recordCheckboxPosition' => $recordCheckboxPosition,'records' => $records,'selectionEnabled' => $isSelectionEnabled]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::summary'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['actions' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasRecordActionsForAnyRecord),'actions-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($recordActionsPosition),'all-table-summary' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasAllTableSummary),'columns' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($columns),'group-column' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($groupColumn),'groups-only' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isGroupsOnly),'page-summary' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasPageSummary),'plural-model-label' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($pluralModelLabel),'record-checkbox-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($recordCheckboxPosition),'records' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($records),'selection-enabled' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($isSelectionEnabled)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed)): ?>
<?php $attributes = $__attributesOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed; ?>
<?php unset($__attributesOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed)): ?>
<?php $component = $__componentOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed; ?>
<?php unset($__componentOriginalf1b432a383fa7f8e0911f754cc748004ad00b6da5b8fe64bb8249ceaadc7d6ed); ?>
<?php endif; ?>
                                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                </tbody>
                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(($records !== null) && count($records) && $contentFooter): ?>
                                <tfoot>
                                    <tr>
                                        <?php echo e($contentFooter->with([
                                                'columns' => $columns,
                                                'records' => $records,
                                            ])); ?>

                                    </tr>
                                </tfoot>
                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                        </table>
                    <?php elseif($records === null): ?>
                        <div
                            role="status"
                            aria-busy="true"
                            aria-live="polite"
                            class="fi-ta-table-loading-ctn"
                        >
                            <?php echo e(\Filament\Support\generate_loading_indicator_html(size: IconSize::TwoExtraLarge)); ?>


                            <span class="fi-sr-only">
                                <?php echo e(__('filament-tables::table.loading')); ?>

                            </span>
                        </div>
                    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                </div>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasEmptyState): ?>
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($emptyState = $getEmptyState()): ?>
                    <?php echo e($emptyState); ?>

                <?php else: ?>
                    <div class="fi-ta-empty-state" role="status">
                        <div class="fi-ta-empty-state-content">
                            <div class="fi-ta-empty-state-icon-bg">
                                <?php echo e(\Filament\Support\generate_icon_html($getEmptyStateIcon(), size: IconSize::Large)); ?>

                            </div>

                            <<?php echo e($secondLevelHeadingTag); ?>

                                class="fi-ta-empty-state-heading"
                            >
                                <?php echo e($getEmptyStateHeading()); ?>

                            </<?php echo e($secondLevelHeadingTag); ?>>

                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(filled($emptyStateDescription = $getEmptyStateDescription())): ?>
                                <p class="fi-ta-empty-state-description">
                                    <?php echo e($emptyStateDescription); ?>

                                </p>
                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($emptyStateActions = array_filter(
                                     $getEmptyStateActions(),
                                     fn (\Filament\Actions\Action | ActionGroup $action): bool => $action->isVisible(),
                                 )): ?>
                                <div
                                    class="fi-ta-actions fi-align-center fi-wrapped"
                                >
                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $emptyStateActions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $action): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                        <?php echo e($action); ?>

                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                </div>
                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                        </div>
                    </div>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasPagination): ?>
              <?phphp
                    $hasExtremePaginationLinks = $hasExtremePaginationLinks();
                    $paginationPageOptions = $getPaginationPageOptions();
              ?>hp

                <?php if (isset($component)) { $__componentOriginal20e7055c0fc94f2c69f4b801366f3804dad18bd1ec2bfa2bac261a8ae75d0dca = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal20e7055c0fc94f2c69f4b801366f3804dad18bd1ec2bfa2bac261a8ae75d0dca = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.pagination.index','data' => ['extremeLinks' => $hasExtremePaginationLinks,'pageOptions' => $paginationPageOptions,'paginator' => $records]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::pagination'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['extreme-links' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasExtremePaginationLinks),'page-options' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($paginationPageOptions),'paginator' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($records)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal20e7055c0fc94f2c69f4b801366f3804dad18bd1ec2bfa2bac261a8ae75d0dca)): ?>
<?php $attributes = $__attributesOriginal20e7055c0fc94f2c69f4b801366f3804dad18bd1ec2bfa2bac261a8ae75d0dca; ?>
<?php unset($__attributesOriginal20e7055c0fc94f2c69f4b801366f3804dad18bd1ec2bfa2bac261a8ae75d0dca); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal20e7055c0fc94f2c69f4b801366f3804dad18bd1ec2bfa2bac261a8ae75d0dca)): ?>
<?php $component = $__componentOriginal20e7055c0fc94f2c69f4b801366f3804dad18bd1ec2bfa2bac261a8ae75d0dca; ?>
<?php unset($__componentOriginal20e7055c0fc94f2c69f4b801366f3804dad18bd1ec2bfa2bac261a8ae75d0dca); ?>
<?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasFiltersBelowContent): ?>
                <?php if (isset($component)) { $__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.filters','data' => ['applyAction' => $filtersApplyAction,'form' => $filtersForm,'headingTag' => $secondLevelHeadingTag,'class' => 'fi-ta-filters-below-content','resetActionPosition' => $filtersResetActionPosition]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::filters'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['apply-action' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersApplyAction),'form' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersForm),'heading-tag' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($secondLevelHeadingTag),'class' => 'fi-ta-filters-below-content','reset-action-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersResetActionPosition)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b)): ?>
<?php $attributes = $__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b; ?>
<?php unset($__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b)): ?>
<?php $component = $__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b; ?>
<?php unset($__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b); ?>
<?php endif; ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        </div>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasFiltersAfterContent): ?>
            <div
                wire:ignore.self
                x-ref="filtersContentContainer"
                x-transition:enter-start="fi-opacity-0"
                x-transition:leave-end="fi-opacity-0"
                x-bind:class="{ 'fi-open': areFiltersOpen }"
                class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                    'fi-ta-filters-after-content-ctn',
                    'lg:fi-open' => ! $hasCollapsibleFilters,
                    (($filtersFormWidth ??= Width::ExtraSmall) instanceof Width) ? "fi-width-{$filtersFormWidth->value}" : (is_string($filtersFormWidth) ? $filtersFormWidth : null),
                ]); ?>"
            >
                <?php if (isset($component)) { $__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-tables::components.filters','data' => ['applyAction' => $filtersApplyAction,'form' => $filtersForm,'headingTag' => $secondLevelHeadingTag,'class' => 'fi-ta-filters-after-content','resetActionPosition' => $filtersResetActionPosition]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-tables::filters'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['apply-action' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersApplyAction),'form' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersForm),'heading-tag' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($secondLevelHeadingTag),'class' => 'fi-ta-filters-after-content','reset-action-position' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($filtersResetActionPosition)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b)): ?>
<?php $attributes = $__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b; ?>
<?php unset($__attributesOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b); ?>
<?php endif; ?>
<?php if (isset($__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b)): ?>
<?php $component = $__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b; ?>
<?php unset($__componentOriginale0d43856f75e1944bbd9ee324c260f689cf5a0b24a5167f18f10e9b182cc5e6b); ?>
<?php endif; ?>
            </div>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    </div>

    <?php if (isset($component)) { $__componentOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament-actions::components.modals','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament-actions::modals'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15)): ?>
<?php $attributes = $__attributesOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15; ?>
<?php unset($__attributesOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15)): ?>
<?php $component = $__componentOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15; ?>
<?php unset($__componentOriginalf7d5d7a06189d05716605d7969f96c2014762aee4bf6e6fce8767a8ea6537c15); ?>
<?php endif; ?>
</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\tables\resources\views/index.blade.php ENDPATH**/ ?>