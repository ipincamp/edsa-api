<?php

namespace App\Filament\Resources;

use App\Enums\PermissionEnum;
use App\Enums\RoleEnum;
use App\Filament\Resources\StudentResource\Pages;
use App\Filament\Resources\StudentResource\RelationManagers;
use App\Models\Group;
use App\Models\User;
use App\Traits\Api\AuthorizeTrait;
use Filament\Forms;
use Filament\Forms\Form;
use Filament\Resources\Resource;
use Filament\Tables;
use Filament\Tables\Table;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\SoftDeletingScope;

class StudentResource extends Resource
{
    use AuthorizeTrait;

    protected static ?string $model = User::class;

    protected static ?string $navigationGroup = 'Users';
    protected static ?string $navigationLabel = 'Students';
    protected static ?int $navigationSort = 2;
    protected static ?string $label = 'Student';
    protected static ?string $pluralLabel = 'Data Students';

    public static function form(Form $form): Form
    {
        return $form
            ->schema([
                // name
                Forms\Components\TextInput::make('name')
                    ->label('Name')
                    ->required()
                    ->maxLength(50),
                // username
                Forms\Components\TextInput::make('username')
                    ->label('Username')
                    ->required()
                    ->maxLength(50)
                    ->unique(ignoreRecord: true),
                // password
                Forms\Components\TextInput::make('password')
                    ->label('Password')
                    ->required()
                    ->password()
                    ->maxLength(32)
                    ->revealable()
                    ->columnSpanFull()
                    ->hiddenOn(['edit', 'view']),
                // group in course
                Forms\Components\Select::make('groups')
                    ->label('Groups')
                    ->relationship('groups', 'name')
                    ->options(
                        Group::with('course')
                            ->get()
                            ->groupBy('course.name')
                            ->mapWithKeys(function ($groups, $courseName) {
                                return [
                                    $courseName => $groups->pluck('name', 'id')->map(function ($name) use ($courseName) {
                                        return $name . ' (' . $courseName . ')';
                                    })->toArray(),
                                ];
                            })
                            ->toArray()
                    )
                    ->columnSpanFull(),
            ]);
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                // name
                Tables\Columns\TextColumn::make('name')
                    ->label('Name')
                    ->searchable()
                    ->limit(50),
                // username
                Tables\Columns\TextColumn::make('username')
                    ->label('Username')
                    ->searchable()
                    ->limit(50),
                // group
                Tables\Columns\TextColumn::make('groups')
                    ->label('Groups')
                    ->formatStateUsing(function ($record) {
                        return $record->groups->map(function ($group) {
                            return $group->name . ' (' . $group->course->name . ')';
                        })->join(', ');
                    })
                    ->searchable(false)
                    ->limit(50),
            ])
            ->filters([
                Tables\Filters\TrashedFilter::make(),
            ])
            ->actions([
                Tables\Actions\ActionGroup::make([
                    Tables\Actions\ViewAction::make()
                        ->color('success')
                        ->label('View')
                        ->icon('heroicon-o-eye'),
                    Tables\Actions\EditAction::make()
                        ->color('warning')
                        ->label('Details')
                        ->icon('heroicon-o-pencil')
                        ->closeModalByClickingAway(false),
                    Tables\Actions\Action::make('password')
                        ->color('warning')
                        ->label('Password')
                        ->icon('heroicon-o-key')
                        ->form([
                            Forms\Components\TextInput::make('password')
                                ->label('New Password')
                                ->required()
                                ->password()
                                ->maxLength(32)
                                ->revealable(),
                        ])
                        ->action(function (User $record, array $data) {
                            $record->password = bcrypt($data['password']);
                            $record->save();
                        })
                        ->authorize(fn() => static::grant(PermissionEnum::CHANGE_PASSWORD_STUDENT->value))
                        ->closeModalByClickingAway(false),
                    Tables\Actions\DeleteAction::make(),
                    Tables\Actions\ForceDeleteAction::make(),
                    Tables\Actions\RestoreAction::make(),
                ])
            ])
            ->bulkActions([
                Tables\Actions\BulkActionGroup::make([
                    Tables\Actions\DeleteBulkAction::make(),
                    Tables\Actions\ForceDeleteBulkAction::make(),
                    Tables\Actions\RestoreBulkAction::make(),
                ]),
            ]);
    }

    public static function getPages(): array
    {
        return [
            'index' => Pages\ManageStudents::route('/'),
        ];
    }

    public static function getEloquentQuery(): Builder
    {
        return parent::getEloquentQuery()
            ->with(['groups', 'roles'])
            ->whereHas('roles', function (Builder $query) {
                $query->where('name', RoleEnum::STUDENT->value);
            })
            ->whereHas('groups')
            ->withoutGlobalScopes([
                SoftDeletingScope::class,
            ]);
    }

    // Custom Permissions
    public static function canViewAny(): bool
    {
        return static::can('viewAnyStudent');
    }

    public static function canCreate(): bool
    {
        return static::can('createStudent');
    }

    public static function canEdit(Model $record): bool
    {
        return static::can('updateStudent', $record);
    }

    public static function canDelete(Model $record): bool
    {
        return static::can('deleteStudent', $record);
    }

    public static function canDeleteAny(): bool
    {
        return static::can('deleteAnyStudent');
    }

    public static function canForceDelete(Model $record): bool
    {
        return static::can('forceDeleteStudent', $record);
    }

    public static function canForceDeleteAny(): bool
    {
        return static::can('forceDeleteAnyStudent');
    }

    public static function canReorder(): bool
    {
        return static::can('reorderStudent');
    }

    public static function canReplicate(Model $record): bool
    {
        return static::can('replicateStudent', $record);
    }

    public static function canRestore(Model $record): bool
    {
        return static::can('restoreStudent', $record);
    }

    public static function canRestoreAny(): bool
    {
        return static::can('restoreAnyStudent');
    }

    public static function canView(Model $record): bool
    {
        return static::can('viewStudent', $record);
    }

    public static function authorizeViewAny(): void
    {
        static::authorize('viewAnyStudent');
    }

    public static function authorizeCreate(): void
    {
        static::authorize('createStudent');
    }

    public static function authorizeEdit(Model $record): void
    {
        static::authorize('updateStudent', $record);
    }

    public static function authorizeView(Model $record): void
    {
        static::authorize('viewStudent', $record);
    }
}
