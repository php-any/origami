<?php if (isset($component)) { $__componentOriginal506ce054ba5b2a394e1ac723cfadb34107316d45fd2d9e78042889e2aca2a309 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal506ce054ba5b2a394e1ac723cfadb34107316d45fd2d9e78042889e2aca2a309 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.layout','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::layout'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

    <?php if (isset($component)) { $__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.section-container','data' => ['class' => 'px-6 py-0 sm:py-0']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::section-container'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'px-6 py-0 sm:py-0']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

        <?php if (isset($component)) { $__componentOriginal62b0f2a022ffce8f785e726c7fb0a8aee9cc76bf869b82a9bf99ff4c69582e75 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal62b0f2a022ffce8f785e726c7fb0a8aee9cc76bf869b82a9bf99ff4c69582e75 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.topbar','data' => ['title' => $exception->title(),'markdown' => $exceptionAsMarkdown]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::topbar'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['title' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exception->title()),'markdown' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exceptionAsMarkdown)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal62b0f2a022ffce8f785e726c7fb0a8aee9cc76bf869b82a9bf99ff4c69582e75)): ?>
<?php $attributes = $__attributesOriginal62b0f2a022ffce8f785e726c7fb0a8aee9cc76bf869b82a9bf99ff4c69582e75; ?>
<?php unset($__attributesOriginal62b0f2a022ffce8f785e726c7fb0a8aee9cc76bf869b82a9bf99ff4c69582e75); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal62b0f2a022ffce8f785e726c7fb0a8aee9cc76bf869b82a9bf99ff4c69582e75)): ?>
<?php $component = $__componentOriginal62b0f2a022ffce8f785e726c7fb0a8aee9cc76bf869b82a9bf99ff4c69582e75; ?>
<?php unset($__componentOriginal62b0f2a022ffce8f785e726c7fb0a8aee9cc76bf869b82a9bf99ff4c69582e75); ?>
<?php endif; ?>
     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db)): ?>
<?php $attributes = $__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db; ?>
<?php unset($__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db)): ?>
<?php $component = $__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db; ?>
<?php unset($__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db); ?>
<?php endif; ?>

    <?php if (isset($component)) { $__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.separator','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::separator'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1)): ?>
<?php $attributes = $__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1; ?>
<?php unset($__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1)): ?>
<?php $component = $__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1; ?>
<?php unset($__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1); ?>
<?php endif; ?>

    <?php if (isset($component)) { $__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.section-container','data' => ['class' => 'flex flex-col gap-8 py-0 sm:py-0']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::section-container'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'flex flex-col gap-8 py-0 sm:py-0']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

        <?php if (isset($component)) { $__componentOriginalbbfed4420714067f54c690b4910bb8460e64526510b4537ac3e61ce1b2f35564 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalbbfed4420714067f54c690b4910bb8460e64526510b4537ac3e61ce1b2f35564 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.header','data' => ['exception' => $exception]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::header'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['exception' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exception)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalbbfed4420714067f54c690b4910bb8460e64526510b4537ac3e61ce1b2f35564)): ?>
<?php $attributes = $__attributesOriginalbbfed4420714067f54c690b4910bb8460e64526510b4537ac3e61ce1b2f35564; ?>
<?php unset($__attributesOriginalbbfed4420714067f54c690b4910bb8460e64526510b4537ac3e61ce1b2f35564); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalbbfed4420714067f54c690b4910bb8460e64526510b4537ac3e61ce1b2f35564)): ?>
<?php $component = $__componentOriginalbbfed4420714067f54c690b4910bb8460e64526510b4537ac3e61ce1b2f35564; ?>
<?php unset($__componentOriginalbbfed4420714067f54c690b4910bb8460e64526510b4537ac3e61ce1b2f35564); ?>
<?php endif; ?>
     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db)): ?>
