<?php

namespace App\Traits;

trait AuthorizeTrait
{
    /**
     * Check if the user is authorized to perform the given action.
     *
     * @param string $ability
     * @param array $arguments
     * @return bool
     */
    public static function grant($ability, $arguments = []): bool
    {
        if (auth()->user() && auth()->user()->can($ability, $arguments)) {
            return true;
        }

        return false;
    }
}
