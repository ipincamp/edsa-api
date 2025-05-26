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
                                            Forms\Components\Repeater::make('images')
                                                ->label('Random Words')
                                                ->relationship('images')
                                                ->schema([
                                                    Forms\Components\Section::make(fn($record) => $record->name)
                                                        ->schema([
                                                            Forms\Components\FileUpload::make('path')
                                                                ->label(false)
                                                                ->avatar()
                                                                ->alignment('center'),
                                                        ]),
                                                ])
                                                ->grid(4)
                                                ->columnSpanFull(),
                                        ]),
                                ]),
                        ]),
                    Tables\Actions\EditAction::make('editRandomWord')
                        ->color('warning')
                        ->label('Random Words')
                        ->icon('heroicon-o-pencil')
                        ->form([
                            Forms\Components\Fieldset::make()
                                ->label(false)
                                ->schema([
                                    Forms\Components\Grid::make(4)
                                        ->schema([
                                            Forms\Components\Repeater::make('images')
                                                ->label('Random Words')
                                                ->relationship('images')
                                                ->schema([
                                                    Forms\Components\FileUpload::make('path')
                                                        ->required()
                                                        ->label('Image')
                                                        ->image()
                                                        ->visibility('public')
                                                        ->placeholder('Upload an image for the word'),
                                                    Forms\Components\TextInput::make('name')
                                                        ->required()
                                                        ->label('Word')
                                                        ->maxLength(15)
                                                        ->placeholder('Enter the word'),
                                                ])
                                                ->grid(4)
                                                ->columnSpanFull()
                                                ->addActionLabel('Add More Words'),
                                        ]),
                                ]),
                        ])
                        ->action(function (array $data, Book $record) {
                            \Filament\Notifications\Notification::make()
                                ->title('Random Word Updated')
                                ->body('Random words for ' . $record->title . ' have been updated.')
                                ->success()
                                ->send();
                        }),
                    /*
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
            ->with('images')
            ->withoutGlobalScopes([
                SoftDeletingScope::class,
            ]);
    }
}
