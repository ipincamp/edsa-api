<?php

namespace App\Http\Resources;

use App\Enums\RoleEnum;
use App\Models\User;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/**
 * @mixin User
 */
class UserResource extends JsonResource
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
             * The unique identifier for the user.
             * @var string
             * @format uuid
             */
            'id' => $this->id,
            /**
             * The name of the user.
             * @var string
             */
            'name' => $this->name,
            /**
             * The username of the user.
             * @var string
             */
            'username' => $this->username,
            /**
             * The role of the user.
             * @var string
             */
            'role' => $this->getRole(),
            /**
             * The courses the user is enrolled in.
             * If the user is a student, only the first course will be returned.
             * @var array
             * @example [{"id": "1", "name": "Course 1", "description": "Description of Course 1"}]
             */
            'course' => $this->getCourses(),
            /**
             * The group the user is enrolled in.
             * If the user is a student, only the first group will be returned.
             * @var array
             * @example [{"id": "1", "name": "Group 1", "description": "Description of Group 1"}]
             */
            'class' => $this->getClasses(),
            /**
             * Indicates if the user is active or not.
             * @var bool
             */
            'is_active' => $this->isActive(),
            /**
             * When the user joined the system.
             * @var string
             */
            'joined_at' => $this->formatDate($this->created_at),
            /**
             * When the user was last updated.
             * @var string
             */
            'last_update' => $this->formatDate($this->updated_at),
        ];
    }

    /**
     * Get the role of the user.
     *
     * @return string
     */
    private function getRole(): string
    {
        return $this->roles()->pluck('name')->first();
    }

    /**
     * Get the courses of the user.
     *
     * @return array|string[]
     */
    private function getCourses(): array|string
    {
        $courses = $this->groups->map(fn($group) => [
            'id' => $group->course->id,
            'name' => $group->course->name,
            'description' => $group->course->description,
        ]);

        return $this->getRole() === RoleEnum::STUDENT->value ? ($courses->first() ?? []) : $courses->toArray();
    }

    /**
     * Get the classes of the user.
     *
     * @return array|string[]
     */
    private function getClasses(): array|string
    {
        $classes = $this->groups->map(fn($group) => [
            'id' => $group->id,
            'name' => $group->name,
            'description' => $group->description,
        ]);

        return $this->getRole() === RoleEnum::STUDENT->value ? ($classes->first() ?? []) : $classes;
    }

    /**
     * Check if the user is active.
     *
     * @return bool
     */
    private function isActive(): bool
    {
        return $this->deleted_at === null;
    }

    /**
     * Format a date to 'Y-m-d H:i:s'.
     *
     * @param \Illuminate\Support\Carbon|null $date
     * @return string|null
     */
    private function formatDate($date): ?string
    {
        return $date ? $date->format('Y-m-d H:i:s') : null;
    }
}
