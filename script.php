namespace tests\func;

class TestInit {
    public static $title;
    public $author;
}

// 不经过构造函数的创建方式
TestInit::title = "123"

$data = TestInit::title;

echo $data;