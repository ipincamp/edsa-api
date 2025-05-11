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
                // TODO: Add relation to summaries
                \App\Models\User::query()->withCount('groups'),
            )
            ->columns(
                collect(range(1, 8))->map(
                    fn($i) => Tables\Columns\TextColumn::make("b_$i")
                )->prepend(
                    Tables\Columns\TextColumn::make('name')
                        ->label('Student Name')
                        ->limit(20)
                )->push(
                    Tables\Columns\TextColumn::make('average')
                        ->label('Average')
                        ->formatStateUsing(fn($state) => number_format($state, 2))
                )->toArray()
            );
    }

    public static function canView(): bool
    {
        return static::grant(PermissionEnum::VIEW_SUMMARY_WIDGET->value);
    }
}
