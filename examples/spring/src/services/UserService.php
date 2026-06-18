<?php

namespace Spring\Service;

use Container\Singleton;
use Database\DB;
use Spring\Model\Entity\UserEntity;

#[Singleton]
class UserService {

    private function db(): DB {
        return DB<UserEntity>();
    }

    public function findAll(): array {
        return $this->db()->orderBy("id ASC")->get();
    }

    public function findById(int $id): ?UserEntity {
        return $this->db()->where("id = ?", $id)->first();
    }

    public function findByEmail(string $email): ?UserEntity {
        return $this->db()->where("email = ?", $email)->first();
    }

    public function create(array $data): UserEntity {
        $entity = new UserEntity();
        $entity->name = $data['name'] ?? '';
        $entity->email = $data['email'] ?? '';
        $entity->age = (int)($data['age'] ?? 0);

        $result = DB::insert($entity);
        $entity->id = $result->insertId;

        return $entity;
    }

    public function update(int $id, array $data): ?UserEntity {
        $existing = $this->findById($id);
        if (!$existing) {
            return null;
        }

        $entity = new UserEntity();
        if (isset($data['name'])) {
            $entity->name = $data['name'];
        }
        if (isset($data['email'])) {
            $entity->email = $data['email'];
        }
        if (isset($data['age'])) {
            $entity->age = (int)$data['age'];
        }

        $this->db()->where("id = ?", $id)->update($entity);
        return $this->findById($id);
    }

    public function delete(int $id): bool {
        if (!$this->findById($id)) {
            return false;
        }
        $this->db()->where("id = ?", $id)->delete();
        return true;
    }

    public function search(string $keyword, string $field = 'name'): array {
        if ($field === 'email') {
            return $this->db()->where("email LIKE ?", "%" . $keyword . "%")->get();
        }
        return $this->db()->where("name LIKE ?", "%" . $keyword . "%")->get();
    }
}
