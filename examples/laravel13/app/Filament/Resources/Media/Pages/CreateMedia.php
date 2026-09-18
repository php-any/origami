<?php

namespace App\Filament\Resources\Media\Pages;

use App\Filament\Resources\Media\MediaResource;
use App\Models\Media;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\Storage;
use Filament\Resources\Pages\CreateRecord;

class CreateMedia extends CreateRecord
{
    protected static string $resource = MediaResource::class;

    /**
     * @param  array<string, mixed>  $data
     */
    protected function handleRecordCreation(array $data): Model
    {
        $path = $data['upload'] ?? null;
        $disk = 'public';

        if (! is_string($path) || $path === '') {
            throw new \RuntimeException('请上传文件。');
        }

        $name = filled($data['name'] ?? null)
            ? (string) $data['name']
            : basename($path);

        return Media::query()->create([
            'name' => $name,
            'disk' => $disk,
            'path' => $path,
            'mime' => Storage::disk($disk)->mimeType($path) ?: null,
            'size' => Storage::disk($disk)->size($path) ?: 0,
            'uploaded_by' => auth('admin')->id(),
        ]);
    }
}
