<?php

use App\Http\Controllers\Api\Auth\AuthController;
use App\Http\Controllers\Api\Book\BookController;
use App\Http\Controllers\Api\Course\CourseController;
use App\Http\Controllers\Api\Course\EnrollmentController;
use App\Http\Controllers\Api\Group\GroupController;
use App\Http\Controllers\Api\User\ProgressController;
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

    // Course
    Route::middleware('role:admin')->group(function () {
        Route::apiResource('courses', CourseController::class);
    });
    Route::middleware('role:student')->group(function () {
        Route::post('courses/enroll', [EnrollmentController::class, 'enroll'])->name('courses.enroll');
    });

    // Group
    Route::middleware('role:admin|teacher')->group(function () {
        Route::get('courses/{course}/groups', [GroupController::class, 'index'])->name('courses.groups.index');

        Route::apiResource('groups', GroupController::class)->except(['index']);

        Route::post('groups/{group}/assign-student', [GroupController::class, 'assignStudent'])->name('groups.assign-student');
        Route::post('groups/{group}/remove-student', [GroupController::class, 'removeStudent'])->name('groups.remove-student');
    });

    // User Progress
    Route::middleware(['role:student'])->prefix('progress')->group(function () {
        Route::post('/start-book', [ProgressController::class, 'startOrContinueBook'])
            ->name('progress.start-book');
        Route::patch('/update-page', [ProgressController::class, 'updateLastPage'])
            ->name('progress.update-last-page');
        Route::post('/submit-interaction', [ProgressController::class, 'submitInteraction'])
            ->name('progress.submit-interaction');
        Route::post('/submit-post-activity', [ProgressController::class, 'submitPostActivity'])
            ->name('progress.submit-post-activity');
        Route::post('/complete-book', [ProgressController::class, 'completeBook'])
            ->name('progress.complete-book');
    });
    Route::get('/students/{user}/progress', [ProgressController::class, 'viewStudentProgress'])
        ->middleware(['role:admin|teacher'])
        ->name('progress.view-student-progress');
});
