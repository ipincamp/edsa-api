<?php

namespace App\Http\Resources;

use App\Models\Book;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/**
 * @mixin Book
 */
class BookResource extends JsonResource
{
    /**
     * Transform the resource into an array.
     *
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        return [
            /**
             * The unique identifier for the book.
             * @var int
             */
            'id' => $this->id,
            /**
             * The title of the book.
             * @var string
             */
            'title' => $this->title,
            /**
             * The author of the book.
             * @var string
             */
            'author' => $this->author,
            /**
             * The publication year of the book.
             * @var int
             * @example 2025
             */
            'year' => $this->year,
            /**
             * The genre of the book.
             * @var string
             */
            'genre' => $this->genre,
            /**
             * The focus of the book.
             * @var string
             */
            'focus' => $this->focus,
            /**
             * The status of the book.
             * If true the book is unlock, if false the book is locked.
             * @var bool
             */
            'status' => boolval($this->status),
            /**
             * The image URL of the book cover.
             * @var string
             */
            'image' => $this->image,
            /**
             * The date when the book was created.
             * @var string
             * @format date-time
             * @example 2025-05-19 14:19:02
             */
            'created_at' => $this->created_at->format('Y-m-d H:i:s'),
            /**
             * The date when the book was last updated.
             * @var string
             * @format date-time
             * @example 2025-05-19 14:20:00
             */
            'updated_at' => $this->updated_at->format('Y-m-d H:i:s'),
            /**
             * The date when the book was deleted.
             * @var string|null
             * @format date-time
             */
            'deleted_at' => $this->deleted_at ? $this->deleted_at->format('Y-m-d H:i:s') : null,
        ];
    }
}
