<?php

namespace App\Enums;

enum PermissionEnum: string
{
    case CHANGE_PASSWORD_SELF = 'Change Password Self';
    case CHANGE_PASSWORD_STUDENT = 'Change Password Student';
    case CHANGE_PASSWORD_TEACHER = 'Change Password Teacher';
    case UPDATE_PROFILE_SELF = 'Update Profile Self';
    case VIEW_PROFILE_WIDGET = 'View Profile Widget';
    case VIEW_SUMMARY_WIDGET = 'View Summary Widget';

    /**
     * Get all permissions for the enum.
     *
     * @return string
     */
    public static function getAllPermissions(): array
    {
        return array_map(fn($enum) => $enum->value, self::cases());
    }
}
