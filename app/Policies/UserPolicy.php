<?php

namespace App\Policies;

use Illuminate\Auth\Access\Response;

use App\Models\User;

class UserPolicy
{
    /** Teacher */
    /**
     * Determine whether the user can view any models.
     */
    public function viewAnyTeacher(User $user): bool
    {
        return $user->checkPermissionTo('View any User');
    }

    /**
     * Determine whether the user can view the model.
     */
    public function viewTeacher(User $user, User $model): bool
    {
        return $user->checkPermissionTo('View User');
    }

    /**
     * Determine whether the user can create models.
     */
    public function createTeacher(User $user): bool
    {
        return $user->checkPermissionTo('Create User');
    }

    /**
     * Determine whether the user can update the model.
     */
    public function updateTeacher(User $user, User $model): bool
    {
        return $user->checkPermissionTo('Update User');
    }

    /**
     * Determine whether the user can delete the model.
     */
    public function deleteTeacher(User $user, User $model): bool
    {
        return $user->checkPermissionTo('Delete User');
    }

    /**
     * Determine whether the user can delete any models.
     */
    public function deleteAnyTeacher(User $user): bool
    {
        return $user->checkPermissionTo('Delete any User');
    }

    /**
     * Determine whether the user can restore the model.
     */
    public function restoreTeacher(User $user, User $model): bool
    {
        return $user->checkPermissionTo('Restore User');
    }

    /**
     * Determine whether the user can restore any models.
     */
    public function restoreAnyTeacher(User $user): bool
    {
        return $user->checkPermissionTo('Restore any User');
    }

    /**
     * Determine whether the user can replicate the model.
     */
    public function replicateTeacher(User $user, User $model): bool
    {
        return $user->checkPermissionTo('Replicate User');
    }

    /**
     * Determine whether the user can reorder the models.
     */
    public function reorderTeacher(User $user): bool
    {
        return $user->checkPermissionTo('Reorder User');
    }

    /**
     * Determine whether the user can permanently delete the model.
     */
    public function forceDeleteTeacher(User $user, User $model): bool
    {
        return $user->checkPermissionTo('Permanently delete User');
    }

    /**
     * Determine whether the user can permanently delete any models.
     */
    public function forceDeleteAnyTeacher(User $user): bool
    {
        return $user->checkPermissionTo('Permanently delete any User');
    }

    /** Student */
    /**
     * Determine whether the user can view any models.
     */
    public function viewAnyStudent(User $user): bool
    {
        return $user->checkPermissionTo('View any User');
    }

    /**
     * Determine whether the user can view the model.
     */
    public function viewStudent(User $user, User $model): bool
    {
        return $user->checkPermissionTo('View User');
    }

    /**
     * Determine whether the user can create models.
     */
    public function createStudent(User $user): bool
    {
        return $user->checkPermissionTo('Create User');
    }

    /**
     * Determine whether the user can update the model.
     */
    public function updateStudent(User $user, User $model): bool
    {
        return $user->checkPermissionTo('Update User');
    }

    /**
     * Determine whether the user can delete the model.
     */
    public function deleteStudent(User $user, User $model): bool
    {
        return $user->checkPermissionTo('Delete User');
    }

    /**
     * Determine whether the user can delete any models.
     */
    public function deleteAnyStudent(User $user): bool
    {
        return $user->checkPermissionTo('Delete any User');
    }

    /**
     * Determine whether the user can restore the model.
     */
    public function restoreStudent(User $user, User $model): bool
    {
        return $user->checkPermissionTo('Restore User');
    }

    /**
     * Determine whether the user can restore any models.
     */
    public function restoreAnyStudent(User $user): bool
    {
        return $user->checkPermissionTo('Restore any User');
    }

    /**
     * Determine whether the user can replicate the model.
     */
    public function replicateStudent(User $user, User $model): bool
    {
        return $user->checkPermissionTo('Replicate User');
    }

    /**
     * Determine whether the user can reorder the models.
     */
    public function reorderStudent(User $user): bool
    {
        return $user->checkPermissionTo('Reorder User');
    }

    /**
     * Determine whether the user can permanently delete the model.
     */
    public function forceDeleteStudent(User $user, User $model): bool
    {
        return $user->checkPermissionTo('Permanently delete User');
    }

    /**
     * Determine whether the user can permanently delete any models.
     */
    public function forceDeleteAnyStudent(User $user): bool
    {
        return $user->checkPermissionTo('Permanently delete any User');
    }
}
