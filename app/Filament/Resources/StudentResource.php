<?php

namespace App\Filament\Resources;

use App\Enums\PermissionEnum as PE;
use App\Filament\Resources\StudentResource\Pages;
use App\Filament\Resources\StudentResource\RelationManagers;
use App\Models\Group;
use App\Models\Student;
use App\Traits\Api\AuthorizeTrait;
use Filament\Forms;
use Filament\Forms\Form;
use Filament\Resources\Resource;
use Filament\Tables;
use Filament\Tables\Table;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\SoftDeletingScope;
use Rmsramos\Activitylog\RelationManagers\ActivitylogRelationManager;

class StudentResource extends Resource
{
    use AuthorizeTrait;

    protected static ?string $model = Student::class;

    protected static ?string $navigationGroup = 'Users';
    protected static ?string $navigationLabel = 'Students';
    protected static ?int $navigationSort = 2;
    protected static ?string $label = 'Student';
    protected static ?string $pluralLabel = 'Data Students';

    public static function form(Form $form): Form
    {
        return $form
            ->schema([
                Forms\Components\Wizard::make()
                    ->steps([
                        Forms\Components\Wizard\Step::make('Identity')
                            ->schema([
                                Forms\Components\TextInput::make('name')
                                    ->columnSpanFull()
                                    ->required()
                                    ->maxLength(255),
                                Forms\Components\TextInput::make('username')
                                    ->columnSpanFull()
                                    ->required()
                                    ->unique(ignoreRecord: true)
                                    ->maxLength(255),
                                Forms\Components\TextInput::make('password')
                                    ->columnSpanFull()
                                    ->password()
                                    ->required()
                                    ->maxLength(255)
                                    ->revealable()
                                    ->visibleOn('create'),
                            ]),
                        Forms\Components\Wizard\Step::make('Groups')
                            ->schema([
                                Forms\Components\Select::make('groups')
                                    ->columnSpanFull()
                                    ->relationship('groups', 'name')
                                    ->required()
                                    ->multiple(false)
                                    ->label('Groups')
                                    ->afterStateUpdated(function (callable $set, $state, $record) {
                                        $set('groups', $state);

                                        if ($record) {
                                            activity('student')
                                                ->performedOn($record)
                                                ->event('updated')
                                                ->withProperties([
                                                    'attributes' => [
                                                        'name' => $record->name,
                                                        'groups' => is_array($state) ? Group::whereIn('id', $state)->pluck('name')->toArray() : [],
                                                    ],
                                                    'old' => [
                                                        'name' => $record->name,
                                                        'groups' => $record->groups->pluck('name')->toArray(),
                                                    ],
                                                ])
                                                ->log('Updated groups');
                                        }
                                    }),
                            ]),
                    ])
                    ->columnSpanFull(),
            ]);
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                Tables\Columns\TextColumn::make('name')
                    ->searchable(),
                Tables\Columns\TextColumn::make('username')
                    ->searchable()
                    ->copyable()
                    ->copyMessage('Copied to clipboard'),
                Tables\Columns\TextColumn::make('groups')
                    ->label('Groups')
                    ->limit(30)
                    ->getStateUsing(fn(Student $record): string => $record->groups ? $record->groups->pluck('name')->implode(', ') : ''),
            ])
            ->filters([
                Tables\Filters\TrashedFilter::make(),
            ])
            ->actions([
                Tables\Actions\ActionGroup::make([
                    Tables\Actions\EditAction::make()
                        ->color('warning')
                        ->icon('heroicon-o-pencil')
                        ->label('Details'),
                    Tables\Actions\Action::make('Password')
                        ->color('success')
                        ->icon('heroicon-o-key')
                        ->authorize(fn() => static::grant(PE::CHANGE_PASSWORD_STUDENT->value))
                        ->url(fn(Student $record): string =>  self::getUrl('password', ['record' => $record->id])),
                    Tables\Actions\DeleteAction::make(),
                    Tables\Actions\ForceDeleteAction::make(),
                    Tables\Actions\RestoreAction::make(),
                ]),
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
            'password' => Pages\ChangePasswordStudent::route('/{record}/password'),
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
