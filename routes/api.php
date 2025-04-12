<?php

use App\Http\Controllers\V1\Auth\AuthController;
use Illuminate\Support\Facades\Route;

Route::prefix('v1')->group(function () {
    // Auth
    Route::prefix('auth')->controller(AuthController::class)->group(function () {
        Route::post('/login', 'login')->name('auth.signIn');
        // change password
        Route::middleware('auth:sanctum')->group(function () {
            Route::post('/change-password', 'changePassword')->name('auth.changePassword');
            Route::post('/update-profile', 'updateProfile')->name('auth.updateProfile');
            Route::post('/logout', 'logout')->name('auth.signOut');
        });
    });
});
