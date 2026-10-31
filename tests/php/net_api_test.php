<?php
function expectInvalid(callable $fn): void { $caught=false;try{$fn();}catch(Throwable $e){$caught=true;}if(!$caught){throw new RuntimeException('invalid option accepted');} }
$http=new Net\Http\Server('127.0.0.1',0);
$http->configure(['maxBodyBytes'=>16,'maxRequests'=>2]);
expectInvalid(function()use($http){$http->configure(['unknown'=>10]);});
expectInvalid(function()use($http){$http->configure(['maxBodyBytes'=>0]);});
$tcp=new Net\Tcp\Server('127.0.0.1:0');
$tcp->configure(['maxConnections'=>2,'readTimeoutMs'=>50]);
expectInvalid(function()use($tcp){$tcp->configure(['readTimeoutMs'=>'50']);});
$path=tempnam(sys_get_temp_dir(),'generic-proto-');
try{
 file_put_contents($path,'syntax="proto3"; import "google/protobuf/wrappers.proto"; service Echo { rpc Say(google.protobuf.StringValue) returns(google.protobuf.StringValue); }');
 $schema=new Net\Grpc\Schema([basename($path)],[dirname($path)]);
 $methods=json_decode($schema->methodsJson(),true);
 if(count($methods)!==1||$methods[0]['path']!=='/Echo/Say'||$methods[0]['input']!=='google.protobuf.StringValue'){throw new RuntimeException('schema metadata');}
 $grpc=new Net\Grpc\Server('127.0.0.1:0',$schema);$grpc->configure(['maxCalls'=>2]);$grpc->reflection(true);
 expectInvalid(function()use($grpc){$grpc->configure(['maxCalls'=>-1]);});
}finally{unlink($path);}
echo "net_api_test: OK\n";
