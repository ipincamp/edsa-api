<?php

namespace App\Http\Resources\Book;

use App\Http\Resources\Auth\AuthResource;
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
            'book' => new BookResource($this->whenLoaded('book')),
            'student' => new AuthResource($this->whenLoaded('student')),
        ];
    }
}
