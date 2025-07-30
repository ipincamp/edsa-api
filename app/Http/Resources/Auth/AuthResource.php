<?php

namespace App\Http\Resources\Auth;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class AuthResource extends JsonResource
{
    /**
     * Token otentikasi.
     *
     * @var string|null
     */
    protected ?string $token;

    /**
     * Create a new resource instance.
     *
     * @param  mixed  $resource
     * @param  string|null  $token
     * @return void
     */
    public function __construct($resource, ?string $token = null)
    {
        parent::__construct($resource);
        $this->token = $token;
    }

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
            'email' => $this->email,
            'role' => $this->role,
            'token' => $this->when($this->token, $this->token),
            'joined_at' => $this->created_at->format('Y-m-d H:i:s'),
        ];
    }
}
