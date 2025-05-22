<?php

namespace App\Http\Requests\Auth;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rules\Password;

/**
 * Register Request Form
 *
 * - Use this request to validate the register form.
 * - The regex pattern allowed for the `username` is:
 * > letters (a-z A-Z) , numbers (0-9) , dots (.)
 */
class RegisterRequest extends FormRequest
{
    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, mixed>
     */
    public function rules(): array
    {
        return [
            'name' => [
                'required',
                'string',
                'max:255',
            ],
            'username' => [
                'required',
                'string',
                'regex:/^[a-zA-Z0-9.]+$/',
                'unique:users,username',
                'min:3',
                'max:20',
            ],
            'password' => [
                'required',
                'string',
                Password::min(8)
                    ->max(20)
                    ->letters()
                    ->numbers(),
            ],
        ];
    }
}
