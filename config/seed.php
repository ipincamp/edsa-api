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
            'name' => env('USER_ADMIN_NAME', null),
            'email' => env('USER_ADMIN_EMAIL', null),
            'username' => env('USER_ADMIN_USERNAME', null),
            'password' => env('USER_ADMIN_PASSWORD', null),
        ],
        'teacher' => [
            'name' => env('USER_TEACHER_NAME', null),
            'email' => env('USER_TEACHER_EMAIL', null),
            'username' => env('USER_TEACHER_USERNAME', null),
            'password' => env('USER_TEACHER_PASSWORD', null),
        ],
        'student' => [
            'name' => env('USER_STUDENT_NAME', null),
            'username' => env('USER_STUDENT_USERNAME', null),
            'password' => env('USER_STUDENT_PASSWORD', null),
        ],
    ],

];
