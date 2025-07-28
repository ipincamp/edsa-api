<?php

use App\Http\Controllers\Api\Auth\AuthController;
use Illuminate\Support\Facades\Route;

Route::prefix('auth')->controller(AuthController::class)->group(function () {
    // Public routes
    Route::post('/login', 'login')->name('auth.login');
    Route::post('/register', 'register')->name('auth.register');

    // Protected routes
    Route::middleware(['auth:sanctum'])->group(function () {
        Route::get('/profile', 'profile')->name('auth.profile');
        Route::patch('/profile/update-details', 'updateDetail')->name('auth.profile.update-details');
        Route::patch('/profile/update-password', 'updatePassword')->name('auth.profile.update-password');
        Route::post('/logout', 'logout')->name('auth.logout');
    });
});
