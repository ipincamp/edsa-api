<?php

namespace Database\Seeders;

use App\Enums\RolesEnum;
use App\Models\Course;
use App\Models\Group;
use App\Models\User;
use Illuminate\Database\Seeder;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Str;

class CourseGroupSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        Group::query()->delete();
        Course::query()->delete();

        // Buat Pengguna (Guru, Siswa)
        $teacher1 = User::create([
            'name' => 'Budi Doremi',
            'email' => 'first.' . config('seed.users.teacher'),
            'password' => Hash::make(config('seed.users.password')),
        ]);
        $teacher1->assignRole(RolesEnum::T->value);

        $teacher2 = User::create([
            'name' => 'Siti Aminah',
            'email' => 'second.' . config('seed.users.teacher'),
            'password' => Hash::make(config('seed.users.password')),
        ]);
        $teacher2->assignRole(RolesEnum::T->value);

        // Buat 15 siswa menggunakan factory
        $students = User::factory(15)->unverified()->create();
        $students->each(fn($student) => $student->assignRole(RolesEnum::S->value));

        // Buat Course
        $course1 = Course::create([
            'name' => 'Beginner English - Level 1',
            'description' => 'Kelas untuk pemula yang berfokus pada alfabet dan angka.',
            'enrollment_code' => Str::random(6),
        ]);

        $course2 = Course::create([
            'name' => 'Storytelling - Level 2',
            'description' => 'Kelas untuk melatih pemahaman cerita dan kosakata.',
            'enrollment_code' => Str::random(6),
        ]);

        // Buat Grup dan tugaskan guru
        $groupA = Group::create([
            'name' => 'Grup A - Pagi',
            'course_id' => $course1->id,
            'teacher_id' => $teacher1->id,
        ]);

        $groupB = Group::create([
            'name' => 'Grup B - Siang',
            'course_id' => $course1->id,
            'teacher_id' => $teacher2->id,
        ]);

        $groupC = Group::create([
            'name' => 'Grup C - Lanjutan',
            'course_id' => $course2->id,
            'teacher_id' => $teacher1->id,
        ]);

        // Masukkan siswa ke dalam grup
        // Ambil 8 siswa pertama untuk Grup A
        $groupA->students()->attach($students->take(8)->pluck('id'));

        // Ambil 7 siswa sisanya untuk Grup B
        $groupB->students()->attach($students->skip(8)->take(7)->pluck('id'));

        // Masukkan 5 siswa secara acak ke Grup C
        $groupC->students()->attach($students->random(5)->pluck('id'));
    }
}
