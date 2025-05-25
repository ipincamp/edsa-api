<?php

namespace App\Http\Controllers\Api\V1\Book;

use App\Http\Controllers\Controller;
use App\Http\Requests\Book\SaveProgressRequest;
use App\Http\Resources\BookResource;
use App\Models\Book;
use Dedoc\Scramble\Attributes\Group;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

#[Group('Books')]
class BookController extends Controller
{
    /**
     * Random Word
     *
     * Retrieve a random word from the book's settings.
     *
     * @param Book $book
     * @param string $key
     * @return array
     */
    private function randomWord(Book $book, string $key = 'random-word-focus'): array
    {
        $words = $book->settings()
            ->where('key', $key)
            ->pluck('value')
            ->toArray();

        if (empty($words)) {
            return [];
        }

        $words = array_map('trim', explode(',', $words[0]));
        $original = $words[array_rand($words)];
        $random = str_split($original);
        shuffle($random);
        $random = implode('', $random);
        $random = str_split($random);

        return [
            'original' => $original,
            'shuffle' => $random,
        ];
    }

    /**
     * All books
     *
     * Retrieve all books stored.
     *
     * @operationId getAllBooks
     * @authenticated
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(): JsonResponse
    {
        try {
            return $this->json(
                message: 'All books retrieved successfully',
                data: BookResource::collection(Book::all()),
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Save progress
     *
     * Store the progress of a book. Then retrieve a random word from the book's settings.
     *
     * @operationId saveBookProgress
     * @authenticated
     * @param SaveProgressRequest $request
     * @param Book $book
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function store(SaveProgressRequest $request, Book $book): JsonResponse
    {
        try {
            $user = auth()->user();
            $group = $user->groups->first();
            $input = $request->validated();

            // See how many takes the user has
            $taken = $user
                ->progress()
                ->where('book_id', $book->id)
                ->count();

            $user->progress()->updateOrCreate(
                ['book_id' => $book->id, 'taken' => $taken + 1], // Ensure unique for each entry
                [
                    'group_id' => $group->id ?? null,
                    'score_correct' => $input['correct'],
                    'score_incorrect' => $input['incorrect'],
                    'time_start' => $input['start_at'],
                    'time_finish' => $input['finish_at'],
                ]
            );

            $word = $this->randomWord(book: $book);

            return $this->json(
                message: 'Book progress updated successfully',
                data: [
                    'word' => $word,
                ]
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Random word other
     *
     * Retrieve a random word from the book's settings, specifically for other words.
     *
     * @operationId getRandomWordOther
     * @param Book $book
     * @return \Illuminate\Http\JsonResponse
     */
    public function show(Book $book): JsonResponse
    {
        try {
            $word = $this->randomWord(book: $book, key: 'random-word-other');

            return $this->json(
                message: 'Random word retrieved successfully',
                data: [
                    'word' => $word,
                ]
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Update the specified resource in storage.
     */
    public function update(Request $request, Book $book)
    {
        //
    }

    /**
     * Remove the specified resource from storage.
     */
    public function destroy(Book $book)
    {
        //
    }
}
