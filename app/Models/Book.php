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
     * Get the postActivity associated with the Book
     *
     * @return \Illuminate\Database\Eloquent\Relations\HasOne
     */
    public function postActivity(): HasOne
    {
        return $this->hasOne(PostActivity::class);
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