<?php $attributes = $__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db; ?>
<?php unset($__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db)): ?>
<?php $component = $__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db; ?>
<?php unset($__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db); ?>
<?php endif; ?>

    <?php if (isset($component)) { $__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.separator','data' => ['class' => '-mt-5 -z-10']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::separator'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => '-mt-5 -z-10']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1)): ?>
<?php $attributes = $__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1; ?>
<?php unset($__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1)): ?>
<?php $component = $__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1; ?>
<?php unset($__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1); ?>
<?php endif; ?>

    <?php if (isset($component)) { $__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.section-container','data' => ['class' => 'flex flex-col gap-8 pt-14']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::section-container'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'flex flex-col gap-8 pt-14']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

        <?php if (isset($component)) { $__componentOriginale1dd2cb05e578b261b016bb6b139e8f5cc3cec94595cc311f8d698eb819ed2fb = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginale1dd2cb05e578b261b016bb6b139e8f5cc3cec94595cc311f8d698eb819ed2fb = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.trace','data' => ['exception' => $exception]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::trace'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['exception' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exception)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginale1dd2cb05e578b261b016bb6b139e8f5cc3cec94595cc311f8d698eb819ed2fb)): ?>
<?php $attributes = $__attributesOriginale1dd2cb05e578b261b016bb6b139e8f5cc3cec94595cc311f8d698eb819ed2fb; ?>
<?php unset($__attributesOriginale1dd2cb05e578b261b016bb6b139e8f5cc3cec94595cc311f8d698eb819ed2fb); ?>
<?php endif; ?>
<?php if (isset($__componentOriginale1dd2cb05e578b261b016bb6b139e8f5cc3cec94595cc311f8d698eb819ed2fb)): ?>
<?php $component = $__componentOriginale1dd2cb05e578b261b016bb6b139e8f5cc3cec94595cc311f8d698eb819ed2fb; ?>
<?php unset($__componentOriginale1dd2cb05e578b261b016bb6b139e8f5cc3cec94595cc311f8d698eb819ed2fb); ?>
<?php endif; ?>

        <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if($exception->previousExceptions()->isNotEmpty()): ?>
            <?php if (isset($component)) { $__componentOriginal58f0cd4b9f3c2378de470092adb9f30d7b3fbb6062798fcbaa0feb87811a1ac3 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal58f0cd4b9f3c2378de470092adb9f30d7b3fbb6062798fcbaa0feb87811a1ac3 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.previous-exceptions','data' => ['exception' => $exception]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::previous-exceptions'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['exception' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exception)]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal58f0cd4b9f3c2378de470092adb9f30d7b3fbb6062798fcbaa0feb87811a1ac3)): ?>
<?php $attributes = $__attributesOriginal58f0cd4b9f3c2378de470092adb9f30d7b3fbb6062798fcbaa0feb87811a1ac3; ?>
<?php unset($__attributesOriginal58f0cd4b9f3c2378de470092adb9f30d7b3fbb6062798fcbaa0feb87811a1ac3); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal58f0cd4b9f3c2378de470092adb9f30d7b3fbb6062798fcbaa0feb87811a1ac3)): ?>
<?php $component = $__componentOriginal58f0cd4b9f3c2378de470092adb9f30d7b3fbb6062798fcbaa0feb87811a1ac3; ?>
<?php unset($__componentOriginal58f0cd4b9f3c2378de470092adb9f30d7b3fbb6062798fcbaa0feb87811a1ac3); ?>
<?php endif; ?>
        <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>

        <?php if (isset($component)) { $__componentOriginal3b71b426e298d07dadc87603f63bb051c48ed7bed7276f017a9ac19a5f6bdde5 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal3b71b426e298d07dadc87603f63bb051c48ed7bed7276f017a9ac19a5f6bdde5 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.query','data' => ['queries' => $exception->applicationQueries()]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::query'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['queries' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exception->applicationQueries())]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal3b71b426e298d07dadc87603f63bb051c48ed7bed7276f017a9ac19a5f6bdde5)): ?>
