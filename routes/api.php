<?php

use App\Http\Controllers\Api\Auth\AuthController;
use App\Http\Controllers\Api\Book\BookController;
use Illuminate\Support\Facades\Route;

// Public routes
Route::post('auth/login', [AuthController::class, 'login'])->name('auth.login');
Route::post('auth/register', [AuthController::class, 'register'])->name('auth.register');

// Protected routes
Route::middleware('auth:sanctum')->group(function () {
    // Auth
    Route::prefix('auth')->controller(AuthController::class)->group(function () {
        Route::get('/profile', 'profile')->name('auth.profile');
        Route::patch('/profile/update-details', 'updateDetail')->name('auth.profile.update-details');
        Route::patch('/profile/update-password', 'updatePassword')->name('auth.profile.update-password');
        Route::post('/logout', 'logout')->name('auth.logout');
    });

    // Book
    Route::prefix('books')->controller(BookController::class)->group(function () {
        Route::get('/', 'index')->name('books.index');
        Route::get('/{book}', 'show')->name('books.show');

        Route::middleware('role:admin')->group(function () {
            Route::post('/', 'store')->name('books.store');
            Route::patch('/{book}', 'update')->name('books.update');
            Route::delete('/{book}', 'destroy')->name('books.destroy');
        });
    });
});
