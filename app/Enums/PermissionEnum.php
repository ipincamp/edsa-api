<?php

namespace App\Enums;

enum PermissionEnum: string
{
/** Personals */
    case CHANGE_PASSWORD_SELF = 'Change Password';
    case UPDATE_PROFILE_SELF = 'Update Profile';
/** Students */
    case CHANGE_PASSWORD_STUDENT = 'Change Password Student';
/** Teachers */
    case CHANGE_PASSWORD_TEACHER = 'Change Password Teacher';
}
