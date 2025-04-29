<?php

namespace Database\Seeders;

use App\Enums\RoleEnum;
use App\Models\Course;
use App\Models\Group;
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

        // courses
        $course = Course::create([
            'name' => 'English',
            'description' => 'English Course',
        ]);

        // groups
        $group = Group::create([
            'name' => 'A',
            'description' => 'Class A',
            'course_id' => $course->id,
        ]);

        // users
        if (config('seed.user.admin.email') !== null) {
            $admin = User::create([
                'name' => config('seed.user.admin.name'),
                'email' => config('seed.user.admin.email'),
                'email_verified_at' => now(),
                'username' => config('seed.user.admin.username'),
                'password' => bcrypt(config('seed.user.admin.password')),
            ]);
            $admin->roles()->attach($adminRole->id);
        }

        if (config('seed.user.teacher.email') !== null) {
            $teacher = User::create([
                'name' => config('seed.user.teacher.name'),
                'email' => config('seed.user.teacher.email'),
                'username' => config('seed.user.teacher.username'),
                'password' => bcrypt(config('seed.user.teacher.password')),
            ]);
            $teacher->roles()->attach($teacherRole->id);
            $group->participants()->attach($teacher->id);
        }

        if (config('seed.user.student.username') !== null) {
            $student = User::create([
                'name' => config('seed.user.student.name'),
                'username' => config('seed.user.student.username'),
                'password' => bcrypt(config('seed.user.student.password')),
            ]);
            $student->roles()->attach($studentRole->id);
            $group->participants()->attach($student->id);
        }

        $studentGuest = User::create([
            'name' => config('seed.user.student.name') . ' Guest',
            'username' => config('seed.user.student.username') . 'guest',
            'password' => config('seed.user.student.password'),
        ]);
        $studentGuest->roles()->attach($studentRole->id);
    }
}
