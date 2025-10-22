<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletes;
use Spatie\Activitylog\Traits\LogsActivity;

class Course extends Model
{
    use LogsActivity,
        SoftDeletes;

    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'name',
        'description',
    ];

    public function getActivitylogOptions(): \Spatie\Activitylog\LogOptions
    {
        return \Spatie\Activitylog\LogOptions::defaults()
            ->logOnly(['name', 'description'])
            ->dontSubmitEmptyLogs()
            ->useLogName('course')
            ->setDescriptionForEvent(fn(string $eventName) => "Course {$this->name} has been {$eventName}")
            ->logOnlyDirty();
    }

    /**
     * Get all of the groups for the Course
     *
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     */
    public function groups(): \Illuminate\Database\Eloquent\Relations\HasMany
    {
        return $this->hasMany(
            Group::class,
            'course_id',
            'id'
        );
    }
}
