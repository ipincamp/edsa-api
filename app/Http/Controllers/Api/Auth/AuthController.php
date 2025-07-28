<?php

namespace App\Http\Controllers\Api\Auth;

use App\Enums\RolesEnum;
use App\Http\Controllers\Controller;
use App\Http\Requests\Auth\LoginRequest;
use App\Http\Requests\Auth\RegisterRequest;
use App\Http\Requests\Auth\UpdateDetailRequest;
use App\Http\Requests\Auth\UpdatePasswordRequest;
use App\Http\Resources\Auth\AuthResource;
use App\Models\User;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Hash;

class AuthController extends Controller
{
    // register
    public function register(RegisterRequest $request)
    {
        try {
            $user = User::create([
                'name' => $request->name,
                'email' => $request->email,
                'password' => Hash::make($request->password),
            ]);
            $user->assignRole(RolesEnum::S->value);
            $token = $user->createToken(
                'auth_up',
                ['*'],
                now()->addDay()
            )->plainTextToken;

            return $this->sendSuccess(
                statusCode: 201,
                message: 'User successfully registered',
                data: new AuthResource($user, $token),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Registration failed: ' . $e->getMessage(),
                statusCode: 500
            );
        }
    }

    // login
    public function login(LoginRequest $request)
    {
        try {
            if (!Auth::attempt($request->only('email', 'password'))) {
                return $this->sendError(
                    message: 'The provided credentials does not match our records.',
                    statusCode: 401
                );
            }

            $user = User::where('email', $request->email)->firstOrFail();
            $token = $user->createToken(
                'auth_in',
                ['*'],
                now()->addDay()
            )->plainTextToken;

            return $this->sendSuccess(
                message: 'Login successful',
                data: new AuthResource($user, $token),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Login failed: ' . $e->getMessage(),
                statusCode: 500
            );
        }
    }

    // profile
    public function profile(Request $request)
    {
        try {
            return $this->sendSuccess(
                message: 'User profile retrieved successfully',
                data: new AuthResource($request->user()),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to retrieve profile: ' . $e->getMessage(),
                statusCode: 500
            );
        }
    }

    // update detail
    public function updateDetail(UpdateDetailRequest $request)
    {
        try {
            $user = $request->user();
            $user->update($request->validated());

            return $this->sendSuccess(
                message: 'Profile updated successfully',
                data: new AuthResource($user),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Profile update failed: ' . $e->getMessage(),
                statusCode: 500
            );
        }
    }

    // update password
    public function updatePassword(UpdatePasswordRequest $request)
    {
        try {
            if (!Hash::check($request->old_password, $request->user()->password)) {
                return $this->sendError(
                    message: 'The provided old password is incorrect.',
                    statusCode: 401
                );
            }

            $request->user()->update([
                'password' => Hash::make($request->new_password),
            ]);

            return $this->sendSuccess(
                message: 'Password updated successfully',
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Password update failed: ' . $e->getMessage(),
                statusCode: 500
            );
        }
    }

    // logout
    public function logout(Request $request)
    {
        try {
            auth()->guard('web')->logout();
            $request->user()->currentAccessToken()->delete();

            return $this->sendSuccess(
                message: 'Logout successful',
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Logout failed: ' . $e->getMessage(),
                statusCode: 500
            );
        }
    }
}
