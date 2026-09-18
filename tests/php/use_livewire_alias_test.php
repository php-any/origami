<?php
namespace App\Test;
use Livewire\Livewire as Livewire;
// This could break Livewire\Features\... expansion if alias is Livewire
class T {
  public static function t() {
    // relative-looking FQN starting with Livewire
    return \Livewire\Features\SupportCompiledWireKeys\SupportCompiledWireKeys::class;
  }
}
echo T::t(), "\n";
