<?php

namespace App\Http\Requests\Book;

use Illuminate\Foundation\Http\FormRequest;
use Illuminate\Validation\Rule;

class StoreBookRequest extends FormRequest
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
            'title' => [
                'required',
                'string',
                'max:255',
                Rule::unique('books', 'title'),
            ],
            'author' => [
                'required',
                'string',
                'max:255',
            ],
            'year' => [
                'required',
                'integer',
                'digits:4',
            ],
            'genre' => [
                'nullable',
                'string',
                'max:255',
            ],
            'focus' => [
                'required',
                'string',
                'max:255',
            ],
            'cover_image' => [
                'nullable',
                'url',
            ],
            'order_sequence' => [
                'required',
                'integer',
            ],
        ];
    }
}
