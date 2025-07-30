<?php

namespace App\Http\Resources\Book;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class StudentProgressResource extends JsonResource
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
            'status' => $this->status,
            'total_points' => $this->total_points,
            'last_page' => $this->last_page,
            'latest_page' => $this->latest_page,
            'total_pages' => $this->book->pages_count,
            'grade' => $this->total_points > 0 ? (int) ceil($this->total_points / 20) : null,
        ];
    }
}
