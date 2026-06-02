<?php

namespace Spring\Model;

class Product {
    private $id;
    private $name;
    private $price;
    private $category;
    private $description;
    
    public function __construct($id = null, $name = '', $price = 0.0, $category = '', $description = '') {
        $this->id = $id;
        $this->name = $name;
        $this->price = $price;
        $this->category = $category;
        $this->description = $description;
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
    
    public function getPrice() {
        return $this->price;
    }
    
    public function setPrice($price) {
        $this->price = $price;
    }
    
    public function getCategory() {
        return $this->category;
    }
    
    public function setCategory($category) {
        $this->category = $category;
    }
    
    public function getDescription() {
        return $this->description;
    }
    
    public function setDescription($description) {
        $this->description = $description;
    }
    
    /**
     * 转换为数组
     */
    public function toArray() {
        return [
            'id' => $this->id,
            'name' => $this->name,
            'price' => $this->price,
            'category' => $this->category,
            'description' => $this->description
        ];
    }
    
    /**
     * JSON 序列化
     */
    public function toJson() {
        return json_encode($this->toArray());
    }
    
    /**
     * 获取格式化的价格
     */
    public function getFormattedPrice() {
        return "¥" . number_format($this->price, 2);
    }
}
