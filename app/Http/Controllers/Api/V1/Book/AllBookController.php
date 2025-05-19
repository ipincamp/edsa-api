<?php

namespace App\Http\Controllers\Api\V1\Book;

use App\Http\Controllers\Controller;
use App\Models\Book;
use Dedoc\Scramble\Attributes\Group;
use Illuminate\Http\JsonResponse;

#[Group('Books')]
class AllBookController extends Controller
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
    public function __invoke(): JsonResponse
    {
        try {
            $books = Book::all();

            return $this->json(
                message: 'All books retrieved successfully',
                data: $books,
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }
}
