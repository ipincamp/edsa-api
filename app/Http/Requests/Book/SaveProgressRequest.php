<?php

namespace App\Http\Requests\Book;

use App\Enums\RoleEnum;
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
     * Determine if the user is authorized to make this request.
     *
     * @return bool
     */
    public function authorize(): bool
    {
        return $this->user()->hasRole(RoleEnum::STUDENT->value);
    }

    /**
     * Get the validation rules that apply to the request.
     *
     * @return array<string, mixed>
     */
    public function rules(): array
    {
        return [
            /**
             * Score correct.
             * @example 10
             */
            'correct' => [
                'required',
                'integer',
                'min:0',
            ],
            /**
             * Score incorrect.
             * @example 2
             */
            'incorrect' => [
                'required',
                'integer',
                'min:0',
            ],
            /**
             * The start and finish date of the book progress.
             * @example 2025-10-20 12:00:00
             */
            'start_at' => [
                'required',
                'date_format:Y-m-d H:i:s',
            ],
            /**
             * The start and finish date of the book progress.
             * @example 2025-10-20 12:30:00
             */
            'finish_at' => [
                'required',
                'date_format:Y-m-d H:i:s',
            ],
        ];
    }
}
