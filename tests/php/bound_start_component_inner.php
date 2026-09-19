<?php

$this->startComponent('cell');
echo 'cell-body';
$cell = $this->renderComponent();
if ($cell !== 'cell') {
    \Log::fatal('include 内 renderComponent 失败: '.var_export($cell, true));
}
