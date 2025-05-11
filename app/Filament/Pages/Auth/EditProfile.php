<?php

namespace App\Filament\Pages\Auth;

use App\Enums\RoleEnum;
use App\Filament\Widgets\ProfileOverview;
use Filament\Forms;
use Filament\Forms\Form;
use Filament\Pages\Auth\EditProfile as BaseEditProfile;

class EditProfile extends BaseEditProfile
{
    protected static bool $isLazy = false;

    public static function getLabel(): string
    {
        return 'Profile';
    }

    protected function getHeaderWidgets(): array
    {
        return [
            ProfileOverview::class,
        ];
    }

    public function form(Form $form): Form
    {
        return $form
            ->schema([
                $this->getNameFormComponent(),
                Forms\Components\TextInput::make('username')
                    ->label('Username')
                    ->required()
                    ->maxLength(50)
                    ->unique(ignoreRecord: true),
                Forms\Components\TextInput::make('email')
                    ->label(__('filament-panels::pages/auth/edit-profile.form.email.label'))
                    ->email()
                    ->required()
                    ->maxLength(255)
                    ->unique(ignoreRecord: true)
                    ->hidden(fn($record) => $record->hasRole(RoleEnum::STUDENT)),
                $this->getPasswordFormComponent(),
                $this->getPasswordConfirmationFormComponent(),
            ]);
    }
}
