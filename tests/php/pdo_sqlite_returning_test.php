<?php
$db = new PDO('sqlite::memory:');
$db->exec('CREATE TABLE records(id INTEGER PRIMARY KEY, value TEXT)');
foreach ([['INSERT INTO records(value) VALUES(?) RETURNING id,value','first'],['UPDATE records SET value=? WHERE id=1 RETURNING id,value','second']] as $case) {
    $statement=$db->prepare($case[0]); $statement->execute([$case[1]]); $row=$statement->fetch(PDO::FETCH_ASSOC);
    if ((int)$row['id'] !== 1 || $row['value'] !== $case[1]) { throw new RuntimeException('RETURNING result lost'); }
}
$statement=$db->prepare('DELETE FROM records WHERE id=? RETURNING value'); $statement->execute([1]);
if ($statement->fetchColumn() !== 'second') { throw new RuntimeException('DELETE RETURNING result lost'); }
echo "pdo returning OK\n";

$db->beginTransaction();
$statement=$db->prepare('INSERT INTO records(id,value) VALUES(?,?) RETURNING id');
$statement->execute([7,'transaction']);
if ($db->lastInsertId() !== '7' || (int)$statement->fetchColumn() !== 7) { throw new RuntimeException('RETURNING lost lastInsertId inside transaction'); }
$db->rollBack();
$statement=$db->prepare('INSERT INTO records(id,value) VALUES(?,?) RETURNING id');
$statement->execute([9,'repeat']);
$statement->execute([10,'repeat']);
if ($db->lastInsertId() !== '10' || (int)$statement->fetchColumn() !== 10) { throw new RuntimeException('Repeated RETURNING execution retained stale metadata'); }
echo "pdo returning metadata OK\n";

$statement=$db->prepare("WITH source(v) AS (SELECT 'cte') INSERT INTO records(id,value) SELECT 11,v FROM source RETURNING id");
$statement->execute();
if ($db->lastInsertId() !== '11' || (int)$statement->fetchColumn() !== 11) { throw new RuntimeException('CTE RETURNING lost metadata'); }
echo "pdo CTE returning OK\n";
