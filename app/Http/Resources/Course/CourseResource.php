<?php

namespace App\Http\Resources\Course;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class CourseResource extends JsonResource
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
            'name' => $this->name,
            'description' => $this->description,
            'created_at' => $this->created_at->toDateTimeString(),
            // Memuat relasi group jika diminta melalui parameter ?include=groups
            // 'groups' => GroupResource::collection($this->whenLoaded('groups')),
            // 'groups_count' => $this->whenCounted('groups'),
        ];
    }
}
