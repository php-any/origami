<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title><?php echo e(config('app.name')); ?></title>
    <?php echo app('Illuminate\Foundation\Vite')(['resources/css/app.css', 'resources/js/app.js']); ?>
</head>
<body>
    <?php echo e($slot); ?>

</body>
</html>
<?php /**PATH D:\gitcode.com\origami\examples\laravel13\resources\views/components/layouts/app.blade.php ENDPATH**/ ?>