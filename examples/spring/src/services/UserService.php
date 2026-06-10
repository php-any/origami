<?php

namespace Spring\Service;

use Container\Singleton;
use Database\DB;
use Spring\Model\Entity\UserEntity;

#[Singleton]
class UserService {

    private function db() {
        return DB::model(UserEntity::class);
    }

    public function findAll() {
        return $this->db()->orderBy("id ASC")->get();
    }

    public function findById($id) {
        return $this->db()->where("id = ?", $id)->first();
    }

    public function findByEmail($email) {
        return $this->db()->where("email = ?", $email)->first();
    }

    public function create($data) {
        $entity = new UserEntity();
        $entity->name = $data['name'] ?? '';
        $entity->email = $data['email'] ?? '';
        $entity->age = (int)($data['age'] ?? 0);

        $result = DB::insert($entity);
        $entity->id = $result->insertId;

        return $entity;
    }

    public function update($id, $data) {
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

    public function delete($id) {
        if (!$this->findById($id)) {
            return false;
        }
        $this->db()->where("id = ?", $id)->delete();
        return true;
    }

    public function search($keyword, $field = 'name') {
        if ($field === 'email') {
            return $this->db()->where("email LIKE ?", "%" . $keyword . "%")->get();
        }
        return $this->db()->where("name LIKE ?", "%" . $keyword . "%")->get();
    }
}
