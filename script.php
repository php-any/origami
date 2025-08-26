namespace tests\func;

class TestInit {
    public $title;
    public $author;
}
// 不经过构造函数的创建方式
$data = TestInit {
    title: "Hello World",
    author: "Tony",
}

echo $data;