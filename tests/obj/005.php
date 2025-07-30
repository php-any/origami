namespace tests\obj;

class Users {
    public $name = "";
}

class DB<T> {
    public $where = {};

    public function __construct() {

    }

    public function where($key, $value) {
        $this->where[$key] = []
        $this->where[$key][] = 1;
        $this->where[$key]->push(2);

        return $this;
    }

    public function get() {
        return [
            new T(),
        ];
    }
}

$list = DB<Users>()->where("name", "张三")->get();

dump($list, "OK");