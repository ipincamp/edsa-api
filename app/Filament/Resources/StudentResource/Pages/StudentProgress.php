<?php

namespace App\Filament\Resources\StudentResource\Pages;

use App\Filament\Resources\StudentResource;
use App\Filament\Widgets\SummaryWidget;
use Filament\Resources\Pages\Page;

class StudentProgress extends Page
{
    protected static string $resource = StudentResource::class;

    protected static string $view = 'filament.resources.student-resource.pages.student-progress';

    protected function getFooterWidgets(): array
    {
        return [
            SummaryWidget::class
        ];
    }
}
