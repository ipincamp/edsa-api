<?php

namespace Database\Seeders;

use App\Enums\RolesEnum;
use App\Models\User;
use Illuminate\Database\Seeder;

class UserSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        $admin = User::updateOrCreate(
            ['email' => config('seed.users.admin')],
            [
                'name' => 'Admin EDSA',
                'email_verified_at' => now(),
                'password' => bcrypt(config('seed.users.password')),
            ]
        );
        $admin->assignRole(RolesEnum::A);

        $student = User::updateOrCreate(
            ['email' => config('seed.users.student')],
            [
                'name' => 'Student EDSA',
                'password' => bcrypt(config('seed.users.password')),
            ]
        );
        $student->assignRole(RolesEnum::S);
    }
}
