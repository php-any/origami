<?php

declare(strict_types=1);

namespace app\optional;

use RuntimeException;
use Throwable;

if (!defined('APP_ROOT')) {
    exit;
}

final class Cron
{
public static function plugin_cron_sync(array $plugin): void
{
    $plugin_id = (string)$plugin['id'];
    $names = [];
    foreach ((array)($plugin['cron'] ?? []) as $name => $task) {
        $name = (string)$name;
        $names[] = $name;
        $interval = self::cron_task_interval($plugin, $task);
        $callback = (string)$task['callback'];
        $enabled = plugin_enabled($plugin) ? 1 : 0;
        app_db_insert_ignore('app_cron_tasks', [
            'plugin_id' => $plugin_id,
            'task_name' => $name,
            'callback' => $callback,
            'interval_seconds' => $interval,
            'enabled' => $enabled,
            'last_error' => '',
        ], ['plugin_id', 'task_name']);
        q("UPDATE app_cron_tasks SET callback=?,interval_seconds=?,enabled=? WHERE plugin_id=? AND task_name=?", [
            $callback, $interval, $enabled, $plugin_id, $name,
        ]);
    }
    if (!$names) {
        q("DELETE FROM app_cron_tasks WHERE plugin_id=?", [$plugin_id]);
        return;
    }
    $marks = sql_marks(count($names));
    q("DELETE FROM app_cron_tasks WHERE plugin_id=? AND task_name NOT IN ($marks)", array_merge([$plugin_id], $names));
}

public static function cron_task_interval(array $plugin, array $task): int
{
    $interval = $task['interval'] ?? 0;
    if (is_string($interval) && !is_numeric($interval)) {
        plugin_load($plugin);
        $interval = plugin_callback_exists($interval) ? $interval($plugin, $task) : throw new RuntimeException('计划任务间隔函数不存在');
    }
    return min(31536000, max(60, (int)$interval));
}

public static function cron_log_start(string $plugin_id, string $task_name, int $started_at): int
{
    q("INSERT INTO app_cron_logs(plugin_id,task_name,status,message,started_at,finished_at) VALUES(?,?,?,?,?,0)", [$plugin_id, $task_name, 'running', '', $started_at]);
    return app_db_last_insert_id('app_cron_logs');
}

public static function cron_log_finish(int $id, string $status, string $message, int $finished_at): void
{
    if ($id <= 0) return;
    q("UPDATE app_cron_logs SET status=?,message=?,finished_at=? WHERE id=?", [$status, $message, $finished_at, $id]);
}

public static function cron_log_prune(): void
{
    $now = time();
    if ($now - (int)setting('cron_logs_pruned_at', '0') < 86400) return;
    q("DELETE FROM app_cron_logs WHERE started_at<?", [time() - CRON_LOG_RETENTION_SECONDS]);
    save_settings_values(['cron_logs_pruned_at' => (string)$now]);
}

public static function cron_task_claim(int $now): ?array
{
    for ($attempt = 0; $attempt < 3; $attempt++) {
        $claimed = tx(function () use ($now): ?array {
            $task = one("SELECT plugin_id,task_name,callback,interval_seconds,last_success_at,status,failure_count,retry_limit,pause_seconds,pause_until FROM app_cron_tasks WHERE enabled=1 AND available_at<=? ORDER BY available_at,plugin_id,task_name LIMIT 1", [$now]);
            if (!$task) return null;
            $token = bin2hex(random_bytes(16));
            $lease_until = $now + CRON_LEASE_SECONDS;
            $resuming = (int)$task['pause_until'] > 0 && (int)$task['pause_until'] <= $now;
            $updated = q("UPDATE app_cron_tasks SET lease_token=?,lease_until=?,available_at=?,last_started_at=?,status='running',attempts=attempts+1,failure_count=?,pause_until=0 WHERE plugin_id=? AND task_name=? AND enabled=1 AND available_at<=?", [
                $token, $lease_until, $lease_until, $now, $resuming ? 0 : (int)$task['failure_count'], (string)$task['plugin_id'], (string)$task['task_name'], $now,
            ])->rowCount();
            if ($updated !== 1) return [];
            if ((string)$task['status'] === 'running') {
                q("UPDATE app_cron_logs SET status='failed',message='任务租约超时',finished_at=? WHERE plugin_id=? AND task_name=? AND status='running' AND finished_at=0", [$now, (string)$task['plugin_id'], (string)$task['task_name']]);
            }
            $task['lease_token'] = $token;
            $task['last_started_at'] = $now;
            $task['touched_at'] = $now;
            if ($resuming) {
                $task['failure_count'] = 0;
                $task['pause_until'] = 0;
            }
            return $task;
        });
        if ($claimed === null) return null;
        if ($claimed) return $claimed;
    }
    return null;
}

public static function cron_lease_touch(): void
{
    $lease = $GLOBALS['__cron_active_lease'] ?? null;
    if (!is_array($lease)) return;
    $now = time();
    if ($now - (int)($lease['touched_at'] ?? 0) < 60) return;
    $until = $now + CRON_LEASE_SECONDS;
    $updated = q("UPDATE app_cron_tasks SET lease_until=?,available_at=? WHERE plugin_id=? AND task_name=? AND lease_token=?", [
        $until, $until, (string)$lease['plugin_id'], (string)$lease['task_name'], (string)$lease['lease_token'],
    ])->rowCount();
    if ($updated === 1) $GLOBALS['__cron_active_lease']['touched_at'] = $now;
}

public static function cron_task_finish(array $task, string $status, string $error, int $finished_at): bool
{
    $interval = max(60, (int)(val('SELECT interval_seconds FROM app_cron_tasks WHERE plugin_id=? AND task_name=?', [(string)$task['plugin_id'], (string)$task['task_name']]) ?: $task['interval_seconds']));
    $failure_count = (int)$task['failure_count'];
    $pause_until = 0;
    if ($status === 'success') {
        $failure_count = 0;
        $available_at = $finished_at + $interval;
    } else {
        $failure_count++;
        $available_at = $finished_at + min(300, $interval);
        if ($failure_count >= max(1, (int)$task['retry_limit'])) {
            $pause_until = $finished_at + max(60, (int)$task['pause_seconds']);
            $available_at = $pause_until;
        }
    }
    return q("UPDATE app_cron_tasks SET available_at=?,lease_token='',lease_until=0,last_finished_at=?,last_success_at=?,status=?,failure_count=?,pause_until=?,last_error=? WHERE plugin_id=? AND task_name=? AND lease_token=?", [
        $available_at,
        $finished_at,
        $status === 'success' ? $finished_at : (int)$task['last_success_at'],
        $status,
        $failure_count,
        $pause_until,
        $status === 'failed' ? cut($error, 500) : '',
        (string)$task['plugin_id'],
        (string)$task['task_name'],
        (string)$task['lease_token'],
    ])->rowCount() === 1;
}

public static function cron_run(): array
{
    $result = ['due' => 0, 'success' => 0, 'failed' => 0, 'tasks' => []];
    try {
        try { self::cron_log_prune(); } catch (Throwable $e) { debug_log_write('[cron] log prune failed', $e); }
        $plugins = plugins();
        while ($task = self::cron_task_claim(time())) {
            $plugin_id = (string)$task['plugin_id'];
            $task_name = (string)$task['task_name'];
            $key = $plugin_id . ':' . $task_name;
            $plugin = $plugins[$plugin_id] ?? null;
            $status = 'failed';
            $error = '';
            $message = '';
            $log_id = 0;
            $result['due']++;
            $GLOBALS['__cron_active_lease'] = $task;
            try { $log_id = self::cron_log_start($plugin_id, $task_name, (int)$task['last_started_at']); }
            catch (Throwable $e) { debug_log_write('[cron] ' . $key . ' log failed', $e); }
            try {
                if (!$plugin || !plugin_enabled($plugin)) throw new RuntimeException('插件未启用或不存在');
                plugin_load($plugin);
                $callback = (string)$task['callback'];
                if (!plugin_callback_exists($callback)) throw new RuntimeException('计划任务函数不存在');
                $definition = (array)($plugin['cron'][$task_name] ?? []);
                $definition['interval'] = (int)$task['interval_seconds'];
                $value = $callback($plugin, $definition);
                $message = is_scalar($value) ? trim((string)$value) : '';
                $status = 'success';
            } catch (Throwable $e) {
                $error = trim($e->getMessage());
                debug_log_write('[cron] ' . $key . ' failed', $e);
            } finally {
                $finished_at = time();
                try {
                    if (!self::cron_task_finish($task, $status, $error, $finished_at)) {
                        $status = 'failed';
                        $error = '任务租约已失效，运行结果未写入';
                    }
                } catch (Throwable $e) {
                    $status = 'failed';
                    $error = '任务状态写入失败：' . $e->getMessage();
                    debug_log_write('[cron] ' . $key . ' state failed', $e);
                }
                $result[$status === 'success' ? 'success' : 'failed']++;
                $result['tasks'][$key] = $status;
                try { self::cron_log_finish($log_id, $status, $status === 'failed' ? cut($error, 500) : $message, $finished_at); }
                catch (Throwable $e) { debug_log_write('[cron] ' . $key . ' log failed', $e); }
                unset($GLOBALS['__cron_active_lease']);
            }
        }
        return $result;
    } finally { unset($GLOBALS['__cron_active_lease']); }
}

public static function cron_route(): void
{
    set_time_limit(0);
    ignore_user_abort(true);
    header('Content-Type: text/plain; charset=utf-8');
    header('Cache-Control: no-store');
    $result = self::cron_run();
    echo 'cron finished: due ' . (int)$result['due'] . ', success ' . (int)$result['success'] . ', failed ' . (int)$result['failed'] . "\n";
    foreach ($result['tasks'] as $task => $status) echo $task . ': ' . $status . "\n";
}
}
