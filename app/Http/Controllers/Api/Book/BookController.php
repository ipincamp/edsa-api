<?php

namespace App\Http\Controllers\Api\Book;

use App\Enums\RolesEnum;
use App\Http\Controllers\Controller;
use App\Http\Requests\Book\StoreBookRequest;
use App\Http\Requests\Book\UpdateBookRequest;
use App\Http\Resources\Book\BookResource;
use App\Models\Book;
use App\Models\StudentProgress;
use Illuminate\Http\Response;
use Illuminate\Support\Facades\Auth;

class BookController extends Controller
{
    // Get books list with pagination
    public function index()
    {
        try {
            $user = Auth::user();
            $allBooks = Book::orderBy('order_sequence', 'asc')->get();

            if ($user->hasRole(RolesEnum::A->value) || $user->hasRole(RolesEnum::T->value)) {
                return $this->sendSuccess(
                    message: 'Books retrieved successfully.',
                    data: BookResource::collection($allBooks),
                );
            }

            $studentProgress = StudentProgress::where('student_id', $user->id)
                ->where('status', 'completed')
                ->pluck('book_id');

            $lastCompletedBookOrder = Book::whereIn('id', $studentProgress)
                ->max('order_sequence') ?? 0;

            $unlockedBooks = $allBooks->filter(function ($book) use ($lastCompletedBookOrder) {
                return $book->order_sequence <= $lastCompletedBookOrder + 1;
            });

            $booksWithLockStatus = $allBooks->map(function ($book) use ($unlockedBooks) {
                $book->is_locked = !$unlockedBooks->contains('id', $book->id);
                return $book;
            });

            return $this->sendSuccess(
                message: 'Books retrieved successfully.',
                data: BookResource::collection($booksWithLockStatus),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'An error occurred while fetching the books.',
                statusCode: 500
            );
        }
    }

    // create a new book
    public function store(StoreBookRequest $request)
    {
        try {
            $book = Book::create($request->validated());

            return $this->sendSuccess(
                message: 'Book created successfully.',
                data: new BookResource($book),
                statusCode: Response::HTTP_CREATED
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'An error occurred while creating the book.',
                statusCode: 500
            );
        }
    }

    // Get book detail
    public function show(Book $book)
    {
        try {
            // If the book is not 'published', only admins can view it
            if ($book->status !== 'published' && !Auth::user()->hasRole(RolesEnum::S->value)) {
                return $this->sendError(
                    statusCode: Response::HTTP_NOT_FOUND,
                    message: 'Book not found or not accessible.',
                );
            }

            $book->load(['pages.interaction', 'postActivities']);

            return $this->sendSuccess(
                message: 'Book retrieved successfully.',
                data: new BookResource($book),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'An error occurred while retrieving the book.',
                statusCode: 500
            );
        }
    }

    // Update book
    public function update(UpdateBookRequest $request, Book $book)
    {
        try {
            $book->update($request->validated());

            return $this->sendSuccess(
                message: 'Book updated successfully.',
                data: new BookResource($book),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'An error occurred while updating the book.',
                statusCode: 500
            );
        }
    }

    // Delete book
    public function destroy(Book $book)
    {
        try {
            $book->delete();

            return $this->sendSuccess(
                statusCode: Response::HTTP_NO_CONTENT,
                message: 'Book deleted successfully.',
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'An error occurred while deleting the book.',
                statusCode: 500
            );
        }
    }
}
