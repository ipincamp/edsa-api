<?php

namespace App\Policies;

use App\Enums\RolesEnum;
use App\Models\Group;
use App\Models\User;

class GroupPolicy
{
    // Admin bisa melakukan segalanya
    public function before(User $user, string $ability): bool|null
    {
        return $user->hasRole(RolesEnum::A->value) ? true : null;
    }

    // Guru hanya bisa melihat grup yang dia ajar
    public function view(User $user, Group $group): bool
    {
        return $user->id === $group->teacher_id;
    }

    // Guru hanya bisa mengupdate grup yang dia ajar
    public function update(User $user, Group $group): bool
    {
        return $user->id === $group->teacher_id;
    }

    // Guru hanya bisa menghapus grup yang dia ajar (jika kosong)
    public function delete(User $user, Group $group): bool
    {
        return $user->id === $group->teacher_id;
    }

    // Guru hanya bisa menambah/hapus siswa di grup yang dia ajar
    public function manageStudents(User $user, Group $group): bool
    {
        return $user->id === $group->teacher_id;
    }
}
