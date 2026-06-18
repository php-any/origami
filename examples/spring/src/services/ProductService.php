<?php

namespace Spring\Service;

use Container\Singleton;
use Database\DB;
use Spring\Model\Entity\ProductEntity;

#[Singleton]
class ProductService {

    private function db(): DB {
        return DB<ProductEntity>();
    }

    public function findAll(): array {
        $entities = $this->db()->orderBy("id ASC")->get();
        return QueryDemoService::entitiesToArray($entities);
    }

    public function findById(int $id): ?array {
        $entity = $this->db()->where("id = ?", $id)->first();
        if (!$entity) {
            return null;
        }
        return $entity->toArray();
    }

    public function create(array $data): array {
        $entity = new ProductEntity();
        $entity->name = $data['name'] ?? '';
        $entity->price = (float)($data['price'] ?? 0);
        $entity->category = $data['category'] ?? '未分类';
        $entity->description = $data['description'] ?? '';

        $result = DB::insert($entity);
        $entity->id = $result->insertId;

        return $entity->toArray();
    }

    public function update(int $id, array $data): ?array {
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

    public function delete(int $id): bool {
        if (!$this->db()->where("id = ?", $id)->first()) {
            return false;
        }
        $this->db()->where("id = ?", $id)->delete();
        return true;
    }

    public function search(string $keyword = '', string $category = ''): array {
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

    public function findByCategory(string $category): array {
        $entities = $this->db()->where("category = ?", $category)->orderBy("price DESC")->get();
        return QueryDemoService::entitiesToArray($entities);
    }

    public function findByPriceRange(float $minPrice, float $maxPrice): array {
        $entities = $this->db()
            ->where("price >= ? AND price <= ?", $minPrice, $maxPrice)
            ->orderBy("price ASC")
            ->get();
        return QueryDemoService::entitiesToArray($entities);
    }
}
