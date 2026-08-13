<?php

namespace Spring\DTO\Request;

use Net\Http\UploadedFile;
use Validation\Annotation\NotBlank;
use Validation\Annotation\Size;

/**
 * POST /api/upload (multipart/form-data)
 * 字段：title=说明文字, file=文件
 */
class UploadFileRequest
{
    #[NotBlank(message: 'title 不能为空')]
    #[Size(min: 1, max: 200)]
    public string $title;

    #[NotBlank(message: '请上传文件')]
    public UploadedFile $file;
}
