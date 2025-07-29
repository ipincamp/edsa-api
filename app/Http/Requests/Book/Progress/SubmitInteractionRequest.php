<?php

namespace App\Http\Requests\Book\Progress;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class SubmitInteractionRequest extends FormRequest
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
            'interaction_id' => [
                'required',
                Rule::exists('interactions', 'id'),
            ],
            'is_correct' => [
                'required',
                'boolean',
            ],
        ];
    }
}
