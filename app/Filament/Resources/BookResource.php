<?php

namespace App\Filament\Resources;

use App\Filament\Resources\BookResource\Pages;
use App\Filament\Resources\BookResource\RelationManagers;
use App\Models\Book;
use Filament\Forms;
use Filament\Forms\Form;
use Filament\Resources\Resource;
use Filament\Tables;
use Filament\Tables\Table;
use Illuminate\Database\Eloquent\Builder;
use Illuminate\Database\Eloquent\SoftDeletingScope;

class BookResource extends Resource
{
    protected static ?string $model = Book::class;

    protected static ?string $navigationGroup = 'Managements';
    protected static ?string $navigationLabel = 'Books';
    protected static ?int $navigationSort = 3;
    protected static ?string $label = 'Book';
    protected static ?string $pluralLabel = 'Data Books';

    public static function form(Form $form): Form
    {
        return $form
            ->schema([
                //
            ]);
    }

    public static function table(Table $table): Table
    {
        return $table
            ->columns([
                // cover image
                Tables\Columns\ImageColumn::make('image')
                    ->label('Cover')
                    ->circular()
                    ->size(64)
                    ->default(config('app.url') . '/assets/default.jpg'),
                // title
                Tables\Columns\TextColumn::make('title')
                    ->searchable()
                    ->sortable()
                    ->limit(50),

                // author
                Tables\Columns\TextColumn::make('author')
                    ->searchable()
                    ->sortable()
                    ->limit(50),

                // year
                Tables\Columns\TextColumn::make('year')
                    ->searchable()
                    ->sortable()
                    ->limit(4),

                // status
                Tables\Columns\ToggleColumn::make('status')
                    ->onColor('success')
                    ->offColor('danger')
                    ->tooltip(fn($record) => $record->status ? 'The book is open for everyone' : 'The book is locked')
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
                            $data['status'] = $data['status'] ? 'Open' : 'Locked';

                            return $data;
                        })
                        ->form([
                            Forms\Components\Fieldset::make()
                                ->label(false)
                                ->schema([
                                    Forms\Components\Grid::make(4)
                                        ->schema([
                                            Forms\Components\Fieldset::make()
                                                ->label('Book Cover')
                                                ->schema([
                                                    Forms\Components\FileUpload::make('image')
                                                        ->label(false)
                                                        ->image()
                                                        ->columnSpanFull(),
                                                ])
                                                ->columnSpan(2),
                                            Forms\Components\Fieldset::make()
                                                ->label(fn($state) => 'The book is ' . $state['status'])
                                                ->schema([
                                                    Forms\Components\TextInput::make('title')
                                                        ->label('Title')
                                                        ->maxLength(100)
                                                        ->columnSpanFull(),
                                                    Forms\Components\TextInput::make('author')
                                                        ->label('Author')
                                                        ->maxLength(20)
                                                        ->columnSpan(1),
                                                    Forms\Components\TextInput::make('year')
                                                        ->label('Year')
                                                        ->maxLength(4)
                                                        ->columnSpan(1),
                                                    Forms\Components\TextInput::make('genre')
                                                        ->label('Genre')
                                                        ->maxLength(20)
                                                        ->columnSpan(1),
                                                    Forms\Components\TextInput::make('focus')
                                                        ->label('Focus')
                                                        ->maxLength(20)
                                                        ->columnSpan(1),
                                                ])
                                                ->columnSpan(2),
                                            Forms\Components\Fieldset::make()
                                                ->label('Settings')
                                                ->schema([
                                                    Forms\Components\Repeater::make('settings')
                                                        ->label(false)
                                                        ->relationship('settings')
                                                        ->schema([
                                                            Forms\Components\TextInput::make('key')
                                                                ->label('Key')
                                                                ->maxLength(50),
                                                            Forms\Components\Textarea::make('value')
                                                                ->label('Value')
                                                                ->rows(3),
                                                        ])
                                                        ->columnSpanFull(),
                                                ]),
                                        ]),
                                ]),
                        ]),
                    Tables\Actions\Action::make('settings')
                        ->color('primary')
                        ->label('Settings')
                        ->icon('heroicon-o-cog')
                        ->form([
                            Forms\Components\Fieldset::make()
                                ->label(false)
                                ->schema([
                                    Forms\Components\Select::make('key')
                                        ->label(false)
                                        ->options([
                                            'ramdom-word' => 'Random Word',
                                        ])
                                        ->required()
                                        ->placeholder('Select a key to set')
                                        ->columnSpanFull(),
                                    Forms\Components\Textarea::make('value')
                                        ->label('Value')
                                        ->required()
                                        ->rows(3)
                                        ->placeholder('Separated by commas')
                                        ->columnSpanFull(),
                                ]),
                        ])
                        ->action(function (array $data, Book $record) {
                            $record->settings()->updateOrCreate(
                                ['key' => $data['key']],
                                ['value' => $data['value']]
                            );

                            \Filament\Notifications\Notification::make()
                                ->title('Settings Updated')
                                ->body('The settings for ' . $record->title . ' have been updated.')
                                ->success()
                                ->send();
                        }),
                    /*
                    Tables\Actions\EditAction::make()
                        ->color('warning')
                        ->label('Details')
                        ->icon('heroicon-o-pencil'),
                    Tables\Actions\DeleteAction::make(),
                    Tables\Actions\ForceDeleteAction::make(),
                    Tables\Actions\RestoreAction::make(),
                    */
                ]),
            ])
            ->bulkActions([
                // Tables\Actions\BulkActionGroup::make([
                //     Tables\Actions\DeleteBulkAction::make(),
                //     Tables\Actions\ForceDeleteBulkAction::make(),
                //     Tables\Actions\RestoreBulkAction::make(),
                // ]),
            ]);
    }

    public static function getPages(): array
    {
        return [
            'index' => Pages\ManageBooks::route('/'),
        ];
    }

    public static function getEloquentQuery(): Builder
    {
        return parent::getEloquentQuery()
            ->with('settings')
            ->withoutGlobalScopes([
                SoftDeletingScope::class,
            ]);
    }
}
