<?php

namespace App\Http\Resources;

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
            'role' => $this->roles()->pluck('name')->first(),
            /**
             * The course of the user.
             * @var string
             */
            'course' => $this->groups->pluck('course.name')->first(),
            /**
             * The class of the user.
             * @var string
             */
            'class' => $this->groups()->pluck('name')->first(),
            /**
             * Indicates if the user active.
             * @var bool
             */
            'is_active' => $this->deleted_at === null,
            /**
             * Since when the user registered.
             * @var string
             * @format Y-m-d H:i:s
             * @example 2025-10-01 12:00:00
             */
            'joined_at' => $this->created_at->format('Y-m-d H:i:s'),
            /**
             * The last time the user updated their profile.
             * @var string
             * @format Y-m-d H:i:s
             * @example 2025-10-01 12:30:00
             */
            'last_update' => $this->updated_at->format('Y-m-d H:i:s'),
        ];
    }
}
