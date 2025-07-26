<?php

namespace App\Traits\Api;

trait JsonResponseTrait
{
    /**
     * Send a JSON response.
     *
     * @param string $message
     * @param int $code
     * @param mixed $data
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function json(
        $message,
        $code = 200,
        $data = null,
    ) {
        $success = $code >= 200 && $code < 300;
        return response()
            ->json([
                'status' => $success,
                'meta' => [
                    'code' => $code,
                    'message' => $message,
                    'data' => $data,
                ],
                'timestamp' => now()->timestamp,
            ], $code);
    }
}
