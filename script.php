$users = [
    {"name": "Jane Doe", "email": "jane@example.com"},
    {"name": "John Doe", "email": "john@example.com"}
];

// 应该循环输出多个div和name
$html = <div for="$k, $v in $users">
    {$v.name}
</div>


echo $html;