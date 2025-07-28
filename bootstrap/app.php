<?php

use Illuminate\Foundation\Application;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;
use Illuminate\Http\Request;
use Illuminate\Validation\ValidationException;
use Illuminate\Auth\AuthenticationException;
use Spatie\Permission\Exceptions\UnauthorizedException;
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;
use Symfony\Component\HttpKernel\Exception\HttpExceptionInterface;

return Application::configure(basePath: dirname(__DIR__))
    ->withRouting(
        web: __DIR__ . '/../routes/web.php',
        api: __DIR__ . '/../routes/api.php',
        commands: __DIR__ . '/../routes/console.php',
        health: '/up',
    )
    ->withMiddleware(function (Middleware $middleware): void {
        $middleware->alias([
            'role' => \Spatie\Permission\Middleware\RoleMiddleware::class,
            'permission' => \Spatie\Permission\Middleware\PermissionMiddleware::class,
            'role_or_permission' => \Spatie\Permission\Middleware\RoleOrPermissionMiddleware::class,
        ]);
        $middleware->api(
            prepend: [
                \App\Http\Middleware\EnsureJsonResponse::class,
            ],
        );
    })
    ->withExceptions(function (Exceptions $exceptions): void {
        $exceptions->render(function (Throwable $e, Request $request) {
            if ($request->is('api/*') || $request->expectsJson()) {
                $statusCode = 500;
                $response = [
                    'success' => false,
                    'message' => 'There was an error processing your request.',
                    'data' => null,
                ];

                switch (true) {
                    case $e instanceof ValidationException:
                        $statusCode = 422;
                        $response['message'] = 'The provided data is invalid.';
                        $response['errors'] = $e->errors();
                        break;
                    case $e instanceof AuthenticationException:
                        $statusCode = 401;
                        $response['message'] = 'Unauthenticated. Please log in first.';
                        break;
                    case $e instanceof UnauthorizedException:
                        $statusCode = 403;
                        $response['message'] = 'You do not have permission to perform this action.';
                        break;
                    case $e instanceof NotFoundHttpException:
                        $statusCode = 404;
                        $response['message'] = 'The endpoint or data you are looking for was not found.';
                        break;
                    default:
                        if ($e instanceof HttpExceptionInterface) {
                            $statusCode = $e->getStatusCode();
                        } else {
                            $statusCode = 500;
                        }

                        if (config('app.debug')) {
                            $response['message'] = $e->getMessage();
                            $response['exception'] = get_class($e);
                        }
                        break;
                }

                if ($statusCode < 100 || $statusCode >= 600) {
                    $statusCode = 500;
                }

                return response()->json(
                    data: $response,
                    status: $statusCode,
                );
            }
        });
    })->create();
