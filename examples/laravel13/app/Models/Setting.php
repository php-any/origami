<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Attributes\Fillable;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\Cache;
use Spatie\Activitylog\LogOptions;
use Spatie\Activitylog\Traits\LogsActivity;

#[Fillable(['key', 'value', 'group', 'label'])]
class Setting extends Model
{
    use LogsActivity;

    protected $table = 'settings';

    protected function casts(): array
    {
        return [
            'value' => 'json',
        ];
    }

    public function getActivitylogOptions(): LogOptions
    {
        return LogOptions::defaults()
            ->logFillable()
            ->logOnlyDirty()
            ->dontSubmitEmptyLogs();
    }

    public static function getValue(string $key, mixed $default = null): mixed
    {
        $setting = static::query()->where('key', $key)->first();

        if (! $setting) {
            return $default;
        }

        return $setting->value ?? $default;
    }

    public static function setValue(string $key, mixed $value, string $group = 'general', ?string $label = null): self
    {
        $setting = static::query()->updateOrCreate(
            ['key' => $key],
            [
                'value' => $value,
                'group' => $group,
                'label' => $label,
            ],
        );

        Cache::forget('setting.'.$key);

        return $setting;
    }
}
