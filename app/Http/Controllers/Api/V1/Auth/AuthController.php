<?php

namespace App\Http\Controllers\Api\V1\Auth;

use App\Enums\RoleEnum;
use App\Http\Controllers\Controller;
use App\Http\Requests\Auth\LoginRequest;
use App\Http\Requests\Auth\RegisterRequest;
use App\Http\Requests\Auth\UpdatePasswordRequest;
use App\Http\Requests\Auth\UpdateProfileRequest;
use App\Http\Resources\UserResource;
use App\Models\Role;
use Dedoc\Scramble\Attributes\Group;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Hash;

#[Group('Auth')]
class AuthController extends Controller
{
    /**
     * Register
     *
     * Register a new user and return the user data along with an access token.
     *
     * @operationId signUp
     * @unauthenticated
     * @param RegisterRequest $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function register(RegisterRequest $request): JsonResponse
    {
        try {
            $inputs = $request->validated();
            $user = \App\Models\User::create([
                'name' => $inputs['name'],
                'username' => $inputs['username'],
                'password' => bcrypt($inputs['password']),
            ]);

            $user->roles()->attach(
                Role::firstWhere('name', RoleEnum::STUDENT->value)->id
            );

            $token = $user->createToken(
                'auth_up',
                ['*'],
                now()->addDay()
            )->plainTextToken;

            activity('auth api')
                ->performedOn($user)
                ->event('register')
                ->withProperties([
                    'attributes' => [
                        'name' => $user->name,
                        'ip' => $request->ip(),
                        'user_agent' => $request->userAgent(),
                        'device' => $request->header('User-Agent'),
                    ],
                ])
                ->log('Register');

            /* Successfully */
            return $this->json(
                message: 'Register successfully',
                data: [
                    'token' => $token,
                    'user' => new UserResource($user),
                ],
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Login
     *
     * Authenticate a user and return the user data along with an access token.
     *
     * @operationId signIn
     * @unauthenticated
     * @param LoginRequest $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function login(LoginRequest $request): JsonResponse
    {
        try {
            if (! Auth::attempt($request->validated())) {
                /**
                 * Invalid credentials
                 *
                 * @status 401
                 * @body {"status": false, "message": "The provided credentials do not match our records.", "data": null}
                 */
                return $this->json(
                    message: 'The provided credentials do not match our records.',
                    code: 401,
                );
            }

            $token = $request->user()->createToken(
                'auth_in',
                ['*'],
                now()->addDay()
            )->plainTextToken;
            $user = $request->user()->load(['groups.course', 'roles']);

            activity('auth api')
                ->performedOn($request->user())
                ->event('login')
                ->withProperties([
                    'attributes' => [
                        'name' => $request->user()->name,
                        'ip' => $request->ip(),
                        'user_agent' => $request->userAgent(),
                        'device' => $request->header('User-Agent'),
                    ],
                ])
                ->log('Login');

            /* Successfully */
            return $this->json(
                message: 'Login successfully',
                data: [
                    'token' => $token,
                    'user' => new UserResource($user),
                ],
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Profile
     *
     * Get the authenticated user profile.
     *
     * @operationId getProfile
     * @authenticated
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function profile(Request $request): JsonResponse
    {
        try {
            $user = $request->user()->load(['groups.course', 'roles']);

            /* Successfully */
            return $this->json(
                message: 'User profile',
                data: new UserResource($user),
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Update Password
     *
     * Update the authenticated user's password.
     *
     * @operationId updatePassword
     * @authenticated
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function updatePassword(UpdatePasswordRequest $request): JsonResponse
    {
        try {
            $inputs = $request->validated();
            $user = $request()->user();
            // dd($inputs);

            if (! Hash::check($inputs['old_password'], $user->password)) {
                /**
                 * Invalid current password
                 *
                 * @status 401
                 * @body {"status": false, "message": "Your old password is incorrect.", "data": null}
                 */
                return $this->json(
                    message: 'The provided current password is incorrect.',
                    code: 401,
                );
            }

            $user->update([
                'password' => bcrypt($inputs['new_password']),
                'updated_at' => now(),
            ]);
            $user->save();

            activity('auth api')
                ->performedOn($request->user())
                ->event('change password')
                ->withProperties([
                    'attributes' => [
                        'name' => $request->user()->name,
                        'ip' => $request->ip(),
                        'user_agent' => $request->userAgent(),
                        'device' => $request->header('User-Agent'),
                    ],
                ])
                ->log('Change password');

            /* Successfully */
            return $this->json(
                message: 'Password changed successfully',
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Update Profile
     *
     * Update the authenticated user's profile.
     *
     * @operationId updateProfile
     * @authenticated
     * @param UpdateProfileRequest $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function updateProfile(UpdateProfileRequest $request): JsonResponse
    {
        try {
            $inputs = $request->validated();
            $user = $request->user()->load(['groups.course', 'roles']);

            $user->update(array_filter([
                'name' => $inputs['name'] ?? null,
                'username' => $inputs['username'] ?? null,
                'updated_at' => now(),
            ]));
            $user->save();

            activity('auth api')
                ->performedOn($request->user())
                ->event('update profile')
                ->withProperties([
                    'attributes' => [
                        'name' => $request->user()->name,
                        'username' => $request->user()->username,
                        'ip' => $request->ip(),
                        'user_agent' => $request->userAgent(),
                        'device' => $request->header('User-Agent'),
                    ],
                    'old' => [
                        'name' => $request->user()->getOriginal('name'),
                        'username' => $request->user()->getOriginal('username'),
                    ],
                ])
                ->log('Update profile');

            /* Successfully */
            return $this->json(
                message: 'Profile updated successfully',
                data: new UserResource($user),
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Logout
     *
     * Revoke the user's access token.
     *
     * @operationId signOut
     * @authenticated
     * @param Request $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function logout(Request $request): JsonResponse
    {
        try {
            $request->user()->currentAccessToken()->delete();

            activity('auth api')
                ->performedOn($request->user())
                ->event('logout')
                ->withProperties([
                    'attributes' => [
                        'name' => $request->user()->name,
                        'ip' => $request->ip(),
                        'user_agent' => $request->userAgent(),
                        'device' => $request->header('User-Agent'),
                    ],
                ])
                ->log('Logout');

            /* Successfully */
            return $this->json(
                message: 'Logout successfully',
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }
}
