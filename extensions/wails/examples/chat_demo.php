<?php
/**
 * Wails v3 聊天应用示例 — 完整前后端分离 Demo
 *
 * 演示一个类 Slack/Discord 的多频道聊天应用:
 *   - 前端目录 (examples/chat/) 提供完整 UI
 *   - PHP 后端维护频道、消息、在线用户状态
 *   - 双向事件通信 (chat:join / chat:send / chat:switch ↔ chat:state / chat:message)
 *   - 斜杠命令 (/help /nick /me /clear /channels)
 *   - 机器人 Origami Bot 自动回复
 *
 * 运行:
 *   go run ./cmd examples/chat_demo.php
 */

use Wails\Application;
use Wails\Options\App;
use Wails\Runtime\Window;
use Wails\Runtime\Events;
use Wails\Runtime\Log;

// ══════════════════════════════════════════════════════════════
//  后端状态
// ══════════════════════════════════════════════════════════════

$nextMsgId = 1;
$nextUserId = 1;

// 频道定义
$channels = [
    'general' => ['id' => 'general', 'name' => '综合',   'topic' => '欢迎新人，随便聊聊', 'count' => 0],
    'random'  => ['id' => 'random',  'name' => '随机',   'topic' => '什么都行，开心就好', 'count' => 0],
    'dev'     => ['id' => 'dev',     'name' => '开发',   'topic' => 'Origami + Wails 技术讨论', 'count' => 0],
];

// 消息存储: channelId => [ msg, msg, ... ]
$messages = [
    'general' => [],
    'random'  => [],
    'dev'     => [],
];

// 在线用户: userId => ['id'=>, 'name'=>, 'type'=>'user'|'bot', 'channel'=>]
$onlineUsers = [];

// 当前用户（本客户端）
$myNick    = '访客';
$myChannel = 'general';

// ── 辅助函数 ──

function makeId(&$counter) {
    return 'm' . ($counter++);
}

function makeUserId(&$counter) {
    return 'u' . ($counter++);
}

function addMessage(&$messages, &$channels, $channel, $author, $text, $type, $time) {
    global $nextMsgId;
    $msg = [
        'id'      => makeId($nextMsgId),
        'channel' => $channel,
        'author'  => $author,
        'text'    => $text,
        'type'    => $type,   // user | bot | system
        'time'    => $time,
    ];
    $messages[$channel][] = $msg;
    if (isset($channels[$channel])) {
        $channels[$channel]['count'] = count($messages[$channel]);
    }
    return $msg;
}

function channelList($channels) {
    return array_values($channels);
}

function userList($onlineUsers) {
    return array_values($onlineUsers);
}

function broadcastState($channel, $nick, $channels, $onlineUsers, $messages) {
    Events::emit("chat:state", [
        'nick'     => $nick,
        'channel'  => $channel,
        'channels' => channelList($channels),
        'users'    => userList($onlineUsers),
        'messages' => $messages[$channel] ?? [],
    ]);
}

function broadcastUsers($onlineUsers) {
    Events::emit("chat:users", userList($onlineUsers));
}

function pushMessage($msg) {
    Events::emit("chat:message", $msg);
}

function pushError($text) {
    Events::emit("chat:error", ['text' => $text]);
}

// ── 初始化种子消息 ──
addMessage($messages, $channels, 'general', 'Origami Bot', '欢迎来到 Origami Chat！输入 /help 查看可用命令。', 'bot', '09:00');
addMessage($messages, $channels, 'general', 'Alice', '大家好 👋', 'user', '09:01');
addMessage($messages, $channels, 'general', 'Bob', '有人试过 Wails v3 吗？', 'user', '09:03');
addMessage($messages, $channels, 'random',  'Origami Bot', '这里是 #随机 频道，随便聊～', 'bot', '09:00');
addMessage($messages, $channels, 'dev',     'Origami Bot', '技术讨论区。有问题尽管问！', 'bot', '09:00');
addMessage($messages, $channels, 'dev',     'Charlie', 'ObjectValue 和 ArrayValue 的区别有人清楚吗？', 'user', '10:15');

// 注册机器人
$onlineUsers['bot'] = ['id' => 'bot', 'name' => 'Origami Bot', 'type' => 'bot', 'channel' => 'general'];

// ══════════════════════════════════════════════════════════════
//  斜杠命令处理
// ══════════════════════════════════════════════════════════════

