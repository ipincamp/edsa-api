<?php

namespace App\Http\Controllers\Api\Book;

use App\Enums\RolesEnum;
use App\Http\Controllers\Controller;
use App\Http\Requests\Book\StoreBookRequest;
use App\Http\Requests\Book\UpdateBookRequest;
use App\Http\Resources\Book\BookResource;
use App\Models\Book;
use Illuminate\Http\Response;
use Illuminate\Support\Facades\Auth;

class BookController extends Controller
{
    // Get books list with pagination
    public function index()
    {
        try {
            $user = Auth::user();
            $query = Book::query()->orderBy('order_sequence', 'asc');

            // If the user is not an admin, filter out unpublished books
            if ($user->role !== RolesEnum::A->value) {
                $query->where('status', 'published');
            }

            return $this->sendSuccess(
                message: 'Books retrieved successfully.',
                data: BookResource::collection($query->paginate(10)),
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
            if ($book->status !== 'published' && Auth::user()->role !== RolesEnum::A->value) {
                return response()->json(['message' => 'Not Found.'], 404);
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