<?php $attributes = $__attributesOriginal3b71b426e298d07dadc87603f63bb051c48ed7bed7276f017a9ac19a5f6bdde5; ?>
<?php unset($__attributesOriginal3b71b426e298d07dadc87603f63bb051c48ed7bed7276f017a9ac19a5f6bdde5); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal3b71b426e298d07dadc87603f63bb051c48ed7bed7276f017a9ac19a5f6bdde5)): ?>
<?php $component = $__componentOriginal3b71b426e298d07dadc87603f63bb051c48ed7bed7276f017a9ac19a5f6bdde5; ?>
<?php unset($__componentOriginal3b71b426e298d07dadc87603f63bb051c48ed7bed7276f017a9ac19a5f6bdde5); ?>
<?php endif; ?>
     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db)): ?>
<?php $attributes = $__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db; ?>
<?php unset($__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db)): ?>
<?php $component = $__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db; ?>
<?php unset($__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db); ?>
<?php endif; ?>

    <?php if (isset($component)) { $__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.separator','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::separator'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1)): ?>
<?php $attributes = $__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1; ?>
<?php unset($__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1)): ?>
<?php $component = $__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1; ?>
<?php unset($__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1); ?>
<?php endif; ?>

    <?php if (isset($component)) { $__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.section-container','data' => ['class' => 'flex flex-col gap-12']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::section-container'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'flex flex-col gap-12']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

        <?php if (isset($component)) { $__componentOriginal63ef56e865b52c7214b1b152335ebdb945c6dacc939aea2c9d6ca5092c7deecc = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal63ef56e865b52c7214b1b152335ebdb945c6dacc939aea2c9d6ca5092c7deecc = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.request-header','data' => ['headers' => $exception->requestHeaders()]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::request-header'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['headers' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exception->requestHeaders())]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal63ef56e865b52c7214b1b152335ebdb945c6dacc939aea2c9d6ca5092c7deecc)): ?>
<?php $attributes = $__attributesOriginal63ef56e865b52c7214b1b152335ebdb945c6dacc939aea2c9d6ca5092c7deecc; ?>
<?php unset($__attributesOriginal63ef56e865b52c7214b1b152335ebdb945c6dacc939aea2c9d6ca5092c7deecc); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal63ef56e865b52c7214b1b152335ebdb945c6dacc939aea2c9d6ca5092c7deecc)): ?>
<?php $component = $__componentOriginal63ef56e865b52c7214b1b152335ebdb945c6dacc939aea2c9d6ca5092c7deecc; ?>
<?php unset($__componentOriginal63ef56e865b52c7214b1b152335ebdb945c6dacc939aea2c9d6ca5092c7deecc); ?>
<?php endif; ?>

        <?php if (isset($component)) { $__componentOriginal169bdd07bc929ec35f0aadfc0756c1bcc202a8a0a003d3f01e2419579f7f4208 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal169bdd07bc929ec35f0aadfc0756c1bcc202a8a0a003d3f01e2419579f7f4208 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.request-body','data' => ['body' => $exception->requestBody()]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::request-body'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['body' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exception->requestBody())]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal169bdd07bc929ec35f0aadfc0756c1bcc202a8a0a003d3f01e2419579f7f4208)): ?>
<?php $attributes = $__attributesOriginal169bdd07bc929ec35f0aadfc0756c1bcc202a8a0a003d3f01e2419579f7f4208; ?>
<?php unset($__attributesOriginal169bdd07bc929ec35f0aadfc0756c1bcc202a8a0a003d3f01e2419579f7f4208); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal169bdd07bc929ec35f0aadfc0756c1bcc202a8a0a003d3f01e2419579f7f4208)): ?>
<?php $component = $__componentOriginal169bdd07bc929ec35f0aadfc0756c1bcc202a8a0a003d3f01e2419579f7f4208; ?>
<?php unset($__componentOriginal169bdd07bc929ec35f0aadfc0756c1bcc202a8a0a003d3f01e2419579f7f4208); ?>
<?php endif; ?>

        <?php if (isset($component)) { $__componentOriginalf8ccdb95bc45a542612cfa137a9b5bde68a535ccf5ad7d972dba5c52c14bb3cf = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginalf8ccdb95bc45a542612cfa137a9b5bde68a535ccf5ad7d972dba5c52c14bb3cf = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.routing','data' => ['routing' => $exception->applicationRouteContext()]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::routing'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['routing' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exception->applicationRouteContext())]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginalf8ccdb95bc45a542612cfa137a9b5bde68a535ccf5ad7d972dba5c52c14bb3cf)): ?>
