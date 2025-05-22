<?php

namespace App\Http\Requests\Book;

use Illuminate\Foundation\Http\FormRequest;

/**
 * Save Progress Book Request Form
 *
 * - Use this request to validate the progress book form.
 * - The `book_id` is required and must be exists in the `books` table.
 * - The `group_id` is optional and must be exists in the `groups` table.
 * - The `correct` and `incorrect` are required and must be integers.
 * - The `start_at` and `finish_at` are required and must be in the format of `Y-m-d H:i:s`.
 */
class SaveProgressRequest extends FormRequest
{
    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, mixed>
     */
    public function rules(): array
    {
        return [
            'book_id' => [
                'required',
                'integer',
                'exists:books,id',
            ],
            'group_id' => [
                'nullable',
                'exists:groups,id',
            ],
            'correct' => [
                'required',
                'integer',
                'min:0',
            ],
            'incorrect' => [
                'required',
                'integer',
                'min:0',
            ],
            'start_at' => [
                'required',
                'date_format:Y-m-d H:i:s',
            ],
            'finish_at' => [
                'required',
                'date_format:Y-m-d H:i:s',
            ],
        ];
    }
}
