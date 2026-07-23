<?php

namespace Bootstrap\Auth;

use App\Models\ApiToken;
use App\Models\User;
use Illuminate\Contracts\Auth\Authenticatable;
use Illuminate\Contracts\Auth\Guard;
use Illuminate\Auth\GuardHelpers;

/**
 * API Token Guard：从 Authorization 头解析 token，查 api_tokens → users。
 * （Net\Http\Request 与 Illuminate\Http\Request 不同，故用自定义 Guard）
 */
class ApiTokenGuard implements Guard
{
    use GuardHelpers;

    /** @var callable(): string */
    private $tokenResolver;

    public function __construct(callable $tokenResolver)
    {
        $this->tokenResolver = $tokenResolver;
    }

    public function user(): ?Authenticatable
    {
        if (!is_null($this->user)) {
            return $this->user;
        }

        $token = (string) call_user_func($this->tokenResolver);
        if ($token === '') {
            return $this->user = null;
        }

        $row = ApiToken::query()->where('token', $token)->first();
        if ($row === null) {
            return $this->user = null;
        }

        $user = User::query()->find((int) $row->user_id);
        if ($user === null) {
            return $this->user = null;
        }

        return $this->user = $user;
    }

    public function validate(array $credentials = []): bool
    {
        if (isset($credentials['token'])) {
            $prev = $this->user;
            $this->user = null;
            $resolver = $this->tokenResolver;
            $this->tokenResolver = static fn () => (string) $credentials['token'];
            $ok = $this->check();
            $this->tokenResolver = $resolver;
            $this->user = $prev;

            return $ok;
        }

        return $this->check();
    }
}
