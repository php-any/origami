<?php
class T {
  public function toEmbeddedHtml(): string {
    return self::w(fn (): string => $this->renderEmbeddedHtml());
  }
  public static function w(callable $c): mixed {
    return $c();
  }
  protected function renderEmbeddedHtml(): string {
    $components = [1];
    $has = false;
    $mapped = array_map(function ($c) use (&$has): array {
      $has = true;
      return [$c, true];
    }, $components);
    if (!$has) return '';
    ob_start(); ?>
    <div><?= 'x' ?></div>
    <?php return ob_get_clean();
  }
}
$t = new T;
echo gettype($t->toEmbeddedHtml()), "\n";
echo $t->toEmbeddedHtml(), "\n";
