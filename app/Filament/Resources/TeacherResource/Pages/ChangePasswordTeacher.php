<?php

namespace App\Filament\Resources\TeacherResource\Pages;

use App\Filament\Resources\TeacherResource;
use Filament\Actions\Action;
use Filament\Forms\Concerns\InteractsWithForms;
use Filament\Notifications\Notification;
use Filament\Resources\Pages\Concerns\InteractsWithRecord;
use Filament\Resources\Pages\Page;
use Filament\Support\Exceptions\Halt;
use Illuminate\Support\Facades\Auth;

class ChangePasswordTeacher extends Page
{
    use InteractsWithForms, InteractsWithRecord;

    protected static string $resource = TeacherResource::class;

    protected static string $view = 'filament.resources.teacher-resource.pages.change-password-teacher';

    public ?array $data = [];

    public function mount(int | string $record): void
    {
        $this->record = TeacherResource::getModel()::findOrFail($record);

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
            \Filament\Forms\Components\TextInput::make('data.your_password')
                ->password()
                ->required()
                ->maxLength(255)
                ->label('Your Password')
                ->revealable(),

            // new_password
            \Filament\Forms\Components\TextInput::make('data.new_password')
                ->password()
                ->required()
                ->maxLength(255)
                ->label('New Password')
                ->revealable(),

            // confirm_new_password
            \Filament\Forms\Components\TextInput::make('data.confirm_new_password')
                ->password()
                ->required()
                ->maxLength(255)
                ->label('Confirm New Password')
                ->revealable(),
        ];
    }

    protected function getFormAction(): array
    {
        return [
            // save
            Action::make('save')
                ->label(__('filament-panels::resources/pages/edit-record.form.actions.save.label'))
                ->submit('save'),

            // cancel
            Action::make('cancel')
                ->label(__('filament-panels::resources/pages/edit-record.form.actions.cancel.label'))
                ->url(fn(): string => TeacherResource::getUrl('index'))
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

            $this->redirect($this->getRedirectUrl());
        } catch (Halt $exception) {
            return;
        }
    }

    protected function getRedirectUrl(): string
    {
        return $this->getResource()::getUrl('index');
    }
}
