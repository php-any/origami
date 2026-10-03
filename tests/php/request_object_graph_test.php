<?php
class RequestGraphChild {
    public int $count = 0;
}
class RequestGraphOwner {
    public array $data = ['nested' => ['values' => [10]], 'gone' => 2];
    public RequestGraphChild $child;
    public RequestGraphChild $alias;
    public $callback;
    public function __construct() {
        $this->child = new RequestGraphChild();
        $this->alias = $this->child;
        $child = $this->child;
        $this->callback = function () use ($child) { return ++$child->count; };
    }
    public function advance(): string {
        $values =& $this->data['nested']['values'];
        $values[] = 20;
        unset($this->data['gone']);
        $this->data['gone'] = 3;
        $callback = $this->callback;
        $count = $callback();
        return $count . ':' . $this->alias->count . ':' . count($values) . ':' . implode(',', array_keys($this->data));
    }
}
$owner = new RequestGraphOwner();
if ($owner->advance() !== '1:1:2:nested,gone') { throw new Exception('object graph alias or nested reference failed'); }
return new RequestGraphOwner();
