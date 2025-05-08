<?php

namespace Database\Seeders;

use App\Enums\RoleEnum;
use App\Models\Course;
use App\Models\Group;
use App\Models\Role;
use App\Models\User;
use Illuminate\Database\Seeder;

class DevelopmentSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
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
            $admin->roles()->attach(
                Role::firstWhere('name', RoleEnum::ADMIN->value)->id,
            );
        }

        if (config('seed.user.teacher.email') !== null) {
            $teacher = User::create([
                'name' => config('seed.user.teacher.name'),
                'email' => config('seed.user.teacher.email'),
                'username' => config('seed.user.teacher.username'),
                'password' => bcrypt(config('seed.user.teacher.password')),
            ]);
            $teacher->roles()->attach(
                Role::firstWhere('name', RoleEnum::TEACHER->value)->id,
            );
            $group->participants()->attach($teacher->id);
        }

        if (config('seed.user.student.username') !== null) {
            $student = User::create([
                'name' => config('seed.user.student.name'),
                'username' => config('seed.user.student.username'),
                'password' => bcrypt(config('seed.user.student.password')),
            ]);
            $student->roles()->attach(
                Role::firstWhere('name', RoleEnum::STUDENT->value)->id,
            );
            $group->participants()->attach($student->id);
        }

        $studentGuest = User::create([
            'name' => config('seed.user.student.name') . ' Guest',
            'username' => config('seed.user.student.username') . 'guest',
            'password' => config('seed.user.student.password'),
        ]);
        $studentGuest->roles()->attach(
            Role::firstWhere('name', RoleEnum::STUDENT->value)->id,
        );
    }
}
