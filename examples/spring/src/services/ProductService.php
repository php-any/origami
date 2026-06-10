<?php

namespace Spring\Service;

use Database\DB;
use Spring\Model\Entity\ProductEntity;

class ProductService {

    private function db() {
        return DB::model(ProductEntity::class);
    }

    public function findAll() {
        $entities = $this->db()->orderBy("id ASC")->get();
        return QueryDemoService::entitiesToArray($entities);
    }

    public function findById($id) {
        $entity = $this->db()->where("id = ?", $id)->first();
        if (!$entity) {
            return null;
        }
        return $entity->toArray();
    }

    public function create($data) {
        $entity = new ProductEntity();
        $entity->name = $data['name'] ?? '';
        $entity->price = (float)($data['price'] ?? 0);
        $entity->category = $data['category'] ?? '未分类';
        $entity->description = $data['description'] ?? '';

        $result = DB::insert($entity);
        $entity->id = $result->insertId;

        return $entity->toArray();
    }

    public function update($id, $data) {
        if (!$this->db()->where("id = ?", $id)->first()) {
            return null;
        }

        $entity = new ProductEntity();
        if (isset($data['name'])) {
            $entity->name = $data['name'];
        }
        if (isset($data['price'])) {
            $entity->price = (float)$data['price'];
        }
        if (isset($data['category'])) {
            $entity->category = $data['category'];
        }
        if (isset($data['description'])) {
            $entity->description = $data['description'];
        }

        $this->db()->where("id = ?", $id)->update($entity);
        return $this->findById($id);
    }

    public function delete($id) {
        if (!$this->db()->where("id = ?", $id)->first()) {
            return false;
        }
        $this->db()->where("id = ?", $id)->delete();
        return true;
    }

    public function search($keyword = '', $category = '') {
        $query = $this->db();

        if (!empty($keyword) && !empty($category)) {
            $query = $query->where(
                "(name LIKE ? OR description LIKE ?) AND category = ?",
                "%" . $keyword . "%",
                "%" . $keyword . "%",
                $category
            );
        } elseif (!empty($keyword)) {
            $query = $query->where("name LIKE ? OR description LIKE ?", "%" . $keyword . "%", "%" . $keyword . "%");
        } elseif (!empty($category)) {
            $query = $query->where("category = ?", $category);
        }

        return QueryDemoService::entitiesToArray($query->orderBy("id ASC")->get());
    }

    public function findByCategory($category) {
        $entities = $this->db()->where("category = ?", $category)->orderBy("price DESC")->get();
        return QueryDemoService::entitiesToArray($entities);
    }

    public function findByPriceRange($minPrice, $maxPrice) {
        $entities = $this->db()
            ->where("price >= ? AND price <= ?", $minPrice, $maxPrice)
            ->orderBy("price ASC")
            ->get();
        return QueryDemoService::entitiesToArray($entities);
    }
}
