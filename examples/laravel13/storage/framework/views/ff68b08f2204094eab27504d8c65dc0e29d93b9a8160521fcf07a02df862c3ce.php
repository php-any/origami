                    <?php $layout->viewContext->mergeIntoNewEnvironment($__env); ?>

                    <?php $__env->startComponent($layout->view, $layout->params); ?>
                        <?php $__env->slot($layout->slotOrSection); ?>
                            <?php echo $content; ?>

                        <?php $__env->endSlot(); ?>

                        <?php
                        // Manually forward slots defined in the Livewire template into the layout component...
                        foreach ($layout->viewContext->slots[-1] ?? [] as $name => $slot) {
                            $__env->slot($name, attributes: $slot->attributes->getAttributes());
                            echo $slot->toHtml();
                            $__env->endSlot();
                        }
                        ?>
                    <?php echo $__env->renderComponent(); ?><?php /**PATH D:\gitcode.com\origami\examples\laravel13\storage\framework\views/fd139a0a48f9507dccc481e6b8a9c1f3c63ccdac54a7957eea7472336491c5c4.blade.php ENDPATH**/ ?>