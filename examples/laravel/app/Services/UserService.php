<?php

namespace App\Services;

use App\Models\User;
use Container\Annotation\Singleton;
use Database\DB;

#[Singleton]
class UserService
{
    private function query(): DB
    {
        return DB::model(User::class);
    }

    public function all(): array
    {
        return $this->query()->orderBy('id ASC')->get();
    }

    public function find(int $id): ?User
    {
        return $this->query()->where('id = ?', $id)->first();
    }

    /**
     * @param int[] $ids
     * @return array<int, User>
     */
    public function findMany(array $ids): array
    {
        $ids = array_values(array_unique(array_filter($ids, fn ($id) => $id > 0)));
        if ($ids === []) {
            return [];
        }

        $placeholders = implode(', ', array_fill(0, count($ids), '?'));
        $users = $this->query()->where('id IN (' . $placeholders . ')', $ids)->get();

        $indexed = [];
        foreach ($users as $user) {
            $indexed[$user->id] = $user;
        }

        return $indexed;
    }

    public function findByEmail(string $email): ?User
    {
        return $this->query()->where('email = ?', $email)->first();
    }

    public function toArray(User $user): array
    {
        return [
            'id' => $user->id,
            'name' => $user->name,
            'email' => $user->email,
            'created_at' => $user->created_at,
        ];
    }
}
