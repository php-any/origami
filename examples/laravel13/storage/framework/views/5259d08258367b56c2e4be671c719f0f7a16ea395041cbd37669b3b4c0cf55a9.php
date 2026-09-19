<?php
    use Filament\Notifications\View\NotificationsIconAlias;
    use Filament\Support\Enums\Alignment;
    use Filament\Support\Icons\Heroicon;
    use Filament\Support\View\ComponentAttributeBag as FilamentComponentAttributeBag;
    use Filament\Support\View\Components\BadgeComponent;
    use Illuminate\Contracts\Pagination\Paginator;

    $notifications = $this->getNotifications();
    $unreadNotificationsCount = $this->getUnreadNotificationsCount();
    $hasNotifications = $notifications->count();
    $isPaginated = $notifications instanceof Paginator && $notifications->hasPages();
    $pollingInterval = $this->getPollingInterval();
?>

<div class="fi-no-database">
    
    <?php if (isset($component)) { $__componentOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalac14c6129a52092b4cf2a9c777954f31327584c9a32747b71f4c8283ba12223c = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.modal.index','data' => ['alignment' => $hasNotifications ? null : Alignment::Center,'ariaLabelledby' => 'database-notifications.heading','closeButton' => true,'description' => $hasNotifications ? null : __('filament-notifications::database.modal.empty.description'),'extraModalWindowAttributeBag' => 
            new FilamentComponentAttributeBag([
                'autofocus' => true,
                'tabindex' => '-1',
            ])
        ,'heading' => $hasNotifications ? null : __('filament-notifications::database.modal.empty.heading'),'icon' => $hasNotifications ? null : Heroicon::OutlinedBellSlash,'iconAlias' => 
            $hasNotifications
            ? null
            : NotificationsIconAlias::DATABASE_MODAL_EMPTY_STATE
        ,'iconColor' => $hasNotifications ? null : 'gray','id' => 'database-notifications','slideOver' => true,'stickyHeader' => $hasNotifications,'teleport' => 'body','width' => 'md','class' => 'fi-no-database','attributes' => 
            new FilamentComponentAttributeBag([
                'wire:poll.' . $pollingInterval => $pollingInterval ? '' : false,
            ])
        ]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::modal'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['alignment' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasNotifications ? null : Alignment::Center),'aria-labelledby' => 'database-notifications.heading','close-button' => true,'description' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasNotifications ? null : __('filament-notifications::database.modal.empty.description')),'extra-modal-window-attribute-bag' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
            new FilamentComponentAttributeBag([
                'autofocus' => true,
                'tabindex' => '-1',
            ])
        ),'heading' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasNotifications ? null : __('filament-notifications::database.modal.empty.heading')),'icon' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasNotifications ? null : Heroicon::OutlinedBellSlash),'icon-alias' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
            $hasNotifications
            ? null
            : NotificationsIconAlias::DATABASE_MODAL_EMPTY_STATE
        ),'icon-color' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasNotifications ? null : 'gray'),'id' => 'database-notifications','slide-over' => true,'sticky-header' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($hasNotifications),'teleport' => 'body','width' => 'md','class' => 'fi-no-database','attributes' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute(
            new FilamentComponentAttributeBag([
                'wire:poll.' . $pollingInterval => $pollingInterval ? '' : false,
            ])
        )]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($trigger = $this->getTrigger()): ?>
             <?php $__env->slot('trigger', null, []); ?> 
                <?php echo e($trigger->with(['unreadNotificationsCount' => $unreadNotificationsCount])); ?>

             <?php $__env->endSlot(); ?>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($hasNotifications): ?>
             <?php $__env->slot('header', null, []); ?> 
                <div>
                    <h2
                        id="database-notifications.heading"
                        class="fi-modal-heading"
                    >
                        <?php echo e(__('filament-notifications::database.modal.heading')); ?>


                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($unreadNotificationsCount): ?>
                            <span
                                <?php echo e((new FilamentComponentAttributeBag)->color(BadgeComponent::class, 'primary')->class([
                                        'fi-badge fi-size-xs',
                                    ])); ?>

                            >
                                <?php echo e($unreadNotificationsCount); ?>

                            </span>
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                    </h2>

                    <div class="fi-ac">
                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($unreadNotificationsCount && $this->markAllNotificationsAsReadAction?->isVisible()): ?>
                            <?php echo e($this->markAllNotificationsAsReadAction); ?>

                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($this->clearNotificationsAction?->isVisible()): ?>
                            <?php echo e($this->clearNotificationsAction); ?>

                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
                    </div>
                </div>
             <?php $__env->endSlot(); ?>

            <div
                aria-label="<?php echo e(__('filament-notifications::database.modal.heading')); ?>"
                role="list"
                class="fi-no-notifications"
            >
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $notifications; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $notification): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
                    <div
                        role="listitem"
                        <?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::$currentLoop['key'] = ''.e($notification->getKey()).'.database-notifications.ctn'; ?>wire:key="<?php echo e($notification->getKey()); ?>.database-notifications.ctn"
                        class="<?php echo \Illuminate\Support\Arr::toCssClasses([
                            'fi-no-notification-read-ctn' => ! $notification->unread(),
                            'fi-no-notification-unread-ctn' => $notification->unread(),
                        ]); ?>"
                    >
                        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($notification->unread()): ?>
                            <span class="fi-sr-only">
                                <?php echo e(__('filament-notifications::database.modal.unread_label')); ?>

                            </span>
                        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

                        <?php echo e($this->getNotification($notification)->inline()); ?>

                    </div>
                <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
            </div>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($broadcastChannel = $this->getBroadcastChannel()): ?>
                                <?php
                    $__scriptKey = '2077465202-0';
                    ob_start();
                ?>
                    <script>
                        window.addEventListener('EchoLoaded', () => {
                            window.Echo.private(<?php echo \Illuminate\Support\Js::from($broadcastChannel)->toHtml() ?>).listen(
                                '.database-notifications.sent',
                                () => {
                                    setTimeout(
                                        () => $wire.call('$refresh'),
                                        500,
                                    )
                                },
                            )
                        })

                        if (window.Echo) {
                            window.dispatchEvent(new CustomEvent('EchoLoaded'))
                        }
                    </script>
                                <?php
                    $__output = ob_get_clean();

                    \Livewire\store($this)->push('scripts', $__output, $__scriptKey)
                ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

            <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($isPaginated): ?>
                 <?php $__env->slot('footer', null, []); ?> 
                    <?php if (isset($component)) { $__componentOriginal20e7055c0fc94f2c69f4b801366f3804dad18bd1ec2bfa2bac261a8ae75d0dca = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal20e7055c0fc94f2c69f4b801366f3804dad18bd1ec2bfa2bac261a8ae75d0dca = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'filament::components.pagination.index','data' => ['paginator' => $notifications]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('filament::pagination'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['paginator' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($notifications)]); ?>
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
                 <?php $__env->endSlot(); ?>
            <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
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
</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\notifications\resources\views/database-notifications.blade.php ENDPATH**/ ?>