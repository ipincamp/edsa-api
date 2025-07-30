<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Support\Facades\Storage;

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
        'author',
        'year',
        'genre',
        'focus',
        'cover_image',
        'order_sequence',
    ];

    /**
     * Accessor untuk mendapatkan URL lengkap dari cover image.
     *
     * @return string
     */
    public function getCoverImageUrlAttribute(): string
    {
        if ($this->cover_image && Storage::disk('public')->exists($this->cover_image)) {
            return asset('assets/' . $this->cover_image);
        }

        return 'https://placehold.co/800x400/cccccc/FFFFFF/png?text=No+Image';
    }

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
