<?php

namespace App\Enums;

enum PermissionEnum: string
{
    /** Personals */
    case CHANGE_PASSWORD_SELF = 'Change Password';
    case UPDATE_PROFILE_SELF = 'Update Profile';
    case CHANGE_PASSWORD_OTHER = 'Change Password Other';
    case UPDATE_PROFILE_OTHER = 'Update Profile Other';
}
