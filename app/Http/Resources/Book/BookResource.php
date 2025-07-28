<?php

namespace App\Http\Resources\Book;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

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
            'id' => $this->id,
            'title' => $this->title,
            'author' => $this->author,
            'year' => $this->year,
            'genre' => $this->genre,
            'focus' => $this->focus,
            'status' => $this->status,
            'cover_image' => $this->cover_image,
            'order_sequence' => $this->order_sequence,
            'pages' => PageResource::collection($this->whenLoaded('pages')),
            'post_activities' => PostActivityResource::collection($this->whenLoaded('postActivities')),
        ];
    }
}
