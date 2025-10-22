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
            ->columns([
                Tables\Columns\TextColumn::make('name')
                    ->label('Student Name')
                    ->limit(20)
                    ->grow(false),
                ...collect(range(1, 8))->map(
                    fn($book) =>
                    Tables\Columns\TextColumn::make("progress." . ($book - 1) . ".score_correct")
                        ->label("Book $book")
                        ->default(0),
                )->toArray(),
                Tables\Columns\TextColumn::make('average_score')
                    ->label('Avg')
                    ->getStateUsing(fn($record) => number_format(
                        num: collect($record->progress)->sum('score_correct') / 8 ?? 0,
                        decimals: 2,
                    )),
            ])
            ->striped();
    }

    public static function canView(): bool
    {
        return static::grant(PermissionEnum::VIEW_SUMMARY_WIDGET->value);
    }
}
