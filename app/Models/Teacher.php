<?php

namespace App\Models;

use App\Enums\RoleEnum;
use Illuminate\Database\Eloquent\Builder;
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
        ));

        static::addGlobalScope('teacher', fn(Builder $builder) => $builder
            ->whereNull('username')
            ->whereNull('email_verified_at'));
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