function handleCommand($cmd, $args, &$myNick, &$myChannel, &$messages, &$channels, &$onlineUsers, $time) {
  global $nextMsgId;

    switch ($cmd) {
        case 'help':
            $help = "可用命令:\n"
                . "  /help      — 显示此帮助\n"
                . "  /nick 名字 — 修改昵称\n"
                . "  /me 动作   — 发送动作消息\n"
                . "  /clear     — 清空当前频道消息\n"
                . "  /channels  — 列出所有频道\n"
                . "  /users     — 列出在线用户";
            $msg = addMessage($messages, $channels, $myChannel, 'Origami Bot', $help, 'bot', $time);
            pushMessage($msg);
            return true;

        case 'nick':
            $newNick = trim($args);
            if ($newNick === '') {
                pushError('用法: /nick 新昵称');
                return true;
            }
            $old = $myNick;
            $myNick = $newNick;
            // 更新在线列表
            foreach ($onlineUsers as $uid => &$u) {
                if ($u['type'] === 'user' && $u['name'] === $old) {
                    $u['name'] = $newNick;
                }
            }
            unset($u);
            $msg = addMessage($messages, $channels, $myChannel, 'system', "{$old} 改名为 {$newNick}", 'system', $time);
            pushMessage($msg);
            broadcastUsers($onlineUsers);
            return true;

        case 'me':
            $action = trim($args);
            if ($action === '') {
                pushError('用法: /me 做了什么');
                return true;
            }
            $msg = addMessage($messages, $channels, $myChannel, $myNick, "* {$myNick} {$action}", 'user', $time);
            pushMessage($msg);
            botReply($myChannel, $myNick, $action, $messages, $channels, $time);
            return true;

        case 'clear':
            $messages[$myChannel] = [];
            $channels[$myChannel]['count'] = 0;
            $msg = addMessage($messages, $channels, $myChannel, 'system', "{$myNick} 清空了频道消息", 'system', $time);
            pushMessage($msg);
            broadcastState($myChannel, $myNick, $channels, $onlineUsers, $messages);
            return true;

        case 'channels':
            $lines = [];
            foreach ($channels as $ch) {
                $lines[] = "  #{$ch['name']} — {$ch['topic']} ({$ch['count']} 条消息)";
            }
            $msg = addMessage($messages, $channels, $myChannel, 'Origami Bot', "频道列表:\n" . implode("\n", $lines), 'bot', $time);
            pushMessage($msg);
            return true;

        case 'users':
            $names = [];
            foreach ($onlineUsers as $u) {
                $tag = $u['type'] === 'bot' ? '🤖' : '👤';
                $names[] = "  {$tag} {$u['name']}";
            }
            $msg = addMessage($messages, $channels, $myChannel, 'Origami Bot', "在线用户 (" . count($onlineUsers) . "):\n" . implode("\n", $names), 'bot', $time);
            pushMessage($msg);
            return true;

        default:
            pushError("未知命令: /{$cmd}，输入 /help 查看帮助");
            return true;
    }
}

// ── 机器人自动回复 ──
function botReply($channel, $nick, $text, &$messages, &$channels, $time) {
    global $nextMsgId;

    $lower = strtolower($text);

  // 关键词匹配
    $reply = null;
    if (str_contains($lower, 'hello') || str_contains($lower, '你好') || str_contains($lower, 'hi')) {
        $reply = "你好 {$nick}！很高兴见到你 😊";
    } elseif (str_contains($lower, 'wails')) {
        $reply = "Wails v3 是用 Go 构建跨平台桌面应用的框架，Origami 让它可以用 PHP 编写！";
    } elseif (str_contains($lower, 'origami')) {
        $reply = "Origami 是一个 Go 实现的 PHP 类脚本语言运行时，支持类、闭包、命名空间等特性。";
    } elseif (str_contains($lower, 'php')) {
        $reply = "没错，你正在用 PHP 语法驱动这个原生桌面应用的后端逻辑 🎉";
    } elseif (str_contains($lower, '帮助') || str_contains($lower, 'help')) {
        $reply = "输入 /help 可以查看所有斜杠命令哦！";
    } elseif (str_contains($lower, '?') || str_contains($lower, '？') || str_contains($lower, '怎么')) {
        $reply = "好问题！你可以在 #开发 频道讨论技术问题，或者输入 /help 查看命令。";
    }

    // #random 频道有 30% 概率随机回复
    if ($reply === null && $channel === 'random' && random_int(1, 100) <= 30) {
        $quotes = [
            "生活就像一盒巧克力 🍫",
            "今天天气不错，适合写代码 ☀️",
            "你知道吗？Origami 的 VM 是用 Go 写的",
            "随机一言：代码写得好，头发掉得少",
            "试试输入 /me 做个动作？",
        ];
        $reply = $quotes[array_rand($quotes)];
    }

    if ($reply !== null) {
        $msg = addMessage($messages, $channels, $channel, 'Origami Bot', $reply, 'bot', $time);
        pushMessage($msg);
    }
}

