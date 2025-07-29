<?php

namespace App\Http\Resources\Group;

use App\Http\Resources\Auth\AuthResource;
use App\Http\Resources\Course\CourseResource;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class GroupResource extends JsonResource
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
            'course' => new CourseResource($this->whenLoaded('course')),
            'teacher' => new AuthResource($this->whenLoaded('teacher')),
            'students' => AuthResource::collection($this->whenLoaded('students')),
            'students_count' => $this->whenCounted('students', $this->students_count),
        ];
    }
}
