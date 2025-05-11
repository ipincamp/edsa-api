<?php

namespace App\Filament\Resources\TeacherResource\Pages;

use App\Enums\RoleEnum;
use App\Filament\Resources\TeacherResource;
use App\Models\Role;
use Filament\Actions;
use Filament\Resources\Pages\ManageRecords;

class ManageTeachers extends ManageRecords
{
    protected static string $resource = TeacherResource::class;

    protected function getHeaderActions(): array
    {
        return [
            Actions\CreateAction::make()
                ->after(function ($record) {
                    $record->roles()->attach(
                        Role::firstWhere('name', RoleEnum::TEACHER->value)->id,
                    );

                    activity('teacher')
                        ->performedOn($record)
                        ->event('created')
                        ->log('Created teacher');
                }),
        ];
    }
}
