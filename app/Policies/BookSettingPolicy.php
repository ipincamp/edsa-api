<?php

namespace App\Policies;

use Illuminate\Auth\Access\Response;
use App\Models\BookSetting;
use App\Models\User;

class BookSettingPolicy
{
    /**
     * Determine whether the user can view any models.
     */
    public function viewAny(User $user): bool
    {
        return $user->checkPermissionTo('View any BookSetting');
    }

    /**
     * Determine whether the user can view the model.
     */
    public function view(User $user, BookSetting $booksetting): bool
    {
        return $user->checkPermissionTo('View BookSetting');
    }

    /**
     * Determine whether the user can create models.
     */
    public function create(User $user): bool
    {
        return $user->checkPermissionTo('Create BookSetting');
    }

    /**
     * Determine whether the user can update the model.
     */
    public function update(User $user, BookSetting $booksetting): bool
    {
        return $user->checkPermissionTo('Update BookSetting');
    }

    /**
     * Determine whether the user can delete the model.
     */
    public function delete(User $user, BookSetting $booksetting): bool
    {
        return $user->checkPermissionTo('Delete BookSetting');
    }

    /**
     * Determine whether the user can delete any models.
     */
    public function deleteAny(User $user): bool
    {
        return $user->checkPermissionTo('Delete any BookSetting');
    }

    /**
     * Determine whether the user can restore the model.
     */
    public function restore(User $user, BookSetting $booksetting): bool
    {
        return $user->checkPermissionTo('Restore BookSetting');
    }

    /**
     * Determine whether the user can restore any models.
     */
    public function restoreAny(User $user): bool
    {
        return $user->checkPermissionTo('Restore any BookSetting');
    }

    /**
     * Determine whether the user can replicate the model.
     */
    public function replicate(User $user, BookSetting $booksetting): bool
    {
        return $user->checkPermissionTo('Replicate BookSetting');
    }

    /**
     * Determine whether the user can reorder the models.
     */
    public function reorder(User $user): bool
    {
        return $user->checkPermissionTo('Reorder BookSetting');
    }

    /**
     * Determine whether the user can permanently delete the model.
     */
    public function forceDelete(User $user, BookSetting $booksetting): bool
    {
        return $user->checkPermissionTo('Permanently delete BookSetting');
    }

    /**
     * Determine whether the user can permanently delete any models.
     */
    public function forceDeleteAny(User $user): bool
    {
        return $user->checkPermissionTo('Permanently delete any BookSetting');
    }
}
