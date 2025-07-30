<?php

namespace App\Filament\Resources\StudentResource\Pages;

use App\Enums\RolesEnum;
use App\Filament\Resources\StudentResource;
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
                    $record->assignRole(RolesEnum::S->value);
                }),
        ];
    }
}
