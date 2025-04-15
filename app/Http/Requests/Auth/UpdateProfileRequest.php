<?php

namespace App\Http\Requests\Auth;

use Illuminate\Foundation\Http\FormRequest;

/**
 * Update Profile Request Form
 *
 * - Use this request to validate the update profile form.
 * - The regex pattern allowed for `username`:
 * > letters (a-z A-Z) , numbers (0-9) , dots (.) , at (@)
 */
class UpdateProfileRequest extends FormRequest
{
    /**
     * Determine if the user is authorized to make this request.
     *
     * @return bool
     */
    public function authorize(): bool
    {
        return $this->user()->checkPermissionTo('update-profile-self');
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, mixed>
     */
    public function rules(): array
    {
        $id = optional($this->user())->id;

        return [
            'name' => ['required', 'string', 'max:50'],
            'username' => ['required', 'string', 'min:3', 'max:20', 'regex:/^[a-zA-Z0-9.@]+$/', 'unique:users,username,' . ($id ?? 'null')],
        ];
    }
}
