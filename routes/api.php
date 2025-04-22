<?php

use App\Http\Controllers\Api\V1\Auth\AuthController;
use Illuminate\Support\Facades\Route;

Route::prefix('v1/auth')->controller(AuthController::class)->group(function () {
    Route::post('login', 'login')->name('auth.signIn');
    Route::middleware('auth:sanctum')->group(function () {
        Route::get('profile', 'profile')->name('auth.profile');
        Route::post('update-password', 'updatePassword')->name('auth.updatePassword');
        Route::post('update-profile', 'updateProfile')->name('auth.updateProfile');
        Route::post('logout', 'logout')->name('auth.signOut');
    });
});
