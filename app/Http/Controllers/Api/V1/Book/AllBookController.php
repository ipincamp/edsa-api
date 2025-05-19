<?php

namespace App\Http\Controllers\Api\V1\Book;

use App\Http\Controllers\Controller;
use App\Http\Resources\BookResource;
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
            return $this->json(
                message: 'All books retrieved successfully',
                data: BookResource::collection(Book::all()),
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }
}
