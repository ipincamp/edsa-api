<?php

namespace App\Models;

use App\Enums\RoleEnum;
use Spatie\Activitylog\LogOptions;
use Spatie\Activitylog\Traits\LogsActivity;

class Teacher extends User
{
    use LogsActivity;

    protected $table = 'users';

    /*
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'name',
        'email',
        'password',
    ];

    protected static function boot()
    {
        parent::boot();

        static::created(fn($student) => $student->roles()->attach(
            Role::findByName(RoleEnum::TEACHER->value)->id,
            ['model_type' => 'App\Models\Teacher', 'model_id' => $student->id],
        ));
    }

    public function getActivitylogOptions(): LogOptions
    {
        return LogOptions::defaults()
            ->logOnly(['name', 'email', 'password'])
            ->dontSubmitEmptyLogs()
            ->useLogName('teacher')
            ->setDescriptionForEvent(fn(string $eventName) => "Teacher has been {$eventName}")
            ->logOnlyDirty();
    }
}
