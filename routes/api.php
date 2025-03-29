<?php

use App\Http\Controllers\Api\V1\Auth\AuthController;
use Illuminate\Support\Facades\Route;

Route::prefix('v1')->group(function () {
    // Auth
    Route::prefix('auth')->controller(AuthController::class)->group(function () {
        Route::post('/login', 'signIn')->name('auth.signIn');
        Route::post('/register', 'signUp')->name('auth.signUp');
        Route::post('/logout', 'signOut')->name('auth.signOut')->middleware('auth:sanctum');
    });
});
