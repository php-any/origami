<?php
    use Filament\Support\Enums\Alignment;
    use Filament\Support\Enums\VerticalAlignment;
?>

<div>
    <div
        class="<?php echo \Illuminate\Support\Arr::toCssClasses([
            'fi-no',
            'fi-align-' . static::$alignment->value,
            'fi-vertical-align-' . static::$verticalAlignment->value,
        ]); ?>"
        role="status"
        aria-atomic="false"
    >
        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::openLoop(); ?><?php endif; ?><?php $__currentLoopData = $notifications; $__env->addLoop($__currentLoopData); foreach($__currentLoopData as $notification): $__env->incrementLoopIndices(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::startLoopIteration(); ?><?php endif; ?>
            <?php echo e($notification); ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::endLoop(); ?><?php endif; ?><?php endforeach; $__env->popLoop(); $loop = $__env->getLastLoop(); ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::closeLoop(); ?><?php endif; ?>
    </div>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($broadcastChannel = $this->getBroadcastChannel()): ?>
                        <?php
                    $__scriptKey = '2322055652-0';
                    ob_start();
                ?>
            <script>
                window.addEventListener('EchoLoaded', () => {
                    window.Echo.private(<?php echo \Illuminate\Support\Js::from($broadcastChannel)->toHtml() ?>).notification(
                        (notification) => {
                            setTimeout(
                                () =>
                                    $wire.handleBroadcastNotification(
                                        notification,
                                    ),
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
</div>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\vendor\filament\notifications\resources\views/notifications.blade.php ENDPATH**/ ?>