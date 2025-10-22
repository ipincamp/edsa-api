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
        'email' => [
            'admin' => env('USER_ADMIN_EMAIL', null),
            'teacher' => env('USER_TEACHER_EMAIL', null),
            'student' => env('USER_STUDENT_EMAIL', null),
        ],
        'password' => [
            'admin' => env('USER_ADMIN_PASSWORD', null),
            'teacher' => env('USER_TEACHER_PASSWORD', null),
            'student' => env('USER_STUDENT_PASSWORD', null),
        ],
    ],

];
