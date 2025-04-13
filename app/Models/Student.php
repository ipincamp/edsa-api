<?php

namespace App\Models;

use App\Enums\RoleEnum;
use Illuminate\Database\Eloquent\Builder;
use Spatie\Permission\Models\Role;

class Student extends User
{
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
}
