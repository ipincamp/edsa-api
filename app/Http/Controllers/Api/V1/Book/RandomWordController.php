<?php

namespace App\Http\Controllers\Api\V1\Book;

use App\Http\Controllers\Controller;
use App\Models\Book;
use Dedoc\Scramble\Attributes\Group;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

#[Group('Books')]
class RandomWordController extends Controller
{
    /**
     * Random Word
     *
     * Retrieve a random word from the book's settings.
     *
     * @param Request $request
     * @param Book $book
     * @return \Illuminate\Http\JsonResponse
     */
    public function __invoke(Request $request, Book $book): JsonResponse
    {
        try {
            $words = $book->settings()
                ->where('key', 'ramdom-word')
                ->pluck('value')
                ->toArray();

            if (empty($words)) {
                return $this->json(
                    message: 'No random words found. Please add some words to the book settings.',
                    data: [],
                );
            }

            $words = array_map('trim', explode(',', $words[0]));
            $original = $words[array_rand($words)];
            $random = str_split($original);
            shuffle($random);
            $random = implode('', $random);
            $random = str_split($random);

            return $this->json(
                message: 'Random word retrieved successfully',
                data: [
                    'original' => $original,
                    'shuffle' => $random,
                ],
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }
}
