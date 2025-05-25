<?php

use App\Http\Controllers\Api\V1\Auth\AuthController;
use App\Http\Controllers\Api\V1\Book\BookController;
use Illuminate\Support\Facades\Route;

Route::prefix('v1')->group(function () {
    Route::prefix('auth')->controller(AuthController::class)->middleware('auth:sanctum')->group(function () {
        Route::post('login', 'login')->withoutMiddleware('auth:sanctum')->name('auth.signIn');
        Route::post('register', 'register')->withoutMiddleware('auth:sanctum')->name('auth.signUp');
        Route::get('profile', 'profile')->name('auth.profile');
        Route::post('update-password', 'updatePassword')->name('auth.updatePassword');
        Route::post('update-profile', 'updateProfile')->name('auth.updateProfile');
        Route::post('logout', 'logout')->name('auth.signOut');
    });

    Route::prefix('books')->controller(BookController::class)->middleware('auth:sanctum')->group(function () {
        Route::get('/', 'index')->name('books.all');
        Route::post('/{book}/progress', 'store')->name('books.saveProgress');
        Route::get('/{book}/random-word', 'show')->name('books.randomWordOther');
    });
});
