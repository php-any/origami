@props(['trigger' => null])
<div>
@if($trigger)
TRIG={{ $trigger }}|
ATTR={{ $trigger->attributes->get('class') }}|
@else
TRIG_NULL
@endif
SLOT={{ $slot }}
</div>