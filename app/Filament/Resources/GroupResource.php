<?php

namespace App\Filament\Resources;

use App\Filament\Resources\GroupResource\Pages;
use App\Filament\Resources\GroupResource\RelationManagers;
use App\Models\Group;
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
            // edit data relationship users with group_participants table
            ->schema([
                Forms\Components\TextInput::make('name')
                    ->required()
                    ->columns($columns)
                    ->columnSpanFull()
                    ->label('Name')
                    ->maxLength(255)
                    ->unique(ignoreRecord: true)
                    ->afterStateUpdated(function (callable $set, $state, $record) {
                        if ($record) {
                            activity('groups')
                                ->performedOn($record)
                                ->event('updated')
                                ->withProperties([
                                    'attributes' => [
                                        'name' => $state,
                                    ],
                                    'old' => [
                                        'name' => $record->name,
                                    ],
                                ])
                                ->log('Updated group name');
                        }
                    }),
                Forms\Components\Textarea::make('description')
                    ->columns($columns)
                    ->columnSpanFull()
                    ->label('Description')
                    ->maxLength(255)
                    ->afterStateUpdated(function (callable $set, $state, $record) {
                        if ($record) {
                            activity('groups')
                                ->performedOn($record)
                                ->event('updated')
                                ->withProperties([
                                    'attributes' => [
                                        'description' => $state,
                                    ],
                                    'old' => [
                                        'description' => $record->description,
                                    ],
                                ])
                                ->log('Updated group description');
                        }
                    }),
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
            ])
            ->filters([
                Tables\Filters\TrashedFilter::make(),
            ])
            ->actions([
                Tables\Actions\ViewAction::make(),
                Tables\Actions\EditAction::make(),
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
