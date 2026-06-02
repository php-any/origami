<?php

namespace Spring\Service;

use Spring\Model\Product;

class ProductService {
    
    private $products = [];
    
    public function __construct() {
        // 初始化示例数据
        $this->products = [
            1 => new Product(1, "iPhone 15 Pro", 7999.00, "电子产品", "Apple 最新旗舰手机"),
            2 => new Product(2, "MacBook Pro 14", 14999.00, "电子产品", "专业级笔记本电脑"),
            3 => new Product(3, "AirPods Pro", 1899.00, "电子产品", "无线降噪耳机"),
            4 => new Product(4, "iPad Air", 4799.00, "电子产品", "轻薄平板电脑"),
            5 => new Product(5, "机械键盘", 599.00, "配件", "Cherry 轴机械键盘")
        ];
    }
    
    /**
     * 获取所有商品
     */
    public function findAll() {
        return array_values($this->products);
    }
    
    /**
     * 根据 ID 查找商品
     */
    public function findById($id) {
        return $this->products[$id] ?? null;
    }
    
    /**
     * 创建新商品
     */
    public function create($data) {
        $id = count($this->products) + 1;
        
        $name = $data['name'] ?? '';
        $price = (float)($data['price'] ?? 0);
        $category = $data['category'] ?? '未分类';
        $description = $data['description'] ?? '';
        
        $product = new Product($id, $name, $price, $category, $description);
        $this->products[$id] = $product;
        
        return $product;
    }
    
    /**
     * 更新商品信息
     */
    public function update($id, $data) {
        if (!isset($this->products[$id])) {
            return null;
        }
        
        $product = $this->products[$id];
        
        if (isset($data['name'])) {
            $product->setName($data['name']);
        }
        if (isset($data['price'])) {
            $product->setPrice((float)$data['price']);
        }
        if (isset($data['category'])) {
            $product->setCategory($data['category']);
        }
        if (isset($data['description'])) {
            $product->setDescription($data['description']);
        }
        
        return $product;
    }
    
    /**
     * 删除商品
     */
    public function delete($id) {
        if (!isset($this->products[$id])) {
            return false;
        }
        
        unset($this->products[$id]);
        return true;
    }
    
    /**
     * 搜索商品
     */
    public function search($keyword = '', $category = '') {
        $results = [];
        
        foreach ($this->products as $product) {
            $matchKeyword = empty($keyword) || 
                           stripos($product->getName(), $keyword) !== false ||
                           stripos($product->getDescription(), $keyword) !== false;
            
            $matchCategory = empty($category) || 
                            $product->getCategory() === $category;
            
            if ($matchKeyword && $matchCategory) {
                $results[] = $product;
            }
        }
        
        return $results;
    }
    
    /**
     * 按分类获取商品
     */
    public function findByCategory($category) {
        $results = [];
        
        foreach ($this->products as $product) {
            if ($product->getCategory() === $category) {
                $results[] = $product;
            }
        }
        
        return $results;
    }
    
    /**
     * 获取价格区间内的商品
     */
    public function findByPriceRange($minPrice, $maxPrice) {
        $results = [];
        
        foreach ($this->products as $product) {
            if ($product->getPrice() >= $minPrice && $product->getPrice() <= $maxPrice) {
                $results[] = $product;
            }
        }
        
        return $results;
    }
}
