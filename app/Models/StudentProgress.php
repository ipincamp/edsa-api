<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\BelongsToMany;

class StudentProgress extends Model
{
    use HasFactory;

    protected $table = 'student_progresses';

    /**
     * The attributes that are mass assignable.
     *
     * @var list<string>
     */
    protected $fillable = [
        'student_id',
        'book_id',
        'last_page',
        'latest_page',
        'total_points',
        'status',
    ];

    /**
     * Get the student that owns the StudentProgress
     *
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     */
    public function student(): BelongsTo
    {
        return $this->belongsTo(User::class, 'student_id');
    }

    /**
     * Get the book that owns the StudentProgress
     *
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     */
    public function book(): BelongsTo
    {
        return $this->belongsTo(Book::class);
    }

    /**
     * The completedInteractions that belong to the StudentProgress
     *
     * @return \Illuminate\Database\Eloquent\Relations\BelongsToMany
     */
    public function completedInteractions(): BelongsToMany
    {
        return $this->belongsToMany(Interaction::class, 'interaction_progress')
            ->withPivot(['correct', 'total'])
            ->withTimestamps();
    }

    /**
     * The completedPostActivities that belong to the StudentProgress
     *
     * @return \Illuminate\Database\Eloquent\Relations\BelongsToMany
     */
    public function completedPostActivities(): BelongsToMany
    {
        return $this->belongsToMany(PostActivity::class, 'post_activity_progress')
            ->withTimestamps();
    }
}
