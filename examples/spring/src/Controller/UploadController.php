<?php

namespace Spring\Controller;

use Net\Annotation\Controller;
use Net\Annotation\PostMapping;
use Net\Annotation\Route;
use Net\Annotation\Middleware;
use Net\Http\Response;
use Net\Http\UploadedFile;
use Spring\DTO\Request\UploadFileRequest;
use Spring\Middleware\LogInterceptor;
use Validation\Annotation\NotBlank;

#[Middleware(LogInterceptor::class)]
#[Controller]
#[Route(prefix: "/api")]
class UploadController {

    private string $uploadDir;

    public function __construct() {
        $this->uploadDir = dirname(dirname(__DIR__)) . '/storage/uploads';
        if (!is_dir($this->uploadDir)) {
            mkdir($this->uploadDir, 0755, true);
        }
    }

    /**
     * multipart/form-data 上传（DTO 自动绑定 title + file）
     *
     * curl 示例：
     * curl -X POST http://localhost:8080/api/upload \
     *   -F "title=示例文档" \
     *   -F "file=@/path/to/hello.txt"
     */
    #[PostMapping(path: "/upload")]
    public function upload(UploadFileRequest $request, Response $response): void {
        if (!$request->file->isValid()) {
            $response->error('文件无效或为空', 400);
            return;
        }

        $savedPath = $request->file->store($this->uploadDir);

        $response->success([
            'title' => $request->title,
            'original_name' => $request->file->originalName(),
            'size' => $request->file->size(),
            'mime_type' => $request->file->mimeType(),
            'saved_path' => $savedPath,
        ], '上传成功', 201);
    }

    /**
     * 单文件形参自动注入（字段名须与形参名一致：avatar）
     *
     * curl 示例：
     * curl -X POST http://localhost:8080/api/upload/avatar \
     *   -F "avatar=@/path/to/photo.png"
     */
    #[PostMapping(path: "/upload/avatar")]
    public function uploadAvatar(
        #[NotBlank(message: '请上传 avatar 文件')]
        UploadedFile $avatar,
        Response $response
    ): void {
        if (!$avatar->isValid()) {
            $response->error('文件无效或为空', 400);
            return;
        }

        $savedPath = $avatar->store($this->uploadDir, 'avatar.png');

        $response->success([
            'original_name' => $avatar->originalName(),
            'size' => $avatar->size(),
            'mime_type' => $avatar->mimeType(),
            'saved_path' => $savedPath,
        ], '头像上传成功', 201);
    }
}
