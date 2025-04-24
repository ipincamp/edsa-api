<?php

namespace App\Models;

use App\Enums\RoleEnum;
use Illuminate\Database\Eloquent\Builder;
use Spatie\Permission\Models\Role;
use Spatie\Activitylog\LogOptions;
use Spatie\Activitylog\Traits\LogsActivity;

class Student extends User
{
    use LogsActivity;

    protected $table = 'users';

    /**
     * The attributes that aren't mass assignable.
     *
     * @var list<string>
     */
    protected $guarded = [
        'email',
        'avatar',
    ];

    protected static function boot()
    {
        parent::boot();

        static::created(fn($student) => $student->roles()->attach(
            Role::findByName(RoleEnum::STUDENT->value)->id,
            ['model_type' => 'App\Models\Student', 'model_id' => $student->id],
        ));
    }

    public function getActivitylogOptions(): LogOptions
    {
        return LogOptions::defaults()
            ->logOnly(['name', 'username', 'password'])
            ->dontSubmitEmptyLogs()
            ->useLogName('student')
            ->setDescriptionForEvent(fn(string $eventName) => "Student has been {$eventName}")
            ->logOnlyDirty();
    }
}
