<?php

namespace App\Filament\Resources;

use App\Filament\Resources\GroupResource\Pages;
use App\Filament\Resources\GroupResource\RelationManagers;
use App\Models\Group;
use App\Models\Student;
use App\Models\Teacher;
use Filament\Forms;
use Filament\Forms\Form;
use Filament\Resources\Resource;
use Filament\Tables;
use Filament\Tables\Table;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\SoftDeletingScope;

class GroupResource extends Resource
{
    protected static ?string $model = Group::class;

    protected static ?string $navigationGroup = 'Managements';
    protected static ?string $navigationLabel = 'Groups';
    protected static ?int $navigationSort = 1;
    protected static ?string $label = 'Group';
    protected static ?string $pluralLabel = 'Group of Class';
    protected static ?string $slug = 'groups';

    public static function form(Form $form): Form
    {
        $columns = [
            'default' => 1,
            'sm' => 2,
            'md' => 3,
        ];

        return $form
            ->schema([
                Forms\Components\Wizard::make()
                    ->steps([
                        Forms\Components\Wizard\Step::make('Details')
                            ->schema([
                                Forms\Components\TextInput::make('name')
                                    ->columnSpanFull()
                                    ->columns($columns)
                                    ->required()
                                    ->label('Name')
                                    ->maxLength(255)
                                    ->unique(ignoreRecord: true),
                                Forms\Components\Textarea::make('description')
                                    ->columnSpanFull()
                                    ->columns($columns)
                                    ->label('Description')
                                    ->maxLength(255),
                            ]),
                        Forms\Components\Wizard\Step::make('Users')
                            ->schema([
                                Forms\Components\CheckboxList::make('teachers')
                                    ->columns($columns)
                                    ->relationship('teachers', 'name')
                                    ->required()
                                    ->label('Teachers')
                                    ->bulkToggleable()
                                    ->afterStateUpdated(function (callable $set, $state, $record) {
                                        $set('teachers', $state);

                                        if ($record) {
                                            activity('group')
                                                ->performedOn($record)
                                                ->event('updated')
                                                ->withProperties([
                                                    'attributes' => [
                                                        'name' => $record->name,
                                                        'teachers' => Teacher::whereIn('id', $state)->pluck('name')->toArray(),
                                                    ],
                                                    'old' => [
                                                        'name' => $record->name,
                                                        'teachers' => $record->teachers->pluck('name')->toArray(),
                                                    ],
                                                ])
                                                ->log('Updated teacher in group');
                                        }
                                    }),
                                Forms\Components\CheckboxList::make('students')
                                    ->columns($columns)
                                    ->relationship('students', 'name')
                                    ->required()
                                    ->label('Students')
                                    ->bulkToggleable()
                                    ->afterStateUpdated(function (callable $set, $state, $record) {
                                        $set('students', $state);

                                        if ($record) {
                                            activity('group')
                                                ->performedOn($record)
                                                ->event('updated')
                                                ->withProperties([
                                                    'attributes' => [
                                                        'name' => $record->name,
                                                        'students' => Student::whereIn('id', $state)->pluck('name')->toArray(),
                                                    ],
                                                    'old' => [
                                                        'name' => $record->name,
                                                        'students' => $record->students->pluck('name')->toArray(),
                                                    ],
                                                ])
                                                ->log('Updated students in group');
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
                    ->label('Group Name')
                    ->searchable()
                    ->sortable(),
                Tables\Columns\TextColumn::make('description')
                    ->label('Description')
                    ->limit(20)
                    ->searchable(),
                Tables\Columns\TextColumn::make('participants')
                    ->label('Participants')
                    ->getStateUsing(function ($record) {
                        return $record->teachers()->count() + $record->students()->count();
                    }),
            ])
            ->filters([
                Tables\Filters\TrashedFilter::make(),
            ])
            ->actions([
                Tables\Actions\ActionGroup::make([
                    Tables\Actions\ViewAction::make(),
                    Tables\Actions\EditAction::make(),
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
            'index' => Pages\ManageGroups::route('/'),
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
