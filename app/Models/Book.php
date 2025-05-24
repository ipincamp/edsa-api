<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;
use Illuminate\Database\Eloquent\SoftDeletes;

class Book extends Model
{
    use SoftDeletes;

    /**
     * The attributes that are mass assignable.
     *
     * @var array
     */
    protected $fillable = [
        'title',
        'author',
        'year',
        'genre',
        'focus',
        'status',
        'image',
    ];

    /**
     * Get all of the settings for the Book
     *
     * @return \Illuminate\Database\Eloquent\Relations\HasMany
     */
    public function settings(): HasMany
    {
        return $this->hasMany(
            BookSetting::class,
            'book_id',
            'id',
        );
    }
}
