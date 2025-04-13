<?php

namespace App\Http\Requests\Auth;

use Illuminate\Foundation\Http\FormRequest;

/**
 * Login Request Form
 *
 * - Use this request to validate the login form.
 * - The regex pattern allowed for the `username` is:
 * > letters (a-z A-Z) , numbers (0-9) , dots (.) , at (@)
 */
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
            'username' => ['required', 'string', 'regex:/^[a-zA-Z0-9.@]+$/', 'min:3', 'max:20'],
            'password' => ['required', 'string', 'min:8', 'max:20'],
            'remember' => ['nullable', 'boolean'],
        ];
    }
}
