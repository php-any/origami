<?php

namespace App\Http\Controllers\Telescope;

use Laravel\Telescope\EntryType;
use Laravel\Telescope\Storage\EntryQueryOptions;
use Laravel\Telescope\Watchers\LogWatcher;
use Laravel\Telescope\Watchers\QueryWatcher;
use Laravel\Telescope\Watchers\RequestWatcher;
use Net\Http\Request;
use Net\Http\Response;

/**
 * Telescope SPA API 代理（Origami Request/Response → EntriesRepository）。
 */
class ApiController
{
    private static function optionsFromRequest(Request $request): EntryQueryOptions
    {
        $opts = new EntryQueryOptions();
        // Origami Request::input 缺省键可能返回非 PHP null，避免 take 被当成 0
        $opts->limit = 50;
        try {
            $take = $request->input('take');
            if (is_numeric($take) && (int) $take > 0) {
                $opts->limit = (int) $take;
            }
        } catch (\Throwable $e) {
        }
        $batchId = $request->input('batch_id');
        if ($batchId !== null && $batchId !== '') {
            $opts->batchId = (string) $batchId;
        }
        $tag = $request->input('tag');
        if ($tag !== null && $tag !== '') {
            $opts->tag = (string) $tag;
        }
        $family = $request->input('family_hash');
        if ($family !== null && $family !== '') {
            $opts->familyHash = (string) $family;
        }
        $before = $request->input('before');
        if ($before !== null && $before !== '') {
            $opts->beforeSequence = $before;
        }

        return $opts;
    }

    private static function json(Response $response, array $payload): void
    {
        $flags = JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES;
        if (defined('JSON_INVALID_UTF8_SUBSTITUTE')) {
            $flags |= JSON_INVALID_UTF8_SUBSTITUTE;
        }
        $json = json_encode($payload, $flags);
        if ($json === false) {
            \Log::error('Telescope JSON encode failed: ' . json_last_error_msg());
            $json = json_encode([
                'entries' => [],
                'status' => 'enabled',
                'error' => 'json_encode failed',
            ], JSON_UNESCAPED_UNICODE);
        }
        $response->header('Content-Type', 'application/json; charset=utf-8');
        $response->write($json);
    }

    private static function listEntries(Response $response, string $type, string $watcher, Request $request): void
    {
        bootstrap_telescope_http();
        $repo = telescope_entries_repository();
        $entries = [];
        $opts = self::optionsFromRequest($request);
        $dbPath = '';
        $sqlCount = -1;
        try {
            $dbPath = (string) eloquent_capsule()->getConnection()->getDatabaseName();
            $sqlCount = (int) eloquent_capsule()->getConnection()->table('telescope_entries')->where('type', $type)->count();
        } catch (\Throwable $e) {
            $dbPath = $e->getMessage();
        }
        $raw = $repo->get($type, $opts);
        $rawCount = 0;
        foreach ($raw as $row) {
            $rawCount++;
            $entries[] = telescope_entry_to_array($row);
        }
        self::json($response, [
            'entries' => $entries,
            'status' => telescope_watcher_status($watcher),
            'debug' => [
                'type' => $type,
                'limit' => $opts->limit,
                'raw' => $rawCount,
                'n' => count($entries),
                'db' => $dbPath,
                'sql' => $sqlCount,
            ],
        ]);
    }

    private static function showEntry(Response $response, string $id): void
    {
        bootstrap_telescope_http();
        $repo = telescope_entries_repository();
        $entry = $repo->find($id);
        $batchOpts = EntryQueryOptions::forBatchId($entry->batchId);
        $batchOpts->limit = 100;
        $batch = [];
        foreach ($repo->get(null, $batchOpts) as $row) {
            $batch[] = telescope_entry_to_array($row);
        }
        self::json($response, [
            'entry' => telescope_entry_to_array($entry),
            'batch' => $batch,
        ]);
    }

    public function requestsIndex(Request $request, Response $response): void
    {
        self::listEntries($response, EntryType::REQUEST, RequestWatcher::class, $request);
    }

    public function requestsShow(Request $request, Response $response, string $id): void
    {
        self::showEntry($response, $id);
    }

    public function logsIndex(Request $request, Response $response): void
    {
        self::listEntries($response, EntryType::LOG, LogWatcher::class, $request);
    }

    public function logsShow(Request $request, Response $response, string $id): void
    {
        self::showEntry($response, $id);
    }

    public function queriesIndex(Request $request, Response $response): void
    {
        self::listEntries($response, EntryType::QUERY, QueryWatcher::class, $request);
    }

    public function queriesShow(Request $request, Response $response, string $id): void
    {
        self::showEntry($response, $id);
    }

    public function exceptionsIndex(Request $request, Response $response): void
    {
        self::listEntries($response, EntryType::EXCEPTION, \Laravel\Telescope\Watchers\ExceptionWatcher::class, $request);
    }

    public function exceptionsShow(Request $request, Response $response, string $id): void
    {
        self::showEntry($response, $id);
    }

    /** 其它类型：空列表，保证 SPA 侧栏可点开 */
    public function emptyIndex(Request $request, Response $response): void
    {
        self::json($response, ['entries' => [], 'status' => 'off']);
    }

    public function emptyShow(Request $request, Response $response, string $id): void
    {
        self::json($response, ['entry' => null, 'batch' => []]);
    }

    public function monitoredTags(Response $response): void
    {
        bootstrap_telescope_http();
        $tags = telescope_entries_repository()->monitoring();
        $list = [];
        foreach ($tags as $tag) {
            $list[] = $tag;
        }
        self::json($response, ['tags' => $list]);
    }

    public function monitoredTagsStore(Request $request, Response $response): void
    {
        bootstrap_telescope_http();
        $tag = (string) $request->input('tag');
        if ($tag !== '') {
            telescope_entries_repository()->monitor([$tag]);
        }
        self::json($response, ['ok' => true]);
    }

    public function monitoredTagsDelete(Request $request, Response $response): void
    {
        bootstrap_telescope_http();
        $tag = (string) $request->input('tag');
        if ($tag !== '') {
            telescope_entries_repository()->stopMonitoring([$tag]);
        }
        self::json($response, ['ok' => true]);
    }

    public function toggleRecording(Response $response): void
    {
        bootstrap_telescope_http();
        $cache = app('cache');
        if ($cache->get('telescope:pause-recording')) {
            $cache->forget('telescope:pause-recording');
        } else {
            $cache->put('telescope:pause-recording', true, now()->addDays(30));
        }
        self::json($response, ['recording' => !cache('telescope:pause-recording')]);
    }

    public function clearEntries(Response $response): void
    {
        bootstrap_telescope_http();
        telescope_entries_repository()->clear();
        self::json($response, ['ok' => true]);
    }
}
