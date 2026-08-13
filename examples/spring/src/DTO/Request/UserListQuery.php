<?php

namespace Spring\DTO\Request;

use Validation\Annotation\Max;
use Validation\Annotation\Min;

/**
 * GET /api/queries/simple/users?min_age=25&limit=5
 */
class UserListQuery
{
    #[Min(value: 0)]
    public int $min_age = 0;

    #[Min(value: 1)]
    #[Max(value: 100)]
    public int $limit = 10;
}
