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

    /*
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'name',
        'username',
        'password',
    ];

    protected static function boot()
    {
        parent::boot();

        static::created(fn($student) => $student->roles()->attach(
            Role::findByName(RoleEnum::STUDENT->value)->id,
        ));

        static::addGlobalScope('student', function (Builder $builder) {
            $builder->whereNotNull('username');
        });
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
