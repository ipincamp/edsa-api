<?php

namespace App\Http\Requests\Auth;

use Illuminate\Foundation\Http\FormRequest;

/**
 * Change Password Request Form
 *
 * - Use this request to validate the change password form.
 * - The regex pattern allows:
 * > letters (a-z A-Z) , numbers (0-9) , dots (.) , underscores (_) , hyphens (-)
 */
class ChangePasswordRequest extends FormRequest
{
    /**
     * Determine if the user is authorized to make this request.
     *
     * @return bool
     */
    public function authorize(): bool
    {
        return $this->user()->checkPermissionTo('change-password-self');
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, mixed>
     */
    public function rules(): array
    {
        return [
            'current_password' => ['required', 'string', 'min:8', 'max:20'],
            'new_password' => ['required', 'string', 'regex:/^[a-zA-Z0-9._-]+$/', 'min:8', 'max:20', 'confirmed'],
        ];
    }
}
