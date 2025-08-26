namespace tests\func;

function dump($data) {
    $data->author = "dump";
    return $data;
}

class TestInit {
    public $title;
    public $author;
}

// 不经过构造函数的创建方式
$data = dump(TestInit {
    title: "Hello World",
    author: "Tony",
})

echo $data;