<?php
if (ignore_user_abort() !== 0 || ignore_user_abort(true) !== 0 || ignore_user_abort() !== 1 || ini_get('ignore_user_abort') !== '1') { throw new Exception('ignore abort enable'); }
if (ignore_user_abort(false) !== 1 || ignore_user_abort(null) !== 0 || connection_status() !== 0 || connection_aborted() !== 0) { throw new Exception('CLI connection status'); }
echo "connection state OK\n";
