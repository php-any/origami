<?php

namespace Spring\DTO\Request;

use Validation\Annotation\Size;

/**
 * GET /api/products/search?keyword=手机&category=电子
 */
class SearchProductsQuery
{
    #[Size(max: 100)]
    public string $keyword = '';

    #[Size(max: 100)]
    public string $category = '';
}
