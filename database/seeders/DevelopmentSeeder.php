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
        if (config('seeder.user.email.admin') !== null) {
            $admin = User::create([
                'name' => 'Admin EDSA',
                'email' => config('seeder.user.email.admin'),
                'email_verified_at' => now(),
                'username' => 'admin.edsa',
                'password' => bcrypt(config('seeder.user.password.admin')),
            ]);
            $admin->roles()->attach(
                Role::firstWhere('name', RoleEnum::ADMIN->value)->id,
            );
        }

        if (config('seeder.user.email.teacher') !== null) {
            $teacher = User::create([
                'name' => 'Teacher EDSA',
                'email' => config('seeder.user.email.teacher'),
                'username' => 'teacher.edsa',
                'password' => bcrypt(config('seeder.user.password.teacher')),
            ]);
            $teacher->roles()->attach(
                Role::firstWhere('name', RoleEnum::TEACHER->value)->id,
            );
            $group->participants()->attach($teacher->id);
        }

        if (config('seeder.user.email.student') !== null) {
            $student = User::create([
                'name' => 'Student EDSA',
                'username' => 'student.edsa',
                'password' => bcrypt(config('seeder.user.password.student')),
            ]);
            $student->roles()->attach(
                Role::firstWhere('name', RoleEnum::STUDENT->value)->id,
            );
            $group->participants()->attach($student->id);
        }

        $studentGuest = User::create([
            'name' => 'Guest EDSA',
            'username' => 'guest.edsa',
            'password' => bcrypt(config('seeder.user.password.student')),
        ]);
        $studentGuest->roles()->attach(
            Role::firstWhere('name', RoleEnum::STUDENT->value)->id,
        );
    }
}
