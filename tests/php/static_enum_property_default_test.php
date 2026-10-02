<?php

namespace StaticEnumPropertyDefault;

enum Alignment: string {
    case Start = 'start';
    case Center = 'center';
}

abstract class BasePage {
    public static string | Alignment $formActionsAlignment = Alignment::Start;

    public function getFormActionsAlignment(): string | Alignment {
        return static::$formActionsAlignment;
    }
}

class Login extends BasePage {}
class CustomLogin extends BasePage {
    public static string | Alignment $formActionsAlignment = Alignment::Center;
}

abstract class LazyBasePage {
    public static string | \StaticEnumPropertyDefault\LateAlignment $formActionsAlignment = LateAlignment::Start;
    public static string | \StaticEnumPropertyDefault\LateAlignment $selfAlignment = LateAlignment::Start;
    public static string | \StaticEnumPropertyDefault\LateAlignment $explicitAlignment = self::DEFAULT_ALIGNMENT;
    public const DEFAULT_ALIGNMENT = LateAlignment::Start;

    public function getFormActionsAlignment(): string | \StaticEnumPropertyDefault\LateAlignment {
        return static::$formActionsAlignment;
    }

    public function getSelfAlignment(): string | \StaticEnumPropertyDefault\LateAlignment {
        return self::$selfAlignment;
    }
}

class LazyLogin extends LazyBasePage {}

require __DIR__ . '/static_enum_property_default_fixtures/alignment.php';

$login = new Login();
if ($login->getFormActionsAlignment() !== Alignment::Start) {
    throw new \Exception('Inherited static enum default was not preserved');
}
if ((new CustomLogin())->getFormActionsAlignment() !== Alignment::Center) {
    throw new \Exception('Redeclared static enum default was not preserved');
}
$lazyLogin = new LazyLogin();
if ($lazyLogin->getFormActionsAlignment() !== LateAlignment::Start) {
    throw new \Exception('Runtime-loaded enum default was not preserved');
}
if ($lazyLogin->getSelfAlignment() !== LateAlignment::Start || LazyBasePage::$explicitAlignment !== LateAlignment::Start) {
    throw new \Exception('self:: and explicit class access lost the enum default');
}
echo "static enum property default: OK\n";
