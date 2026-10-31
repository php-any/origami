<?php
$input = '<a title="Bob\'s">&amp; &#34;</a>';
$expected = '&lt;a title=&quot;Bob&#039;s&quot;&gt;&amp;amp; &amp;#34;&lt;/a&gt;';
if (htmlspecialchars($input) !== $expected) { throw new RuntimeException('default quote entities'); }
if (htmlspecialchars($input, ENT_QUOTES | ENT_SUBSTITUTE, 'UTF-8') !== $expected) { throw new RuntimeException('explicit quote entities'); }
class HTMLQuoteString {
    public function __toString() { return '"\''; }
}
if (htmlspecialchars(new HTMLQuoteString()) !== '&quot;&#039;') { throw new RuntimeException('object quote entities'); }
echo "htmlspecialchars_quote_entities_test: OK\n";
