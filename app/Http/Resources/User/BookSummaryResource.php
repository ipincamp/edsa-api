<?php

namespace App\Http\Resources\User;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class BookSummaryResource extends JsonResource
{
    /**
     * Transform the resource into an array.
     *
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        $progressPercentage = 0;
        if ($this->pages_count > 0) {
            // (halaman terjauh / total halaman) * 100
            $progressPercentage = round(($this->latest_page / $this->pages_count) * 100, 2);
        }

        return [
            'id' => $this->id,
            'title' => $this->title,
            'cover_image' => $this->cover_image_url,
            'order_sequence' => $this->order_sequence,
            'is_locked' => $this->when(isset($this->is_locked), $this->is_locked, false),
            'progress' => $progressPercentage,
            'grade' => (int) ceil($this->when(isset($this->total_points), $this->total_points, 0) / 20),
        ];
    }
}
