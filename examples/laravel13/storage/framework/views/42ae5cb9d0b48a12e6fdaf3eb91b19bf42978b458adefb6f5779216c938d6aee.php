<?php
    use Filament\Support\Facades\FilamentView;
    use Filament\Support\Icons\Heroicon;
    use Filament\View\PanelsIconAlias;
    use Filament\View\PanelsRenderHook;
    use Illuminate\Support\Arr;

    $debounce = filament()->getGlobalSearchDebounce();
    $keyBindings = filament()->getGlobalSearchKeyBindings();
    $suffix = filament()->getGlobalSearchFieldSuffix();
?>

<div class="fi-global-search-ctn">
    <?php echo e(FilamentView::renderHook(PanelsRenderHook::GLOBAL_SEARCH_START)); ?>


    <div
        x-on:focus-first-global-search-result.stop="$el.querySelector('.fi-global-search-result-link')?.focus()"
        class="fi-global-search"
    >
        <div x-id="['input']" class="fi-global-search-field">
            <label x-bind:for="$id('input')" class="fi-sr-only">
                <?php echo e(__('filament-panels::global-search.field.label')); ?>

            </label>

            <?php if (isset($component)) { $__componentOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal264db04203c004b4f9e51b1b2c6c7c4047ced98999c21746d8dafd0d23736bd7 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.input.wrapper','data' => ['prefixIcon' => Heroicon::MagnifyingGlass,'prefixIconAlias' => PanelsIconAlias::GLOBAL_SEARCH_FIELD,'inlinePrefix' => true,'suffix' => $suffix,'inlineSuffix' => true,'wire:target' => 'search']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::input.wrapper'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['prefix-icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(Heroicon::MagnifyingGlass),'prefix-icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(PanelsIconAlias::GLOBAL_SEARCH_FIELD),'inline-prefix' => true,'suffix' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($suffix),'inline-suffix' => true,'wire:target' => 'search']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

                <input
                    autocomplete="off"
                    maxlength="1000"
                    placeholder="<?php echo e(__('filament-panels::global-search.field.placeholder')); ?>"
                    type="search"
                    <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = 'global-search.field.input'; ?>wire:key="global-search.field.input"
                    x-bind:id="$id('input')"
                    x-on:keydown.down.prevent.stop="$dispatch('focus-first-global-search-result')"
                    wire:model.live.debounce.<?php echo e($debounce); ?>="search"
                    x-mousetrap.global.<?php echo e(collect($keyBindings)->map(fn (string $keyBinding): string => str_replace('+', '-', $keyBinding))->implode('.')); ?>="document.getElementById($id('input'))?.focus()"
                    class="fi-input fi-input-has-inline-prefix"
                />
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
        </div>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($results !== null): ?>
            <div
                x-data="{
                    isOpen: false,

                    open(event) {
                        this.isOpen = true
                    },

                    close(event) {
                        this.isOpen = false
                    },
                }"
                x-init="$nextTick(() => open())"
                x-on:click.away="close()"
                x-on:keydown.escape.window="close()"
                x-on:keydown.up.prevent="$focus.wrap().previous()"
                x-on:keydown.down.prevent="$focus.wrap().next()"
                x-on:open-global-search-results.window="$nextTick(() => open())"
                x-show="isOpen"
                x-transition:enter-start="fi-transition-enter-start"
                x-transition:leave-end="fi-transition-leave-end"
                class="fi-global-search-results-ctn"
            >
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($results->getCategories()->isEmpty()): ?>
                    <p class="fi-global-search-no-results-message">
                        <?php echo e(__('filament-panels::global-search.no_results_message')); ?>

                    </p>
                <?php else: ?>
                    <ul class="fi-global-search-results">
                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $results->getCategories(); $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $group => $groupedResults): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                            <li class="fi-global-search-result-group">
                                <h3
                                    class="fi-global-search-result-group-header"
                                >
                                    <?php echo e($group); ?>

                                </h3>

                                <ul
                                    class="fi-global-search-result-group-results"
                                >
                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $groupedResults; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $result): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                        <?php
                                            $resultVisibleActions = $result->getVisibleActions();
                                        ?>

                                        <li
                                            class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                                                'fi-global-search-result',
                                                'fi-global-search-result-has-actions' => $resultVisibleActions,
                                            ]); ?>"
                                        >
                                            <a
                                                <?php echo e(\Filament\Support\generate_href_html($result->url)); ?>

                                                x-on:click="close()"
                                                class="fi-global-search-result-link"
                                            >
                                                <h4
                                                    class="fi-global-search-result-heading"
                                                >
                                                    <?php echo e($result->title); ?>

                                                </h4>

                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($result->details): ?>
                                                    <dl
                                                        class="fi-global-search-result-details"
                                                    >
                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $result->details; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $label => $value): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                            <div
                                                                class="fi-global-search-result-detail"
                                                            >
                                                                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isAssoc ??= Arr::isAssoc($result->details)): ?>
                                                                    <dt
                                                                        class="fi-global-search-result-detail-label"
                                                                    >
                                                                        <?php echo e($label); ?>:
                                                                    </dt>
                                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                                                                <dd
                                                                    class="fi-global-search-result-detail-value"
                                                                >
                                                                    <?php echo e($value); ?>

                                                                </dd>
                                                            </div>
                                                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                    </dl>
                                                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                            </a>

                                            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($resultVisibleActions): ?>
                                                <div
                                                    class="fi-global-search-result-actions"
                                                >
                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $resultVisibleActions; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $action): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                                                        <?php echo e($action); ?>

                                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                                </div>
                                            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                                        </li>
                                    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                                </ul>
                            </li>
                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
                    </ul>
                <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
            </div>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
    </div>

    <?php echo e(FilamentView::renderHook(PanelsRenderHook::GLOBAL_SEARCH_END)); ?>

</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\filament\resources\views/livewire/global-search.blade.php ENDPATH**/ ?>