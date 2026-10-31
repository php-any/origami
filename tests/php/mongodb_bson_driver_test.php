<?php

use MongoDB\BSON\ObjectId;
use MongoDB\BSON\UTCDateTime;
use MongoDB\BSON\Binary;
use MongoDB\BSON\Regex;
use MongoDB\BSON\Decimal128;
use MongoDB\BSON\Timestamp;
use MongoDB\Driver\BulkWrite;
use MongoDB\Driver\Command;
use MongoDB\Driver\Manager;
use MongoDB\Driver\Query;
use MongoDB\Driver\WriteConcern;

function checkMongo($condition, $message) {
    if (!$condition) { throw new RuntimeException($message); }
}
function invalidMongo($callback) {
    try { $callback(); } catch (MongoDB\Driver\Exception\InvalidArgumentException $e) {
        checkMongo(strlen($e->getMessage()) > 0, 'exception message');
        return;
    }
    throw new RuntimeException('expected InvalidArgumentException');
}

checkMongo(extension_loaded('mongodb'), 'extension registration');
checkMongo(phpversion('mongodb') === MONGODB_VERSION, 'honest module version');
checkMongo(in_array('mongodb', get_loaded_extensions()), 'loaded extensions');
$id = new ObjectId('507f1f77bcf86cd799439011');
$other = new ObjectId('000000000000000000000001');
checkMongo((string) $id === '507f1f77bcf86cd799439011', 'ObjectId isolation');
checkMongo($id->getTimestamp() === 1350508407, 'ObjectId timestamp');
checkMongo(json_encode($id) === '{"$oid":"507f1f77bcf86cd799439011"}', 'BSON JSON serialization');
checkMongo($id instanceof MongoDB\BSON\Type, 'BSON marker interface');
$document = [
    '_id' => $id,
    'list' => [1, true, null, 'text'],
    'sparse' => [2 => 'two', 5 => 'five'],
    'emptyArray' => [],
    'emptyObject' => new stdClass(),
    'date' => new UTCDateTime(-1234),
    'binary' => new Binary("a\0b", Binary::TYPE_GENERIC),
    'regex' => new Regex('abc', 'mi'),
    'decimal' => new Decimal128('123.45'),
    'timestamp' => new Timestamp(7, 42),
    'min' => new MongoDB\BSON\MinKey(),
    'max' => new MongoDB\BSON\MaxKey(),
];
$raw = MongoDB\BSON\fromPHP($document);
$decoded = MongoDB\BSON\toPHP($raw);
checkMongo($decoded instanceof stdClass, 'root defaults to object');
checkMongo((string) $decoded->_id === (string) $id, 'ObjectId roundtrip');
checkMongo($decoded->list === [1, true, null, 'text'], 'array roundtrip');
checkMongo($decoded->emptyArray === [], 'empty BSON array');
checkMongo($decoded->emptyObject instanceof stdClass, 'empty BSON document');
checkMongo((string) $decoded->date === '-1234', 'negative datetime');
checkMongo($decoded->binary->getData() === "a\0b", 'binary NUL');
checkMongo($decoded->regex->getFlags() === 'im', 'regex flags sorted');
checkMongo((string) $decoded->decimal === '123.45', 'decimal roundtrip');
checkMongo($decoded->timestamp->getIncrement() === 7, 'timestamp increment');
checkMongo($decoded->timestamp->getTimestamp() === 42, 'timestamp seconds');
$array = MongoDB\BSON\toPHP($raw, ['root' => 'array', 'document' => 'array']);
checkMongo($array['sparse'] === [2 => 'two', 5 => 'five'], 'sparse keys');
$json = MongoDB\BSON\toCanonicalExtendedJSON($raw);
$again = MongoDB\BSON\toPHP(MongoDB\BSON\fromJSON($json));
checkMongo((string) $again->_id === (string) $id, 'extended JSON roundtrip');
invalidMongo(function () { new ObjectId('bad'); });
invalidMongo(function () { MongoDB\BSON\toPHP('broken'); });
invalidMongo(function () { MongoDB\BSON\fromPHP(["a\0b" => 1]); });
invalidMongo(function () { new Binary('short', Binary::TYPE_UUID); });
invalidMongo(function () { new Query([], ['unknownOption' => true]); });
invalidMongo(function () { new Command([]); });

$bulk = new BulkWrite();
$source = ['name' => 'before'];
$generated = $bulk->insert($source);
checkMongo($generated instanceof ObjectId, 'insert generates id');
checkMongo(!isset($source['_id']), 'insert must not mutate input');
checkMongo($bulk->insert(['_id' => 'string-id']) === null, 'non-ObjectId insert return');
$bulk->update(['name' => 'before'], ['$set' => ['name' => 'after']], ['multi' => true]);
$bulk->delete(['name' => 'after'], ['limit' => 1]);
checkMongo(count($bulk) === 4, 'bulk operation count');
checkMongo(count(new BulkWrite()) === 0, 'bulk instance isolation');
invalidMongo(function () use ($bulk) { $bulk->update([], ['field' => 1], ['multi' => true]); });
invalidMongo(function () use ($bulk) { $bulk->update([], [['$set' => ['field' => 1]]]); });
checkMongo(count($bulk) === 4, 'invalid update must not enqueue');
$wc = new WriteConcern(WriteConcern::MAJORITY, 0, true);
checkMongo($wc->getW() === 'majority' && $wc->getJournal() === true, 'write concern');

// Manager construction must succeed without a listening MongoDB server.
$manager = new Manager('mongodb://127.0.0.1:1/', ['serverSelectionTimeoutMS' => 30, 'connectTimeoutMS' => 30]);
try {
    $manager->executeCommand('test', new Command(['ping' => 1]));
    throw new RuntimeException('connection unexpectedly succeeded');
} catch (MongoDB\Driver\Exception\ConnectionTimeoutException $e) {
    checkMongo(strlen($e->getMessage()) > 0, 'connection error message');
}
echo "mongodb_bson_driver_test: OK\n";