<?php $attributes = $__attributesOriginalf8ccdb95bc45a542612cfa137a9b5bde68a535ccf5ad7d972dba5c52c14bb3cf; ?>
<?php unset($__attributesOriginalf8ccdb95bc45a542612cfa137a9b5bde68a535ccf5ad7d972dba5c52c14bb3cf); ?>
<?php endif; ?>
<?php if (isset($__componentOriginalf8ccdb95bc45a542612cfa137a9b5bde68a535ccf5ad7d972dba5c52c14bb3cf)): ?>
<?php $component = $__componentOriginalf8ccdb95bc45a542612cfa137a9b5bde68a535ccf5ad7d972dba5c52c14bb3cf; ?>
<?php unset($__componentOriginalf8ccdb95bc45a542612cfa137a9b5bde68a535ccf5ad7d972dba5c52c14bb3cf); ?>
<?php endif; ?>

        <?php if (isset($component)) { $__componentOriginaleaa46d40b07a44d39cbca18f8aa9392d92ef99d732b9a6425905af8b44ffc56f = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginaleaa46d40b07a44d39cbca18f8aa9392d92ef99d732b9a6425905af8b44ffc56f = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.routing-parameter','data' => ['routeParameters' => $exception->applicationRouteParametersContext()]] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::routing-parameter'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['routeParameters' => \Illuminate\View\Compilers\BladeCompiler::sanitizeComponentAttribute($exception->applicationRouteParametersContext())]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginaleaa46d40b07a44d39cbca18f8aa9392d92ef99d732b9a6425905af8b44ffc56f)): ?>
<?php $attributes = $__attributesOriginaleaa46d40b07a44d39cbca18f8aa9392d92ef99d732b9a6425905af8b44ffc56f; ?>
<?php unset($__attributesOriginaleaa46d40b07a44d39cbca18f8aa9392d92ef99d732b9a6425905af8b44ffc56f); ?>
<?php endif; ?>
<?php if (isset($__componentOriginaleaa46d40b07a44d39cbca18f8aa9392d92ef99d732b9a6425905af8b44ffc56f)): ?>
<?php $component = $__componentOriginaleaa46d40b07a44d39cbca18f8aa9392d92ef99d732b9a6425905af8b44ffc56f; ?>
<?php unset($__componentOriginaleaa46d40b07a44d39cbca18f8aa9392d92ef99d732b9a6425905af8b44ffc56f); ?>
<?php endif; ?>
     <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db)): ?>
