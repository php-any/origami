<?php

namespace Spring\Model;

class User {
    private $id;
    private $name;
    private $email;
    private $age;
    
    public function __construct($id = null, $name = '', $email = '', $age = 0) {
        $this->id = $id;
        $this->name = $name;
        $this->email = $email;
        $this->age = $age;
    }
    
    public function getId() {
        return $this->id;
    }
    
    public function setId($id) {
        $this->id = $id;
    }
    
    public function getName() {
        return $this->name;
    }
    
    public function setName($name) {
        $this->name = $name;
    }
    
    public function getEmail() {
        return $this->email;
    }
    
    public function setEmail($email) {
        $this->email = $email;
    }
    
    public function getAge() {
        return $this->age;
    }
    
    public function setAge($age) {
        $this->age = $age;
    }
    
    /**
     * 转换为数组
     */
    public function toArray() {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'email' => $this->email,
            'age' => $this->age
        ];
    }
    
    /**
     * JSON 序列化
     */
    public function toJson() {
        return json_encode($this->toArray());
    }
}
