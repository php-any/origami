$users = [
    {"name": "Jane Doe 1", "email": "jane@example.com"},
    {"name": "John Doe 2", "email": "john@example.com"}
];

// 应该循环输出多个div和name
$html = <div for="$k, $v in $users">
    {$k} => {$v.name}
</div>


echo $html;