<?php $attributes = $__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db; ?>
<?php unset($__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db)): ?>
<?php $component = $__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db; ?>
<?php unset($__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db); ?>
<?php endif; ?>

    <?php if (isset($component)) { $__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1 = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1 = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.separator','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::separator'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1)): ?>
<?php $attributes = $__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1; ?>
<?php unset($__attributesOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1)): ?>
<?php $component = $__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1; ?>
<?php unset($__componentOriginal6b0b50b5e1498f3f0d850f6d860fdd75e4117e177e794bbca0e87ee93252e4d1); ?>
<?php endif; ?>

    <?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if BLOCK]><![endif]--><?php endif; ?><?php if(! app()->runningUnitTests() && ! app()->runningInConsole()): ?>
        <?php if (isset($component)) { $__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.section-container','data' => ['class' => 'pb-0 sm:pb-0']] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::section-container'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes(['class' => 'pb-0 sm:pb-0']); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

            <?php if (isset($component)) { $__componentOriginal1f5d31b38ddf700d8d742567bba35f333d76ed96acabdf924ea10cfce13d5e8c = $component; } ?>
<?php if (isset($attributes)) { $__attributesOriginal1f5d31b38ddf700d8d742567bba35f333d76ed96acabdf924ea10cfce13d5e8c = $attributes; } ?>
<?php $component = Illuminate\View\AnonymousComponent::resolve(['view' => 'laravel-exceptions-renderer::components.laravel-ascii-spotlight','data' => []] + (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag ? $attributes->all() : [])); ?>
<?php $component->withName('laravel-exceptions-renderer::laravel-ascii-spotlight'); ?>
<?php if ($component->shouldRender()): ?>
<?php $__env->startComponent($component->resolveView(), $component->data()); ?>
<?php if (isset($attributes) && $attributes instanceof Illuminate\View\ComponentAttributeBag): ?>
<?php $attributes = $attributes->except(\Illuminate\View\AnonymousComponent::ignoredParameterNames()); ?>
<?php endif; ?>
<?php $component->withAttributes([]); ?>
<?php \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::processComponentKey($component); ?>

<?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal1f5d31b38ddf700d8d742567bba35f333d76ed96acabdf924ea10cfce13d5e8c)): ?>
<?php $attributes = $__attributesOriginal1f5d31b38ddf700d8d742567bba35f333d76ed96acabdf924ea10cfce13d5e8c; ?>
<?php unset($__attributesOriginal1f5d31b38ddf700d8d742567bba35f333d76ed96acabdf924ea10cfce13d5e8c); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal1f5d31b38ddf700d8d742567bba35f333d76ed96acabdf924ea10cfce13d5e8c)): ?>
<?php $component = $__componentOriginal1f5d31b38ddf700d8d742567bba35f333d76ed96acabdf924ea10cfce13d5e8c; ?>
<?php unset($__componentOriginal1f5d31b38ddf700d8d742567bba35f333d76ed96acabdf924ea10cfce13d5e8c); ?>
<?php endif; ?>
         <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db)): ?>
<?php $attributes = $__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db; ?>
<?php unset($__attributesOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db)): ?>
<?php $component = $__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db; ?>
<?php unset($__componentOriginal7ab96d2a82d233c97d6028ec78358f11018fbe90abae91186a2bbe8e9719d1db); ?>
<?php endif; ?>
    <?php endif; ?><?php if(\Livewire\Mechanisms\ExtendBlade\ExtendBlade::isRenderingLivewireComponent()): ?><!--[if ENDBLOCK]><![endif]--><?php endif; ?>
 <?php echo $__env->renderComponent(); ?>
<?php endif; ?>
<?php if (isset($__attributesOriginal506ce054ba5b2a394e1ac723cfadb34107316d45fd2d9e78042889e2aca2a309)): ?>
<?php $attributes = $__attributesOriginal506ce054ba5b2a394e1ac723cfadb34107316d45fd2d9e78042889e2aca2a309; ?>
<?php unset($__attributesOriginal506ce054ba5b2a394e1ac723cfadb34107316d45fd2d9e78042889e2aca2a309); ?>
<?php endif; ?>
<?php if (isset($__componentOriginal506ce054ba5b2a394e1ac723cfadb34107316d45fd2d9e78042889e2aca2a309)): ?>
<?php $component = $__componentOriginal506ce054ba5b2a394e1ac723cfadb34107316d45fd2d9e78042889e2aca2a309; ?>
<?php unset($__componentOriginal506ce054ba5b2a394e1ac723cfadb34107316d45fd2d9e78042889e2aca2a309); ?>
<?php endif; ?>
<?php /**PATH d:\gitcode.com\origami\examples\laravel13\vendor\laravel\framework\src\illuminate\foundation\providers/../resources/exceptions/renderer/show.blade.php ENDPATH**/ ?>