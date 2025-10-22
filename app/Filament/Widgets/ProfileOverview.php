<?php

namespace App\Filament\Widgets;

use App\Enums\PermissionEnum;
use App\Traits\Api\AuthorizeTrait;
use Filament\Widgets\StatsOverviewWidget as BaseWidget;
use Filament\Widgets\StatsOverviewWidget\Stat;

class ProfileOverview extends BaseWidget
{
    use AuthorizeTrait;

    protected static ?string $pollingInterval = '10s';
    protected static bool $isLazy = false;

    protected int | string | array $columnSpan = [
        'md' => 2,
        'xl' => 3,
    ];

    protected function getStats(): array
    {
        $user = auth()->user();

        return [
            Stat::make('Welcome', $user->name)
                ->description("You're logged in as " . ucfirst($user->getRoleNames()->first()))
                ->icon('heroicon-o-user-circle')
                ->color('secondary'),
            Stat::make('Joined since', $user->created_at->diffForHumans())
                ->description($user->created_at->format('F j, Y \| H:i:s'))
                ->icon('heroicon-o-calendar')
                ->color('success'),
            Stat::make('Last update', $user->updated_at->diffForHumans())
                ->description($user->updated_at->format('F j, Y \| H:i:s'))
                ->icon('heroicon-o-clock')
                ->color('warning'),
        ];
    }

    public static function canView(): bool
    {
        return static::grant(PermissionEnum::VIEW_PROFILE_WIDGET->value);
    }
}
