<?php

namespace App\Filament\Resources\CourseResource\Pages;

use App\Filament\Resources\CourseResource;
use Filament\Actions;
use Filament\Resources\Pages\ManageRecords;

class ManageCourses extends ManageRecords
{
    protected static string $resource = CourseResource::class;

    protected function getHeaderActions(): array
    {
        return [
            Actions\CreateAction::make()
                ->after(function ($record) {
                    $record->groups()->create([
                        'name' => 'Default',
                        'description' => 'Default group for ' . $record->name,
                        'course_id' => $record->id,
                    ]);

                    activity('course')
                        ->performedOn($record)
                        ->event('created')
                        ->log('Course created');
                }),
        ];
    }
}
