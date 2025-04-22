<?php

namespace Database\Seeders;

use App\Enums\RoleEnum;
use App\Models\User;
use Illuminate\Database\Seeder;
use Spatie\Permission\Models\Role;

class UserSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        // roles
        $adminRole = Role::create([
            'name' => RoleEnum::ADMIN->value,
            'guard_name' => 'web',
        ]);
        $teacherRole = Role::create([
            'name' => RoleEnum::TEACHER->value,
            'guard_name' => 'web',
        ]);
        $studentRole = Role::create([
            'name' => RoleEnum::STUDENT->value,
            'guard_name' => 'web',
        ]);

        // users
        if (config('seed.user.admin.email') !== null) {
            $admin = User::create([
                'name' => config('seed.user.admin.name'),
                'email' => config('seed.user.admin.email'),
                'email_verified_at' => now(),
                'password' => bcrypt(config('seed.user.admin.password')),
            ]);
            $admin->assignRole($adminRole);
        }

        if (config('seed.user.teacher.email') !== null) {
            $teacher = User::create([
                'name' => config('seed.user.teacher.name'),
                'email' => config('seed.user.teacher.email'),
                'password' => bcrypt(config('seed.user.teacher.password')),
            ]);
            $teacher->assignRole($teacherRole);
        }

        if (config('seed.user.student.username') !== null) {
            $student = User::create([
                'name' => config('seed.user.student.name'),
                'username' => config('seed.user.student.username'),
                'password' => bcrypt(config('seed.user.student.password')),
            ]);
            $student->assignRole($studentRole);
        }
    }
}
