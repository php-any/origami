$html =
<html >
    <header>
        <title>
            {$title}
        </title>
    </header>
    <body>
        <div for="k, v in $list">
            这里会循环输出列表内容：{$v.name}
        </div>
    </body>
</html>

echo $html;