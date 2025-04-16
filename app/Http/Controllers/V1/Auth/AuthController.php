<?php

namespace App\Http\Controllers\V1\Auth;

use App\Http\Controllers\Controller;
use App\Http\Requests\Auth\ChangePasswordRequest;
use App\Http\Requests\Auth\LoginRequest;
use App\Http\Requests\Auth\UpdateProfileRequest;
use Dedoc\Scramble\Attributes\Group;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Hash;

#[Group('Auth')]
class AuthController extends Controller
{
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
    public function login(LoginRequest $request)
    {
        try {
            $i = $request->safe()->only(['username', 'password']);
            $iEmail = ['email' => $i['username'], 'password' => $i['password']];
            $iUsername = ['username' => $i['username'], 'password' => $i['password']];

            if (! Auth::attempt($iEmail) && ! Auth::attempt($iUsername)) {
                /**
                 * Invalid credentials
                 *
                 * @status 401
                 * @body {"status": false, "message": "The provided credentials do not match our records.", "data": null}
                 */
                return response()->json([
                    'status' => false,
                    'message' => 'The provided credentials do not match our records.',
                    'data' => null,
                ], 401);
            }

            $token = $request->user()->createToken('auth_in', ['*'], now()->addDay())->plainTextToken;
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
                ->log('User logged in');

            /** Expired token issued for 1 day since creation */
            return response()->json([
                'status' => true,
                'message' => 'Login successful',
                'data' => [
                    'user' => $request->user()->only(['id', 'name', 'created_at', 'updated_at']),
                    'role' => $request->user()->getRoleNames()->first(),
                    'token' => $token,
                ],
            ], 200);
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Change Password
     *
     * Change the password of the authenticated user.
     *
     * @operationId changePassword
     * @authenticated
     * @param ChangePasswordRequest $request
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function changePassword(ChangePasswordRequest $request)
    {
        try {
            $user = $request->user();
            $password = $request->safe()->only(['current_password']);
            $newPassword = $request->safe()->only(['new_password']);

            if (! Hash::check($password['current_password'], $user->password)) {
                /**
                 * Invalid current password
                 *
                 * @status 401
                 * @body {"status": false, "message": "The provided current password is incorrect.", "data": null}
                 */
                return response()->json([
                    'status' => false,
                    'message' => 'The provided current password is incorrect.',
                    'data' => null,
                ], 401);
            }

            $user->update([
                'password' => Hash::make($newPassword['new_password']),
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
                ->log('User changed password');

            /**
             * Password changed successfully
             *
             * @status 200
             * @body {"status": true, "message": "Password changed successfully", "data": null}
             */
            return response()->json([
                'status' => true,
                'message' => 'Password changed successfully',
                'data' => null,
            ], 200);
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
    public function updateProfile(UpdateProfileRequest $request)
    {
        try {
            $user = $request->user();
            $data = $request->safe()->only(['name', 'username']);

            $user->update($data);
            $user->save();

            activity('auth api')
                ->performedOn($request->user())
                ->event('update profile')
                ->withProperties([
                    'attributes' => [
                        'name' => $request->user()->name,
                        'username' => $request->user()->username,
                        'email' => $request->user()->email,
                        'ip' => $request->ip(),
                        'user_agent' => $request->userAgent(),
                        'device' => $request->header('User-Agent'),
                    ],
                    'old' => [
                        'name' => $user->getOriginal('name'),
                        'username' => $user->getOriginal('username'),
                        'email' => $user->getOriginal('email'),
                    ],
                ])
                ->log('User updated profile');

            /* Successfully */
            return response()->json([
                'status' => true,
                'message' => 'Profile updated successfully',
                'data' => null,
            ], 200);
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
    public function logout(Request $request)
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
                ->log('User logged out');

            /* Successfully */
            return response()->json([
                'status' => true,
                'message' => 'Logout successful',
                'data' => null,
            ], 200);
        } catch (\Exception $e) {
            throw $e;
        }
    }
}
