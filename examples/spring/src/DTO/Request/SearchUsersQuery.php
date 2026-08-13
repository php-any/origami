<?php

namespace Spring\DTO\Request;

use Validation\Annotation\Size;

/**
 * GET /api/queries/users/search?keyword=张
 */
class SearchUsersQuery
{
    #[Size(max: 100)]
    public string $keyword = '';
}
