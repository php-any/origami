@php
    use Illuminate\Contracts\Support\Htmlable;
    use Illuminate\Support\Arr;

    echo "L1 isset=".(isset($attributes)?'yes':'no')." type=".gettype($attributes ?? null)."\n";

    $brandName = filament()->getBrandName();
    echo "L2 isset=".(isset($attributes)?'yes':'no')."\n";
    $brandLogo = filament()->getBrandLogo();
    $brandLogoHeight = filament()->getBrandLogoHeight() ?? '1.5rem';
    $darkModeBrandLogo = filament()->getDarkModeBrandLogo();
    $hasDarkModeBrandLogo = filled($darkModeBrandLogo);

    $getLogoClasses = fn (bool $isDarkMode): string => Arr::toCssClasses([
        'fi-logo',
        'fi-logo-light' => $hasDarkModeBrandLogo && (! $isDarkMode),
        'fi-logo-dark' => $isDarkMode,
    ]);

    $logoStyles = 'height: ' . e($brandLogoHeight);
    echo "L3 isset=".(isset($attributes)?'yes':'no')."\n";
    $gdv = get_defined_vars();
    echo "L3 gdv_attr=".(isset($gdv['attributes'])?'yes':'no')." type=".gettype($gdv['attributes'] ?? null)."\n";
@endphp

@capture($content, $logo, $isDarkMode = false)
    @php
        echo "CAP isset=".(isset($attributes)?'yes':'no')." type=".gettype($attributes ?? null)."\n";
    @endphp
    <div>x</div>
@endcapture

{{ $content($brandLogo) }}
