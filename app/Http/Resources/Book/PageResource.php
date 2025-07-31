<?php

namespace App\Http\Resources\Book;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class PageResource extends JsonResource
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
            'page_number' => $this->page_number,
            'content' => $this->when($this->page_number === 0, $this->content),
            'image' => $this->when($this->page_number === 0, $this->book->cover_image_url),
            'interaction' => new InteractionResource($this->whenLoaded('interaction')),
        ];
    }
}
