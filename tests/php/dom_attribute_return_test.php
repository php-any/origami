<?php
$document = new DOMDocument();
$document->loadHTML('<div id="sample" class=""></div>');
$element = $document->getElementsByTagName('div')->item(0);
function domString(DOMElement $element, string $attribute): string {
    return $element->getAttribute($attribute);
}
if (domString($element, 'missing') !== '' || domString($element, 'id') !== 'sample' || domString($element, 'class') !== '') {
    throw new Exception('DOM attributes must return strings');
}
echo "DOM attribute return: PASS\n";
