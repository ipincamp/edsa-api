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
            'cover_image' => $this->cover_image_url,
            'order_sequence' => $this->order_sequence,
            'is_locked' => $this->when(isset($this->is_locked), $this->is_locked),
            'latest_page' => $this->latest_page,
            'last_page' => $this->when(isset($this->last_page), $this->last_page),
            'total_pages' => $this->pages_count,
            'pages' => PageResource::collection($this->whenLoaded('pages')),
            'post_activities' => PostActivityResource::collection($this->whenLoaded('postActivities')),
        ];
    }
}
