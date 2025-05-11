<?php

namespace App\Filament\Resources;

use App\Filament\Resources\CourseResource\Pages;
use App\Filament\Resources\CourseResource\RelationManagers;
use App\Models\Course;
use Filament\Forms;
use Filament\Forms\Form;
use Filament\Resources\Resource;
use Filament\Tables;
use Filament\Tables\Table;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\SoftDeletingScope;

class CourseResource extends Resource
{
    protected static ?string $model = Course::class;

    protected static ?string $navigationGroup = 'Managements';
    protected static ?string $navigationLabel = 'Courses';
    protected static ?int $navigationSort = 1;
    protected static ?string $label = 'Course';
    protected static ?string $pluralLabel = 'Data Courses';

    public static function form(Form $form): Form
    {
        return $form
            ->schema([
                // name
                Forms\Components\TextInput::make('name')
                    ->columnSpanFull()
                    ->label('Name')
                    ->required()
                    ->maxLength(255),

                // description
                Forms\Components\Textarea::make('description')
                    ->columnSpanFull()
                    ->label('Description')
                    ->rows(3)
                    ->maxLength(65535),
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
                    ->limit(10),

                // description
                Tables\Columns\TextColumn::make('description')
                    ->label('Description')
                    ->limit(20)
                    ->toggleable(),

                // total groups
                Tables\Columns\TextColumn::make('groups_count')
                    ->label('Total Groups')
                    ->counts('groups')
                    ->toggleable(),
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

                            // Get the course with the groups and participants
                            $course = Course::with(['groups', 'groups.participants'])->find($data['record']->id);
                            if (!$course) {
                                return $data;
                            }

                            // Get the groups and participants
                            $data['groups'] = $course->groups->map(function ($group) {
                                return [
                                    'name' => $group->name,
                                    'participants' => $group->participants->map(function ($participant) {
                                        return [
                                            'name' => $participant->name,
                                            'role' => $participant->roles->pluck('name')->first(),
                                        ];
                                    }),
                                ];
                            });

                            // Get the participants
                            $data['participants'] = $course->groups->flatMap(function ($group) {
                                return $group->participants;
                            });

                            // Get the total groups
                            $data['total_groups'] = $course->groups->count();

                            return $data;
                        })
                        ->form([
                            // course
                            Forms\Components\Fieldset::make()
                                ->label('Course Details')
                                ->columnSpanFull()
                                ->schema([
                                    // name
                                    Forms\Components\TextInput::make('name')
                                        ->label('Name')
                                        ->columnSpanFull()
                                        ->maxLength(255)
                                        ->disabled(),

                                    // description
                                    Forms\Components\Textarea::make('description')
                                        ->label('Description')
                                        ->columnSpanFull()
                                        ->rows(3)
                                        ->maxLength(65535)
                                ]),
                            // groups
                            Forms\Components\Fieldset::make()
                                ->label('Groups')
                                ->schema([
                                    Forms\Components\Repeater::make('groups')
                                        ->label(false)
                                        ->relationship('groups')
                                        ->columnSpanFull()
                                        ->schema([
                                            Forms\Components\Section::make(fn($record) => $record->name)
                                                ->description(fn($record) => $record->description . '. Total ' . count($record->participants) . ' participants.')
                                                ->schema([
                                                    Forms\Components\CheckboxList::make('participants')
                                                        ->label(false)
                                                        ->columns(2)
                                                        ->options(fn($record) => $record->participants->mapWithKeys(function ($participant) {
                                                            return [$participant->id => $participant->name . ' - ' . $participant->roles->pluck('name')->first()];
                                                        }))
                                                        ->disabled(),
                                                ])
                                                ->collapsed(),
                                        ]),
                                ]),
                        ]),
                    Tables\Actions\EditAction::make()
                        ->color('warning')
                        ->label('Details')
                        ->icon('heroicon-o-pencil')
                        ->closeModalByClickingAway(false),
                    Tables\Actions\Action::make('new_group')
                        ->color('success')
                        ->label('New Group')
                        ->icon('heroicon-o-plus')
                        ->modalHeading(fn($record) => 'New Group for ' . $record->name . ' Course')
                        ->form([
                            Forms\Components\TextInput::make('name')
                                ->label('Group Name')
                                ->required()
                                ->maxLength(50),
                            Forms\Components\Textarea::make('description')
                                ->label('Group Description')
                                ->rows(2)
                                ->maxLength(255),
                        ])
                        ->action(function (array $data, Course $record) {
                            $record->groups()->create($data);
                        }),
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
            'index' => Pages\ManageCourses::route('/'),
        ];
    }

    public static function getEloquentQuery(): Builder
    {
        return parent::getEloquentQuery()
            ->with(['groups', 'groups.participants'])
            ->withoutGlobalScopes([
                SoftDeletingScope::class,
            ]);
    }
}
