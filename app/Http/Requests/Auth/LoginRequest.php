<?php

namespace App\Http\Requests\Auth;

use Illuminate\Foundation\Http\FormRequest;

class LoginRequest extends FormRequest
{
    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, mixed>
     */
    public function rules(): array
    {
        return [
            'username' => ['required', 'string', 'regex:/^[a-zA-Z0-9._-]+$/', 'min:3', 'max:20'],
            'password' => ['required', 'string', 'min:8', 'max:20'],
            'remember' => ['nullable', 'boolean'],
        ];
    }
}
