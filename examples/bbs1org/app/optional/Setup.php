<?php

declare(strict_types=1);

namespace app\optional;

use PDO;
use InvalidArgumentException;
use RuntimeException;
use Throwable;

if (!defined('APP_ROOT')) exit;

define('INSTALL_DATA_DIR', DATA_DIR);
define('INSTALL_DB_CONFIG_FILE', DB_CONFIG_FILE);
define('INSTALL_DEFAULT_DB_FILE', DATA_DIR . '/forum.sqlite');
define('UPDATE_DATA_DIR', DATA_DIR);
define('UPDATE_DB_CONFIG_FILE', DB_CONFIG_FILE);
define('UPDATE_INSTALL_LOCK_FILE', INSTALL_LOCK_FILE);
define('UPDATE_RUN_LOCK_FILE', DATA_DIR . '/update.lock');
define('UPDATE_REPOSITORY', 'bbs1org/bbs1org');
define('UPDATE_BRANCH', 'main');
define('UPDATE_MAX_ARCHIVE_BYTES', 52428800);
define('UPDATE_NOTICE_CHECK_INTERVAL', 21600);
define('UPDATE_PROTECTED_DIRS', ['app/data', 'app/plugins', 'app/avatars', 'app/upload', 'app/assets/plugins.css', 'app/assets/plugins.js', '.git']);

