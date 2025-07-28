<?php

namespace Database\Seeders;

use App\Enums\RolesEnum;
use Illuminate\Database\Seeder;
use Spatie\Permission\Models\Role;

class RoleSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        Role::updateOrCreate(
            ['name' => RolesEnum::A->value, 'guard_name' => 'web'],
            []
        );
        Role::updateOrCreate(
            ['name' => RolesEnum::T->value, 'guard_name' => 'web'],
            []
        );
        Role::updateOrCreate(
            ['name' => RolesEnum::S->value, 'guard_name' => 'web'],
            []
        );
    }
}
