<?php

namespace App\Filament\Resources;

use App\Enums\RoleEnum;
use App\Filament\Resources\RoleResource\Pages;
use App\Filament\Resources\RoleResource\RelationManagers;
use App\Models\Role;
use Filament\Forms;
use Filament\Forms\Form;
use Filament\Resources\Resource;
use Filament\Tables;
use Filament\Tables\Table;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\SoftDeletingScope;
use Rmsramos\Activitylog\RelationManagers\ActivitylogRelationManager;
use Spatie\Permission\Models\Permission;

class RoleResource extends Resource
{
    protected static ?string $model = Role::class;

    protected static ?string $navigationIcon = 'heroicon-o-shield-check';
    protected static ?string $navigationGroup = 'Settings';
    protected static ?string $navigationLabel = 'Permissions';
    protected static ?string $label = 'Permission';
    protected static ?string $pluralLabel = 'Manage Permissions';
    protected static ?string $slug = 'permissions';

    public static function form(Form $form): Form
    {
        return $form
            ->schema([
                Forms\Components\CheckboxList::make('permissions')
                    ->relationship('permissions', 'name')
                    ->columns([
                        'default' => 1,
                        'sm' => 2,
                        'md' => 3,
                    ])
                    ->columnSpanFull()
                    ->label('Permissions')
                    ->bulkToggleable()
                    ->searchable()
                    ->afterStateUpdated(function (callable $set, $state, $record) {
                        $set('permissions', $state);

                        if ($record) {
                            activity('permissions')
                                ->performedOn($record)
                                ->event('updated')
                                ->withProperties([
                                    'attributes' => [
                                        'permissions' => Permission::whereIn('id', $state)->pluck('name')->toArray(),
                                    ],
                                    'old' => [
                                        'permissions' => $record->permissions->pluck('name')->toArray(),
                                    ],
                                ])
                                ->log('Updated permissions');
                        }
                    }),
            ]);
    }

    public static function getEloquentQuery(): Builder
    {
        return parent::getEloquentQuery()
            ->withoutGlobalScopes([SoftDeletingScope::class])
            ->where('name', '!=', RoleEnum::ADMIN->value);
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                Tables\Columns\TextColumn::make('name')
                    ->label('Role Name')
                    ->getStateUsing(fn($record) => ucwords($record->name))
                    ->copyable()
                    ->copyMessage('Copied!')
                    ->searchable(),
                Tables\Columns\TextColumn::make('permission_count')
                    ->label('Permissions')
                    ->getStateUsing(fn($record) => $record->permissions->count()),
            ])
            ->actions([
                Tables\Actions\EditAction::make()
                    ->label('Permissions')
                    ->icon('heroicon-o-pencil'),
            ]);
    }

    public static function getRelations(): array
    {
        return [
            ActivitylogRelationManager::class,
        ];
    }

    public static function getPages(): array
    {
        return [
            'index' => Pages\ManageRoles::route('/'),
        ];
    }
}
