{{-- Mirror Filament logo @capture extract --}}
@php
    echo "TOP isset=".(isset($attributes)?'yes':'no')." type=".gettype($attributes ?? null)."\n";
    $gdvTop = get_defined_vars();
    echo "TOP gdv_has=".(isset($gdvTop['attributes'])?'yes':'no')."\n";
    $brandName = 'Brand';
    $getLogoClasses = fn (bool $isDarkMode): string => 'fi-logo';
    $logoStyles = 'height: 1.5rem';
    echo "MID isset=".(isset($attributes)?'yes':'no')."\n";
@endphp

@capture($content, $logo, $isDarkMode = false)
    @php
        echo "CAP attrs_isset=".(isset($attributes)?'yes':'no')." type=".gettype($attributes ?? null)."\n";
    @endphp
    <div>logo</div>
@endcapture

{!! $content(null) !!}
