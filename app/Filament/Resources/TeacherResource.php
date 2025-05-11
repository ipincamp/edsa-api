<?php

namespace App\Filament\Resources;

use App\Enums\RoleEnum;
use App\Filament\Resources\TeacherResource\Pages;
use App\Filament\Resources\TeacherResource\RelationManagers;
use App\Models\Group;
use App\Models\User;
use Filament\Forms;
use Filament\Forms\Form;
use Filament\Resources\Resource;
use Filament\Tables;
use Filament\Tables\Table;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\SoftDeletingScope;

class TeacherResource extends Resource
{
    protected static ?string $model = User::class;

    protected static ?string $navigationGroup = 'Users';
    protected static ?string $navigationLabel = 'Teachers';
    protected static ?int $navigationSort = 1;
    protected static ?string $label = 'Teacher';
    protected static ?string $pluralLabel = 'Data Teachers';

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
                // email
                Forms\Components\TextInput::make('email')
                    ->label('Email')
                    ->required()
                    ->email()
                    ->maxLength(255)
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
                // email
                Tables\Columns\TextColumn::make('email')
                    ->label('Email')
                    ->searchable()
                    ->limit(50)
                    ->formatStateUsing(fn($state) => substr($state, 0, 2) . '***@' . substr($state, strpos($state, '@') + 1)),
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
                        // TODO: Authorize this action
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
            'index' => Pages\ManageTeachers::route('/'),
        ];
    }

    public static function getEloquentQuery(): Builder
    {
        return parent::getEloquentQuery()
            ->with(['groups', 'roles'])
            ->whereHas('roles', function (Builder $query) {
                $query->where('name', RoleEnum::TEACHER->value);
            })
            ->withoutGlobalScopes([
                SoftDeletingScope::class,
            ]);
    }
}
