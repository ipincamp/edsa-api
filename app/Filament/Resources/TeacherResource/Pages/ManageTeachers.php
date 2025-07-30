<?php

namespace App\Filament\Resources\TeacherResource\Pages;

use App\Enums\RolesEnum;
use App\Filament\Resources\TeacherResource;
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
                    $record->assignRole(RolesEnum::T->value);
                }),
        ];
    }
}
