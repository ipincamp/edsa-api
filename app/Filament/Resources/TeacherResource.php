<?php

namespace App\Filament\Resources;

use App\Enums\PermissionEnum as PE;
use App\Filament\Resources\TeacherResource\Pages;
use App\Filament\Resources\TeacherResource\RelationManagers;
use App\Models\Teacher;
use App\Traits\AuthorizeTrait;
use Filament\Forms;
use Filament\Forms\Form;
use Filament\Resources\Resource;
use Filament\Tables;
use Filament\Tables\Table;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\SoftDeletingScope;

class TeacherResource extends Resource
{
    use AuthorizeTrait;

    protected static ?string $model = Teacher::class;

    protected static ?string $navigationIcon = 'heroicon-o-users';
    protected static ?string $navigationGroup = 'Users';
    protected static ?string $navigationLabel = 'Teachers';
    protected static ?string $label = 'Teacher';
    protected static ?string $pluralLabel = 'Data Teachers';

    public static function form(Form $form): Form
    {
        return $form
            ->schema([
                Forms\Components\TextInput::make('name')
                    ->required()
                    ->maxLength(255),
                Forms\Components\TextInput::make('email')
                    ->required()
                    ->unique(ignoreRecord: true)
                    ->maxLength(255)
                    ->email(),
                Forms\Components\TextInput::make('password')
                    ->columnSpanFull()
                    ->password()
                    ->required()
                    ->maxLength(255)
                    ->revealable()
                    ->visibleOn('create'),
            ]);
    }

    public static function table(Table $table): Table
    {
        // sembunyikan teacher dengan role admin

        return $table
            ->columns([
                Tables\Columns\TextColumn::make('name')
                    ->searchable(),
                Tables\Columns\TextColumn::make('email')
                    ->searchable()
                    ->copyable()
                    ->copyMessage('Copied to clipboard'),
                // TODO: group name
            ])
            ->filters([
                Tables\Filters\TrashedFilter::make(),
            ])
            ->actions([
                Tables\Actions\EditAction::make(),
                Tables\Actions\Action::make('Password')
                    ->color('success')
                    ->icon('heroicon-o-key')
                    ->authorize(fn() => static::grant(PE::CHANGE_PASSWORD_TEACHER->value))
                    ->url(fn(Teacher $record): string =>  self::getUrl('password', ['record' => $record->id])),
                Tables\Actions\DeleteAction::make(),
                Tables\Actions\ForceDeleteAction::make(),
                Tables\Actions\RestoreAction::make(),
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
            'password' => Pages\ChangePasswordTeacher::route('/{record}/password'),
        ];
    }

    public static function getEloquentQuery(): Builder
    {
        return parent::getEloquentQuery()
            ->withoutGlobalScopes([
                SoftDeletingScope::class,
            ]);
    }
}
