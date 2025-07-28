<?php

namespace App\Traits;

use Illuminate\Http\JsonResponse;

trait ApiResponseTrait
{
    /**
     * JSON response for successful operations.
     *
     * @param mixed $data
     * @param string $message
     * @param int $statusCode
     * @return \Illuminate\Http\JsonResponse
     */
    public function sendSuccess($data = null, string $message = 'Success', int $statusCode = 200): JsonResponse
    {
        return response()->json([
            'status'  => $statusCode,
            'message' => $message,
            'data'    => $data,
        ], $statusCode);
    }

    /**
     * JSON response for error operations.
     *
     * @param string $message
     * @param int $statusCode
     * @return \Illuminate\Http\JsonResponse
     */
    public function sendError(string $message = 'Error', int $statusCode = 400): JsonResponse
    {
        return response()->json([
            'status'  => $statusCode,
            'message' => $message,
            'data'    => null,
        ], $statusCode);
    }
}
