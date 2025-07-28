<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Database\Eloquent\Relations\HasOne;

class Book extends Model
{
    use HasFactory;

    /**
     * The attributes that are mass assignable.
     *
     * @var list<string>
     */
    protected $fillable = [
        'title',
        'cover_image',
        'order_sequence',
    ];

    /**
     * Get all of the pages for the Book
     *
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     */
    public function pages(): HasMany
    {
        return $this->hasMany(Page::class)
            ->orderBy('page_number', 'asc');
    }

    /**
     * Get all of the postActivities for the Book
     *
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     */
    public function postActivities(): HasMany
    {
        return $this->hasMany(PostActivity::class)
            ->orderBy('order', 'asc');
    }

    /**
     * Get all of the progresses for the Book
     *
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     */
    public function progresses(): HasMany
    {
        return $this->hasMany(StudentProgress::class);
    }
}
