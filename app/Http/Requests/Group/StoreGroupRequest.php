<?php

namespace App\Http\Requests\Group;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class StoreGroupRequest extends FormRequest
{
    /**
     * Determine if the user is authorized to make this request.
     */
    public function authorize(): bool
    {
        return true;
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, \Illuminate\Contracts\Validation\ValidationRule|array<mixed>|string>
     */
    public function rules(): array
    {
        return [
            'name' => [
                'required',
                'string',
                'max:255',
            ],
            'course_id' => [
                'required',
                Rule::exists('courses', 'id'),
            ],
            'teacher_ids' => [
                'required',
                'array',
            ],
            'teacher_ids.*' => [
                'required',
                Rule::exists('users', 'id'),
            ],
        ];
    }
}
