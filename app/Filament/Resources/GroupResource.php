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
    protected static ?int $navigationSort = 2;
    protected static ?string $label = 'Group';
    protected static ?string $pluralLabel = 'Data Groups';

    public static function form(Form $form): Form
    {
        return $form
            ->schema([
                // group
                Forms\Components\Fieldset::make()
                    ->label('Group Information')
                    ->schema([
                        // name
                        Forms\Components\TextInput::make('name')
                            ->columnSpanFull()
                            ->label('Name')
                            ->required()
                            ->maxLength(50),

                        // description
                        Forms\Components\Textarea::make('description')
                            ->columnSpanFull()
                            ->label('Description')
                            ->rows(2)
                            ->maxLength(255),

                        // course
                        Forms\Components\Select::make('course_id')
                            ->columnSpanFull()
                            ->label('Course')
                            ->options(\App\Models\Course::query()->pluck('name', 'id'))
                            ->required(),
                    ]),
                // participants
                Forms\Components\Fieldset::make()
                    ->label('Participants')
                    ->columnSpanFull()
                    ->schema([
                        // teachers
                        Forms\Components\CheckboxList::make('teachers')
                            ->label('Teachers')
                            ->relationship('teachers', 'name')
                            ->options(
                                \App\Models\User::query()
                                    ->whereHas('roles', function ($query) {
                                        $query->whereIn('name', ['teacher']);
                                    })
                                    ->pluck('name', 'id')
                            )
                            ->columns(2)
                            ->helperText('Select one or more teachers'),

                        // students
                        Forms\Components\CheckboxList::make('students')
                            ->label('Students')
                            ->relationship('students', 'name')
                            ->options(
                                \App\Models\User::query()
                                    ->whereHas('roles', function ($query) {
                                        $query->whereIn('name', ['student']);
                                    })
                                    ->whereDoesntHave('groups')
                                    ->pluck('name', 'id')
                            )
                            ->columns(2)
                            ->helperText('Select one or more students'),
                    ]),
            ]);
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                Tables\Columns\TextColumn::make('name')
                    ->label('Group')
                    ->searchable(),
                Tables\Columns\TextColumn::make('participants_count')
                    ->label('Participants')
                    ->counts('participants'),
            ])
            ->groups([
                Tables\Grouping\Group::make('course.name')
                    ->titlePrefixedWithLabel(true),
            ])
            ->filters([
                Tables\Filters\TrashedFilter::make(),
            ])
            ->actions([
                Tables\Actions\ActionGroup::make([
                    Tables\Actions\ViewAction::make()
                        ->color('success')
                        ->label('View')
                        ->icon('heroicon-o-eye')
                        ->mutateRecordDataUsing(function (array $data): array {
                            if (!isset($data['record']) || !$data['record']) {
                                return $data;
                            }

                            // Get the group with the participants
                            $group = Group::with(['participants'])->find($data['record']->id);
                            if (!$group) {
                                return $data;
                            }

                            // Add the participants to the data
                            $data['participants'] = $group->participants;

                            return $data;
                        })
                        ->form([
                            // fieldset
                            Forms\Components\Fieldset::make()
                                ->label('Group Information')
                                ->schema([
                                    // name
                                    Forms\Components\TextInput::make('name')
                                        ->label('Name')
                                        ->columns(3)
                                        ->maxLength(50)
                                        ->disabled(),
                                    // course name
                                    Forms\Components\Radio::make('course.name')
                                        ->label('Course')
                                        ->columns(1)
                                        ->options(function ($record) {
                                            return \App\Models\Course::query()
                                                ->where('id', $record->course_id ?? null)
                                                ->pluck('name', 'id');
                                        })
                                        ->disabled(),
                                    // description
                                    Forms\Components\Textarea::make('description')
                                        ->label('Description')
                                        ->columnSpanFull()
                                        ->rows(2)
                                        ->maxLength(255)
                                        ->disabled(),
                                ]),
                            // participants
                            Forms\Components\Fieldset::make()
                                ->label('Participants')
                                ->schema([
                                    // teachers
                                    Forms\Components\CheckboxList::make('teachers')
                                        ->label('Teachers')
                                        ->relationship('teachers', 'name')
                                        ->options(
                                            \App\Models\User::query()
                                                ->whereHas('roles', function ($query) {
                                                    $query->whereIn('name', ['teacher']);
                                                })
                                                ->pluck('name', 'id')
                                        )
                                        ->columns(2),

                                    // students
                                    Forms\Components\CheckboxList::make('students')
                                        ->label('Students')
                                        ->relationship('students', 'name')
                                        ->options(
                                            \App\Models\User::query()
                                                ->whereHas('roles', function ($query) {
                                                    $query->whereIn('name', ['student']);
                                                })
                                                ->whereDoesntHave('groups')
                                                ->pluck('name', 'id')
                                        )
                                        ->columns(2),
                                ]),
                        ]),
                    Tables\Actions\EditAction::make()
                        ->color('warning')
                        ->label('Details')
                        ->icon('heroicon-o-pencil')
                        ->closeModalByClickingAway(false)
                        ->form([
                            // fieldset
                            Forms\Components\Fieldset::make()
                                ->label('Group Information')
                                ->schema([
                                    // name
                                    Forms\Components\TextInput::make('name')
                                        ->label('Name')
                                        ->columns(3)
                                        ->maxLength(50)
                                        ->required(),
                                    // course name
                                    Forms\Components\Radio::make('course.name')
                                        ->label('Course')
                                        ->columns(1)
                                        ->options(function ($record) {
                                            return \App\Models\Course::query()
                                                ->where('id', $record->course_id ?? null)
                                                ->pluck('name', 'id');
                                        })
                                        ->disabled(),
                                    // description
                                    Forms\Components\Textarea::make('description')
                                        ->label('Description')
                                        ->columnSpanFull()
                                        ->rows(2)
                                        ->maxLength(255),
                                ]),
                            // participants
                            Forms\Components\Fieldset::make()
                                ->label('Participants')
                                ->schema([
                                    // teachers
                                    Forms\Components\CheckboxList::make('teachers')
                                        ->label('Teachers')
                                        ->relationship('teachers', 'name')
                                        ->options(
                                            \App\Models\User::query()
                                                ->whereHas('roles', function ($query) {
                                                    $query->whereIn('name', ['teacher']);
                                                })
                                                ->pluck('name', 'id')
                                        )
                                        ->columns(2),

                                    // students
                                    Forms\Components\CheckboxList::make('students')
                                        ->label('Students')
                                        ->relationship('students', 'name')
                                        ->options(
                                            \App\Models\User::query()
                                                ->whereHas('roles', function ($query) {
                                                    $query->whereIn('name', ['student']);
                                                })
                                                ->whereDoesntHave('groups')
                                                ->pluck('name', 'id')
                                        )
                                        ->columns(2),
                                ]),
                        ]),
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
            'index' => Pages\ManageGroups::route('/'),
        ];
    }

    public static function getEloquentQuery(): Builder
    {
        return parent::getEloquentQuery()
            ->with('course')
            ->withoutGlobalScopes([
                SoftDeletingScope::class,
            ]);
    }
}
