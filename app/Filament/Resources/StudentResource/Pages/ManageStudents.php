<?php

namespace App\Filament\Resources\StudentResource\Pages;

use App\Enums\RoleEnum;
use App\Filament\Resources\StudentResource;
use App\Models\Role;
use Filament\Actions;
use Filament\Resources\Pages\ManageRecords;

class ManageStudents extends ManageRecords
{
    protected static string $resource = StudentResource::class;

    protected function getHeaderActions(): array
    {
        return [
            Actions\CreateAction::make()
                ->after(function ($record) {
                    $record->roles()->attach(
                        Role::firstWhere('name', RoleEnum::STUDENT->value)->id,
                    );

                    activity('student')
                        ->performedOn($record)
                        ->event('created')
                        ->log('Created student');
                }),
        ];
    }
}
