<?php

namespace App\Filament\Resources\StudentResource\Pages;

use App\Filament\Resources\StudentResource;
use Filament\Actions\Action;
use Filament\Forms;
use Filament\Forms\Concerns\InteractsWithForms;
use Filament\Forms\Contracts\HasForms;
use Filament\Notifications\Notification;
use Filament\Resources\Pages\Concerns\InteractsWithRecord;
use Filament\Resources\Pages\Page;
use Filament\Support\Exceptions\Halt;
use Illuminate\Support\Facades\Auth;

class ChangePasswordStudent extends Page implements HasForms
{
    use InteractsWithForms, InteractsWithRecord;

    protected static string $resource = StudentResource::class;

    protected static string $view = 'filament.resources.student-resource.pages.change-password-student';

    public ?array $data = [];

    public function mount(int | string $record): void
    {
        $this->record = StudentResource::getModel()::findOrFail($record);

        $this->form->fill();
    }

    public function getTitle(): string
    {
        return __('Change Password for :name', [
            'name' => $this->record->name,
        ]);
    }

    protected function getFormSchema(): array
    {
        return [

            // your_password
            Forms\Components\TextInput::make('data.your_password')
                ->password()
                ->required()
                ->maxLength(255)
                ->label('Your Password')
                ->revealable(),

            // new_password
            Forms\Components\TextInput::make('data.new_password')
                ->password()
                ->required()
                ->maxLength(255)
                ->label('New Password')
                ->revealable(),

            // confirm_new_password
            Forms\Components\TextInput::make('data.confirm_new_password')
                ->password()
                ->required()
                ->maxLength(255)
                ->label('Confirm New Password')
                ->revealable()
                ->same('data.new_password'),

        ];
    }

    protected function getFormActions(): array
    {
        return [
            // save
            Action::make('save')
                ->label(__('filament-panels::resources/pages/edit-record.form.actions.save.label'))
                ->submit('save'),

            // cancel
            Action::make('cancel')
                ->label(__('filament-panels::resources/pages/edit-record.form.actions.cancel.label'))
                ->url(fn(): string => StudentResource::getUrl('index'))
                ->color('secondary'),
        ];
    }

    public function save(): void
    {
        try {
            $data = $this->form->getState();

            if (!password_verify($data['data']['your_password'], Auth::user()->password)) {
                Notification::make()
                    ->danger()
                    ->title(__('Your password is incorrect'))
                    ->send();
                return;
            }

            $this->record->update([
                'password' => bcrypt($data['data']['new_password']),
            ]);

            Notification::make()
                ->success()
                ->title(__('Password updated successfully'))
                ->send();
        } catch (Halt $exception) {
            return;
        }
    }
}