final class Setup
{
public static function debug_log_write(string $message, ?Throwable $e = null): void
{
    $exception_text = $e ? exception_detail($e) : '';
    $fingerprint = hash('sha256', $message . "\n" . $exception_text);
    $now = time();
    $throttle = @fopen(DEBUG_LOG_FILE . '.throttle', 'c+');
    if (is_resource($throttle) && @flock($throttle, LOCK_EX)) {
        rewind($throttle);
        $recent = json_decode((string)stream_get_contents($throttle), true);
        $recent = is_array($recent) ? $recent : [];
        foreach ($recent as $key => $timestamp) {
            if ($now - (int)$timestamp >= DEBUG_LOG_DEDUP_SECONDS) unset($recent[$key]);
        }
        if (isset($recent[$fingerprint])) {
            @flock($throttle, LOCK_UN);
            @fclose($throttle);
            return;
        }
        $recent[$fingerprint] = $now;
        @ftruncate($throttle, 0);
        rewind($throttle);
        @fwrite($throttle, json_encode($recent, JSON_UNESCAPED_SLASHES));
        @flock($throttle, LOCK_UN);
        @fclose($throttle);
    } elseif (is_resource($throttle)) {
        @fclose($throttle);
    }
    $line = '[' . date('Y-m-d H:i:s') . '] ' . trim($message);
    $ip = ip_addr();
    if ($ip !== '') $line .= "\nIP: " . $ip;
    $uri = trim((string)($_SERVER['REQUEST_METHOD'] ?? '') . ' ' . (string)($_SERVER['REQUEST_URI'] ?? ''));
    if ($uri !== '') $line .= "\n" . $uri;
    if ($exception_text !== '') $line .= "\n" . $exception_text;
    $line .= "\n\n";
    @file_put_contents(DEBUG_LOG_FILE, $line, FILE_APPEND | LOCK_EX);
}

public static function setup_html(string $title, string $body, bool $project_modal = false): never
{
    if (PHP_SAPI === 'cli') {
        $text = preg_replace('#<(?:style|script)\b[^>]*>.*?</(?:style|script)>#is', '', $body) ?? $body;
        $text = preg_replace('#<\s*/?\s*(?:br|p|div|section|article|header|footer|h[1-6]|li|tr|ul|ol)\b[^>]*>#i', "\n", $text) ?? $text;
        $text = html_entity_decode(strip_tags($text), ENT_QUOTES | ENT_HTML5, 'UTF-8');
        $text = preg_replace('/[ \t]+\n/u', "\n", $text) ?? $text;
        $text = preg_replace('/\n{3,}/', "\n\n", $text) ?? $text;
        $text = trim($text);
        $output = $text === '' || $text === $title || str_starts_with($text, $title . "\n") ? $text : $title . "\n\n" . $text;
        fwrite(STDOUT, trim($output) . "\n");
        exit;
    }
    $meta = '<meta name="viewport" content="width=device-width,initial-scale=1"><meta charset="utf-8">';
    $project_assets = $project_modal ? '<link rel="stylesheet" href="' . h(app_url('app/assets/index.css')) . '?v=' . h(APP_VERSION) . '">' : '';
    $project_scripts = $project_modal ? \project_modal_html() . '<script src="' . h(app_url('app/assets/index.js')) . '?v=' . h(APP_VERSION) . '" defer></script>' : '';
    echo '<!doctype html><html lang="zh-CN"><head>' . $meta . '<title>' . h($title) . '</title><link rel="icon" type="image/svg+xml" href="app/assets/index.svg">' . $project_assets . '<style>
    :root{--bg:#eef2f7;--panel:#fff;--line:#dfe6ee;--line2:#edf1f5;--text:#1f2937;--muted:#6b7280;--brand:#2563eb;--brand2:#1d4ed8;--ok:#059669;--warn:#b45309;--danger:#dc2626;--line-soft:var(--line2);--text-muted:var(--muted);--text-subtle:var(--muted);--text-disabled:var(--muted);--brand-hover:var(--brand2);--brand-soft:#eff6ff;--inverse:#111827;--inverse-text:#fff;--color-dark-rgb:31,41,55;--backdrop:rgba(var(--color-dark-rgb),.42);--shadow-medium:rgba(15,23,42,.18);--font-size-sm:12px;--font-size-md:14px;--font-size-lg:15px;--radius:10px;--radius-sm:8px;--focus-ring:rgba(37,99,235,.2)}
    *{box-sizing:border-box}body{margin:0;color:var(--text);font:14px/1.6 -apple-system,BlinkMacSystemFont,"PingFang SC","Microsoft YaHei",sans-serif}
    a{color:var(--brand);text-decoration:none}a:hover{color:var(--brand2)}.wrap{max-width:1060px;margin:0 auto;padding:24px 16px 40px}
    .hero{display:grid;gap:8px;margin-bottom:18px}.hero h1{margin:0;font-size:28px;line-height:1.2}.hero p{margin:0;color:var(--muted)}
    .grid{display:grid;grid-template-columns:minmax(0,1.15fr) minmax(0,.85fr);gap:16px;align-items:start}.card{background:var(--panel);border:1px solid var(--line);border-radius:var(--radius);box-shadow:0 10px 24px rgba(15,23,42,.05)}
    .card .hd{padding:16px 18px;border-bottom:1px solid var(--line2)}.card .hd h2{margin:0;font-size:16px}.card .bd{padding:16px 18px}
    .note{margin:10px 0;padding:12px 14px;border:1px solid #dbeafe;background:#eff6ff;color:#1e3a8a;border-radius:8px}.warn{border-color:#fde68a;background:#fffbeb;color:#92400e}.ok{border-color:#bbf7d0;background:#f0fdf4;color:#166534}
    .form{display:grid;gap:12px}.row{display:grid;gap:6px}.row label{font-size:12px;color:var(--muted)}.row small{color:var(--muted);font-size:11px;line-height:1.4}.row.compact{grid-template-columns:minmax(0,1fr) minmax(0,1fr);gap:10px}.row.compact .field{display:grid;gap:6px}.db-fields[hidden]{display:none}input[type=text],input[type=password],input[type=email],select,textarea{width:100%;border:1px solid #d6dbe3;border-radius:8px;padding:10px 12px;font:inherit;background:#fff;color:var(--text)}textarea{min-height:128px;resize:vertical}
    input:focus,select:focus,textarea:focus{outline:0;border-color:#93c5fd;box-shadow:0 0 0 3px rgba(59,130,246,.12)}.db-fields{display:grid;gap:12px;padding:12px;border:1px solid var(--line2);border-radius:8px;background:#fafcff}.checks{display:grid;gap:10px}.check{display:flex;gap:10px;align-items:flex-start;padding:10px 12px;border:1px solid var(--line2);border-radius:8px;background:#fafcff}.check input{margin-top:3px}
    .actions{display:flex;gap:10px;align-items:center;justify-content:flex-end}.btn{display:inline-flex;align-items:center;justify-content:center;min-height:38px;padding:0 16px;border:0;border-radius:8px;background:var(--brand);color:#fff;cursor:pointer;font:inherit;font-weight:600}.btn:hover{background:var(--brand2);color:#fff}.btn.alt{background:#fff;color:#374151;border:1px solid #d1d5db}.btn.alt:hover{background:#f8fafc;color:#111;border-color:#cbd5e1}
    .list{margin:0;padding-left:18px;color:#374151}.mono{font-family:ui-monospace,SFMono-Regular,Menlo,Consolas,monospace;overflow-wrap:anywhere;word-break:break-word}.kv{display:grid;width:100%;min-width:0;grid-template-columns:120px minmax(0,1fr);gap:8px 12px;font-size:13px}.kv div{min-width:0;max-width:100%;overflow-wrap:anywhere;word-break:break-word}.kv div:nth-child(odd){color:var(--muted)}.admin-pass{padding:14px;border:1px solid #fecaca;background:#fff1f2;color:#991b1b;border-radius:8px;word-break:break-all}.footer{margin-top:16px;color:var(--muted);font-size:12px;text-align:center}
    @media (max-width:860px){.grid{grid-template-columns:1fr}.hero h1{font-size:24px}.wrap{padding:18px 12px 30px}}
    </style></head><body><main class="wrap">' . $body . '</main>' . $project_scripts . '</body></html>';
    exit;
}

public static function app_db_config_source(array $config): string
{
    $saved = $config;
    unset($saved['path']);
    if (($saved['driver'] ?? '') === 'sqlite') $saved['db_file'] = $saved['database'];
    return "<?php\nif (!defined('APP_ROOT')) exit;\nreturn " . var_export($saved, true) . ";\n";
}

public static function app_db_schema(string $driver): array
{
    $types = app_db_types($driver);
    $id = $types['id'];
    $uint = $types['uint'];
    $short = $types['string'];
    $key = $types['key'];
    $long = $types['text'];
    $tables = [
        'app_groups' => "CREATE TABLE app_groups(id $id,name $key NOT NULL UNIQUE,allow_manage INTEGER NOT NULL DEFAULT 0,allow_admin INTEGER NOT NULL DEFAULT 0)",
        'app_users' => "CREATE TABLE app_users(id $id,username $key NOT NULL UNIQUE,password $short NOT NULL,email $short NOT NULL DEFAULT '',bio $long NOT NULL,avatar_style $short NOT NULL DEFAULT '',avatar_seed $short NOT NULL DEFAULT '',group_id $uint NOT NULL DEFAULT 2,points INTEGER NOT NULL DEFAULT 0,is_banned INTEGER NOT NULL DEFAULT 0,is_muted INTEGER NOT NULL DEFAULT 0,unread_notifications $uint NOT NULL DEFAULT 0,last_post_at $uint NOT NULL DEFAULT 0,created_at $uint NOT NULL)",
        'app_notifications' => "CREATE TABLE app_notifications(id $id,recipient_id $uint NOT NULL,sender_id $uint DEFAULT NULL,kind $short NOT NULL DEFAULT 'direct',content $long NOT NULL,topic_id $uint DEFAULT NULL,reply_id $uint DEFAULT NULL,read_at $uint NOT NULL DEFAULT 0,created_at $uint NOT NULL)",
        'app_forums' => "CREATE TABLE app_forums(id $id,name $short NOT NULL,description $long NOT NULL,sort $uint NOT NULL DEFAULT 0,allow_view_groups $short NOT NULL DEFAULT '',allow_post_groups $short NOT NULL DEFAULT '',allow_reply_groups $short NOT NULL DEFAULT '')",
        'app_topics' => "CREATE TABLE app_topics(id $id,forum_id $uint NOT NULL,user_id $uint NOT NULL,title $short NOT NULL,body $long NOT NULL,highlight_style $short NOT NULL DEFAULT '',reply_order INTEGER NOT NULL DEFAULT 0,reply_count $uint NOT NULL DEFAULT 0,view_count $uint NOT NULL DEFAULT 0,last_reply_at $uint NOT NULL DEFAULT 0,last_reply_user_id $uint NOT NULL DEFAULT 0,created_at $uint NOT NULL)",
        'app_replies' => "CREATE TABLE app_replies(id $id,topic_id $uint NOT NULL,user_id $uint NOT NULL,body $long NOT NULL,created_at $uint NOT NULL,updated_at $uint NOT NULL)",
        'app_attachments' => "CREATE TABLE app_attachments(id $id,user_id $uint NOT NULL,hash $short NOT NULL,file_name $short NOT NULL,original_name $short NOT NULL DEFAULT '',ext $short NOT NULL DEFAULT '',mime $short NOT NULL DEFAULT '',size $uint NOT NULL DEFAULT 0,is_image INTEGER NOT NULL DEFAULT 0,created_at $uint NOT NULL)",
        'app_cron_logs' => "CREATE TABLE app_cron_logs(id $id,plugin_id $short NOT NULL,task_name $short NOT NULL,status $short NOT NULL,message $long NOT NULL,started_at $uint NOT NULL,finished_at $uint NOT NULL DEFAULT 0)",
        'app_plugins' => "CREATE TABLE app_plugins(id $key PRIMARY KEY,name $short NOT NULL,version $short NOT NULL DEFAULT '',file $short NOT NULL,code_hash $short NOT NULL,manifest_json $long NOT NULL,config_json $long NOT NULL,entries_json $long NOT NULL,enabled $uint NOT NULL DEFAULT 0,status $short NOT NULL DEFAULT '',disabled_reason $long NOT NULL,installed_at $uint NOT NULL,updated_at $uint NOT NULL)",
        'app_cron_tasks' => "CREATE TABLE app_cron_tasks(plugin_id $key NOT NULL,task_name $key NOT NULL,callback $short NOT NULL,interval_seconds $uint NOT NULL,enabled $uint NOT NULL DEFAULT 1,available_at $uint NOT NULL DEFAULT 0,lease_token $short NOT NULL DEFAULT '',lease_until $uint NOT NULL DEFAULT 0,last_started_at $uint NOT NULL DEFAULT 0,last_finished_at $uint NOT NULL DEFAULT 0,last_success_at $uint NOT NULL DEFAULT 0,status $short NOT NULL DEFAULT '',attempts $uint NOT NULL DEFAULT 0,failure_count $uint NOT NULL DEFAULT 0,retry_limit $uint NOT NULL DEFAULT 3,pause_seconds $uint NOT NULL DEFAULT 1800,pause_until $uint NOT NULL DEFAULT 0,last_error $long NOT NULL,PRIMARY KEY(plugin_id,task_name))",
        'app_settings' => "CREATE TABLE app_settings(name $key PRIMARY KEY,value $long NOT NULL)",
    ];
    if ($driver === 'mysql') {
        foreach ($tables as &$sql) $sql .= ' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci';
        unset($sql);
    }
    $indexes = [
        'idx_users_group' => 'app_users(group_id)', 'idx_forums_sort' => 'app_forums(sort,id)',
        'idx_replies_topic_time' => 'app_replies(topic_id,created_at,id)',
        'idx_replies_user_time' => 'app_replies(user_id,created_at DESC,id DESC)',
        'idx_attachments_user' => 'app_attachments(user_id,created_at DESC,id DESC)',
        'idx_notifications_recipient_unread' => 'app_notifications(recipient_id,read_at)',
        'idx_notifications_recipient_time' => 'app_notifications(recipient_id,created_at DESC,id DESC)',
        'idx_cron_logs_plugin_time' => 'app_cron_logs(plugin_id,started_at DESC,id DESC)',
        'idx_cron_logs_time' => 'app_cron_logs(started_at DESC,id DESC)',
        'idx_cron_tasks_due' => 'app_cron_tasks(enabled,available_at,plugin_id,task_name)',
        'idx_topics_created' => 'app_topics(created_at DESC,id DESC)', 'idx_topics_last_reply' => 'app_topics(last_reply_at DESC,id DESC)',
        'idx_topics_user_created' => 'app_topics(user_id,created_at DESC,id DESC)', 'idx_topics_forum_created' => 'app_topics(forum_id,created_at DESC,id DESC)',
        'idx_topics_forum_last_reply' => 'app_topics(forum_id,last_reply_at DESC,id DESC)',
    ];
    foreach ($indexes as $name => &$target) $target = 'CREATE INDEX ' . $name . ' ON ' . $target;
    unset($target);
    $indexes['idx_attachments_user_hash'] = 'CREATE UNIQUE INDEX idx_attachments_user_hash ON app_attachments(user_id,hash)';
    if ($driver === 'mysql') {
        $indexes['idx_topics_search_title'] = 'CREATE FULLTEXT INDEX idx_topics_search_title ON app_topics(title) WITH PARSER ngram';
        $indexes['idx_topics_search_body'] = 'CREATE FULLTEXT INDEX idx_topics_search_body ON app_topics(body) WITH PARSER ngram';
        $indexes['idx_replies_search_body'] = 'CREATE FULLTEXT INDEX idx_replies_search_body ON app_replies(body) WITH PARSER ngram';
    } elseif ($driver === 'pgsql') {
        $indexes['idx_topics_search_title'] = 'CREATE INDEX idx_topics_search_title ON app_topics USING gin (title gin_trgm_ops)';
        $indexes['idx_topics_search_body'] = 'CREATE INDEX idx_topics_search_body ON app_topics USING gin (body gin_trgm_ops)';
        $indexes['idx_replies_search_body'] = 'CREATE INDEX idx_replies_search_body ON app_replies USING gin (body gin_trgm_ops)';
    }
    return [$tables, $indexes];
}

public static function app_db_prepare_search(PDO $db, string $driver): void
{
    if ($driver === 'pgsql') $db->exec('CREATE EXTENSION IF NOT EXISTS pg_trgm');
}

public static function app_db_create_schema_index(PDO $db, string $driver, string $index, string $sql): bool
{
    try {
        $db->exec($sql);
        return true;
    } catch (Throwable $e) {
        $optional = array_keys(mysql_search_index_definitions());
        if ($driver === 'mysql' && in_array($index, $optional, true)) return false;
        throw $e;
    }
}

public static function app_db_search_fallback_notice(): string
{
    return '当前 MySQL 环境无法创建 ngram 全文索引，已自动跳过；站内搜索将使用 LIKE，功能可用，但数据量较大时速度可能较慢。';
}

public static function app_db_index_table(string $sql): string
{
    return preg_match('/\bON\s+[`"]?([A-Za-z_][A-Za-z0-9_]*)/i', $sql, $match) ? $match[1] : '';
}

public static function i_db_name(): string
{
    if (is_file(INSTALL_DB_CONFIG_FILE)) {
        $config = include INSTALL_DB_CONFIG_FILE;
        $name = is_array($config) ? basename((string)($config['db_file'] ?? '')) : '';
        if ($name !== '' && preg_match('/^[A-Za-z0-9][A-Za-z0-9._-]*\.sqlite$/', $name)) return $name;
    }
    if (is_file(INSTALL_DEFAULT_DB_FILE)) return basename(INSTALL_DEFAULT_DB_FILE);
    return 'forum-' . bin2hex(random_bytes(8)) . '.sqlite';
}
public static function i_save_db_config(array $config): void
{
    if (!is_dir(INSTALL_DATA_DIR)) mkdir(INSTALL_DATA_DIR, 0755, true);
    if (file_put_contents(INSTALL_DB_CONFIG_FILE, self::app_db_config_source($config), LOCK_EX) === false) self::i_install_error('安装失败', '数据库配置文件写入失败。');
}

public static function i_require_writable_dirs(): void
{
    $dirs = [
        INSTALL_DATA_DIR => 'app/data/',
    ];
    foreach ($dirs as $dir => $label) {
        if (!is_dir($dir) && !mkdir($dir, 0755, true) && !is_dir($dir)) {
            self::i_install_error('安装环境检查未通过', $label . '目录无法创建，请检查目录权限。');
        }
        if (!is_writable($dir)) {
            self::i_install_error('安装环境检查未通过', $label . '目录不可写，请赋予 PHP 写入权限。');
        }
        $probe = tempnam($dir, '.install-');
        if ($probe === false || file_put_contents($probe, '1', LOCK_EX) === false) {
            if ($probe !== false) @unlink($probe);
            self::i_install_error('安装环境检查未通过', $label . '目录无法写入文件，请检查目录权限。');
        }
        @unlink($probe);
    }
}

public static function i_install_error(string $title, string $message): void
{
    if (PHP_SAPI === 'cli') {
        fwrite(STDERR, $title . PHP_EOL . $message . PHP_EOL);
        exit(1);
    }
    self::setup_html($title, '<div class="hero"><h1>' . h($title) . '</h1><p>安装环境检查未通过。</p></div><div class="card"><div class="bd"><div class="note warn">' . h($message) . '</div></div></div>');
}
public static function i_db(array $config): PDO
{
    try {
        return app_db_connect($config);
    } catch (Throwable $e) {
        self::i_install_error('数据库初始化失败', '数据库连接失败：' . $e->getMessage());
    }
}
public static function i_database_install_state(PDO $db, string $driver): string
{
    if (!app_db_table_exists($db, $driver, 'app_users') || $db->query('SELECT id FROM app_users ORDER BY id LIMIT 1')->fetchColumn() === false) return 'empty';
    foreach (['app_settings', 'app_groups', 'app_forums', 'app_topics', 'app_replies'] as $table) if (!app_db_table_exists($db, $driver, $table)) return 'partial';
    $site = $db->query("SELECT value FROM app_settings WHERE name='site_name' LIMIT 1")->fetchColumn();
    return $site === false ? 'partial' : 'installed';
}
public static function i_restore_existing_install(array $config): never
{
    self::i_save_db_config($config);
    if (file_put_contents(INSTALL_LOCK_FILE, (string)now(), LOCK_EX) === false) self::i_install_error('安装失败', '安装锁文件写入失败。');
    header('Location: index.php?a=login', true, 303);
    exit;
}
public static function i_result(string $title, string $admin_user, string $admin_pass, string $admin_email, string $site_name, string $database, string $notice = ''): void
{
    $warning = $notice === '' ? '' : '<div class="note warn">' . h($notice) . '</div>';
    self::setup_html($title, '<div class="hero"><h1>安装完成</h1><p>站点已初始化，管理员账号已创建。</p></div><div class="grid"><section class="card"><div class="hd"><h2>安装结果</h2></div><div class="bd"><div class="note ok">可以直接进入论坛使用，建议立即登录后台修改密码。</div>' . $warning . '<div style="height:12px"></div><div class="kv"><div>站点名</div><div>' . h($site_name) . '</div><div>数据库</div><div class="mono">' . h($database) . '</div><div>管理员用户名</div><div class="mono">' . h($admin_user) . '</div><div>管理员邮箱</div><div class="mono">' . h($admin_email) . '</div><div>管理员密码</div><div class="admin-pass mono">' . h($admin_pass) . '</div></div><div style="height:14px"></div><div class="actions"><a class="btn alt" href="index.php">进入首页</a><a class="btn" href="index.php?a=admin">进入后台</a></div></div></section><aside class="card"><div class="hd"><h2>已完成内容</h2></div><div class="bd"><ul class="list"><li>创建数据库结构和可用索引</li><li>创建默认版块</li><li>创建第一个管理员</li><li>生成缓存文件</li><li>数据库密码仅保存在 app/data/db.php</li></ul></div></aside></div><div class="footer">请立即保存本页显示的管理员密码，离开后无法再次查看。</div>');
}
public static function i_locked(): void
{
    self::setup_html('安装已锁定', '<div class="hero"><h1>安装已锁定</h1><p>安装入口当前不可访问。</p></div><div class="card"><div class="bd"><div class="note warn">如需重新安装，请先删除安装锁文件后再访问。</div><div style="height:14px"></div><div class="actions"><a class="btn" href="index.php">进入首页</a></div></div></div>');
}
public static function i_form(string $site_name, string $admin_user, string $admin_email, string $admin_pass, string $default_forum, array $values = []): void
{
    $type = in_array((string)($values['db_type'] ?? 'sqlite'), ['sqlite', 'mysql', 'pgsql'], true) ? (string)($values['db_type'] ?? 'sqlite') : 'sqlite';
    $option = fn(string $value, string $label): string => '<option value="' . $value . '"' . ($type === $value ? ' selected' : '') . '>' . $label . '</option>';
    $v = fn(string $name, string $default = ''): string => h((string)($values[$name] ?? $default));
    $default_host = $type === 'pgsql' ? 'postgres' : 'mysql';
    $default_port = $type === 'pgsql' ? '5432' : '3306';
    $db_fields = '<div class="db-fields" id="server-db-fields"' . ($type === 'sqlite' ? ' hidden' : '') . '><div class="row compact"><div class="field"><label>数据库地址</label><input type="text" name="db_host" value="' . $v('db_host', $default_host) . '"></div><div class="field"><label>端口</label><input type="text" name="db_port" value="' . $v('db_port', $default_port) . '"></div></div><div class="row"><label>数据库名</label><input type="text" name="db_name" value="' . $v('db_name') . '"><small>数据库需要提前创建，安装器会创建其中的数据表。</small></div><div class="row compact"><div class="field"><label>数据库用户</label><input type="text" name="db_user" value="' . $v('db_user') . '"></div><div class="field"><label>数据库密码</label><input type="password" name="db_password" value="' . $v('db_password') . '"></div></div></div>';
    $body = '<div class="hero"><h1>安装</h1><p>一页完成初始化，创建管理员和默认版块。</p></div><div class="grid"><section class="card"><div class="hd"><h2>安装配置</h2></div><div class="bd"><form class="form" method="post"><input type="hidden" name="step" value="install"><div class="row"><label>数据库类型</label><select name="db_type" id="db-type">' . $option('sqlite', 'SQLite（默认）') . $option('mysql', 'MySQL') . $option('pgsql', 'PostgreSQL') . '</select></div>' . $db_fields . '<div class="row"><label>站点名称</label><input type="text" name="site_name" value="' . h($site_name) . '" required></div><div class="row"><label>管理员用户名</label><input type="text" name="admin_username" value="' . h($admin_user) . '" required></div><div class="row"><label>管理员邮箱</label><input type="email" name="admin_email" value="' . h($admin_email) . '" required><small>用于找回密码与通知。</small></div><div class="row"><label>管理员密码</label><input type="password" name="admin_password" value="' . h($admin_pass) . '" required></div><div class="row"><label>确认管理员密码</label><input type="password" name="admin_password2" value="' . h($admin_pass) . '" required></div><div class="row"><label>默认版块名称</label><input type="text" name="forum_name" value="' . h($default_forum) . '" required></div><div class="checks"><label class="check"><input type="checkbox" name="confirm_clean" value="1" required><span>我确认这是全新安装，数据将被清理。</span></label><label class="check"><input type="checkbox" name="confirm_admin" value="1" required><span>我确认需要手工设置第一个管理员密码。</span></label></div><div class="actions"><button class="btn" type="submit">开始安装</button></div></form></div></section><aside class="card"><div class="hd"><h2>安装说明</h2></div><div class="bd"><ul class="list"><li>SQLite 无需填写连接信息</li><li>MySQL/PostgreSQL 数据库需提前创建</li><li>第一个管理员将拥有全部权限</li><li>管理员邮箱可用于找回密码</li></ul></div></aside></div><script>const type=document.getElementById("db-type"),fields=document.getElementById("server-db-fields"),host=fields.querySelector("[name=db_host]"),port=fields.querySelector("[name=db_port]"),defaults={mysql:["mysql","3306"],pgsql:["postgres","5432"]};function toggleDb(change){const values=defaults[type.value]||null;fields.hidden=!values;if(change&&values){host.value=values[0];port.value=values[1]}}type.addEventListener("change",()=>toggleDb(true));addEventListener("pageshow",()=>toggleDb(false));toggleDb(false);</script>';
    self::setup_html('安装', $body);
}

public static function setup_install_run(): never
{
    if (is_file(INSTALL_LOCK_FILE)) {
        self::i_locked();
    }
    self::i_require_writable_dirs();
    $step = (string)($_POST['step'] ?? '');
    if ($step !== 'install') {
        self::i_form('我的论坛', 'admin', '', '', '默认版块');
    }
    $form_values = $_POST;
    if (!isset($_POST['confirm_clean'], $_POST['confirm_admin'])) self::i_form('我的论坛', 'admin', '', '', '默认版块', $form_values);
    $driver = in_array((string)($_POST['db_type'] ?? 'sqlite'), ['sqlite', 'mysql', 'pgsql'], true) ? (string)($_POST['db_type'] ?? 'sqlite') : 'sqlite';
    $db_name = trim((string)($_POST['db_name'] ?? ''));
    $sqlite_name = $driver === 'sqlite' ? self::i_db_name() : '';
    $config = $driver === 'sqlite' ? [
        'driver' => 'sqlite', 'database' => $sqlite_name, 'path' => INSTALL_DATA_DIR . '/' . $sqlite_name,
    ] : [
        'driver' => $driver,
        'host' => trim((string)($_POST['db_host'] ?? '127.0.0.1')),
        'port' => max(1, (int)($_POST['db_port'] ?? ($driver === 'mysql' ? 3306 : 5432))),
        'database' => $db_name,
        'username' => (string)($_POST['db_user'] ?? ''),
        'password' => (string)($_POST['db_password'] ?? ''),
    ];
    if ($driver !== 'sqlite' && ($config['host'] === '' || $config['database'] === '' || $config['username'] === '')) self::i_form('我的论坛', 'admin', '', '', '默认版块', $form_values);
    $site_name = trim((string)($_POST['site_name'] ?? '我的论坛'));
    $admin_username = trim((string)($_POST['admin_username'] ?? 'admin'));
    $admin_email = trim((string)($_POST['admin_email'] ?? ''));
    $admin_password = (string)($_POST['admin_password'] ?? '');
    $admin_password2 = (string)($_POST['admin_password2'] ?? '');
    $forum_name = trim((string)($_POST['forum_name'] ?? '默认版块'));
    if ($site_name === '' || $admin_username === '' || $admin_email === '' || $admin_password === '' || $forum_name === '') self::i_form($site_name ?: '我的论坛', $admin_username ?: 'admin', $admin_email, $admin_password, $forum_name ?: '默认版块', $form_values);
    if ($admin_password !== $admin_password2) self::i_form($site_name, $admin_username, $admin_email, $admin_password, $forum_name, $form_values);
    if (is_file(INSTALL_LOCK_FILE)) self::i_locked();
    $db = self::i_db($config);
    $install_state = self::i_database_install_state($db, $driver);
    if ($install_state === 'installed') self::i_restore_existing_install($config);
    if ($install_state === 'partial') self::i_install_error('检测到已有数据', '数据库中已有用户数据，但安装记录不完整。请恢复原程序文件或使用空数据库安装。');
    self::i_save_db_config($config);
    [$tables, $indexes] = self::app_db_schema($driver);
    foreach ($tables as $table => $sql) if (!app_db_table_exists($db, $driver, $table)) $db->exec($sql);
    if ($driver === 'sqlite') {
        app_db_create_fts5_table($db, 'app_topics_fts', 'title, body');
        app_db_create_fts5_table($db, 'app_replies_fts', 'body');
    }
    self::app_db_prepare_search($db, $driver);
    $search_fallback = false;
    foreach ($indexes as $index => $sql) {
        if (app_db_index_exists($db, $driver, $index, self::app_db_index_table($sql))) continue;
        if (!self::app_db_create_schema_index($db, $driver, $index, $sql)) $search_fallback = true;
    }
    $seed = $db->prepare(app_db_upsert_sql($driver, 'app_groups', ['id', 'name', 'allow_manage', 'allow_admin'], ['id']));
    $seed->execute([1, '管理员', 1, 1]); $seed->execute([2, '会员', 0, 0]);
    $seed = $db->prepare(app_db_upsert_sql($driver, 'app_forums', ['id', 'name', 'description', 'sort'], ['id']));
    $seed->execute([1, $forum_name, '欢迎发帖', 0]);
    if ($driver === 'pgsql') {
        $db->exec("SELECT setval(pg_get_serial_sequence('app_groups','id'), (SELECT MAX(id) FROM app_groups))");
        $db->exec("SELECT setval(pg_get_serial_sequence('app_forums','id'), (SELECT MAX(id) FROM app_forums))");
    }
    $settings = array_merge(default_settings(), mysql_search_index_settings($db, $driver));
    $settings['site_name'] = $site_name;
    $stmt = $db->prepare(app_db_upsert_sql($driver, 'app_settings', ['name', 'value'], ['name']));
    foreach ($settings as $name => $value) $stmt->execute([$name, $value]);
    $admin_pass = $admin_password;
    $welcome_ts = now();
    $db->prepare("INSERT INTO app_users(username,password,email,bio,avatar_style,avatar_seed,group_id,last_post_at,created_at) VALUES(?,?,?,?,?,?,?,?,?)")->execute([$admin_username, password_hash($admin_pass, PASSWORD_DEFAULT), $admin_email, '站点管理员', '', '', 1, $welcome_ts, $welcome_ts]);
    forums_cache(true);
    groups_cache(true);
    home_stats_record_insert('users', 1);
    Plugin::plugin_registry_sync();
    Plugin::plugin_assets_rebuild();
    if (file_put_contents(INSTALL_LOCK_FILE, (string)now(), LOCK_EX) === false) self::i_install_error('安装失败', '安装锁文件写入失败。');
    $database_label = $driver === 'sqlite' ? 'app/data/' . $config['database'] : strtoupper($driver === 'pgsql' ? 'PostgreSQL' : 'MySQL') . ' / ' . $config['database'];
    self::i_result('安装完成', $admin_username, $admin_pass, $admin_email, $site_name, $database_label, $search_fallback ? self::app_db_search_fallback_notice() : '');
}

public static function us_unlock(): void
{
    if (isset($GLOBALS['update_lock_handle']) && is_resource($GLOBALS['update_lock_handle'])) {
        flock($GLOBALS['update_lock_handle'], LOCK_UN);
        fclose($GLOBALS['update_lock_handle']);
        unset($GLOBALS['update_lock_handle']);
    }
}

public static function us_styles(): string
{
    return '.update-page{min-height:100vh;padding:28px 12px;background:#f6f7f8}.update-card{width:min(720px,100%);margin:auto;padding:28px;border:1px solid #e8e8e8;border-radius:8px;background:#fff;box-shadow:0 18px 45px rgba(16,24,40,.08)}.update-title{display:flex;align-items:center;gap:9px;margin:0;color:#111;font-size:22px;line-height:1.3}.update-file-version{padding:2px 6px;border:1px solid #dfe4e1;border-radius:4px;background:#f7f9f8;color:#6d7571;font:12px/1.4 ui-monospace,SFMono-Regular,Menlo,monospace}.update-sub{margin:7px 0 20px;color:#777;font-size:13px;line-height:1.7}.update-grid{display:grid;grid-template-columns:1fr 1fr;gap:12px;margin:0 0 18px}.update-panel{padding:16px;border:1px solid #eee;border-radius:6px;background:#fafafa}.update-panel strong{display:block;margin-bottom:6px;color:#222}.update-panel span{display:block;color:#777;font-size:13px;line-height:1.6}.update-version{font-family:ui-monospace,SFMono-Regular,Menlo,monospace}.update-notice,.update-warning,.update-error{margin:0 0 18px;padding:14px;border:1px solid #dfe8e3;border-radius:6px;background:#f8fcfa;color:#376348;font-size:13px;line-height:1.7;word-break:break-word}.update-warning{border-color:#f3d6a2;background:#fffaf0;color:#8a5a13}.update-error{border-color:#ffd8d8;background:#fff8f8;color:#b42318}.update-list{list-style:none;margin:0 0 20px;padding:0;border:1px solid #e5e7eb;border-radius:6px;overflow:hidden;background:#fff}.update-list li{padding:0;border-bottom:1px solid #edf0f2;color:#444;font-size:13px}.update-list li:last-child{border-bottom:0}.update-list label{display:flex;align-items:center;gap:12px;min-height:48px;padding:10px 14px;cursor:pointer;transition:background .15s ease}.update-list label:hover{background:#f7faf8}.update-list li:has(input:checked){background:#f8fcfa}.update-list input[type=checkbox]{width:18px;height:18px;margin:0;accent-color:#20a45a;cursor:pointer;flex:0 0 18px}.update-file-type{display:inline-flex;align-items:center;justify-content:center;min-width:44px;padding:3px 7px;border-radius:4px;background:#eef8f2;color:#267247;font-size:12px;line-height:1.2}.update-file-path{min-width:0;color:#252b2e;font-family:ui-monospace,SFMono-Regular,Menlo,monospace;overflow-wrap:anywhere}.update-schema-item{background:#fafafa}.update-schema-copy{display:grid;gap:2px;min-width:0}.update-schema-copy strong{color:#30363a;font-size:13px}.update-schema-copy span{color:#7a8185;font-size:12px;line-height:1.5}.update-result-item{padding:11px 14px!important;line-height:1.6}.update-actions{display:flex;align-items:center;justify-content:flex-end;gap:10px;flex-wrap:wrap}.update-actions form{display:flex;margin:0}.update-actions a,.update-actions button{display:inline-flex;align-items:center;justify-content:center;min-height:36px;padding:0 14px;border:1px solid #ddd;border-radius:6px;background:#fff;color:#555;font:inherit;text-decoration:none;cursor:pointer}.update-actions button.primary{border-color:#2ecc71;background:#2ecc71;color:#fff}.update-actions button:disabled{cursor:not-allowed;opacity:.55}@media(max-width:600px){.update-card{padding:20px}.update-grid{grid-template-columns:1fr}.update-list label{gap:10px;padding:10px 12px}.update-file-type{min-width:40px}.update-actions{display:grid;grid-template-columns:repeat(auto-fit,minmax(128px,1fr));align-items:stretch}.update-actions a,.update-actions form,.update-actions button{width:100%;min-width:0}}';
}

public static function us_result_page(string $title, array $changes, string $error = ''): void
{
    $body = '<h1 class="update-title">' . h($title) . '</h1><p class="update-sub">勾选数据库同步时，将根据当前程序幂等同步数据库结构和索引。</p>';
    if ($error !== '') {
        $body .= '<div class="update-error">' . h($error) . '</div>';
    } elseif ($changes) {
        $body .= '<ul class="update-list">';
        foreach ($changes as $change) $body .= '<li class="update-result-item">' . h($change) . '</li>';
        $body .= '</ul>';
    } else {
        $body .= '<div class="update-notice">数据库结构和索引已是最新，无需调整。</div>';
    }
    $body .= '<div class="update-actions"><a href="index.php?a=update">返回升级页</a><a href="index.php">进入首页</a></div>';
    self::us_unlock();
    self::setup_html($title, '<style>' . self::us_styles() . '</style><section class="update-card">' . $body . '</section>');
}

public static function us_need_admin(): void
{
    if (!uid()) self::us_result_page('请先登录', [], '请先登录管理员账号后再执行升级。');
    if (uid() === 1) return;
    if (!app_db_table_exists(db(), db_driver(), 'app_users') || !app_db_table_exists(db(), db_driver(), 'app_groups') || !can_access_admin()) self::us_result_page('无权限', [], '当前账号没有后台管理权限。');
}

public static function us_legacy_upgrade_state(PDO $db, string $driver): string
{
    $legacy = [];
    $current = [];
    foreach (self::migrate_core_table_map() as $table => $target) {
        if ($table === 'topics_fts') continue;
        if (app_db_table_exists($db, $driver, $table)) $legacy[] = $table;
        if (app_db_table_exists($db, $driver, $target)) $current[] = $target;
    }
    if (!$legacy) return '';
    if ($current) return '系统表同时存在旧名称和新名称，请先检查数据库。';
    foreach (['groups', 'users', 'forums', 'topics', 'replies', 'settings'] as $table) {
        if (!app_db_table_exists($db, $driver, $table)) return '旧数据库缺少核心表：' . $table . '。请使用备份重新安装后执行数据迁入。';
    }
    return 'legacy';
}

public static function us_legacy_admin_id(string $username, string $password): int
{
    $user = one('SELECT * FROM users WHERE username=?', [trim($username)]);
    if (!$user || !password_verify($password, (string)$user['password']) || (int)($user['is_banned'] ?? 0) === 1) return 0;
    $user_id = (int)$user['id'];
    if ($user_id === 1) return $user_id;
    $group = one('SELECT * FROM groups WHERE id=?', [(int)($user['group_id'] ?? 0)]);
    return (int)($group['allow_admin'] ?? 0) === 1 ? $user_id : 0;
}

public static function us_acquire_lock(): void
{
    $lock = fopen(UPDATE_RUN_LOCK_FILE, 'c');
    if (!$lock || !flock($lock, LOCK_EX | LOCK_NB)) self::us_result_page('升级失败', [], '升级正在执行，请稍后再试。');
    $GLOBALS['update_lock_handle'] = $lock;
}

public static function us_legacy_upgrade_page(string $error = ''): never
{
    $token = csrf_token();
    $body = '<h1 class="update-title">旧版本数据库升级 <span class="update-file-version">' . h(APP_VERSION) . '</span></h1><p class="update-sub">检测到数据库仍使用无 app_ 前缀的旧系统表。升级将原子改名系统表并同步当前结构。</p>';
    if ($error !== '') $body .= '<div class="update-error">' . h($error) . '</div>';
    $body .= '<div class="update-warning"><strong>操作前必须完整备份数据库。</strong>若旧库结构不完整，请重新安装后使用“数据迁入”。</div><form method="post"><input type="hidden" name="_csrf" value="' . h($token) . '"><input type="hidden" name="legacy_upgrade" value="1"><div class="update-grid"><label class="update-panel"><strong>旧版管理员用户名</strong><input type="text" name="username" required autocomplete="username"></label><label class="update-panel"><strong>旧版管理员密码</strong><input type="password" name="password" required autocomplete="current-password"></label></div><label class="update-warning"><input type="checkbox" name="confirm_backup" value="1" required> 已完成数据库备份</label><div class="update-actions"><button class="primary" type="submit">升级旧数据库</button></div></form>';
    self::setup_html('旧版本数据库升级', '<style>' . self::us_styles() . '</style><section class="update-card">' . $body . '</section>');
}

public static function us_handle_legacy_upgrade(): never
{
    if ($_SERVER['REQUEST_METHOD'] !== 'POST') self::us_legacy_upgrade_page();
    if (!hash_equals(csrf_token(), (string)($_POST['_csrf'] ?? ''))) self::us_legacy_upgrade_page('请求已过期，请返回重试。');
    if (!isset($_POST['confirm_backup'])) self::us_legacy_upgrade_page('请先确认已完成数据库备份。');
    $user_id = self::us_legacy_admin_id((string)($_POST['username'] ?? ''), (string)($_POST['password'] ?? ''));
    if ($user_id <= 0) self::us_legacy_upgrade_page('管理员账号或密码错误。');
    self::us_acquire_lock();
    try {
        $changes = self::us_sync_schema();
        start_cookie_login($user_id);
        self::us_result_page('升级完成', $changes);
    } catch (Throwable $e) {
        self::us_result_page('升级失败', [], $e->getMessage());
    }
}

public static function us_http(string $url): string
{
    $context = stream_context_create(['http' => [
        'method' => 'GET',
        'header' => "Accept: application/vnd.github+json\r\nUser-Agent: bbs1org-updater\r\n",
        'timeout' => 15,
        'follow_location' => 1,
        'max_redirects' => 3,
        'ignore_errors' => true,
    ]]);
    $body = @file_get_contents($url, false, $context);
    $status = 0;
    foreach (($http_response_header ?? []) as $header) {
        if (preg_match('/^HTTP\/\S+\s+(\d{3})/', $header, $m)) $status = (int)$m[1];
    }
    if ($body === false || $status < 200 || $status >= 300) throw new RuntimeException('连接 GitHub 失败（HTTP ' . ($status ?: '未知') . '）。');
    return $body;
}

public static function us_remote_release(): array
{
    $json = json_decode(self::us_http('https://api.github.com/repos/' . UPDATE_REPOSITORY . '/commits/' . UPDATE_BRANCH), true, 512, JSON_THROW_ON_ERROR);
    $sha = (string)($json['sha'] ?? '');
    $tree_url = (string)($json['commit']['tree']['url'] ?? '');
    if (!preg_match('/^[a-f0-9]{40}$/', $sha) || $tree_url === '') throw new RuntimeException('GitHub 返回的版本信息无效。');
    $tree = json_decode(self::us_http($tree_url . '?recursive=1'), true, 512, JSON_THROW_ON_ERROR);
    if (!empty($tree['truncated']) || !is_array($tree['tree'] ?? null)) throw new RuntimeException('GitHub 返回的文件清单不完整。');
    $files = [];
    foreach ($tree['tree'] as $item) {
        if (($item['type'] ?? '') !== 'blob') continue;
        $path = (string)($item['path'] ?? '');
        if (!self::us_update_ignored_path($path)) $files[$path] = (string)($item['sha'] ?? '');
    }
    return [
        'sha' => $sha,
        'short_sha' => substr($sha, 0, 12),
        'date' => (string)($json['commit']['committer']['date'] ?? ''),
        'message' => trim(strtok((string)($json['commit']['message'] ?? ''), "\r\n")),
        'files' => $files,
    ];
}

public static function us_git_blob_sha(string $file): string
{
    $content = (string)file_get_contents($file);
    return sha1('blob ' . strlen($content) . "\0" . $content);
}

public static function update_state_data(): array
{
    if (!is_file(UPDATE_STATE_FILE)) return [];
    $state = json_decode((string)file_get_contents(UPDATE_STATE_FILE), true);
    return is_array($state) ? $state : [];
}

public static function update_state_write(array $state): void
{
    $json = json_encode($state, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES | JSON_PRETTY_PRINT);
    $swap = UPDATE_STATE_FILE . '.tmp-' . bin2hex(random_bytes(4));
    if ($json === false || file_put_contents($swap, $json, LOCK_EX) === false || !rename($swap, UPDATE_STATE_FILE)) {
        @unlink($swap);
        throw new RuntimeException('无法更新升级状态');
    }
}

public static function us_local_changes(array $remote_files): array
{
    $changes = [];
    foreach ($remote_files as $path => $sha) {
        $file = APP_ROOT . '/' . $path;
        if (!is_file($file)) $changes[] = ['path' => $path, 'type' => '新增'];
        elseif (!hash_equals($sha, self::us_git_blob_sha($file))) $changes[] = ['path' => $path, 'type' => '更新'];
    }
    foreach ((array)(self::update_state_data()['files'] ?? []) as $path) {
        if (is_string($path) && !isset($remote_files[$path]) && !self::us_protected_path($path) && is_file(APP_ROOT . '/' . $path)) {
            $changes[] = ['path' => $path, 'type' => '删除'];
        }
    }
    return $changes;
}

public static function us_json(array $data): never
{
    header('Content-Type: application/json; charset=utf-8');
    echo json_encode($data, JSON_UNESCAPED_UNICODE | JSON_UNESCAPED_SLASHES);
    exit;
}

public static function deliver_update_notice(): void
{
    $state = self::update_state_data();
    $notice = is_array($state['update_notice'] ?? null) ? $state['update_notice'] : [];
    $sha = (string)($notice['sha'] ?? '');
    if (!preg_match('/^[a-f0-9]{40}$/', $sha)) return;
    $short_sha = substr($sha, 0, 12);
    $message = trim((string)($notice['message'] ?? ''));
    $content = '检测到系统新版本 ' . $short_sha . '。' . ($message !== '' ? "\n\n" . cut($message, 120) : '') . "\n\n请前往后台设置中的“系统升级”完成升级。";
    if (!one("SELECT 1 FROM app_notifications WHERE recipient_id=? AND kind='system_update' AND content=? LIMIT 1", [uid(), $content])) {
        create_notification(uid(), 0, 'system_update', $content);
    }
    unset($state['update_notice']);
    $state['update_notice_sent_sha'] = $sha;
    self::update_state_write($state);
    unset($GLOBALS['__me_cache']);
}

public static function us_notice_check(): never
{
    $lock = @fopen(UPDATE_RUN_LOCK_FILE, 'c');
    if (!$lock || !flock($lock, LOCK_EX | LOCK_NB)) self::us_json(['ok' => 1, 'pending' => 1]);
    try {
        $state = self::update_state_data();
        $available = is_array($state['update_notice'] ?? null) || preg_match('/^[a-f0-9]{40}$/', (string)($state['update_notice_sent_sha'] ?? '')) === 1;
        $last_checked = strtotime((string)($state['last_notice_checked_at'] ?? '')) ?: 0;
        if (!$available && $last_checked > time() - UPDATE_NOTICE_CHECK_INTERVAL) {
            self::us_json(['ok' => 1, 'update_available' => 0, 'cached' => 1]);
        }
        if (!is_array($state['update_notice'] ?? null)) {
            $state['last_notice_checked_at'] = date(DATE_ATOM);
            self::update_state_write($state);
            $release = self::us_remote_release();
            $changes = self::us_local_changes((array)$release['files']);
            $sha = (string)$release['sha'];
            if ($changes && !hash_equals((string)($state['update_notice_sent_sha'] ?? ''), $sha)) {
                $state['update_notice'] = [
                    'sha' => $sha,
                    'message' => (string)($release['message'] ?? ''),
                    'checked_at' => date(DATE_ATOM),
                ];
                self::update_state_write($state);
                $available = true;
            }
        }
        self::us_json(['ok' => 1, 'update_available' => $available ? 1 : 0]);
    } catch (Throwable $e) {
        self::us_json(['ok' => 0, 'message' => $e->getMessage() ?: '检查升级失败']);
    } finally {
        flock($lock, LOCK_UN);
        fclose($lock);
    }
}

public static function us_update_page(?array $release = null, string $error = ''): void
{
    $token = csrf_token();
    $state = self::update_state_data();
    $local = isset($state['sha']) ? substr((string)$state['sha'], 0, 12) : '未记录';
    $local_time = ($timestamp = strtotime((string)($state['updated_at'] ?? ''))) !== false ? date('Y-m-d H:i', $timestamp) : '';
    $body = '<h1 class="update-title">系统升级 <span class="update-file-version">' . h(APP_VERSION) . '</span></h1><p class="update-sub">检测并安装 ' . h(UPDATE_REPOSITORY) . ' 主分支的最新代码，也可单独同步当前代码对应的数据库结构。</p>';
    if ($error !== '') $body .= '<div class="update-error">' . h($error) . '</div>';
    if ($release) {
        $changes = self::us_local_changes($release['files']);
        $remote_time = ($timestamp = strtotime((string)$release['date'])) !== false ? date('Y-m-d H:i', $timestamp) : (string)$release['date'];
        $body .= '<div class="update-grid"><div class="update-panel"><strong>本地记录</strong><span class="update-version">' . h($local) . '</span>' . ($local_time !== '' ? '<span>更新时间：' . h($local_time) . '</span>' : '') . '</div><div class="update-panel"><strong>远端最新</strong><span class="update-version">' . h($release['short_sha']) . '</span><span>最后提交：' . h($remote_time) . '</span><span>' . h($release['message']) . '</span></div></div>';
        if ($changes) {
            $body .= '<div class="update-notice">检测到 ' . count($changes) . ' 个代码文件需要新增或更新。</div><div class="update-warning"><strong>警告：</strong>勾选文件的本地内容和修改将被 GitHub main 分支版本覆盖。</div><ul class="update-list">';
            foreach ($changes as $change) {
                $path = (string)($change['path'] ?? '');
                $type = (string)($change['type'] ?? '变更');
                $body .= '<li><label><input type="checkbox" name="files[]" value="' . h($path) . '" form="online-update-form" checked><span class="update-file-type">' . h($type) . '</span><span class="update-file-path">' . h($path) . '</span></label></li>';
            }
            $body .= '<li class="update-schema-item"><label><input type="checkbox" name="sync_schema" value="1" form="online-update-form" checked><span class="update-schema-copy"><strong>同步数据库结构</strong><span>文件更新并确认 OPcache 已清理后，由新请求同步缺少的表、字段和索引</span></span></label></li>';
            $body .= '</ul>';
        } else {
            $body .= '<div class="update-notice">当前程序文件已是最新版本。</div>';
        }
    } else {
        $changes = [];
        $body .= '<div class="update-notice">点击“检测更新”连接 GitHub 并逐文件核对当前程序。</div>';
    }
    $body .= '<div class="update-actions"><a href="index.php">返回首页</a><a href="index.php?a=migrate">数据迁入</a><a href="index.php?a=update&amp;check=1">检测更新</a>';
    if (!$release || !$changes) $body .= '<form method="post" data-no-ajax="1"><input type="hidden" name="_csrf" value="' . h($token) . '"><input type="hidden" name="action" value="schema"><button type="submit">同步数据库</button></form>';
    if ($release && $changes) $body .= '<form id="online-update-form" method="post" data-no-ajax="1" data-confirm="确定下载并覆盖已勾选的程序文件？"><input type="hidden" name="_csrf" value="' . h($token) . '"><input type="hidden" name="action" value="online"><input type="hidden" name="sha" value="' . h($release['sha']) . '"><button class="primary" type="submit">在线升级</button></form>';
    $body .= '</div>';
    self::us_unlock();
    self::setup_html('系统升级', '<style>' . self::us_styles() . '</style><section class="update-card">' . $body . '</section>', true);
}

public static function us_protected_path(string $path): bool
{
    $path = trim(str_replace('\\', '/', $path), '/');
    if ($path === '' || str_contains($path, "\0") || preg_match('#(^|/)\.\.(/|$)#', $path)) return true;
    foreach (UPDATE_PROTECTED_DIRS as $protected) if ($path === $protected || str_starts_with($path, $protected . '/')) return true;
    return false;
}

public static function us_update_ignored_path(string $path): bool
{
    $path = trim(str_replace('\\', '/', $path), '/');
    if ($path === '' || str_contains($path, "\0") || preg_match('#(^|/)\.\.(/|$)#', $path)) return true;
    foreach (explode('/', $path) as $name) if (str_starts_with($name, '.')) return true;
    return str_ends_with(strtolower($path), '.md');
}

public static function us_remove_dir(string $dir): void
{
    if (!is_dir($dir)) return;
    $items = scandir($dir);
    if ($items === false) return;
    foreach ($items as $item) {
        if ($item === '.' || $item === '..') continue;
        $path = $dir . '/' . $item;
        is_dir($path) && !is_link($path) ? self::us_remove_dir($path) : @unlink($path);
    }
    @rmdir($dir);
}

public static function us_writable_parent(string $path): bool
{
    $parent = dirname($path);
    while (!is_dir($parent) && $parent !== dirname($parent)) $parent = dirname($parent);
    return is_dir($parent) && is_writable($parent);
}

public static function us_refresh_fpm_opcache_after_update(): string
{
    if (PHP_SAPI !== 'cli') return 'OPcache 已成功清理';

    $lock_file = UPDATE_DATA_DIR . '/opcache-refresh.lock';
    if (file_put_contents($lock_file, '', LOCK_EX) === false) throw new RuntimeException('无法创建 OPcache 刷新锁。');
    try {
        $result = @file_get_contents('http://nginx/?a=opcache_refresh');
        if ($result !== 'ok') throw new RuntimeException('程序文件已更新，但 PHP-FPM OPcache 清理失败；请重启 PHP-FPM 后重试。');
        return 'OPcache 已成功清理';
    } finally {
        @unlink($lock_file);
    }
}

public static function us_refresh_opcache_after_update(array $files): string
{
    $php_updated = (bool)array_filter($files, static fn(string $path): bool => str_ends_with(strtolower($path), '.php'));
    if (!$php_updated) return '本次未更新 PHP 文件，无需清理 OPcache';

    $opcache_enabled = filter_var((string)ini_get('opcache.enable'), FILTER_VALIDATE_BOOL);
    if (function_exists('opcache_get_status')) {
        try {
            $opcache_enabled = $opcache_enabled || is_array(@opcache_get_status(false));
        } catch (Throwable $e) {
            // 无法读取状态时，以配置值判断，随后仍要求 opcache_reset() 成功。
        }
    }
    if (!$opcache_enabled) return self::us_refresh_fpm_opcache_after_update();
    if (!function_exists('opcache_reset')) {
        throw new RuntimeException('程序文件已更新，但 OPcache 已启用且无法调用 opcache_reset()；数据库结构尚未同步，请清理 OPcache 后重试。');
    }
    try {
        $cleared = opcache_reset();
    } catch (Throwable $e) {
        $cleared = false;
    }
    if (!$cleared) {
        throw new RuntimeException('程序文件已更新，但 OPcache 清理失败；数据库结构尚未同步，请清理 OPcache 后重试。');
    }
    return self::us_refresh_fpm_opcache_after_update();
}

public static function us_install_files(string $sha, array $remote_files, array $selected): array
{
    if (!preg_match('/^[a-f0-9]{40}$/', $sha)) throw new RuntimeException('升级版本无效。');
    $temp = UPDATE_DATA_DIR . '/update-' . bin2hex(random_bytes(6));
    if (!mkdir($temp, 0700, true)) throw new RuntimeException('无法创建升级临时目录。');
    try {
        $files = [];
        foreach ($remote_files as $path => $expected_sha) {
            if (!is_string($path) || self::us_update_ignored_path($path)) continue;
            if (!in_array($path, $selected, true)) continue;
            $content = self::us_http('https://raw.githubusercontent.com/' . UPDATE_REPOSITORY . '/' . $sha . '/' . $path);
            if (strlen($content) > UPDATE_MAX_ARCHIVE_BYTES || !hash_equals($expected_sha, sha1('blob ' . strlen($content) . "\0" . $content))) throw new RuntimeException('远端文件校验失败：' . $path);
            $target = $temp . '/' . $path;
            if (!is_dir(dirname($target)) && !mkdir(dirname($target), 0700, true)) throw new RuntimeException('无法创建临时目录：' . dirname($path));
            if (file_put_contents($target, $content, LOCK_EX) === false) throw new RuntimeException('无法保存临时文件：' . $path);
            $files[] = $path;
        }
        if (!$files) throw new RuntimeException('请至少选择一个需要升级的文件。');
        $backups = [];
        foreach ($files as $path) {
            $target = APP_ROOT . '/' . $path;
            if (!self::us_writable_parent($target)) throw new RuntimeException('文件所在目录不可写：' . $path);
            $backup = $temp . '/backup/' . $path;
            $existed = is_file($target);
            if ($existed) {
                if (!is_dir(dirname($backup)) && !mkdir(dirname($backup), 0700, true)) throw new RuntimeException('无法创建备份目录。');
                if (!copy($target, $backup) || !hash_equals((string)hash_file('sha256', $target), (string)hash_file('sha256', $backup))) throw new RuntimeException('备份文件失败：' . $path);
            }
            $backups[$path] = ['file' => $backup, 'existed' => $existed];
        }
        $state_existed = is_file(UPDATE_STATE_FILE);
        $state_backup = $temp . '/update-state.json';
        if ($state_existed && (!copy(UPDATE_STATE_FILE, $state_backup) || !hash_equals((string)hash_file('sha256', UPDATE_STATE_FILE), (string)hash_file('sha256', $state_backup)))) throw new RuntimeException('备份版本记录失败。');
        $replaced = [];
        try {
            foreach ($files as $path) {
                $source = $temp . '/' . $path;
                $target = APP_ROOT . '/' . $path;
                if (!is_dir(dirname($target)) && !mkdir(dirname($target), 0755, true)) throw new RuntimeException('无法创建目录：' . dirname($path));
                $swap = $target . '.update-' . bin2hex(random_bytes(4));
                if (!copy($source, $swap) || !rename($swap, $target)) {
                    @unlink($swap);
                    throw new RuntimeException('更新文件失败：' . $path);
                }
                $replaced[] = $path;
                if (!hash_equals((string)$remote_files[$path], self::us_git_blob_sha($target))) throw new RuntimeException('更新后校验失败：' . $path);
            }
            $state = json_encode(['sha' => $sha, 'updated_at' => date(DATE_ATOM), 'files' => $files], JSON_UNESCAPED_SLASHES | JSON_PRETTY_PRINT);
            $state_swap = UPDATE_STATE_FILE . '.update-' . bin2hex(random_bytes(4));
            if ($state === false || file_put_contents($state_swap, $state, LOCK_EX) === false || !rename($state_swap, UPDATE_STATE_FILE)) {
                @unlink($state_swap);
                throw new RuntimeException('无法写入版本记录。');
            }
        } catch (Throwable $e) {
            $rollback_errors = [];
            foreach (array_reverse($replaced) as $path) {
                $target = APP_ROOT . '/' . $path;
                $backup = $backups[$path];
                if ($backup['existed']) {
                    $swap = $target . '.rollback-' . bin2hex(random_bytes(4));
                    if (!copy($backup['file'], $swap) || !rename($swap, $target)) {
                        @unlink($swap);
                        $rollback_errors[] = $path;
                    }
                } elseif (is_file($target) && !unlink($target)) {
                    $rollback_errors[] = $path;
                }
            }
            if ($state_existed) {
                if (!copy($state_backup, UPDATE_STATE_FILE)) $rollback_errors[] = basename(UPDATE_STATE_FILE);
            } elseif (is_file(UPDATE_STATE_FILE) && !unlink(UPDATE_STATE_FILE)) {
                $rollback_errors[] = basename(UPDATE_STATE_FILE);
            }
            if ($rollback_errors) throw new RuntimeException($e->getMessage() . '；回滚失败：' . implode('、', $rollback_errors), 0, $e);
            throw new RuntimeException($e->getMessage() . '；已恢复升级前文件。', 0, $e);
        }
        clearstatcache();
        return ['count' => count($files), 'opcache' => self::us_refresh_opcache_after_update($files)];
    } finally {
        self::us_remove_dir($temp);
    }
}

public static function us_split_defs(string $body): array
{
    $defs = [];
    $buf = '';
    $depth = 0;
    for ($i = 0, $len = strlen($body); $i < $len; $i++) {
        $ch = $body[$i];
        if ($ch === '(') $depth++;
        if ($ch === ')') $depth--;
        if ($ch === ',' && $depth === 0) {
            $defs[] = trim($buf);
            $buf = '';
            continue;
        }
        $buf .= $ch;
    }
    if (trim($buf) !== '') $defs[] = trim($buf);
    return $defs;
}

public static function us_parse_table_sql(string $sql): array
{
    if (!preg_match('/CREATE\s+TABLE(?:\s+IF\s+NOT\s+EXISTS)?\s+[`"]?([a-zA-Z0-9_]+)[`"]?\s*\((.*)\)\s*(?:ENGINE=.*)?;?\s*$/is', trim($sql), $m)) return [];
    $columns = [];
    foreach (self::us_split_defs($m[2]) as $def) {
        if (preg_match('/^(PRIMARY|UNIQUE|CHECK|FOREIGN|CONSTRAINT)\b/i', $def)) continue;
        if (preg_match('/^([a-zA-Z0-9_]+)\s+(.+)$/s', $def, $cm)) $columns[$cm[1]] = $def;
    }
    return ['name' => $m[1], 'sql' => rtrim(trim($sql), ';') . ';', 'columns' => $columns];
}

public static function us_install_schema(): array
{
    $driver = db_driver();
    [$schema_tables, $schema_indexes] = self::app_db_schema($driver);
    $tables = [];
    foreach ($schema_tables as $sql) {
        $table = self::us_parse_table_sql($sql);
        if ($table) $tables[$table['name']] = $table;
    }
    $virtual_tables = [];
    if ($driver === 'sqlite') {
        $virtual_tables['app_topics_fts'] = 'title, body';
        $virtual_tables['app_replies_fts'] = 'body';
    }
    return [$tables, $virtual_tables, $schema_indexes];
}

public static function us_column_type(PDO $db, string $driver, string $table, string $column): string
{
    if ($driver === 'sqlite') {
        foreach ($db->query('PRAGMA table_info(' . app_db_identifier($driver, $table) . ')')->fetchAll() as $row) {
            if ((string)($row['name'] ?? '') === $column) return strtolower((string)($row['type'] ?? ''));
        }
        return '';
    }
    $sql = 'SELECT ' . ($driver === 'mysql' ? 'column_type' : 'data_type') . ' FROM information_schema.columns WHERE table_schema=' . ($driver === 'mysql' ? 'DATABASE()' : 'current_schema()') . ' AND table_name=? AND column_name=?';
    $stmt = $db->prepare($sql);
    $stmt->execute([$table, $column]);
    return strtolower((string)$stmt->fetchColumn());
}

public static function us_rename_legacy_system_tables(PDO $db, string $driver): array
{
    $changes = [];
    $tables = self::migrate_core_table_map();
    foreach ($tables as $table => $target) {
        if (app_db_table_exists($db, $driver, $table) && app_db_table_exists($db, $driver, $target)) {
            throw new RuntimeException('系统表同时存在旧名称和新名称：' . $table . '、' . $target . '，请先检查数据。');
        }
    }
    if ($driver === 'mysql') {
        $renames = [];
        foreach ($tables as $table => $target) {
            if (!app_db_table_exists($db, $driver, $table)) continue;
            $renames[] = app_db_identifier($driver, $table) . ' TO ' . app_db_identifier($driver, $target);
            $changes[] = '重命名系统表：' . $table . ' -> ' . $target;
        }
        if ($renames) $db->exec('RENAME TABLE ' . implode(',', $renames));
        return $changes;
    }
    foreach ($tables as $table => $target) {
        if (!app_db_table_exists($db, $driver, $table)) continue;
        if ($driver === 'sqlite' && $table === 'topics_fts') {
            $fts_created = app_db_create_fts5_table($db, $target, 'title, body', false);
            if ($fts_created && app_db_table_exists($db, $driver, 'app_topics')) {
                $db->exec('INSERT INTO app_topics_fts(rowid,title,body) SELECT id,title,body FROM app_topics');
            }
            $db->exec('DROP TABLE ' . app_db_identifier($driver, $table));
            $changes[] = $fts_created ? '重命名系统表：' . $table . ' -> ' . $target : '删除不兼容的旧搜索表：' . $table;
            continue;
        } else {
            $db->exec('ALTER TABLE ' . app_db_identifier($driver, $table) . ' RENAME TO ' . app_db_identifier($driver, $target));
        }
        $changes[] = '重命名系统表：' . $table . ' -> ' . $target;
    }
    return $changes;
}

public static function us_migrate_legacy_plugin_settings(): int
{
    $settings = settings_cache();
    $registered = array_fill_keys(array_map('strval', q("SELECT id FROM app_plugins")->fetchAll(PDO::FETCH_COLUMN)), true);
    $files = [];
    foreach (Plugin::plugin_files() as $file) {
        $id = basename(dirname($file));
        if ($id !== 'plugin_market') $files[$id] = $file;
    }
    $files['plugin_market'] = '';
    $suffixes = ['enabled', 'version', 'config', 'disabled_reason', 'entry_feature_links', 'entry_sidebar_cards'];
    $delete = [];
    $count = 0;
    foreach ($files as $id => $file) {
        $names = [];
        foreach ($suffixes as $suffix) {
            $name = 'plugin_' . $id . '_' . $suffix;
            if (array_key_exists($name, $settings)) $names[] = $name;
        }
        if (!$names) continue;
        $delete = array_merge($delete, $names);
        $count++;
        if ($file === '' || isset($registered[$id])) continue;
        $config = plugin_json_decode($settings['plugin_' . $id . '_config'] ?? '') ?? [];
        $enabled = (string)($settings['plugin_' . $id . '_enabled'] ?? '0') === '1' ? 1 : 0;
        app_db_insert_ignore('app_plugins', [
            'id' => $id,
            'name' => $id,
            'version' => (string)($settings['plugin_' . $id . '_version'] ?? ''),
            'file' => ltrim(str_replace(APP_ROOT, '', $file), '/'),
            'code_hash' => '',
            'manifest_json' => plugin_json_encode([]),
            'config_json' => plugin_json_encode($config),
            'entries_json' => plugin_json_encode([
                'feature_links' => (string)($settings['plugin_' . $id . '_entry_feature_links'] ?? '1') === '1',
                'sidebar_cards' => (string)($settings['plugin_' . $id . '_entry_sidebar_cards'] ?? '1') === '1',
            ]),
            'enabled' => $enabled,
            'status' => $enabled ? 'enabled' : 'disabled',
            'disabled_reason' => (string)($settings['plugin_' . $id . '_disabled_reason'] ?? ''),
            'installed_at' => now(),
            'updated_at' => now(),
        ], ['id']);
        $registered[$id] = true;
    }
    if ($delete) {
        $delete = array_values(array_unique($delete));
        q("DELETE FROM app_settings WHERE name IN (" . sql_marks(count($delete)) . ")", $delete);
        db_row_cache_clear();
        if (is_array($GLOBALS['__settings_cache'] ?? null)) foreach ($delete as $name) unset($GLOBALS['__settings_cache'][$name]);
    }
    return $count;
}

public static function us_sync_schema(): array
{
    [$tables, $virtual_tables, $indexes] = self::us_install_schema();
    if (!$tables) throw new RuntimeException('未读取到当前程序的数据表结构。');
    $db = db();
    $transactional = db_driver() !== 'mysql';
    $changes = [];
    try {
        if ($transactional) $db->beginTransaction();
        $changes = array_merge($changes, self::us_rename_legacy_system_tables($db, db_driver()));
        self::app_db_prepare_search($db, db_driver());
        $created_virtual_tables = [];
        foreach ($virtual_tables as $table => $columns) {
            if (!app_db_table_exists($db, db_driver(), $table)) {
                if (!app_db_create_fts5_table($db, $table, $columns)) continue;
                $created_virtual_tables[] = $table;
                $changes[] = '新增虚拟表：' . $table;
            }
        }
        foreach ($tables as $table => $schema) {
            if (!app_db_table_exists($db, db_driver(), $table)) {
                $db->exec($schema['sql']);
                $changes[] = '新增表：' . $table;
                continue;
            }
            $current_columns = app_db_columns($db, db_driver(), $table);
            foreach ($schema['columns'] as $column => $definition) {
                if (isset($current_columns[$column])) continue;
                $db->exec('ALTER TABLE ' . app_db_identifier(db_driver(), $table) . ' ADD COLUMN ' . $definition);
                $changes[] = '新增字段：' . $table . '.' . $column;
            }
        }
        if (in_array('app_topics_fts', $created_virtual_tables, true)) {
            $db->exec('INSERT INTO app_topics_fts(rowid,title,body) SELECT id,title,body FROM app_topics');
            $changes[] = '初始化主题搜索索引';
        }
        if (in_array('app_replies_fts', $created_virtual_tables, true)) {
            $db->exec('INSERT INTO app_replies_fts(rowid,body) SELECT id,body FROM app_replies');
            $changes[] = '初始化回帖搜索索引';
        }
        foreach (['idx_attachments_hash'=>'app_attachments', 'idx_topics_user'=>'app_topics', 'idx_topics_user_updated'=>'app_topics', 'idx_topics_forum_updated'=>'app_topics', 'idx_users_created'=>'app_users', 'idx_replies_user'=>'app_replies', 'idx_replies_user_topic_time'=>'app_replies', 'idx_notifications_recipient_read'=>'app_notifications', 'idx_notifications_sender'=>'app_notifications', 'idx_cron_logs_started'=>'app_cron_logs'] as $index => $table) {
            if (!app_db_index_exists($db, db_driver(), $index, $table)) continue;
            app_db_drop_index($index, $table);
            $changes[] = '删除索引：' . $index;
        }
        $topics_table = 'app_topics';
        $topic_columns = app_db_columns($db, db_driver(), $topics_table);
        if (isset($topic_columns['updated_at'])) {
            $db->exec('ALTER TABLE ' . app_db_identifier(db_driver(), $topics_table) . ' DROP COLUMN ' . app_db_identifier(db_driver(), 'updated_at'));
            $changes[] = '删除字段：topics.updated_at';
        }
        $forums_table = app_db_identifier(db_driver(), 'app_forums');
        $forum_columns = app_db_columns($db, db_driver(), 'app_forums');
        foreach (['last_topic_id', 'last_topic_title'] as $column) {
            if (!isset($forum_columns[$column])) continue;
            $db->exec("ALTER TABLE $forums_table DROP COLUMN " . app_db_identifier(db_driver(), $column));
            $changes[] = '删除字段：forums.' . $column;
        }
        $cron_table = app_db_identifier(db_driver(), 'app_cron_tasks');
        $cron_columns = app_db_columns($db, db_driver(), 'app_cron_tasks');
        foreach (['next_run_at', 'updated_at'] as $column) {
            if (!isset($cron_columns[$column])) continue;
            $db->exec("ALTER TABLE $cron_table DROP COLUMN " . app_db_identifier(db_driver(), $column));
            $changes[] = '删除字段：cron_tasks.' . $column;
        }
        if (!str_contains(self::us_column_type($db, db_driver(), $topics_table, 'reply_order'), 'int')) {
            $table = app_db_identifier(db_driver(), $topics_table);
            $column = app_db_identifier(db_driver(), 'reply_order');
            $db->exec("ALTER TABLE $table DROP COLUMN $column");
            $db->exec("ALTER TABLE $table ADD COLUMN $column INTEGER NOT NULL DEFAULT 0");
            $changes[] = '更新字段：topics.reply_order';
        }
        foreach ($indexes as $index => $sql) {
            if (!app_db_index_exists($db, db_driver(), $index, self::app_db_index_table($sql))) {
                if ($index === 'idx_attachments_user_hash') {
                    $removed = $db->exec('DELETE FROM app_attachments WHERE id NOT IN (SELECT keep_id FROM (SELECT MIN(id) keep_id FROM app_attachments GROUP BY user_id,hash) attachment_dedup)');
                    if ($removed) $changes[] = '清理重复附件：' . $removed . ' 条';
                }
                if (self::app_db_create_schema_index($db, db_driver(), $index, $sql)) $changes[] = '新增索引：' . $index;
                elseif (!in_array(self::app_db_search_fallback_notice(), $changes, true)) $changes[] = self::app_db_search_fallback_notice();
            }
        }
        save_settings_values(mysql_search_index_settings($db, db_driver()));
        if ($transactional) $db->commit();
        $legacy_plugin_count = self::us_migrate_legacy_plugin_settings();
        if ($legacy_plugin_count > 0) $changes[] = '迁移旧插件配置：' . $legacy_plugin_count . ' 个';
        $plugin_count = count(Plugin::plugin_registry_sync());
        Plugin::plugin_assets_rebuild();
        $changes[] = '同步插件注册表：' . $plugin_count . ' 个';
        return $changes;
    } catch (Throwable $e) {
        if ($transactional && $db->inTransaction()) $db->rollBack();
        throw $e;
    }
}

public static function us_defer_schema_after_update(array $changes): never
{
    $nonce = bin2hex(random_bytes(24));
    save_settings_values(['update_schema_pending' => json_encode([
        'nonce' => $nonce,
        'created_at' => time(),
        'status' => 'pending',
        'changes' => array_values($changes),
    ], JSON_UNESCAPED_UNICODE | JSON_THROW_ON_ERROR)]);
    self::us_unlock();
    header('Location: index.php?a=update&schema_after_update=' . rawurlencode($nonce), true, 303);
    exit;
}

public static function us_run_deferred_schema(): never
{
    $nonce = (string)($_GET['schema_after_update'] ?? '');
    $pending = preg_match('/^[a-f0-9]{48}$/D', $nonce) === 1 ? json_decode(setting('update_schema_pending', ''), true) : null;
    $created_at = (int)($pending['created_at'] ?? 0);
    $status = (string)($pending['status'] ?? 'pending');
    if (!is_array($pending)
        || $nonce === ''
        || !hash_equals((string)($pending['nonce'] ?? ''), $nonce)
        || $created_at <= 0
        || time() - $created_at > 300
        || !in_array($status, ['pending', 'completed'], true)
    ) {
        self::us_result_page('升级失败', [], '数据库同步请求无效或已过期，请返回升级页重试。');
    }

    $changes = array_values(array_filter((array)($pending['changes'] ?? []), 'is_string'));
    if ($status === 'completed') self::us_result_page('升级完成', $changes);
    self::us_acquire_lock();
    try {
        $schema_changes = self::us_sync_schema();
        $changes = array_merge($changes, $schema_changes ?: ['数据库结构同步完成，当前结构无需调整']);
        save_settings_values(['update_schema_pending' => json_encode([
            'nonce' => $nonce,
            'created_at' => $created_at,
            'status' => 'completed',
            'changes' => $changes,
        ], JSON_UNESCAPED_UNICODE | JSON_THROW_ON_ERROR)]);
        self::us_result_page('升级完成', $changes);
    } catch (Throwable $e) {
        $prefix = $changes ? implode('；', $changes) . '；' : '';
        self::us_result_page('数据库同步失败', [], $prefix . $e->getMessage());
    }
}

public static function setup_update_run(): never
{
    if (!is_file(UPDATE_INSTALL_LOCK_FILE) || !is_file(UPDATE_DB_CONFIG_FILE)) self::us_result_page('请先安装', [], '请先执行安装操作。');
    if (db_driver() === 'sqlite' && !is_file((string)db_config()['path'])) self::us_result_page('请先安装', [], '请先执行安装操作。');
    $legacy_state = self::us_legacy_upgrade_state(db(), db_driver());
    if ($legacy_state === 'legacy') self::us_handle_legacy_upgrade();
    if ($legacy_state !== '') self::us_result_page('无法自动升级', [], $legacy_state);
    self::us_need_admin();

    if (isset($_GET['schema_after_update'])) self::us_run_deferred_schema();

    if ((string)($_GET['notice_check'] ?? '') === '1') self::us_notice_check();

    if ($_SERVER['REQUEST_METHOD'] !== 'POST') {
        if (!isset($_GET['check'])) self::us_update_page();
        try {
            self::us_update_page(self::us_remote_release());
        } catch (Throwable $e) {
            self::us_update_page(null, $e->getMessage());
        }
    }
    if (!hash_equals(csrf_token(), (string)($_POST['_csrf'] ?? ''))) self::us_result_page('升级失败', [], '请求已过期，请返回重试。');

    self::us_acquire_lock();

    try {
        $action = (string)($_POST['action'] ?? 'schema');
        $changes = [];
        if ($action === 'online') {
            $remote = self::us_remote_release();
            $requested_sha = (string)($_POST['sha'] ?? '');
            if (!hash_equals($remote['sha'], $requested_sha)) throw new RuntimeException('远端版本已变化，请重新检测后再升级。');
            $selected = array_values(array_unique(array_filter((array)($_POST['files'] ?? ''), static fn($path): bool => is_string($path) && isset($remote['files'][$path]))));
            if (in_array('index.php', $selected, true)) {
                foreach (['app/optional/Cron.php', 'app/optional/Plugin.php'] as $dependency) {
                    if (isset($remote['files'][$dependency]) && !in_array($dependency, $selected, true)) $selected[] = $dependency;
                }
            }
            if ($selected) {
                $installed = self::us_install_files($remote['sha'], $remote['files'], $selected);
                $changes[] = '程序代码已更新至 ' . $remote['short_sha'] . '（' . (int)$installed['count'] . ' 个文件）';
                $changes[] = (string)$installed['opcache'];
                if (isset($_POST['sync_schema'])) self::us_defer_schema_after_update($changes);
            } elseif (!isset($_POST['sync_schema'])) {
                throw new RuntimeException('请至少选择一个需要执行的升级操作。');
            }
        } elseif ($action !== 'schema') {
            throw new RuntimeException('未知升级操作。');
        }
        if ($action === 'schema' || isset($_POST['sync_schema'])) $changes = array_merge($changes, self::us_sync_schema());
        self::us_result_page('升级完成', $changes);
    } catch (Throwable $e) {
        self::us_result_page('升级失败', [], $e->getMessage());
    }
}

public static function migrate_driver(string $driver): string
{
    if (!in_array($driver, ['sqlite', 'mysql', 'pgsql'], true)) throw new RuntimeException('不支持的数据库类型。');
    return $driver;
}

public static function migrate_identifier(string $driver, string $name, string $context = ''): string
{
    if ($name === '' || str_contains($name, "\0")) {
        $reason = $name === '' ? '为空' : '包含 NUL 字符';
        $detail = '数据库标识符无效：' . $reason . '；驱动=' . $driver . '；原始值=' . var_export($name, true);
        if ($context !== '') $detail .= '；位置=' . $context;
        throw new InvalidArgumentException($detail . '。');
    }
    return $driver === 'mysql' ? '`' . str_replace('`', '``', $name) . '`' : '"' . str_replace('"', '""', $name) . '"';
}

public static function migrate_source_config(): array
{
    $driver = self::migrate_driver((string)($_POST['source_driver'] ?? 'sqlite'));
    if ($driver === 'sqlite') {
        $path = trim((string)($_POST['source_sqlite'] ?? ''));
        if ($path === '') throw new RuntimeException('SQLite 文件不能为空。');
        if ($path[0] !== DIRECTORY_SEPARATOR) $path = APP_ROOT . '/' . $path;
        return ['driver' => 'sqlite', 'database' => basename($path), 'path' => $path];
    }
    $config = [
        'driver' => $driver,
        'host' => trim((string)($_POST['source_host'] ?? '127.0.0.1')),
        'port' => max(1, (int)($_POST['source_port'] ?? ($driver === 'mysql' ? 3306 : 5432))),
        'database' => trim((string)($_POST['source_database'] ?? '')),
        'username' => (string)($_POST['source_username'] ?? ''),
        'password' => (string)($_POST['source_password'] ?? ''),
    ];
    if ($config['host'] === '' || $config['database'] === '' || $config['username'] === '') throw new RuntimeException('数据库连接信息不完整。');
    return $config;
}

public static function migrate_source_db(array $config): PDO
{
    if ($config['driver'] === 'sqlite' && !is_file((string)$config['path'])) throw new RuntimeException('SQLite 文件不存在：' . $config['path']);
    return app_db_connect($config);
}

public static function migrate_tables(PDO $db, string $driver): array
{
    if ($driver === 'sqlite') {
        $rows = $db->query("SELECT name,sql FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name")->fetchAll();
        $virtual = [];
        foreach ($rows as $row) if (stripos((string)$row['sql'], 'CREATE VIRTUAL TABLE') === 0) $virtual[] = (string)$row['name'];
        $tables = [];
        foreach ($rows as $row) {
            $name = (string)$row['name'];
            foreach ($virtual as $prefix) if ($name === $prefix || str_starts_with($name, $prefix . '_')) continue 2;
            $tables[] = $name;
        }
        return $tables;
    }
    $stmt = $driver === 'mysql'
        ? $db->query("SELECT table_name FROM information_schema.tables WHERE table_schema=DATABASE() AND table_type='BASE TABLE' ORDER BY table_name")
        : $db->query('SELECT tablename FROM pg_tables WHERE schemaname=current_schema() ORDER BY tablename');
    return array_map('strval', $stmt->fetchAll(PDO::FETCH_COLUMN));
}

public static function migrate_columns(PDO $db, string $driver, string $table): array
{
    if ($driver === 'sqlite') {
        $rows = $db->query('PRAGMA table_info(' . self::migrate_identifier($driver, $table, '读取源表字段：' . $table) . ')')->fetchAll();
        return array_values(array_filter(array_map(fn(array $row): string => (string)($row['name'] ?? ''), $rows)));
    }
    return array_keys(app_db_columns($db, $driver, $table));
}

public static function migrate_column_schema(PDO $db, string $driver, string $table): array
{
    if ($driver === 'sqlite') {
        $rows = $db->query('PRAGMA table_info(' . self::migrate_identifier($driver, $table, '读取源表结构：' . $table) . ')')->fetchAll();
        return array_map(fn(array $row): array => [
            'name' => (string)$row['name'], 'type' => (string)$row['type'], 'nullable' => !(bool)$row['notnull'],
            'default' => $row['dflt_value'], 'auto' => false, 'pk' => (int)$row['pk'],
        ], $rows);
    }
    if ($driver === 'mysql') {
        $sql = 'SELECT column_name,data_type,column_type,is_nullable,column_default,extra FROM information_schema.columns WHERE table_schema=DATABASE() AND table_name=? ORDER BY ordinal_position';
    } else {
        $sql = 'SELECT column_name,data_type,udt_name,is_nullable,column_default,is_identity FROM information_schema.columns WHERE table_schema=current_schema() AND table_name=? ORDER BY ordinal_position';
    }
    $stmt = $db->prepare($sql);
    $stmt->execute([$table]);
    $columns = [];
    foreach ($stmt->fetchAll() as $row) {
        $columns[] = [
            'name' => (string)$row['column_name'],
            'type' => (string)($row[$driver === 'mysql' ? 'column_type' : 'udt_name'] ?? $row['data_type']),
            'nullable' => (string)$row['is_nullable'] === 'YES',
            'default' => $row['column_default'],
            'auto' => $driver === 'mysql' ? str_contains((string)$row['extra'], 'auto_increment') : (string)$row['is_identity'] === 'YES' || str_starts_with((string)$row['column_default'], 'nextval('),
            'pk' => 0,
        ];
    }
    return $columns;
}

public static function migrate_db_bool(mixed $value): bool
{
    return in_array($value, [true, 1, '1', 't', 'true'], true);
}

public static function migrate_index_schema(PDO $db, string $driver, string $table, array $columns): array
{
    $indexes = [];
    if ($driver === 'sqlite') {
        $primary = [];
        foreach ($columns as $column) if ($column['pk']) $primary[(int)$column['pk']] = $column['name'];
        if ($primary) {
            ksort($primary);
            $indexes['PRIMARY'] = ['name' => 'PRIMARY', 'unique' => true, 'primary' => true, 'columns' => array_values($primary)];
        }
        foreach ($db->query('PRAGMA index_list(' . self::migrate_identifier($driver, $table, '读取源表索引：' . $table) . ')')->fetchAll() as $index) {
            if ((string)$index['origin'] === 'pk') continue;
            $name = (string)$index['name'];
            $items = $db->query('PRAGMA index_info(' . self::migrate_identifier($driver, $name, '读取源索引字段：表=' . $table . '，索引=' . $name) . ')')->fetchAll();
            $names = array_values(array_filter(array_map(fn(array $row): string => (string)($row['name'] ?? ''), $items)));
            if ($names) $indexes[$name] = ['name' => $name, 'unique' => (bool)$index['unique'], 'primary' => false, 'columns' => $names];
        }
        return array_values($indexes);
    }
    if ($driver === 'mysql') {
        $rows = $db->query('SHOW INDEX FROM ' . self::migrate_identifier($driver, $table, '读取源表索引：' . $table))->fetchAll();
        foreach ($rows as $row) {
            $name = (string)$row['Key_name'];
            $indexes[$name] ??= ['name' => $name, 'unique' => !(bool)$row['Non_unique'], 'primary' => $name === 'PRIMARY', 'columns' => []];
            $indexes[$name]['columns'][(int)$row['Seq_in_index']] = (string)$row['Column_name'];
        }
    } else {
        $stmt = $db->prepare('SELECT oid FROM pg_class WHERE relname=? AND relnamespace=(SELECT oid FROM pg_namespace WHERE nspname=current_schema())');
        $stmt->execute([$table]);
        $table_oid = (int)$stmt->fetchColumn();
        if ($table_oid <= 0) return [];
        $stmt = $db->prepare('SELECT indexrelid,indisunique,indisprimary,indkey::text index_keys FROM pg_index WHERE indrelid=?');
        $stmt->execute([$table_oid]);
        $rows = $stmt->fetchAll();
        $index_ids = array_values(array_unique(array_filter(array_map('intval', array_column($rows, 'indexrelid')))));
        $index_names = [];
        if ($index_ids) {
            $marks = sql_marks(count($index_ids));
            $stmt = $db->prepare("SELECT oid,relname FROM pg_class WHERE oid IN ($marks)");
            $stmt->execute($index_ids);
            foreach ($stmt->fetchAll() as $index) $index_names[(int)$index['oid']] = $index;
        }
        $column_ids = [];
        foreach ($rows as $row) $column_ids = array_merge($column_ids, preg_split('/\s+/', trim((string)$row['index_keys'])) ?: []);
        $column_ids = array_values(array_unique(array_filter(array_map('intval', $column_ids))));
        $column_names = [];
        if ($column_ids) {
            $marks = sql_marks(count($column_ids));
            $stmt = $db->prepare("SELECT attnum,attname FROM pg_attribute WHERE attrelid=? AND attnum IN ($marks)");
            $stmt->execute(array_merge([$table_oid], $column_ids));
            foreach ($stmt->fetchAll() as $column) $column_names[(int)$column['attnum']] = (string)$column['attname'];
        }
        foreach ($rows as $row) {
            $index_id = (int)$row['indexrelid'];
            $name = (string)($index_names[$index_id]['relname'] ?? '');
            if ($name === '') continue;
            $keys = array_values(array_filter(array_map('intval', preg_split('/\s+/', trim((string)$row['index_keys'])) ?: [])));
            $names = array_values(array_filter(array_map(fn(int $key): string => $column_names[$key] ?? '', $keys)));
            if ($names) $indexes[$name] = ['name' => $name, 'unique' => self::migrate_db_bool($row['indisunique']), 'primary' => self::migrate_db_bool($row['indisprimary']), 'columns' => $names];
        }
    }
    foreach ($indexes as &$index) {
        ksort($index['columns']);
        $index['columns'] = array_values($index['columns']);
    }
    unset($index);
    return array_values($indexes);
}

public static function migrate_default_sql(PDO $db, string $source_driver, mixed $default, bool $expression = false): string
{
    if ($default === null) return '';
    $value = trim((string)$default);
    if (($value === '' && $source_driver !== 'mysql') || strcasecmp($value, 'NULL') === 0 || str_starts_with($value, 'nextval(')) return '';
    if ($value === '') $literal = $db->quote('');
    elseif (is_numeric($value)) $literal = $value;
    elseif (strcasecmp($value, 'true') === 0 || strcasecmp($value, 'false') === 0) $literal = strcasecmp($value, 'true') === 0 ? '1' : '0';
    if (preg_match('/^CURRENT_(?:TIMESTAMP|DATE|TIME)(?:\(\))?$/i', $value)) return ' DEFAULT ' . strtoupper(rtrim($value, '()'));
    if (!isset($literal)) {
        if ($source_driver === 'pgsql' && preg_match("/^'(.*)'(?:::[A-Za-z0-9_\[\] ]+)?$/s", $value, $match)) $value = str_replace("''", "'", $match[1]);
        elseif ($source_driver !== 'mysql' && strlen($value) >= 2 && (($value[0] === "'" && $value[-1] === "'") || ($value[0] === '"' && $value[-1] === '"'))) $value = str_replace($value[0] . $value[0], $value[0], substr($value, 1, -1));
        $literal = $db->quote($value);
    }
    return ' DEFAULT ' . ($expression ? '(' . $literal . ')' : $literal);
}

public static function migrate_column_type(string $source_type, string $target_driver, bool $indexed, bool $auto): string
{
    $type = strtolower($source_type);
    if ($auto) return match ($target_driver) {
        'mysql' => 'BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT',
        'pgsql' => 'BIGINT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY',
        default => 'INTEGER PRIMARY KEY AUTOINCREMENT',
    };
    if (preg_match('/int|serial|bool/', $type)) return $target_driver === 'sqlite' ? 'INTEGER' : 'BIGINT' . ($target_driver === 'mysql' && str_contains($type, 'unsigned') ? ' UNSIGNED' : '');
    if (preg_match('/real|float|double|decimal|numeric/', $type)) return match ($target_driver) {'mysql' => 'DOUBLE', 'pgsql' => 'DOUBLE PRECISION', default => 'REAL'};
    if (preg_match('/blob|binary|bytea/', $type)) return match ($target_driver) {'mysql' => 'LONGBLOB', 'pgsql' => 'BYTEA', default => 'BLOB'};
    if ($indexed) return $target_driver === 'mysql' ? 'VARCHAR(191)' : 'TEXT';
    return $target_driver === 'mysql' ? 'LONGTEXT' : 'TEXT';
}

public static function migrate_index_name(PDO $db, string $driver, string $table, string $name, array $columns, bool $unique): string
{
    if ($name === 'PRIMARY' || str_starts_with($name, 'sqlite_autoindex_') || !preg_match('/^[A-Za-z_][A-Za-z0-9_]*$/', $name)) $name = ($unique ? 'uidx_' : 'idx_') . $table . '_' . implode('_', $columns);
    $name = strlen($name) > 55 ? substr($name, 0, 44) . '_' . substr(hash('sha256', $name), 0, 10) : $name;
    if (!app_db_index_exists($db, $driver, $name, $table)) return $name;
    return substr($name, 0, 44) . '_' . substr(hash('sha256', $table . ':' . $name), 0, 10);
}

public static function migrate_install_table(PDO $source, string $source_driver, PDO $target, string $target_driver, string $source_table, string $target_table): void
{
    $columns = self::migrate_column_schema($source, $source_driver, $source_table);
    $indexes = self::migrate_index_schema($source, $source_driver, $source_table, $columns);
    $primary = [];
    $indexed = [];
    foreach ($indexes as $index) {
        if ($index['primary']) $primary = $index['columns'];
        foreach ($index['columns'] as $column) $indexed[$column] = true;
    }
    $auto_column = count($primary) === 1 ? $primary[0] : '';
    $definitions = [];
    $auto_primary = false;
    foreach ($columns as $column) {
        $auto = $column['name'] === $auto_column && ($column['auto'] || preg_match('/int|serial/i', $column['type']));
        if ($auto) $auto_primary = true;
        $type = self::migrate_column_type($column['type'], $target_driver, isset($indexed[$column['name']]), $auto);
        $definition = self::migrate_identifier($target_driver, (string)$column['name'], '创建目标表字段：表=' . $target_table . '，字段=' . (string)$column['name']) . ' ' . $type;
        if (!$auto) $definition .= (!$column['nullable'] ? ' NOT NULL' : '') . self::migrate_default_sql($target, $source_driver, $column['default'], $target_driver === 'mysql' && preg_match('/TEXT|BLOB/', $type));
        $definitions[] = $definition;
    }
    if ($primary && !$auto_primary) $definitions[] = 'PRIMARY KEY(' . implode(',', array_map(fn(string $name): string => self::migrate_identifier($target_driver, $name, '创建目标表主键：表=' . $target_table . '，字段=' . $name), $primary)) . ')';
    $sql = 'CREATE TABLE ' . self::migrate_identifier($target_driver, $target_table, '创建目标表：源表=' . $source_table . '，目标表=' . $target_table) . '(' . implode(',', $definitions) . ')';
    if ($target_driver === 'mysql') $sql .= ' ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci';
    $target->exec($sql);
    foreach ($indexes as $index) {
        if ($index['primary'] || !$index['columns']) continue;
        $name = self::migrate_index_name($target, $target_driver, $target_table, $index['name'], $index['columns'], $index['unique']);
        $target->exec('CREATE ' . ($index['unique'] ? 'UNIQUE ' : '') . 'INDEX ' . self::migrate_identifier($target_driver, $name, '创建目标索引：表=' . $target_table . '，索引=' . $name) . ' ON ' . self::migrate_identifier($target_driver, $target_table, '创建目标索引所属表：' . $target_table) . '(' . implode(',', array_map(fn(string $column): string => self::migrate_identifier($target_driver, $column, '创建目标索引字段：表=' . $target_table . '，索引=' . $name . '，字段=' . $column), $index['columns'])) . ')');
    }
}

public static function migrate_database_identity(PDO $db, array $config): string
{
    if ($config['driver'] === 'sqlite') return (string)realpath((string)$config['path']);
    try {
        if ($config['driver'] === 'mysql') return (string)$db->query("SELECT CONCAT(@@server_uuid,':',DATABASE())")->fetchColumn();
        return (string)$db->query("SELECT current_database()||':'||system_identifier::text FROM pg_control_system()")->fetchColumn();
    } catch (Throwable $e) {
        return '';
    }
}

public static function migrate_same_database(PDO $source_db, array $source, PDO $target_db, array $target): bool
{
    if ($source['driver'] !== $target['driver']) return false;
    $source_identity = self::migrate_database_identity($source_db, $source);
    $target_identity = self::migrate_database_identity($target_db, $target);
    if ($source_identity !== '' && $target_identity !== '') return hash_equals($source_identity, $target_identity);
    return strtolower((string)$source['host']) === strtolower((string)$target['host'])
        && (int)$source['port'] === (int)$target['port']
        && (string)$source['database'] === (string)$target['database'];
}

public static function migrate_core_table_map(): array
{
    return [
        'groups' => 'app_groups',
        'users' => 'app_users',
        'forums' => 'app_forums',
        'topics' => 'app_topics',
        'replies' => 'app_replies',
        'attachments' => 'app_attachments',
        'notifications' => 'app_notifications',
        'settings' => 'app_settings',
        'plugins' => 'app_plugins',
        'cron_tasks' => 'app_cron_tasks',
        'topics_fts' => 'app_topics_fts',
        'replies_fts' => 'app_replies_fts',
    ];
}

public static function migrate_order_tables(array $tables): array
{
    $core = array_flip(array_values(self::migrate_core_table_map()));
    usort($tables, fn(string $a, string $b): int => (($core[self::migrate_target_table($a)] ?? PHP_INT_MAX) <=> ($core[self::migrate_target_table($b)] ?? PHP_INT_MAX)) ?: strcmp($a, $b));
    return $tables;
}

public static function migrate_target_table(string $table): string
{
    return self::migrate_core_table_map()[$table] ?? $table;
}

public static function migrate_value(mixed $value): mixed
{
    if (is_resource($value)) return stream_get_contents($value);
    return is_bool($value) ? (int)$value : $value;
}

public static function migrate_rebuild_search(PDO $db, string $driver): void
{
    if ($driver !== 'sqlite') return;
    if (app_db_table_exists($db, $driver, 'app_topics_fts')) {
        $db->exec('DELETE FROM app_topics_fts');
        $db->exec('INSERT INTO app_topics_fts(rowid,title,body) SELECT id,title,body FROM app_topics');
    }
    if (app_db_table_exists($db, $driver, 'app_replies_fts')) {
        $db->exec('DELETE FROM app_replies_fts');
        $db->exec('INSERT INTO app_replies_fts(rowid,body) SELECT id,body FROM app_replies');
    }
}

public static function migrate_reset_sequences(PDO $db, string $driver, array $tables): void
{
    if ($driver !== 'pgsql') return;
    $sequence = $db->prepare("SELECT pg_get_serial_sequence(?, 'id')");
    $set = $db->prepare('SELECT setval(CAST(? AS regclass),?,?)');
    foreach ($tables as $table) {
        if (!in_array('id', self::migrate_columns($db, $driver, $table), true)) continue;
        $sequence->execute([$table]);
        $name = $sequence->fetchColumn();
        if (!$name) continue;
        $max = (int)$db->query('SELECT COALESCE(MAX(id),0) FROM ' . self::migrate_identifier($driver, $table, '重置序列：表=' . $table))->fetchColumn();
        $set->execute([$name, max(1, $max), $max > 0]);
    }
}

public static function migrate_run(PDO $source, array $source_config): array
{
    $target = db();
    $target_config = db_config();
    if (self::migrate_same_database($source, $source_config, $target, $target_config)) throw new RuntimeException('源数据库不能与当前数据库相同。');
    [$schema_tables] = self::app_db_schema($target_config['driver']);
    $source_tables = [];
    foreach (array_keys($schema_tables) as $target_table) {
        $source_candidates = [$target_table];
        foreach (self::migrate_core_table_map() as $source_table => $mapped_table) if ($mapped_table === $target_table) $source_candidates[] = $source_table;
        foreach (array_unique($source_candidates) as $source_table) if (app_db_table_exists($source, $source_config['driver'], $source_table)) $source_tables[] = $source_table;
    }
    $tables = self::migrate_order_tables($source_tables);
    if (!$tables) throw new RuntimeException('源数据库没有可迁入的系统数据表。');
    $target_tables = [];
    foreach (array_keys($schema_tables) as $table) if (app_db_table_exists($target, $target_config['driver'], $table)) $target_tables[$table] = true;
    $mapped_tables = [];
    foreach ($tables as $source_table) {
        $target_table = self::migrate_target_table($source_table);
        if (isset($mapped_tables[$target_table])) throw new RuntimeException('多个源表映射到同一目标表：' . $mapped_tables[$target_table] . '、' . $source_table);
        $mapped_tables[$target_table] = $source_table;
        if (isset($target_tables[$target_table])) continue;
        self::migrate_install_table($source, $source_config['driver'], $target, $target_config['driver'], $source_table, $target_table);
        $target_tables[$target_table] = true;
    }
    $plans = [];
    foreach ($tables as $source_table) {
        $target_table = self::migrate_target_table($source_table);
        $columns = array_values(array_intersect(self::migrate_columns($target, $target_config['driver'], $target_table), self::migrate_columns($source, $source_config['driver'], $source_table)));
        if ($columns) $plans[] = ['source' => $source_table, 'target' => $target_table, 'columns' => $columns];
    }
    if (!$plans) throw new RuntimeException('没有找到兼容的数据字段。');
    $counts = [];
    if ($source_config['driver'] === 'mysql') $source->exec('SET TRANSACTION ISOLATION LEVEL REPEATABLE READ');
    if ($source_config['driver'] === 'pgsql') $source->exec('BEGIN ISOLATION LEVEL REPEATABLE READ');
    else $source->beginTransaction();
    $target->beginTransaction();
    try {
        foreach (array_reverse($plans) as $plan) $target->exec('DELETE FROM ' . self::migrate_identifier($target_config['driver'], $plan['target'], '清空目标表：' . $plan['target']));
        foreach ($plans as $plan) {
            $source_table = $plan['source'];
            $target_table = $plan['target'];
            $columns = $plan['columns'];
            $source_columns = implode(',', array_map(fn(string $name): string => self::migrate_identifier($source_config['driver'], $name, '读取源字段：表=' . $source_table . '，字段=' . $name), $columns));
            $target_columns = implode(',', array_map(fn(string $name): string => self::migrate_identifier($target_config['driver'], $name, '写入目标字段：表=' . $target_table . '，字段=' . $name), $columns));
            $order = in_array('id', $columns, true) ? ' ORDER BY ' . self::migrate_identifier($source_config['driver'], 'id', '读取源表排序字段：' . $source_table) : '';
            $read = $source->query('SELECT ' . $source_columns . ' FROM ' . self::migrate_identifier($source_config['driver'], $source_table, '读取源表：' . $source_table) . $order);
            $insert = 'INSERT INTO ' . self::migrate_identifier($target_config['driver'], $target_table, '写入目标表：源表=' . $source_table . '，目标表=' . $target_table) . '(' . $target_columns . ') VALUES(' . sql_marks(count($columns)) . ')';
            $insert = $target_config['driver'] === 'mysql' ? str_replace('INSERT INTO', 'INSERT IGNORE INTO', $insert) : $insert . ' ON CONFLICT DO NOTHING';
            $write = $target->prepare($insert);
            $count = 0;
            $attachment_hashes = [];
            while ($row = $read->fetch()) {
                if (in_array($source_table, ['topics', 'app_topics'], true) && array_key_exists('reply_order', $row)) $row['reply_order'] = strtolower(trim((string)$row['reply_order'])) === 'desc' || (string)$row['reply_order'] === '1' ? 1 : 0;
                if (in_array($source_table, ['attachments', 'app_attachments'], true) && isset($row['user_id'], $row['hash'])) {
                    $key = $row['user_id'] . ':' . $row['hash'];
                    if (isset($attachment_hashes[$key])) continue;
                    $attachment_hashes[$key] = true;
                }
                $write->execute(array_map(fn(string $column) => self::migrate_value($row[$column]), $columns));
                if ($write->rowCount() > 0) $count++;
            }
            $counts[$source_table] = $count;
        }
        self::migrate_rebuild_search($target, $target_config['driver']);
        self::migrate_reset_sequences($target, $target_config['driver'], array_column($plans, 'target'));
        $source->commit();
        $target->commit();
    } catch (Throwable $e) {
        if ($source->inTransaction()) $source->rollBack();
        if ($target->inTransaction()) $target->rollBack();
        throw $e;
    }
    return $counts;
}

public static function migrate_refresh_caches(): void
{
    unset($GLOBALS['__settings_cache'], $GLOBALS['__home_stats_cache']);
    forums_cache(true);
    groups_cache(true);
    save_settings_values(['plugin_sync_pending' => '1']);
}

public static function migrate_page(): void
{
    need_admin();
    set_time_limit(0);
    ignore_user_abort(true);
    $error = '';
    $counts = null;
    if ($_SERVER['REQUEST_METHOD'] === 'POST') {
        try {
            if (!isset($_POST['confirm_replace'])) throw new RuntimeException('请确认清空当前新数据库。');
            $source_config = self::migrate_source_config();
            $counts = self::migrate_run(self::migrate_source_db($source_config), $source_config);
        } catch (Throwable $e) {
            $error = $e->getMessage();
        }
    }
    $target = db_config();
    $target_label = $target['driver'] === 'sqlite' ? 'SQLite / ' . basename((string)$target['path']) : strtoupper($target['driver'] === 'pgsql' ? 'PostgreSQL' : 'MySQL') . ' / ' . $target['database'];
    if (is_array($counts)) {
        $rows = '';
        foreach ($counts as $table => $count) $rows .= '<div>' . h($table) . '</div><div>' . (int)$count . '</div>';
        self::migrate_refresh_caches();
        self::setup_html('迁入完成', '<div class="hero"><h1>迁入完成</h1><p>旧数据库数据已写入当前数据库。</p></div><div class="card"><div class="hd"><h2>迁入结果</h2></div><div class="bd"><div class="note ok">共迁入 ' . array_sum($counts) . ' 条数据。</div><div style="height:14px"></div><div class="kv"><div>目标数据库</div><div class="mono">' . h($target_label) . '</div>' . $rows . '</div><div style="height:14px"></div><div class="actions"><a class="btn alt" href="' . h(route_url('update')) . '">返回升级</a><a class="btn" href="' . h(route_url('home')) . '">进入首页</a></div></div></div>');
    }
    $driver = in_array((string)($_POST['source_driver'] ?? 'sqlite'), ['sqlite', 'mysql', 'pgsql'], true) ? (string)($_POST['source_driver'] ?? 'sqlite') : 'sqlite';
    $v = fn(string $name, string $default = ''): string => (string)($_POST[$name] ?? $default);
    $options = '<option value="sqlite"' . ($driver === 'sqlite' ? ' selected' : '') . '>SQLite</option><option value="mysql"' . ($driver === 'mysql' ? ' selected' : '') . '>MySQL</option><option value="pgsql"' . ($driver === 'pgsql' ? ' selected' : '') . '>PostgreSQL</option>';
    $message = $error !== '' ? '<div class="note warn">' . h($error) . '</div>' : '';
    $sqlite_fields = '<div class="db-fields" id="sqlite-fields"><div class="row"><label>SQLite 文件路径</label><input type="text" name="source_sqlite" value="' . h($v('source_sqlite', 'app/data/old.sqlite')) . '"></div></div>';
    $server_fields = '<div class="db-fields" id="server-fields"><div class="row compact"><div class="field"><label>数据库地址</label><input type="text" name="source_host" value="' . h($v('source_host', '127.0.0.1')) . '"></div><div class="field"><label>端口</label><input type="text" name="source_port" value="' . h($v('source_port', $driver === 'pgsql' ? '5432' : '3306')) . '"></div></div><div class="row"><label>数据库名</label><input type="text" name="source_database" value="' . h($v('source_database')) . '"></div><div class="row compact"><div class="field"><label>用户名</label><input type="text" name="source_username" value="' . h($v('source_username')) . '"></div><div class="field"><label>密码</label><input type="password" name="source_password"></div></div></div>';
    $form = '<form class="form" method="post" action="' . h(route_url('migrate')) . '" autocomplete="off">' . form_token() . '<div class="row"><label>旧数据库类型</label><select name="source_driver" id="source-driver">' . $options . '</select></div>' . $sqlite_fields . $server_fields . '<div class="checks"><label class="check"><input type="checkbox" name="confirm_replace" value="1" required><span>确认清空当前数据库中的同名数据表。</span></label></div><div class="actions"><a class="btn alt" href="' . h(route_url('update')) . '">取消</a><button class="btn" type="submit">开始迁入</button></div></form>';
    $body = '<div class="hero"><h1>数据迁入</h1><p>从旧数据库迁入当前已安装数据库。</p></div>' . $message . '<div class="grid"><section class="card"><div class="hd"><h2>旧数据库配置</h2></div><div class="bd">' . $form . '</div></section><aside class="card"><div class="hd"><h2>迁入说明</h2></div><div class="bd"><ul class="list"><li>目标数据库：' . h($target_label) . '</li><li>仅迁入系统定义的数据表</li><li>其他表不读取、不创建、不修改</li><li>缺少的数据表会自动创建</li><li>同名数据表将清空后替换</li><li>附件、头像和插件文件需单独复制</li></ul></div></aside></div><script>const type=document.getElementById("source-driver"),sqlite=document.getElementById("sqlite-fields"),server=document.getElementById("server-fields"),port=document.querySelector("[name=source_port]");function toggle(change){sqlite.hidden=type.value!=="sqlite";server.hidden=type.value==="sqlite";if(change)port.value=type.value==="pgsql"?"5432":"3306"}type.addEventListener("change",()=>toggle(true));toggle(false);</script>';
    self::setup_html('数据迁入', $body);
}
}
