<?php

class UserEntity {
    public string $id;
    public string $username;
    public string $nickname;
}

$entity = new UserEntity();
$fields = ['id' => 'u1', 'username' => 'alice', 'nickname' => 'Alice'];

foreach ($fields as $key => $value) {
    $entity->$key = $value;
}

if ($entity->id == 'u1' && $entity->username == 'alice' && $entity->nickname == 'Alice') {
    Log::info("动态属性赋值测试通过");
} else {
    Log::fatal("动态属性赋值测试失败");
}
