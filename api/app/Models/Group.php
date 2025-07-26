<?php

namespace App\Models;

use App\Enums\RoleEnum;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\BelongsToMany;
use Illuminate\Database\Eloquent\SoftDeletes;
use Spatie\Activitylog\LogOptions;
use Spatie\Activitylog\Traits\LogsActivity;

class Group extends Model
{
    use LogsActivity, SoftDeletes;

    /**
     * The attributes that are mass assignable.
     *
     * @var list<string>
     */
    protected $fillable = [
        'name',
        'description',
        'course_id',
    ];

    public function getActivitylogOptions(): LogOptions
    {
        return LogOptions::defaults()
            ->logOnly(['name', 'description', 'course_id'])
            ->dontSubmitEmptyLogs()
            ->useLogName('group')
            ->setDescriptionForEvent(fn(string $eventName) => "Group has been {$eventName}")
            ->logOnlyDirty();
    }

    /**
     * Get the course that owns the Group
     *
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     */
    public function course(): BelongsTo
    {
        return $this->belongsTo(
            Course::class,
            'course_id',
            'id'
        );
    }

    /**
     * The users that belong to the Group
     *
     * @return \Illuminate\Database\Eloquent\Relations\BelongsToMany
     */
    public function participants(): BelongsToMany
    {
        return $this->belongsToMany(
            User::class,
            'group_participants',
            'group_id',
            'user_id'
        )->withTimestamps();
    }

    public function teachers()
    {
        return $this->participants()->whereHas('roles', function ($query) {
            $query->where('name', RoleEnum::TEACHER->value);
        });
    }

    public function students()
    {
        return $this->participants()->whereHas('roles', function ($query) {
            $query->where('name', RoleEnum::STUDENT->value);
        });
    }
}