// ══════════════════════════════════════════════════════════════
//  应用配置
// ══════════════════════════════════════════════════════════════

$options = new App([
    'Title'     => '💬 Origami Chat',
    'Width'     => 960,
    'Height'    => 680,
    'MinWidth'  => 720,
    'MinHeight' => 500,
    'AssetDir'  => __DIR__ . '/chat',
]);

// ══════════════════════════════════════════════════════════════
//  事件处理
// ══════════════════════════════════════════════════════════════

// 用户加入
Events::on("chat:join", function ($data) use (&$myNick, &$myChannel, &$onlineUsers, &$messages, &$channels, &$nextUserId) {
    if (is_array($data) && isset($data['nick'])) {
        $myNick = trim((string)$data['nick']) ?: '访客';
    }

    // 注册到在线列表
    $uid = makeUserId($nextUserId);
    $onlineUsers[$uid] = [
        'id'      => $uid,
        'name'    => $myNick,
        'type'    => 'user',
        'channel' => $myChannel,
    ];

    Log::info("用户加入: {$myNick}");
    broadcastState($myChannel, $myNick, $channels, $onlineUsers, $messages);
});

// 发送消息
Events::on("chat:send", function ($data) use (&$myNick, &$myChannel, &$messages, &$channels, &$onlineUsers) {
    if (!is_array($data)) {
        return;
    }

    $text    = trim((string)($data['text'] ?? ''));
    $time    = (string)($data['time'] ?? '??:??');
    $channel = (string)($data['channel'] ?? $myChannel);
    $nick    = trim((string)($data['nick'] ?? $myNick));

    if ($text === '') {
        return;
    }
    if (!isset($channels[$channel])) {
        pushError("频道不存在: {$channel}");
        return;
    }

    $myNick    = $nick;
    $myChannel = $channel;

    // 斜杠命令
    if (str_starts_with($text, '/')) {
        $parts = explode(' ', $text, 2);
        $cmd   = strtolower(ltrim($parts[0], '/'));
        $args  = $parts[1] ?? '';
        handleCommand($cmd, $args, $myNick, $myChannel, $messages, $channels, $onlineUsers, $time);
        return;
    }

    // 普通消息
    $msg = addMessage($messages, $channels, $channel, $nick, $text, 'user', $time);
    pushMessage($msg);
    Log::info("[#{$channels[$channel]['name']}] {$nick}: {$text}");

    // 机器人回复
    botReply($channel, $nick, $text, $messages, $channels, $time);
});

// 切换频道
Events::on("chat:switch", function ($data) use (&$myChannel, &$myNick, &$channels, &$onlineUsers, &$messages) {
    if (!is_array($data)) {
        return;
    }
    $ch = (string)($data['channel'] ?? '');
    if (!isset($channels[$ch])) {
        pushError("频道不存在: {$ch}");
        return;
    }

    $myChannel = $ch;
    Log::info("切换到频道: #{$channels[$ch]['name']}");
    broadcastState($myChannel, $myNick, $channels, $onlineUsers, $messages);
});

// 修改昵称
Events::on("chat:nick", function ($data) use (&$myNick, &$myChannel, &$onlineUsers, &$messages, &$channels) {
    if (!is_array($data)) {
        return;
    }
    $newNick = trim((string)($data['nick'] ?? ''));
    if ($newNick === '' || $newNick === $myNick) {
        return;
    }

    $old = $myNick;
    $myNick = $newNick;

    foreach ($onlineUsers as $uid => &$u) {
        if ($u['type'] === 'user' && $u['name'] === $old) {
            $u['name'] = $newNick;
        }
    }
    unset($u);

    $msg = addMessage($messages, $channels, $myChannel, 'system', "{$old} 改名为 {$newNick}", 'system', '??:??');
    pushMessage($msg);
    broadcastUsers($onlineUsers);
    Log::info("昵称变更: {$old} → {$newNick}");
});

// ── 生命周期 ──
$options->onDomReady(function () {
    Window::center();
    Log::info("Origami Chat 已就绪");
});

$options->onShutdown(function () {
    Log::info("聊天应用关闭");
});

// ── 启动 ──
Application::run($options);
