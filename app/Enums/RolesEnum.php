<?php

namespace App\Enums;

enum RolesEnum: string
{
    // admin, teacher, student
    case A = 'admin';
    case T = 'teacher';
    case S = 'student';
}
