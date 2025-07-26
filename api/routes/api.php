<?php

use App\Http\Controllers\Api\V1\Auth\AuthController;
use App\Http\Controllers\Api\V1\Book\BookController;
use App\Http\Controllers\Api\V1\User\StudentProgressController;
use Illuminate\Support\Facades\Route;

Route::prefix('v1')->group(function () {
    // Authentication Routes
    Route::prefix('auth')->controller(AuthController::class)->group(function () {
        // Public
        Route::post('login', 'login')->name('auth.signIn');
        Route::post('register', 'register')->name('auth.signUp');

        // Protected
        Route::middleware('auth:sanctum')->group(function () {
            Route::get('profile', 'profile')->name('auth.profile');
            Route::post('update-password', 'updatePassword')->name('auth.updatePassword');
            Route::post('update-profile', 'updateProfile')->name('auth.updateProfile');
            Route::post('logout', 'logout')->name('auth.signOut');
        });
    });

    Route::middleware('auth:sanctum')->group(function () {
        // Books Routes
        Route::prefix('books')->controller(BookController::class)->group(function () {
            Route::get('/', 'index')->name('books.all');
            Route::post('/{book}/progress', 'store')->name('books.saveProgress');
            Route::get('/{book}/random-word', 'show')->name('books.randomWordOther');
        });

        // User Progress Routes
        Route::prefix('users')->controller(StudentProgressController::class)->group(function () {
            Route::get('/progress', 'index')->name('users.progress.index');
            Route::get('/progress/{user}', 'show')->name('users.progress.show');
        });
    });


    // Route::prefix('auth')->controller(AuthController::class)->middleware('auth:sanctum')->group(function () {
    //     Route::post('login', 'login')->withoutMiddleware('auth:sanctum')->name('auth.signIn');
    //     Route::post('register', 'register')->withoutMiddleware('auth:sanctum')->name('auth.signUp');
    //     Route::get('profile', 'profile')->name('auth.profile');
    //     Route::post('update-password', 'updatePassword')->name('auth.updatePassword');
    //     Route::post('update-profile', 'updateProfile')->name('auth.updateProfile');
    //     Route::post('logout', 'logout')->name('auth.signOut');
    // });

    // Route::prefix('books')->controller(BookController::class)->middleware('auth:sanctum')->group(function () {
    //     Route::get('/', 'index')->name('books.all');
    //     Route::post('/{book}/progress', 'store')->name('books.saveProgress');
    //     Route::get('/{book}/random-word', 'show')->name('books.randomWordOther');
    // });

    // Route::prefix('users')->controller(StudentProgressController::class)->middleware('auth:sanctum')->group(function () {
    //     Route::get('/progress', 'index')->name('users.progress.index');
    //     Route::get('/progress/{user}', 'show')->name('users.progress.show');
    // });
});
