<?php

namespace App\Http\Requests\Auth;

use App\Enums\PermissionEnum as PE;
use App\Traits\Api\AuthorizeTrait;
use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rules\Password;

/**
 * Update Password Request Form
 *
 * - Use this request to validate the change password form.
 * - The regex pattern allows:
 * > letters (a-z A-Z) , numbers (0-9) , dots (.) , underscores (_) , hyphens (-)
 */
class UpdatePasswordRequest extends FormRequest
{
    use AuthorizeTrait;

    /**
     * Determine if the user is authorized to make this request.
     *
     * @return bool
     */
    public function authorize(): bool
    {
        return $this->grant(PE::CHANGE_PASSWORD_SELF->value);
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, mixed>
     */
    public function rules(): array
    {
        return [
            'old_password' => [
                'required',
                'string',
                Password::min(8)
                    ->max(20)
                    ->letters()
                    ->numbers(),
            ],
            'new_password' => [
                'required',
                'string',
                Password::min(8)
                    ->max(20)
                    ->letters()
                    ->numbers(),
                'confirmed'
            ],
        ];
    }
}
