<?php

namespace App\Http\Requests\Book\Progress;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class UpdateLastPageRequest extends FormRequest
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
            'book_id' => [
                'required',
                Rule::exists('books', 'id'),
            ],
            'page_number' => [
                'required',
                'integer',
                Rule::exists('pages', 'page_number')->where(function ($query) {
                    $query->where('book_id', $this->book_id);
                }),
            ],
        ];
    }
}
