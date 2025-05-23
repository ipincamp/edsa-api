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
     * Store the progress of a book.
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

            return $this->json(
                message: 'Book progress updated successfully',
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Display the specified resource.
     */
    public function show(Book $book)
    {
        //
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
