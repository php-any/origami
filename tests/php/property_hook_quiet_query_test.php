<?php
class HookQuietQuery {
    public int $reads = 0;
    public ?string $value { get { $this->reads++; return 'yes'; } }
    public ?string $nothing { get { $this->reads++; return null; } }
    public array $items { get { $this->reads++; return ['key' => 42]; } }
    public string $failure { get { throw new Exception('getter failure'); } }
}
$object = new HookQuietQuery;
if (!isset($object->value) || $object->reads !== 1) { throw new Exception('isset getter count'); }
if (($object->value ?? 'fallback') !== 'yes' || $object->reads !== 2) { throw new Exception('coalesce getter count'); }
if (isset($object->nothing) || $object->reads !== 3) { throw new Exception('nullable isset'); }
if (($object->nothing ?? 'fallback') !== 'fallback' || $object->reads !== 4) { throw new Exception('nullable coalesce'); }
if (($object->items['key'] ?? 0) !== 42 || $object->reads !== 5) { throw new Exception('nested getter count'); }
if (!isset($object->items['key']) || $object->reads !== 6) { throw new Exception('nested isset getter count'); }
$name = 'value';
if (($object->$name ?? 'fallback') !== 'yes' || $object->reads !== 7) { throw new Exception('dynamic getter count'); }
$caught = 0;
try { isset($object->failure); } catch (Exception $e) { if ($e->getMessage() === 'getter failure') { $caught++; } }
try { $object->failure ?? 'fallback'; } catch (Exception $e) { if ($e->getMessage() === 'getter failure') { $caught++; } }
if ($caught !== 2) { throw new Exception('quiet queries must preserve getter exceptions'); }
echo "property hook quiet query OK\n";
