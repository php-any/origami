<?php
trait SyntaxConflictLeft { public function value(): string { return 'left'; } }
trait SyntaxConflictRight { public function value(): string { return 'right'; } }
class SyntaxChooseLeft {
    use SyntaxConflictLeft, SyntaxConflictRight {
        SyntaxConflictLeft::value insteadof SyntaxConflictRight;
        SyntaxConflictRight::value as other;
    }
}
class SyntaxChooseRight {
    use SyntaxConflictLeft, SyntaxConflictRight {
        SyntaxConflictRight::value insteadof SyntaxConflictLeft;
        SyntaxConflictLeft::value as other;
    }
}
$left = new SyntaxChooseLeft;
$right = new SyntaxChooseRight;
if ($left->value() !== 'left' || $left->other() !== 'right'
    || $right->value() !== 'right' || $right->other() !== 'left') {
    throw new Exception('trait conflict selection and excluded method alias');
}

trait SyntaxPrecedenceKeep { public function value(): string { return 'keep'; } }
trait SyntaxPrecedenceDrop { public function value(): string { return 'drop'; } }
trait SyntaxPrecedenceThird { public function value(): string { return 'third'; } }
class SyntaxPrecedenceTarget {
    use SyntaxPrecedenceKeep, SyntaxPrecedenceDrop {
        SyntaxPrecedenceKeep::value insteadof SyntaxPrecedenceDrop;
        SyntaxPrecedenceDrop::value as private;
    }
}
class SyntaxThreeTraits {
    use SyntaxPrecedenceKeep, SyntaxPrecedenceDrop, SyntaxPrecedenceThird {
        SyntaxPrecedenceKeep::value insteadof SyntaxPrecedenceDrop, SyntaxPrecedenceThird;
    }
}
$kept = new SyntaxPrecedenceTarget;
if ($kept->value() !== 'keep' || !(new ReflectionMethod($kept, 'value'))->isPublic()) {
    throw new Exception('same-name visibility change must not replace insteadof winner');
}
if ((new SyntaxThreeTraits)->value() !== 'keep') {
    throw new Exception('insteadof must exclude every listed trait');
}

trait SyntaxAfterLeft { public function value(): string { return 'left'; } }
trait SyntaxAfterRight { public function value(): string { return 'right'; } }
trait SyntaxAfterOwner {
    use SyntaxAfterLeft, SyntaxAfterRight;
    public function value(): string { return 'owner'; }
}
class SyntaxAfterUser { use SyntaxAfterOwner; }
if ((new SyntaxAfterUser)->value() !== 'owner') {
    throw new Exception('trait method declared after use resolves the conflict');
}
echo "trait insteadof OK\n";
