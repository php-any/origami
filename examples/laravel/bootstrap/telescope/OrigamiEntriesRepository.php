<?php

namespace Bootstrap\Telescope;

use Illuminate\Support\Collection;
use Laravel\Telescope\EntryResult;
use Laravel\Telescope\EntryType;
use Laravel\Telescope\Storage\DatabaseEntriesRepository;
use Laravel\Telescope\Storage\EntryQueryOptions;

/**
 * Origami 兼容层：
 * - flatMap()->all() 可能返回空对象而非 []
 * - EntryModel 属性访问触发 Illuminate\Support\Carbon::format parent:: 死循环
 */
class OrigamiEntriesRepository extends DatabaseEntriesRepository
{
    public function store(Collection $entries)
    {
        if ($entries->isEmpty()) {
            return;
        }

        [$exceptions, $entries] = $entries->partition->isException();

        $this->storeExceptions($exceptions);

        $table = $this->table('telescope_entries');

        $entries->chunk($this->chunkSize)->each(function ($chunked) use ($table) {
            $rows = [];
            foreach ($chunked as $entry) {
                $content = $entry->content;
                if (!is_string($content)) {
                    $content = json_encode($content, defined('JSON_INVALID_UTF8_SUBSTITUTE') ? JSON_INVALID_UTF8_SUBSTITUTE : 0);
                }
                $entry->content = $content;
                $rows[] = $entry->toArray();
            }
            if (count($rows) === 0) {
                return;
            }
            $table->insert($rows);
        });

        $this->storeTags($entries->pluck('tags', 'uuid'));
    }

    public function find($id): EntryResult
    {
        $row = $this->table('telescope_entries')->where('uuid', $id)->first();
        if ($row === null) {
            throw (new \Illuminate\Database\Eloquent\ModelNotFoundException())->setModel(
                \Laravel\Telescope\Storage\EntryModel::class,
                [$id]
            );
        }

        $tags = $this->table('telescope_entries_tags')
            ->where('entry_uuid', $id)
            ->pluck('tag');
        $tagList = [];
        foreach ($tags as $tag) {
            $tagList[] = $tag;
        }

        return $this->rowToEntryResult($row, $tagList);
    }

    public function get($type, EntryQueryOptions $options)
    {
        $query = $this->table('telescope_entries')->orderByDesc('sequence');

        if ($type !== null && $type !== '') {
            $query->where('type', $type);
        }
        if (!empty($options->batchId)) {
            $query->where('batch_id', $options->batchId);
        }
        if (!empty($options->familyHash)) {
            $query->where('family_hash', $options->familyHash);
        }
        if (empty($options->batchId) && empty($options->tag) && empty($options->familyHash)) {
            $query->where('should_display_on_index', 1);
        }

        $limit = (int) ($options->limit ?? 50);
        $rows = $query->limit($limit)->get();

        $results = [];
        foreach ($rows as $row) {
            $result = $this->rowToEntryResult($row, []);
            if (is_array($result->content)) {
                $results[] = $result;
            }
        }

        return new Collection($results);
    }

    protected function storeTags($results)
    {
        $results->chunk($this->chunkSize)->each(function ($chunked) {
            $payload = [];
            foreach ($chunked as $uuid => $tags) {
                if ($tags === null) {
                    continue;
                }
                if (!is_array($tags)) {
                    if ($tags instanceof \Traversable) {
                        $tags = iterator_to_array($tags);
                    } else {
                        continue;
                    }
                }
                foreach ($tags as $tag) {
                    $payload[] = [
                        'entry_uuid' => (string) $uuid,
                        'tag' => $tag,
                    ];
                }
            }
            if (count($payload) === 0) {
                return;
            }
            try {
                $this->table('telescope_entries_tags')->insert($payload);
            } catch (\Illuminate\Database\UniqueConstraintViolationException $e) {
                // Ignore duplicate tags.
            }
        });
    }

    protected function rowToEntryResult($row, array $tags): EntryResult
    {
        $get = static function ($row, string $key) {
            if (is_array($row)) {
                return $row[$key] ?? null;
            }

            return $row->{$key} ?? null;
        };

        $content = $get($row, 'content');
        if (is_string($content)) {
            $decoded = json_decode($content, true);
            $content = is_array($decoded) ? $decoded : [];
        } elseif (!is_array($content)) {
            $content = [];
        }

        $createdAt = $get($row, 'created_at');
        try {
            $createdAt = \Carbon\Carbon::parse($createdAt ?: 'now');
        } catch (\Throwable $e) {
            $createdAt = \Carbon\Carbon::now();
        }

        return new EntryResult(
            $get($row, 'uuid'),
            $get($row, 'sequence'),
            $get($row, 'batch_id'),
            $get($row, 'type'),
            $get($row, 'family_hash'),
            $content,
            $createdAt,
            $tags
        );
    }
}
