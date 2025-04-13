<?php

return [

    /*
    |--------------------------------------------------------------------------
    | Users
    |--------------------------------------------------------------------------
    |
    | This value determines the "users" your application is currently
    | running in. This may determine how you prefer to configure various
    | services the application utilizes. Set this in your ".env" file.
    |
    */

    'user' => [
        'admin' => [
            'name' => env('USER_ADMIN_NAME'),
            'email' => env('USER_ADMIN_EMAIL'),
            'password' => env('USER_ADMIN_PASSWORD'),
        ],
        'teacher' => [
            'name' => env('USER_TEACHER_NAME'),
            'email' => env('USER_TEACHER_EMAIL'),
            'password' => env('USER_TEACHER_PASSWORD'),
        ],
        'student' => [
            'name' => env('USER_STUDENT_NAME'),
            'username' => env('USER_STUDENT_USERNAME'),
            'password' => env('USER_STUDENT_PASSWORD'),
        ],
    ],

];
