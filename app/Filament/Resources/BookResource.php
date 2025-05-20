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
                        ->icon('heroicon-o-eye'),
                    /*
                    Tables\Actions\EditAction::make()
                        ->color('warning')
                        ->label('Details')
                        ->icon('heroicon-o-pencil'),
                    */
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
            'index' => Pages\ManageBooks::route('/'),
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
