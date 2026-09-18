<?php
function t(): string {
  ob_start(); ?>
  <div><?= 'x' ?></div>
  <?php return ob_get_clean();
}
$r = t();
echo "type=".gettype($r)."\n";
echo "val=".var_export($r, true)."\n";
