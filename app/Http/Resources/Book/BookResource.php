<?php

namespace App\Http\Resources\Book;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;
use App\Models\Page;

class BookResource extends JsonResource
{
    /**
     * Transform the resource into an array.
     *
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        $pagesCollection = $this->whenLoaded('pages');

        $coverPage = new Page();
        $coverPage->id = 0;
        $coverPage->page_number = 0;
        $coverPage->content = "{'title':'{$this->title}','author':'{$this->author}'}";
        $coverPage->setRelation('book', $this->resource);
        $coverPage->setRelation('interaction', null);

        $allPagesWithCover = $pagesCollection->prepend($coverPage);

        return [
            'id' => $this->id,
            'title' => $this->title,
            'author' => $this->author,
            'cover_image' => $this->cover_image_url,
            'order_sequence' => $this->order_sequence,
            'is_locked' => $this->when(isset($this->is_locked), $this->is_locked),
            'latest_page' => $this->when(isset($this->latest_page), $this->latest_page),
            'last_page' => $this->when(isset($this->last_page), $this->last_page),
            'total_pages' => $this->pages_count + 1,
            'pages' => PageResource::collection($allPagesWithCover),
            'post_activities' => PostActivityResource::collection($this->whenLoaded('postActivities')),
        ];
    }
}
