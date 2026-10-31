<?php
$child=tempnam(sys_get_temp_dir(),'stream-chunks-');
try {
    file_put_contents($child,'<?php $out=fopen("php://stdout","w"); fwrite($out,"first"); usleep(150000); fwrite($out,"second"); fclose($out);');
    $pipes=[]; $process=proc_open([PHP_BINARY,$child],[1=>['pipe','w']],$pipes);
    $result=stream_get_contents($pipes[1],100);
    fclose($pipes[1]); $code=proc_close($process);
    if($result!=='firstsecond' || $code!==0) { throw new RuntimeException('stream_get_contents truncated: '.json_encode($result)); }
} finally { unlink($child); }
echo "stream_get_contents_pipe_length_test: OK\n";

$stream=fopen('php://memory','w+');
fwrite($stream,'abcdef'); rewind($stream);
if(stream_get_contents($stream,0)!=='' || stream_get_contents($stream,3)!=='abc' || stream_get_contents($stream,10)!=='def' || stream_get_contents($stream,1)!=='') { throw new RuntimeException('length boundary or stream position changed'); }
fclose($stream);
echo "stream length boundaries OK\n";
