<?php

namespace App\Filament\Widgets;

use App\Enums\PermissionEnum;
use App\Traits\Api\AuthorizeTrait;
use Filament\Tables;
use Filament\Tables\Table;
use Filament\Widgets\TableWidget as BaseWidget;

class SummaryWidget extends BaseWidget
{
    use AuthorizeTrait;

    protected static ?string $pollingInterval = '60s';
    protected static bool $isLazy = false;

    protected int | string | array $columnSpan = [
        'md' => 2,
        'xl' => 3,
    ];

    public function table(Table $table): Table
    {
        return $table
            ->query(
                \App\Models\User::query()
                    ->whereHas('roles', fn($query) => $query->where('name', 'student'))
                    ->with(['progress' => fn($query) => $query->orderBy('book_id')])
            )
            ->columns(
                [
                    Tables\Columns\TextColumn::make('name')
                        ->label('Student Name')
                        ->limit(20),
                    Tables\Columns\TextColumn::make('progress.0.score_correct')
                        ->label('Book 1')
                        ->default(0),
                    Tables\Columns\TextColumn::make('progress.1.score_correct')
                        ->label('Book 2')
                        ->default(0),
                    Tables\Columns\TextColumn::make('progress.2.score_correct')
                        ->label('Book 3')
                        ->default(0),
                    Tables\Columns\TextColumn::make('progress.3.score_correct')
                        ->label('Book 4')
                        ->default(0),
                    Tables\Columns\TextColumn::make('progress.4.score_correct')
                        ->label('Book 5')
                        ->default(0),
                    Tables\Columns\TextColumn::make('progress.5.score_correct')
                        ->label('Book 6')
                        ->default(0),
                    Tables\Columns\TextColumn::make('progress.6.score_correct')
                        ->label('Book 7')
                        ->default(0),
                    Tables\Columns\TextColumn::make('progress.7.score_correct')
                        ->label('Book 8')
                        ->default(0),
                    Tables\Columns\TextColumn::make('average_score')
                        ->label('Avg')
                        ->getStateUsing(fn($record) => number_format(collect($record->progress)->avg('score_correct') ?? 0, 2))
                ]
            );
    }

    public static function canView(): bool
    {
        return static::grant(PermissionEnum::VIEW_SUMMARY_WIDGET->value);
    }
}
