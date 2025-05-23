<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class StudentProgress extends Model
{
    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'student_id',
        'book_id',
        'group_id',
        'score_correct',
        'score_incorrect',
        'time_start',
        'time_finish',
        'taken',
    ];

    /**
     * Get the book that owns the StudentProgress
     *
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     */
    public function book(): BelongsTo
    {
        return $this->belongsTo(
            Book::class,
            'book_id',
            'id',
        );
    }

    /**
     * Get the student that owns the StudentProgress
     *
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     */
    public function student(): BelongsTo
    {
        return $this->belongsTo(
            User::class,
            'student_id',
            'id',
        )->whereHas('roles', function ($query) {
            $query->where('name', 'student');
        });
    }

    /**
     * Get the group that owns the StudentProgress
     *
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     */
    public function group(): BelongsTo
    {
        return $this->belongsTo(
            Group::class,
            'group_id',
            'id',
        );
    }
}
