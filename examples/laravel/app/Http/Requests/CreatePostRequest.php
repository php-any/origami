<?php

namespace App\Http\Requests;

use Validation\Annotation\NotBlank;
use Validation\Annotation\Size;

/**
 * 创建文章请求 DTO
 */
class CreatePostRequest
{
    #[NotBlank(message: 'title 不能为空')]
    #[Size(min: 1, max: 200)]
    public string $title = '';

    #[NotBlank(message: 'body 不能为空')]
    #[Size(min: 1, max: 5000)]
    public string $body = '';
}
