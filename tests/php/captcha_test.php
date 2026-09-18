<?php

namespace tests\php;

/**
 * Captcha 标准库：生成短语、PNG/JPEG、会话存取与校验。
 */

if (!class_exists(\Captcha::class, false)) {
    Log::fatal('Captcha 类未注册');
}
if (!function_exists('captcha_img') || !function_exists('captcha_src') || !function_exists('captcha_check')) {
    Log::fatal('captcha_* 函数未注册');
}
if (\Captcha::MODE_TEXT !== 'text' || \Captcha::MODE_MATH !== 'math') {
    Log::fatal('Captcha 常量不正确');
}

$c = (new \Captcha())
    ->setCharset('A')
    ->setLength(4)
    ->setNoise(0)
    ->setExpire(0)
    ->build();

if ($c->getPhrase() !== 'AAAA') {
    Log::fatal('charset/length 生成失败: ' . $c->getPhrase());
}
if ($c->getDisplay() !== 'AAAA') {
    Log::fatal('display 应与 phrase 相同');
}
if (!$c->test('AAAA') || !$c->verify('aaaa')) {
    Log::fatal('实例 test/verify 失败');
}
if ($c->test('BBBB')) {
    Log::fatal('错误输入不应通过');
}

$png = $c->png();
if (strlen($png) < 32 || ord(substr($png, 0, 1)) !== 137 || substr($png, 1, 3) !== 'PNG') {
    Log::fatal('png() 不是合法 PNG');
}

$jpeg = $c->jpeg(70);
if (strlen($jpeg) < 16 || ord(substr($jpeg, 0, 1)) !== 255 || ord(substr($jpeg, 1, 1)) !== 216) {
    Log::fatal('jpeg() 不是合法 JPEG');
}

$uri = $c->inline();
$prefix = 'data:image/png;base64,';
if (strpos($uri, $prefix) !== 0) {
    Log::fatal('inline() 前缀错误');
}
$decoded = base64_decode(substr($uri, strlen($prefix)));
if ($decoded === false || $decoded !== $png) {
    Log::fatal('inline/base64 与 png() 不一致');
}

$html = $c->img();
if (strpos($html, '<img ') !== 0 || strpos($html, $prefix) === false) {
    Log::fatal('img() HTML 不正确: ' . $html);
}

$path = sys_get_temp_dir() . '/origami_captcha_test.png';
if (!$c->save($path) || !file_exists($path) || filesize($path) < 32) {
    Log::fatal('save() 写入失败');
}
if (file_exists($path)) {
    unlink($path);
}

$c->store();
if (!isset($_SESSION['captcha']['phrase']) || $_SESSION['captcha']['phrase'] !== 'AAAA') {
    Log::fatal('store() 未写入 $_SESSION');
}
if (!\Captcha::check('aaaa')) {
    Log::fatal('Captcha::check 应通过');
}
if (\Captcha::check('aaaa')) {
    Log::fatal('check 成功后应消费会话');
}

$preset = new \Captcha('XyZ9');
$preset->setIgnoreCase(false)->build();
if ($preset->getPhrase() !== 'XyZ9') {
    Log::fatal('构造短语未保留');
}
if ($preset->test('xyz9')) {
    Log::fatal('ignoreCase=false 时不应忽略大小写');
}
if (!$preset->test('XyZ9')) {
    Log::fatal('精确匹配失败');
}

$math = (new \Captcha())->setMode(\Captcha::MODE_MATH)->setNoise(1)->build();
$ans = $math->getPhrase();
$shown = $math->getDisplay();
if ($shown === '' || $ans === '' || $shown === $ans) {
    Log::fatal('算术验证码 display/phrase 异常: ' . $shown . ' / ' . $ans);
}
if (!preg_match('/^[0-9]+$/', $ans)) {
    Log::fatal('算术答案应为数字: ' . $ans);
}
if (!$math->test($ans) || $math->test('not-a-number')) {
    Log::fatal('算术 test 行为错误');
}

$created = \Captcha::create([
    'charset' => 'B',
    'length' => 3,
    'noise' => 0,
    'expire' => 0,
    'key' => 'login_code',
]);
if ($created->getPhrase() !== 'BBB') {
    Log::fatal('Captcha::create 短语错误: ' . $created->getPhrase());
}
if (!captcha_check('bbb', 'login_code')) {
    Log::fatal('captcha_check 失败');
}

$src = captcha_src(['charset' => 'C', 'length' => 2, 'noise' => 0, 'expire' => 0, 'store' => false]);
if (strpos($src, $prefix) !== 0) {
    Log::fatal('captcha_src 前缀错误');
}

Log::info('Captcha 标准库测试通过');
