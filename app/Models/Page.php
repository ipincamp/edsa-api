<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;
use Illuminate\Database\Eloquent\Relations\HasOne;

class Page extends Model
{
    use HasFactory;

    /**
     * The attributes that are mass assignable.
     *
     * @var list<string>
     */
    protected $fillable = [
        'book_id',
        'page_number',
        'content',
    ];

    /**
     * Get the book that owns the Page
     *
     * @return \Illuminate\Database\Eloquent\Relations\BelongsTo
     */
    public function book(): BelongsTo
    {
        return $this->belongsTo(Book::class);
    }

    /**
     * Get the interaction associated with the Page
     *
     * @return \Illuminate\Database\Eloquent\Relations\HasOne
     */
    public function interaction(): HasOne
    {
        return $this->hasOne(Interaction::class);
    }
}
