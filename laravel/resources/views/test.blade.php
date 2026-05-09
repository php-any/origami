<!DOCTYPE html>
<html>
<head>
    <title>Origami Laravel 测试</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <style>
        body { font-family: sans-serif; max-width: 800px; margin: 2em auto; padding: 0 2em; }
        h1 { color: #333; }
        .info { background: #f5f5f5; padding: 1em; border-radius: 4px; margin: 1em 0; }
        .pass { color: green; font-weight: bold; }
    </style>
</head>
<body>
    <h1>Origami Laravel 测试</h1>

    <div class="info">
        <p><span class="pass">✓</span> 欢迎 <?php echo $name ?? '访客'; ?></p>
        <p><span class="pass">✓</span> <?php echo $greeting; ?></p>
    </div>

    @if (isset($items) && count($items) > 0)
        <h3>项目列表 (<?php echo count($items); ?> 项)</h3>
        <ul>
        <?php foreach ($items as $item): ?>
            <li><?php echo $item; ?></li>
        <?php endforeach; ?>
        </ul>
    @endif

    <div class="info">
        <p><span class="pass">✓</span> Blade 替代语法正常</p>
    </div>
</body>
</html>
