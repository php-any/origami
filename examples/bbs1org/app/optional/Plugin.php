<?php
namespace app\optional;

use RuntimeException;
use Throwable;
use WeakMap;

if (!defined('APP_ROOT')) {
    exit;
}

const PLUGIN_MARKET_ENDPOINT = 'https://bbs1.org/index.php';
const PLUGIN_MARKET_SHARE_MAX = 200000;

final class Plugin
{
private static function plugin_runtime_cache_rows_valid(array $rows): bool
{
    $fields = array_fill_keys(['id', 'file', 'config', 'entries', 'hooks', 'routes', 'admin_tabs', 'assets', 'cron'], true);
    foreach ($rows as $row) {
        if (!is_array($row) || array_diff_key($row, $fields) || array_diff_key($fields, $row)) return false;
        $id = (string)$row['id'];
        $file = str_replace('\\', '/', ltrim((string)$row['file'], '/'));
        if (!plugin_id_valid($id) || $id === 'plugin_market' || $file !== 'app/plugins/' . $id . '/plugin.php') return false;
        foreach (['config', 'entries', 'hooks', 'routes', 'admin_tabs', 'assets', 'cron'] as $field) {
            if (!is_array($row[$field])) return false;
        }
    }
    return true;
}

public static function plugin_runtime_cache_rows(bool $refresh = false): array
{
    $rows = $refresh ? null : json_decode(setting('cache_plugins'), true);
    if (is_array($rows) && self::plugin_runtime_cache_rows_valid($rows)) return $rows;
    $rows = [];
    foreach (q("SELECT id,file,manifest_json,config_json,entries_json FROM app_plugins WHERE enabled=1 ORDER BY id")->fetchAll() as $row) {
        $id = (string)($row['id'] ?? '');
        $file = str_replace('\\', '/', ltrim((string)($row['file'] ?? ''), '/'));
        $manifest = plugin_json_decode($row['manifest_json'] ?? '', null);
        if (!plugin_id_valid($id) || $id === 'plugin_market' || $file !== 'app/plugins/' . $id . '/plugin.php' || $manifest === null) continue;
        $rows[] = [
            'id' => $id,
            'file' => $file,
            'config' => plugin_json_decode($row['config_json'] ?? '') ?? [],
            'entries' => plugin_json_decode($row['entries_json'] ?? '') ?? [],
            'hooks' => is_array($manifest['hooks'] ?? null) ? $manifest['hooks'] : [],
            'routes' => is_array($manifest['routes'] ?? null) ? $manifest['routes'] : [],
            'admin_tabs' => is_array($manifest['admin_tabs'] ?? null) ? $manifest['admin_tabs'] : [],
            'assets' => is_array($manifest['assets'] ?? null) ? $manifest['assets'] : [],
            'cron' => is_array($manifest['cron'] ?? null) ? $manifest['cron'] : [],
        ];
    }
    save_settings_values(['cache_plugins' => plugin_json_encode($rows)]);
    return $rows;
}

public static function plugin_registry(?string $id = null): array
{
    if ($id !== null && !plugin_id_valid($id)) return [];
    $plugins = [];
    $sql = "SELECT id,name,version,file,manifest_json,config_json,entries_json,enabled,disabled_reason,updated_at FROM app_plugins";
    $rows = q($sql . ($id === null ? ' ORDER BY id' : ' WHERE id=?'), $id === null ? [] : [$id])->fetchAll();
    foreach ($rows as $row) {
        if ((string)$row['id'] === 'plugin_market') continue;
        $plugin = plugin_registry_row($row);
        if ($plugin) $plugins[(string)$plugin['id']] = $plugin;
    }
    return $plugins;
}

public static function plugin_market_url(string $action): string
{
    return append_url_query(PLUGIN_MARKET_ENDPOINT, ['a' => $action]);
}

public static function remote_http_request(string $url, int $timeout = 8, array $headers = [], ?array $post_fields = null): array
{
    if (!function_exists('curl_init')) return ['ok' => false, 'status' => 0, 'body' => '', 'error' => '服务器未启用 cURL'];
    $ch = curl_init($url);
    if (!$ch) return ['ok' => false, 'status' => 0, 'body' => '', 'error' => '无法初始化请求'];
    $options = [
        CURLOPT_RETURNTRANSFER => true,
        CURLOPT_FOLLOWLOCATION => true,
        CURLOPT_MAXREDIRS => 3,
        CURLOPT_CONNECTTIMEOUT => min(4, max(1, $timeout)),
        CURLOPT_TIMEOUT => max(1, $timeout),
        CURLOPT_USERAGENT => 'bbs1org/' . APP_VERSION,
    ];
    if (setting('ignore_ssl_errors') === '1') {
        $options[CURLOPT_SSL_VERIFYPEER] = false;
        $options[CURLOPT_SSL_VERIFYHOST] = 0;
    }
    if ($headers) $options[CURLOPT_HTTPHEADER] = $headers;
    if ($post_fields !== null) {
        $options[CURLOPT_POST] = true;
        $options[CURLOPT_POSTFIELDS] = http_build_query($post_fields);
    }
    if (defined('CURLOPT_PROTOCOLS') && defined('CURLPROTO_HTTP') && defined('CURLPROTO_HTTPS')) $options[CURLOPT_PROTOCOLS] = CURLPROTO_HTTP | CURLPROTO_HTTPS;
    if (defined('CURLOPT_REDIR_PROTOCOLS') && defined('CURLPROTO_HTTP') && defined('CURLPROTO_HTTPS')) $options[CURLOPT_REDIR_PROTOCOLS] = CURLPROTO_HTTP | CURLPROTO_HTTPS;
    curl_setopt_array($ch, $options);
    $body = curl_exec($ch);
    $error = curl_error($ch);
    $status = (int)curl_getinfo($ch, CURLINFO_RESPONSE_CODE);
    if ($body === false) return ['ok' => false, 'status' => $status, 'body' => '', 'error' => $error !== '' ? $error : '请求失败'];
    if ($status < 200 || $status >= 300) return ['ok' => false, 'status' => $status, 'body' => (string)$body, 'error' => 'HTTP ' . $status];
    return ['ok' => true, 'status' => $status, 'body' => (string)$body, 'error' => ''];
}

public static function plugin_market_fetch(): array
{
    $response = self::remote_http_request(self::plugin_market_url('plugin_market_feed'), 8, ['Accept: application/json']);
    if (!$response['ok']) return ['ok' => 0, 'message' => '无法连接插件市场' . ((string)$response['error'] !== '' ? '：' . (string)$response['error'] : ''), 'plugins' => []];
    $data = json_decode((string)$response['body'], true);
    if (!is_array($data)) return ['ok' => 0, 'message' => '插件市场返回格式错误', 'plugins' => []];
    $plugins = [];
    foreach ((array)($data['plugins'] ?? []) as $item) {
        if (!is_array($item)) continue;
        $id = (string)($item['id'] ?? '');
        $code = (string)($item['code'] ?? '');
        if (!plugin_id_valid($id) || $id === 'plugin_market' || trim($code) === '') continue;
        $plugins[$id] = [
            'id' => $id,
            'title' => (string)($item['title'] ?? ''),
            'name' => (string)($item['name'] ?? $id),
            'version' => (string)($item['version'] ?? ''),
            'description' => (string)($item['description'] ?? ''),
            'author' => (string)($item['author'] ?? ''),
            'creator' => (string)($item['creator'] ?? ($item['username'] ?? ($item['author'] ?? ''))),
            'creator_id' => (int)($item['creator_id'] ?? 0),
            'topic_id' => (int)($item['topic_id'] ?? 0),
            'updated_at' => (int)($item['updated_at'] ?? 0),
            'sha256' => (string)($item['sha256'] ?? hash('sha256', $code)),
            'url' => clean_site_base_url((string)($item['url'] ?? '')),
            'code' => $code,
        ];
    }
    return ['ok' => (int)($data['ok'] ?? 1), 'message' => (string)($data['message'] ?? ''), 'plugins' => $plugins];
}

public static function plugin_market_install(string $id): void
{
    if (!plugin_id_valid($id)) err('插件不存在');
    if ($id === 'plugin_market') err('该插件 ID 为系统保留');
    require_writable_dir(PLUGIN_DIR, '插件目录不可写，请检查 app/plugins/ 目录权限');
    $market = self::plugin_market_fetch();
    $item = $market['plugins'][$id] ?? null;
    if (!is_array($item)) err((string)($market['message'] ?? '') ?: '插件市场没有返回该插件');
    $code = (string)($item['code'] ?? '');
    if (!str_starts_with(ltrim($code), '<?php')) err('插件代码格式错误');
    if ((string)$item['sha256'] !== '' && !hash_equals((string)$item['sha256'], hash('sha256', $code))) err('插件代码校验失败');
    if (preg_match('/[\'"]id[\'"]\s*=>\s*([\'"])(.*?)\1/s', $code, $match) !== 1 || (string)$match[2] !== $id) err('插件代码 ID 与市场 ID 不一致');
    $dir = PLUGIN_DIR . '/' . $id;
    $file = $dir . '/plugin.php';
    if (!is_dir($dir) && !mkdir($dir, 0755, true)) err('插件目录创建失败');
    require_writable_dir($dir, '插件目录不可写，请检查 app/plugins/ 目录权限');
    if (is_file($file)) {
        $backup_dir = DATA_DIR . '/plugin-backups';
        if (!is_dir($backup_dir) && !mkdir($backup_dir, 0755, true)) err('插件备份目录创建失败');
        if (!copy($file, $backup_dir . '/' . $id . '-' . date('YmdHis') . '.php')) err('现有插件备份失败');
    }
    $tmp = $file . '.tmp.' . bin2hex(random_bytes(4));
    if (file_put_contents($tmp, $code, LOCK_EX) === false) err('插件写入失败');
    if (!rename($tmp, $file)) {
        @unlink($tmp);
        err('插件安装失败');
    }
    if (function_exists('opcache_invalidate')) @opcache_invalidate($file, true);
    plugin_update_row($id, ['enabled' => 0, 'status' => 'disabled', 'disabled_reason' => '']);
    q("UPDATE app_cron_tasks SET enabled=0 WHERE plugin_id=?", [$id]);
    save_settings_values([
        'plugin_' . $id . '_market_sha256' => (string)$item['sha256'],
        'plugin_' . $id . '_market_topic_id' => (string)(int)$item['topic_id'],
        'plugin_sync_pending' => '1',
    ]);
    self::plugin_assets_mark_dirty();
}

public static function plugin_market_install_page(): void
{
    need_admin();
    require_post();
    self::plugin_market_install((string)($_POST['plugin_id'] ?? ''));
    $message = '插件已安装或更新，已自动停用，请启用后使用。';
    if (ajax_request()) {
        header('Content-Type: application/json; charset=utf-8');
        echo json_encode(['ok' => 1, 'message' => $message, 'refresh' => 1], JSON_UNESCAPED_UNICODE);
        exit;
    }
    set_flash($message);
    go(admin_url(['tab' => 'plugins', 'view' => 'market']));
}

public static function plugin_market_share_page(): void
{
    need_admin();
    require_post();
    $id = (string)($_POST['plugin_id'] ?? '');
    $target = plugin_id_valid($id) ? (self::plugin_registry($id)[$id] ?? null) : null;
    if (!is_array($target)) err('插件不存在');
    $file = (string)($target['file'] ?? '');
    $code = $file !== '' && is_file($file) ? file_get_contents($file) : false;
    if (!is_string($code) || trim($code) === '') err('插件文件不存在或为空');
    if (preg_match('/^\s*```\s*[\w-]*\s*$/m', $code)) err('插件代码包含独立的 Markdown 代码块标记，无法安全分享');
    $body = "```php\n" . rtrim($code) . "\n```";
    if (strlen($body) > PLUGIN_MARKET_SHARE_MAX) err('插件代码超过分享长度限制');
    $name = trim((string)($target['name'] ?? '')) ?: $id;
    $share_form = '<form class="post-action-form" method="post" action="' . h(self::plugin_market_url('plugin_share_receive')) . '" data-no-ajax="1" data-plugin-share-auto="1"><input type="hidden" name="title" value="' . h('[' . $id . ']' . $name) . '"><textarea name="body" hidden>' . h($body) . '</textarea><button type="submit" class="plugin-enable">立即继续</button></form>';
    $author = trim((string)($target['author'] ?? ''));
    $head = '<div class="admin-plugin-summary"><strong>正在前往插件市场</strong><span>插件代码已准备好，请在官方站点确认后发布。</span></div>';
    $row = '<li class="admin-list-item admin-object-row plugin-item"><div class="admin-row-main"><div class="plugin-title-line"><strong class="admin-content-title">' . h($name) . '</strong><span class="admin-flag on">分享</span></div><div class="admin-row-meta"><span class="plugin-id">ID ' . h($id) . '</span>' . ((string)($target['version'] ?? '') !== '' ? '<span>版本 ' . h((string)$target['version']) . '</span>' : '') . ($author !== '' ? '<span>插件作者 ' . h($author) . '</span>' : '') . '</div><div class="admin-content-text plugin-desc">' . h((string)($target['description'] ?? '')) . '</div><div class="plugin-file">' . h($file) . '</div></div><div class="admin-inline-ops plugin-ops">' . $share_form . '</div></li>';
    page('分享插件', shell_html('<div class="admin-list-panel plugin-list-panel">' . admin_list_head($head, '') . '<ul class="admin-manage-list plugin-list">' . $row . '</ul></div>', sidebar_stack_html([sidebar_user_card_html()])));
}

public static function plugin_market_update_available(array $plugin, array $item): bool
{
    $remote = trim((string)($item['version'] ?? ''));
    $local = trim((string)($plugin['version'] ?? ''));
    return $remote !== '' && $local !== '' && version_compare($remote, $local, '>');
}

public static function plugin_market_search_form(string $query): string
{
    $hidden = hidden_inputs(['a' => 'admin', 'tab' => 'plugins', 'view' => 'market']);
    $url = admin_url(['tab' => 'plugins', 'view' => 'market']);
    $clear = $query !== '' ? '<a class="admin-search-clear" href="' . h($url) . '">清空</a>' : '';
    return '<form class="admin-table-search" method="get" action="' . h(index_url()) . '">' . $hidden . '<div class="admin-search-field"><input name="q" value="' . h($query) . '" placeholder="搜索标题 / 插件ID / 制作者" minlength="' . search_min_chars() . '"><button class="admin-search-submit" type="submit">搜索</button></div>' . $clear . '</form>';
}

public static function plugin_market_matches(array $item, string $query): bool
{
    if ($query === '') return true;
    foreach (['title', 'name', 'id', 'creator', 'description'] as $key) if (stripos((string)($item[$key] ?? ''), $query) !== false) return true;
    return false;
}

public static function plugin_market_page_html(bool $with_tabs = true): string
{
    $market = self::plugin_market_fetch();
    $items = is_array($market['plugins'] ?? null) ? $market['plugins'] : [];
    $local = self::plugin_registry();
    $updates = [];
    foreach ($items as $id => $item) if (isset($local[$id]) && is_array($item) && self::plugin_market_update_available($local[$id], $item)) $updates[$id] = true;
    uksort($items, fn(string $a, string $b): int => (int)isset($updates[$b]) <=> (int)isset($updates[$a]));
    $query = trim((string)($_GET['q'] ?? ''));
    $url = admin_url(['tab' => 'plugins', 'view' => 'market']);
    $head = '<div class="admin-plugin-summary"><strong>插件市场</strong><span>仅展示官方审核通过的插件，安装后默认仍需手动启用。</span></div>';
    $actions = '<div class="plugin-head-actions">' . self::plugin_market_search_form($query) . '<a class="admin-search-clear" href="' . h($url) . '">刷新</a></div>';
    $html = ($with_tabs ? self::admin_plugins_tabs_html('market') : '') . '<div class="admin-list-panel plugin-list-panel">' . admin_list_head($head, $actions) . '<ul class="admin-manage-list plugin-list">';
    if (!(int)($market['ok'] ?? 0)) return $html . '<li class="empty-state">' . h((string)($market['message'] ?? '插件市场暂不可用')) . '</li></ul></div>';
    $shown = 0;
    foreach ($items as $item) {
        if (!is_array($item)) continue;
        $id = (string)($item['id'] ?? '');
        if (!plugin_id_valid($id) || !self::plugin_market_matches($item, $query)) continue;
        $shown++;
        $installed = isset($local[$id]);
        $needs_update = isset($updates[$id]);
        $label = $installed ? ($needs_update ? '更新' : '重新安装') : '安装';
        $button_class = $installed && !$needs_update ? '' : 'plugin-enable';
        $ops = '<form class="post-action-form" method="post" action="' . h(route_url('plugin_market_install')) . '" data-replace-target=".plugin-list-panel" data-confirm="确定' . h($label) . '该插件？插件代码将写入本地 plugins 目录，完成后插件会自动停用。">' . form_token() . hidden_inputs(['plugin_id' => $id]) . '<button type="submit" data-loading-text="安装中"' . ($button_class !== '' ? ' class="' . h($button_class) . '"' : '') . '>' . h($label) . '</button></form>';
        $meta = [];
        if ((string)($item['version'] ?? '') !== '') $meta[] = '版本 ' . (string)$item['version'];
        $creator = trim((string)($item['creator'] ?? ''));
        $meta[] = '插件制作者 ' . ($creator !== '' ? $creator : '未声明');
        if ((int)($item['updated_at'] ?? 0) > 0) $meta[] = date('Y-m-d H:i', (int)$item['updated_at']);
        $topic_url = (string)($item['url'] ?? '');
        $title = h((string)($item['name'] ?? $id));
        $title = $topic_url !== '' ? '<a class="admin-content-title" href="' . h($topic_url) . '" target="_blank" rel="noopener">' . $title . '</a>' : '<strong class="admin-content-title">' . $title . '</strong>';
        $flag = $installed ? '<span class="admin-flag ' . ($needs_update ? 'update' : 'on') . '">' . h($needs_update ? '可更新' : '已安装') . '</span>' : '<span class="admin-flag">未安装</span>';
        $class = $needs_update ? ' plugin-market-update plugin-update-item' : ($installed ? ' plugin-market-installed' : ' plugin-market-available');
        $html .= '<li class="admin-list-item admin-object-row plugin-item' . $class . '"><div class="admin-row-main"><div class="plugin-title-line">' . $title . $flag . '</div><div class="admin-row-meta"><span class="plugin-id">ID ' . h($id) . '</span><span>' . h(implode(' / ', $meta)) . '</span></div><div class="admin-content-text plugin-desc">' . h((string)($item['description'] ?? '')) . '</div><div class="plugin-file">' . h(substr((string)($item['sha256'] ?? ''), 0, 16)) . '</div></div><div class="admin-inline-ops plugin-ops">' . $ops . '</div></li>';
    }
    if ($shown === 0) $html .= '<li class="empty-state">' . h($query !== '' ? '没有匹配的插件。' : '暂无已审核通过的插件。') . '</li>';
    return $html . '</ul></div>';
}

public static function plugin_market_admin_actions(array $plugin): string
{
    $id = (string)($plugin['id'] ?? '');
    if (!plugin_id_valid($id)) return '';
    return '<form class="post-action-form" method="post" action="' . h(route_url('plugin_market_share')) . '" target="_blank" rel="noopener" data-no-ajax="1">' . form_token() . hidden_inputs(['plugin_id' => $id]) . '<button type="submit">分享</button></form>';
}

public static function admin_plugin_action_form(string $id, string $action, string $label, string $class = '', string $confirm = ''): string
{
    $confirm_attr = $confirm !== '' ? ' data-confirm="' . h($confirm) . '"' : '';
    return '<form class="post-action-form" method="post" action="' . h(admin_url(['tab' => 'plugins'])) . '" data-replace-target=".plugin-list-panel"' . $confirm_attr . '>' . form_token() . hidden_inputs(['plugin_id' => $id, 'plugin_action' => $action]) . '<button type="submit"' . ($class !== '' ? ' class="' . h($class) . '"' : '') . '>' . h($label) . '</button></form>';
}
public static function admin_plugin_uninstall_form(string $id): string
{
    return '<form class="post-action-form" method="post" action="' . h(admin_url(['tab' => 'plugins'])) . '" data-plugin-uninstall="1" data-replace-target=".plugin-list-panel" data-confirm="确定卸载插件？插件目录将被永久删除。">' . form_token() . hidden_inputs(['plugin_id' => $id, 'plugin_action' => 'uninstall', 'keep_plugin_data' => '1']) . '<button type="submit" class="danger">卸载</button></form>';
}
public static function admin_plugin_entry_toggle_form(array $plugin, string $entry, string $label): string
{
    if (!plugin_uses_entry($plugin, $entry)) return '';
    $id = (string)$plugin['id'];
    $checked = plugin_entry_enabled($plugin, $entry);
    return '<form class="post-action-form plugin-entry-form" method="post" action="' . h(admin_url(['tab' => 'plugins'])) . '" data-replace-target=".plugin-list-panel">' . form_token() . hidden_inputs(['plugin_id' => $id, 'plugin_action' => 'entry_toggle', 'entry' => $entry, 'entry_enabled' => '0']) . '<label class="plugin-entry-check"><input type="checkbox" name="entry_enabled" value="1" data-auto-submit' . ($checked ? ' checked' : '') . '><span>' . h($label) . '</span></label></form>';
}
public static function admin_plugins_page_html(bool $with_tabs = true): string
{
    $plugins = self::plugin_registry();
    uasort($plugins, function (array $a, array $b): int {
        $a_time = (int)($a['updated_at'] ?? 0);
        $b_time = (int)($b['updated_at'] ?? 0);
        return ($b_time <=> $a_time) ?: strcmp((string)($a['id'] ?? ''), (string)($b['id'] ?? ''));
    });
    $enabled_count = 0;
    foreach ($plugins as $plugin) if (plugin_enabled($plugin)) $enabled_count++;
    $head_left = '<div class="admin-plugin-summary"><strong>插件</strong><span>已发现 ' . count($plugins) . ' 个，已启用 ' . $enabled_count . ' 个</span></div>';
    $head_right = self::admin_plugin_action_form('', 'sync', '同步插件');
    $html = ($with_tabs ? self::admin_plugins_tabs_html('local') : '') . '<div class="admin-list-panel plugin-list-panel">' . admin_list_head($head_left, $head_right) . '<ul class="admin-manage-list plugin-list">';
    foreach ($plugins as $plugin) {
        $id = (string)$plugin['id'];
        $enabled = plugin_enabled($plugin);
        $manage_url = '';
        if ($enabled && !empty($plugin['admin_tabs']) && is_array($plugin['admin_tabs'])) {
            foreach ($plugin['admin_tabs'] as $key => $fn) {
                if (is_string($key) && is_string($fn)) {
                    $manage_url = admin_url(['tab' => $key]);
                    break;
                }
            }
        }
        $ops = $manage_url !== '' ? '<a class="plugin-manage-link" href="' . h($manage_url) . '">管理</a>' : '';
        $ops .= $enabled
            ? self::admin_plugin_action_form($id, 'disable', '停用', 'danger', '确定停用插件？')
            : self::admin_plugin_action_form($id, 'enable', '启用', 'plugin-enable');
        $entry_ops = self::admin_plugin_entry_toggle_form($plugin, 'feature_links', '快捷功能');
        $entry_ops .= self::admin_plugin_entry_toggle_form($plugin, 'sidebar_cards', '边栏卡片');
        $ops .= self::plugin_market_admin_actions($plugin);
        $ops .= (string)hook('admin.plugin.actions', '', ['plugin' => $plugin]);
        $ops .= self::admin_plugin_uninstall_form($id);
        $meta = [];
        if ((string)($plugin['version'] ?? '') !== '') $meta[] = '版本 ' . (string)$plugin['version'];
        if ((string)($plugin['author'] ?? '') !== '') $meta[] = (string)$plugin['author'];
        $features = [];
        if (!empty($plugin['hooks'])) $features[] = count($plugin['hooks']) . ' 个钩子';
        if (!empty($plugin['routes'])) $features[] = count($plugin['routes']) . ' 个路由';
        if (!empty($plugin['admin_tabs'])) $features[] = count($plugin['admin_tabs']) . ' 个后台页';
        $file = str_replace(APP_ROOT . '/', '', (string)($plugin['file'] ?? ''));
        $disabled_reason = !$enabled ? trim((string)($plugin['disabled_reason'] ?? '')) : '';
        $reason_line = $disabled_reason !== '' ? '<div class="plugin-disabled-reason"><strong>自动停用原因</strong><span>' . h($disabled_reason) . '</span></div>' : '';
        $entry_line = $entry_ops !== '' ? '<div class="plugin-entry-line"><span class="plugin-entry-label">展示位置</span><div class="plugin-entry-options">' . $entry_ops . '</div></div>' : '';
        $local_class = $enabled ? ' plugin-local-enabled' : ' plugin-local-disabled';
        $title = $manage_url !== '' ? '<a class="admin-content-title" href="' . h($manage_url) . '">' . h((string)$plugin['name']) . '</a>' : '<strong class="admin-content-title">' . h((string)$plugin['name']) . '</strong>';
        $html .= '<li class="admin-list-item admin-object-row plugin-item' . $local_class . '"><div class="admin-row-main"><div class="plugin-title-line">' . $title . '<span class="admin-flag' . ($enabled ? ' on' : '') . '">' . h($enabled ? '已启用' : '已停用') . '</span></div><div class="admin-row-meta"><span class="plugin-id">ID ' . h($id) . '</span>' . ($meta ? '<span>' . h(implode(' / ', $meta)) . '</span>' : '') . ($features ? '<span>' . h(implode(' / ', $features)) . '</span>' : '') . '</div><div class="admin-content-text plugin-desc">' . h((string)($plugin['description'] ?? '')) . '</div>' . $reason_line . '<div class="plugin-file">' . h($file) . '</div></div>' . $entry_line . '<div class="admin-inline-ops plugin-ops">' . $ops . '</div></li>';
    }
    if (!$plugins) $html .= '<li class="empty-state">暂无插件，放入 app/plugins/*/plugin.php 后点击“同步插件”。</li>';
    return $html . '</ul></div>';
}
public static function admin_plugins_tabs_html(string $active): string
{
    $items = [
        'local' => ['label' => '本地插件', 'href' => admin_url(['tab' => 'plugins'])],
        'market' => ['label' => '插件市场', 'href' => admin_url(['tab' => 'plugins', 'view' => 'market'])],
    ];
    $hook_items = hook('admin.plugins.tabs', $items, ['active' => $active]);
    if (is_array($hook_items)) $items = $hook_items;
    $items['cron'] = ['label' => '计划任务日志', 'href' => admin_url(['tab' => 'plugins', 'view' => 'cron'])];
    return tab_bar_html($items, $active, 'plugin-tabs');
}
public static function admin_plugins_cron_logs_page_html(): string
{
    $size = 50;
    $page = max(1, (int)($_GET['p'] ?? 1));
    $offset = ($page - 1) * $size;
    $names = [];
    foreach (self::plugin_registry() as $plugin) {
        $names[(string)$plugin['id']] = (string)($plugin['name'] ?? $plugin['id']);
    }
    $rows = q("SELECT plugin_id,task_name,status,message,started_at,finished_at FROM app_cron_logs ORDER BY started_at DESC,id DESC LIMIT ? OFFSET ?", [$size + 1, $offset])->fetchAll();
    $has_next = count($rows) > $size;
    if ($has_next) array_pop($rows);
    $labels = ['success' => '成功', 'failed' => '失败', 'running' => '运行中'];
    $html = self::admin_plugins_tabs_html('cron') . '<div class="admin-list-panel plugin-list-panel">' . admin_list_head('<div class="admin-plugin-summary"><strong>计划任务日志</strong><span>最新运行记录</span></div>', '') . '<ul class="admin-manage-list">';
    foreach ($rows as $row) {
        $status = (string)$row['status'];
        $class = $status === 'success' ? ' on' : ($status === 'failed' ? ' danger' : '');
        $started_at = (int)$row['started_at'];
        $finished_at = (int)$row['finished_at'];
        $duration = $finished_at > 0 ? max(0, $finished_at - $started_at) . ' 秒' : '进行中';
        $message = trim((string)$row['message']);
        $plugin_id = (string)$row['plugin_id'];
        $plugin_name = $names[$plugin_id] ?? $plugin_id;
        $html .= '<li class="admin-list-item"><div class="admin-row-main"><div class="plugin-title-line"><strong class="admin-content-title">' . h($plugin_name) . '</strong><span class="admin-flag' . $class . '">' . h($labels[$status] ?? $status) . '</span></div><div class="admin-row-meta"><span class="plugin-id">' . h($plugin_id) . ' / ' . h((string)$row['task_name']) . '</span><span>' . date('Y-m-d H:i:s', $started_at) . '</span><span>' . h($duration) . '</span>' . ($message !== '' ? '<span title="' . h($message) . '">' . h(cut($message, 160)) . '</span>' : '') . '</div></div></li>';
    }
    if (!$rows) $html .= '<li class="empty-state">暂无计划任务运行记录</li>';
    $html .= '</ul></div>';
    $pagination = simple_paginate($page > 1, $has_next, $page, admin_url(['tab' => 'plugins', 'view' => 'cron']));
    return $html . ($pagination === '' ? '' : '<div class="pagination-bar">' . $pagination . '</div>');
}
public static function admin_plugin_tab_html(string $tab): ?string
{
    foreach (plugins() as $plugin) {
        if (!plugin_enabled($plugin)) continue;
        $fn = $plugin['admin_tabs'][$tab] ?? null;
        if ($fn !== null) {
            $html = plugin_call($plugin, fn(): ?string => plugin_callback_exists($fn) ? (string)$fn($plugin) : null);
            if ($html !== null) return $html;
        }
    }
    return null;
}

public static function plugin_disable_after_exception(string $id, Throwable $e): void
{
    if (!plugin_id_valid($id)) return;
    static $handled;
    $handled ??= new WeakMap();
    if (isset($handled[$e])) return;
    $handled[$e] = true;
    $message = trim((string)preg_replace('/\s+/', ' ', $e->getMessage()));
    $file = str_replace('\\', '/', $e->getFile());
    $root = rtrim(str_replace('\\', '/', APP_ROOT), '/') . '/';
    if (str_starts_with($file, $root)) $file = substr($file, strlen($root));
    $reason = date('Y-m-d H:i:s') . ' ' . get_class($e) . ($message !== '' ? ': ' . cut($message, 500) : '') . ($file !== '' ? ' (' . $file . ':' . $e->getLine() . ')' : '');
    try {
        plugin_update_row($id, ['enabled' => 0, 'status' => 'error', 'disabled_reason' => $reason]);
        q("UPDATE app_cron_tasks SET enabled=0 WHERE plugin_id=?", [$id]);
        plugins(true);
        self::plugin_assets_mark_dirty();
    } catch (Throwable $disable_error) {
        debug_log_write('插件 ' . $id . ' 自动停用失败', $disable_error);
        debug_log_write('插件 ' . $id . ' 运行异常', $e);
        return;
    }
    debug_log_write('插件 ' . $id . ' 运行异常，已自动停用', $e);
}

public static function plugin_manifest_validate(array $plugin, string $file): ?array
{
    $id = (string)($plugin['id'] ?? '');
    if (!plugin_id_valid($id) || $id !== basename(dirname($file))) return null;
    $base = [
        'id' => $id,
        'name' => (string)($plugin['name'] ?? $id),
        'version' => (string)($plugin['version'] ?? ''),
        'description' => (string)($plugin['description'] ?? ''),
        'author' => (string)($plugin['author'] ?? ''),
        'enabled' => !empty($plugin['enabled']),
        'hooks' => is_array($plugin['hooks'] ?? null) ? $plugin['hooks'] : [],
        'routes' => is_array($plugin['routes'] ?? null) ? $plugin['routes'] : [],
        'admin_tabs' => is_array($plugin['admin_tabs'] ?? null) ? $plugin['admin_tabs'] : [],
        'assets' => is_array($plugin['assets'] ?? null) ? $plugin['assets'] : [],
        'cron' => is_array($plugin['cron'] ?? null) ? $plugin['cron'] : [],
        'install' => (string)($plugin['install'] ?? ''),
        'uninstall' => (string)($plugin['uninstall'] ?? ''),
        'file' => $file,
    ];
    foreach (['hooks', 'routes', 'admin_tabs'] as $map) {
        $items = [];
        foreach ($base[$map] as $name => $fn) if (is_string($name) && is_string($fn) && preg_match('/^[A-Za-z_][A-Za-z0-9_]*$/', $fn)) $items[$name] = $fn;
        $base[$map] = $items;
    }
    $assets = [];
    foreach (['css', 'js'] as $type) {
        $fn = $base['assets'][$type] ?? null;
        if (is_string($fn) && preg_match('/^[A-Za-z_][A-Za-z0-9_]*$/', $fn)) $assets[$type] = $fn;
    }
    $base['assets'] = $assets;
    $cron = [];
    foreach ($base['cron'] as $name => $task) {
        if (!is_string($name) || preg_match('/^[a-z0-9][a-z0-9_-]{0,63}$/', $name) !== 1 || !is_array($task)) continue;
        $callback = (string)($task['callback'] ?? '');
        $interval = $task['interval'] ?? 0;
        if (preg_match('/^[A-Za-z_][A-Za-z0-9_]*$/', $callback) !== 1) continue;
        if (is_string($interval) && !is_numeric($interval)) {
            if (preg_match('/^[A-Za-z_][A-Za-z0-9_]*$/', $interval) !== 1) continue;
        } else {
            $interval = (int)$interval;
            if ($interval < 60) continue;
            $interval = min(31536000, $interval);
        }
        $cron[$name] = ['callback' => $callback, 'interval' => $interval];
    }
    $base['cron'] = $cron;
    foreach (['install', 'uninstall'] as $key) if ($base[$key] !== '' && preg_match('/^[A-Za-z_][A-Za-z0-9_]*$/', $base[$key]) !== 1) $base[$key] = '';
    return $base;
}

public static function plugin_files(): array
{
    $files = glob(PLUGIN_DIR . '/*/plugin.php') ?: [];
    sort($files);
    return array_values(array_filter($files, 'is_file'));
}

public static function plugin_assets_mark_dirty(): void
{
    save_settings_values(['plugin_assets_dirty' => '1']);
}

public static function plugin_asset_write(string $file, string $content): void
{
    $tmp = $file . '.tmp.' . bin2hex(random_bytes(4));
    if (is_writable(dirname($file)) && file_put_contents($tmp, $content, LOCK_EX) !== false) {
        if (@rename($tmp, $file)) return;
        @unlink($tmp);
    }
    if (file_put_contents($file, $content, LOCK_EX) === false) throw new RuntimeException('插件资源文件不可写：' . basename($file));
}

public static function plugin_assets_rebuild(): array
{
    $chunks = ['css' => [], 'js' => []];
    foreach (plugins() as $plugin) {
        if (!plugin_enabled($plugin) || empty($plugin['assets'])) continue;
        foreach (['css', 'js'] as $type) {
            $fn = $plugin['assets'][$type] ?? null;
            if (!is_string($fn)) continue;
            $content = trim((string)plugin_call($plugin, function () use ($fn): string {
                return plugin_callback_exists($fn) ? (string)$fn() : '';
            }));
            if ($content === '') continue;
            $chunk = '/* ' . (string)$plugin['id'] . " */\n" . $content;
            $chunks[$type][] = $chunk;
        }
    }
    $files = ['css' => PLUGIN_CSS_FILE, 'js' => PLUGIN_JS_FILE];
    $manifest = [];
    foreach ($files as $key => $file) {
        $content = implode("\n", $chunks[$key]) . ($chunks[$key] ? "\n" : '');
        self::plugin_asset_write($file, $content);
        $manifest[$key] = hash('sha256', $content);
        $manifest[$key . '_size'] = strlen($content);
    }
    save_settings_values([
        'plugin_assets_css_hash' => (string)($manifest['css'] ?? ''),
        'plugin_assets_css_size' => (string)(int)($manifest['css_size'] ?? 0),
        'plugin_assets_js_hash' => (string)($manifest['js'] ?? ''),
        'plugin_assets_js_size' => (string)(int)($manifest['js_size'] ?? 0),
        'plugin_assets_dirty' => '0',
    ]);
    return $manifest;
}

public static function plugin_registry_sync(): array
{
    $existing = [];
    foreach (q("SELECT * FROM app_plugins")->fetchAll() as $row) $existing[(string)$row['id']] = $row;
    $synced = [];
    $disable = function (string $id, string $reason) use (&$existing, &$synced): void {
        if (!isset($existing[$id])) return;
        plugin_update_row($id, ['enabled' => 0, 'status' => 'error', 'disabled_reason' => $reason]);
        q("UPDATE app_cron_tasks SET enabled=0 WHERE plugin_id=?", [$id]);
        $synced[$id] = true;
    };
    foreach (self::plugin_files() as $file) {
        $id = basename(dirname($file));
        if ($id === 'plugin_market') continue;
        try {
            if (function_exists('opcache_invalidate')) @opcache_invalidate($file, true);
            $raw = include $file;
        } catch (Throwable $e) {
            $disable($id, cut($e->getMessage(), 500));
            continue;
        }
        $GLOBALS['__plugin_raw'][$file] = $raw;
        if (!is_array($raw)) {
            $disable($id, '插件定义格式无效');
            continue;
        }
        $plugin = self::plugin_manifest_validate($raw, $file);
        if (!$plugin) {
            $disable($id, '插件定义校验失败');
            continue;
        }
        $id = (string)$plugin['id'];
        $old = $existing[$id] ?? [];
        $installed_at = (int)($old['installed_at'] ?? 0) ?: now();
        $code_hash = hash_file('sha256', $file) ?: '';
        $updated_at = isset($old['updated_at']) && (string)($old['code_hash'] ?? '') === $code_hash
            ? (int)$old['updated_at']
            : now();
        $config = plugin_json_decode($old['config_json'] ?? '') ?? [];
        $entries = isset($old['entries_json']) ? plugin_json_decode($old['entries_json']) ?? [] : [
            'feature_links' => true,
            'sidebar_cards' => true,
        ];
        $enabled = isset($old['enabled']) ? (int)$old['enabled'] : 0;
        app_db_upsert('app_plugins', [
            'id' => $id,
            'name' => (string)$plugin['name'],
            'version' => (string)$plugin['version'],
            'file' => ltrim(str_replace(APP_ROOT, '', $file), '/'),
            'code_hash' => $code_hash,
            'manifest_json' => plugin_json_encode(array_intersect_key($plugin, array_flip(['description', 'author', 'hooks', 'routes', 'admin_tabs', 'assets', 'cron', 'install', 'uninstall']))),
            'config_json' => plugin_json_encode($config),
            'entries_json' => plugin_json_encode($entries),
            'enabled' => $enabled,
            'status' => $enabled ? 'enabled' : 'disabled',
            'disabled_reason' => (string)($old['disabled_reason'] ?? ''),
            'installed_at' => $installed_at,
            'updated_at' => $updated_at,
        ], ['id']);
        $synced[$id] = true;
    }
    foreach (array_keys($existing) as $id) {
        if (isset($synced[$id])) continue;
        q("DELETE FROM app_cron_tasks WHERE plugin_id=?", [$id]);
        q("DELETE FROM app_plugins WHERE id=?", [$id]);
    }
    $plugins = self::plugin_registry();
    plugins(true);
    foreach ($plugins as $plugin) Cron::plugin_cron_sync($plugin);
    return $plugins;
}

public static function admin_plugins_handle_post(): void
{
    $plugin_action = (string)($_POST['plugin_action'] ?? '');
    $plugin_id = (string)($_POST['plugin_id'] ?? '');
    $message = '';
    $refresh_after_response = false;
    if ($plugin_action === 'sync') {
        save_settings_values(['plugin_sync_pending' => '1']);
        $message = '插件已同步';
        $refresh_after_response = true;
    } elseif ($plugin_action === 'enable') {
        self::plugin_set_enabled($plugin_id, true);
        $message = '插件已启用';
    } elseif ($plugin_action === 'disable') {
        self::plugin_set_enabled($plugin_id, false);
        $message = '插件已停用';
    } elseif ($plugin_action === 'uninstall') {
        $keep_data = (string)($_POST['keep_plugin_data'] ?? '1') === '1';
        self::plugin_uninstall($plugin_id, $keep_data);
        $message = $keep_data ? '插件已卸载，目录已删除，数据已保留' : '插件已卸载，目录和数据已删除';
    } elseif ($plugin_action === 'entry_toggle') {
        self::plugin_set_entry_enabled($plugin_id, (string)($_POST['entry'] ?? ''), (string)($_POST['entry_enabled'] ?? '0') === '1');
        $message = '插件入口显示已更新';
    } else err('参数错误');
    $view = (string)($_GET['view'] ?? '');
    if (ajax_request()) {
        if ($refresh_after_response) {
            header('Content-Type: application/json; charset=utf-8');
            echo json_encode(['ok' => 1, 'message' => $message, 'refresh' => 1], JSON_UNESCAPED_UNICODE);
            exit;
        }
        $html = self::admin_plugins_page_html(false);
        header('Content-Type: application/json; charset=utf-8');
        echo json_encode(['ok' => 1, 'message' => $message, 'html' => $html], JSON_UNESCAPED_UNICODE);
        exit;
    }
    set_flash($message);
    go(admin_url(['tab' => 'plugins', 'view' => $view === 'cron' ? $view : null]));
}

public static function plugin_set_entry_enabled(string $id, string $entry, bool $enabled): void
{
    if (!plugin_id_valid($id) || plugin_entry_hook_name($entry) === '') err('参数错误');
    $plugin = self::plugin_registry($id)[$id] ?? null;
    if (!$plugin || !plugin_uses_entry($plugin, $entry)) err('插件未使用该入口');
    $entries = (array)($plugin['entries'] ?? []);
    $entries[$entry] = $enabled;
    plugin_update_row($id, ['entries_json' => $entries]);
    plugins(true);
}

public static function plugin_set_enabled(string $id, bool $enabled): void
{
    if (!plugin_id_valid($id)) err('插件不存在');
    $plugin = self::plugin_registry($id)[$id] ?? null;
    if (!$plugin) err('插件不存在');
    if ($enabled) {
        plugin_call($plugin, function () use ($plugin): void {
            if (!plugin_enabled($plugin) && plugin_callback_exists($plugin['install'] ?? null)) {
                call_user_func((string)$plugin['install'], $plugin);
            }
        });
    }
    plugin_update_row($id, ['enabled' => $enabled ? 1 : 0, 'status' => $enabled ? 'enabled' : 'disabled', 'disabled_reason' => '']);
    q("UPDATE app_cron_tasks SET enabled=? WHERE plugin_id=?", [$enabled ? 1 : 0, $id]);
    $runtime_plugins = plugins(true);
    if ($enabled && isset($runtime_plugins[$id])) Cron::plugin_cron_sync($runtime_plugins[$id]);
    self::plugin_assets_mark_dirty();
}

public static function plugin_uninstall(string $id, bool $keep_data = true): void
{
    if (!plugin_id_valid($id)) err('插件不存在');
    $plugin = self::plugin_registry($id)[$id] ?? null;
    if (!$plugin) err('插件不存在');
    $dir = rtrim(str_replace('\\', '/', PLUGIN_DIR), '/') . '/' . $id;
    if (str_replace('\\', '/', (string)($plugin['file'] ?? '')) !== $dir . '/plugin.php') err('插件目录无效');
    if (!self::plugin_directory_removable($dir)) err('插件目录不可删除，请检查目录权限');
    if (!$keep_data) {
        plugin_call($plugin, function () use ($plugin, $id): void {
            $fn = (string)($plugin['uninstall'] ?? '');
            if (!plugin_callback_exists($fn)) $fn = str_replace('-', '_', $id) . '_uninstall';
            if (plugin_callback_exists($fn)) call_user_func($fn, $plugin);
        });
    }
    self::plugin_remove_directory($dir);
    q("DELETE FROM app_cron_tasks WHERE plugin_id=?", [$id]);
    q("DELETE FROM app_plugins WHERE id=?", [$id]);
    plugins(true);
    plugin_runtime_cache_reset();
    self::plugin_assets_mark_dirty();
}

private static function plugin_directory_removable(string $dir): bool
{
    if (!file_exists($dir) && !is_link($dir)) return true;
    if (is_link($dir)) return is_writable(dirname($dir));
    if (!is_dir($dir)) return false;
    if (!is_readable($dir) || !is_writable($dir)) return false;
    $items = scandir($dir);
    if ($items === false) return false;
    foreach ($items as $item) {
        if ($item === '.' || $item === '..') continue;
        $path = $dir . '/' . $item;
        if (is_dir($path) && !is_link($path) && !self::plugin_directory_removable($path)) return false;
    }
    return is_writable(dirname($dir));
}

private static function plugin_remove_directory(string $dir): void
{
    if (is_link($dir)) {
        if (!unlink($dir)) throw new RuntimeException('无法删除插件目录');
        return;
    }
    if (!is_dir($dir)) return;
    $items = scandir($dir);
    if ($items === false) throw new RuntimeException('无法读取插件目录');
    foreach ($items as $item) {
        if ($item === '.' || $item === '..') continue;
        $path = $dir . '/' . $item;
        if (is_dir($path) && !is_link($path)) {
            self::plugin_remove_directory($path);
        } elseif (!unlink($path)) {
            throw new RuntimeException('无法删除插件文件：' . $item);
        }
    }
    if (!rmdir($dir)) throw new RuntimeException('无法删除插件目录');
}
}
