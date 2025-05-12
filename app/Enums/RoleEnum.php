<?php

namespace App\Enums;

enum RoleEnum: string
{
    case ADMIN = 'admin';
    case TEACHER = 'teacher';
    case STUDENT = 'student';

    /**
     * Get the label of the enum case.
     *
     * @return string
     */
    public function label(): string
    {
        return match ($this) {
            self::ADMIN => 'Admin',
            self::TEACHER => 'Teacher',
            self::STUDENT => 'Student',
        };
    }

    /**
     * Get all roles for the enum.
     *
     * @return string
     */
    public static function getAllRoles(): array
    {
        return array_map(fn($enum) => $enum->value, self::cases());
    }
}
