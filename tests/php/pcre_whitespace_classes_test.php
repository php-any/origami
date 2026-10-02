<?php
function whitespaceCheck($actual, $expected, $label) {
    if ($actual !== $expected) { throw new Exception($label . ': ' . json_encode($actual)); }
}
whitespaceCheck(preg_match('/^\h+$/', " \t"), 1, 'horizontal ASCII');
whitespaceCheck(preg_match('/^[^\d\h]+$/', 'seconds'), 1, 'horizontal in negated class');
whitespaceCheck(preg_match('/^\h$/', "\n"), 0, 'horizontal excludes newline');
whitespaceCheck(preg_match('/^\v+$/', "\r\n\v\f"), 1, 'vertical includes newline');
whitespaceCheck(preg_match('/^[\v]+$/', "\n"), 1, 'vertical in character class');
whitespaceCheck(preg_match('/^\H+$/', 'abc'), 1, 'horizontal complement');
whitespaceCheck(preg_match('/^[\H]+$/', 'abc'), 1, 'horizontal complement in class');
whitespaceCheck(preg_match('/^[\V]+$/', 'abc'), 1, 'vertical complement in class');
whitespaceCheck(preg_match('/^[\H]+$/', ' '), 0, 'horizontal complement excludes space');
whitespaceCheck(preg_match('/^\h$/u', "\u{2003}"), 1, 'horizontal Unicode');
whitespaceCheck(preg_match('/^\v$/u', "\u{2028}"), 1, 'vertical Unicode');
whitespaceCheck(preg_match('/^\Q\h\E$/', '\h'), 1, 'quoted escape remains literal');
whitespaceCheck(preg_match('/^\\\\h$/', '\h'), 1, 'escaped backslash remains literal');
preg_match_all('/(-?\d+(?:\.\d+)?)\h*([^\d\h]*)/i', '1 second', $parts, PREG_SET_ORDER);
whitespaceCheck($parts, [['1 second', '1', 'second']], 'Carbon interval groups');
echo "PCRE whitespace classes: PASS\n";
