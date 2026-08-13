<?php

namespace App\Services;

use App\Models\User;
use Container\Annotation\Singleton;

#[Singleton]
class UserService
{
    public function all(): array
    {
        return User::query()->orderBy('id')->get()->all();
    }

    public function find(int $id): ?User
    {
        return User::query()->find($id);
    }

    /**
     * @param int[] $ids
     * @return array<int, User>
     */
    public function findMany(array $ids): array
    {
        $ids = array_values(array_unique(array_filter($ids, static fn ($id) => $id > 0)));
        if ($ids === []) {
            return [];
        }

        return User::query()
            ->whereIn('id', $ids)
            ->get()
            ->keyBy('id')
            ->all();
    }

    public function findByEmail(string $email): ?User
    {
        return User::query()->where('email', $email)->first();
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
