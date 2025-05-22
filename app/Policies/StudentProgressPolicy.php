<?php

namespace App\Policies;

use Illuminate\Auth\Access\Response;
use App\Models\StudentProgress;
use App\Models\User;

class StudentProgressPolicy
{
    /**
     * Determine whether the user can view any models.
     */
    public function viewAny(User $user): bool
    {
        return $user->checkPermissionTo('View any StudentProgress');
    }

    /**
     * Determine whether the user can view the model.
     */
    public function view(User $user, StudentProgress $studentprogress): bool
    {
        return $user->checkPermissionTo('View StudentProgress');
    }

    /**
     * Determine whether the user can create models.
     */
    public function create(User $user): bool
    {
        return $user->checkPermissionTo('Create StudentProgress');
    }

    /**
     * Determine whether the user can update the model.
     */
    public function update(User $user, StudentProgress $studentprogress): bool
    {
        return $user->checkPermissionTo('Update StudentProgress');
    }

    /**
     * Determine whether the user can delete the model.
     */
    public function delete(User $user, StudentProgress $studentprogress): bool
    {
        return $user->checkPermissionTo('Delete StudentProgress');
    }

    /**
     * Determine whether the user can delete any models.
     */
    public function deleteAny(User $user): bool
    {
        return $user->checkPermissionTo('Delete any StudentProgress');
    }

    /**
     * Determine whether the user can restore the model.
     */
    public function restore(User $user, StudentProgress $studentprogress): bool
    {
        return $user->checkPermissionTo('Restore StudentProgress');
    }

    /**
     * Determine whether the user can restore any models.
     */
    public function restoreAny(User $user): bool
    {
        return $user->checkPermissionTo('Restore any StudentProgress');
    }

    /**
     * Determine whether the user can replicate the model.
     */
    public function replicate(User $user, StudentProgress $studentprogress): bool
    {
        return $user->checkPermissionTo('Replicate StudentProgress');
    }

    /**
     * Determine whether the user can reorder the models.
     */
    public function reorder(User $user): bool
    {
        return $user->checkPermissionTo('Reorder StudentProgress');
    }

    /**
     * Determine whether the user can permanently delete the model.
     */
    public function forceDelete(User $user, StudentProgress $studentprogress): bool
    {
        return $user->checkPermissionTo('Permanently delete StudentProgress');
    }

    /**
     * Determine whether the user can permanently delete any models.
     */
    public function forceDeleteAny(User $user): bool
    {
        return $user->checkPermissionTo('Permanently delete any StudentProgress');
    }
}
