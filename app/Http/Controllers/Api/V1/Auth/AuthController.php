<?php

namespace App\Http\Controllers\Api\V1\Auth;

use App\Http\Controllers\Controller;
use App\Http\Requests\Auth\RequestSignIn;
use App\Http\Requests\Auth\RequestSignUp;
use App\Models\User;
use Dedoc\Scramble\Attributes\Group;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;

#[Group('Auth')]
class AuthController extends Controller
{
    /**
     * Register
     *
     * Create a new user and return the user data along with an access token.
     *
     * @operationId signUp
     * @unauthenticated
     * @param RequestSignUp $req
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function signUp(RequestSignUp $req): \Illuminate\Http\JsonResponse
    {
        try {
            $credentials = $req->safe()->only(['name', 'email', 'password']);
            $credentials['password'] = bcrypt($credentials['password']);
            $user = User::create($credentials);
            $token = $user->createToken('auth_up')->plainTextToken;

            /* Successfully */
            return response()->json([
                'status' => true,
                'message' => 'Registration successful',
                'data' => [
                    'user' => $user,
                    'token' => $token,
                ],
            ], 201);
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
     * @param RequestSignIn $req
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function signIn(RequestSignIn $req): \Illuminate\Http\JsonResponse
    {
        try {
            $credentials = $req->safe()->only(['email', 'password']);

            if (! Auth::attempt($credentials)) {
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

            /* Successfully */
            return response()->json([
                'status' => true,
                'message' => 'Login successful',
                'data' => [
                    'user' => $req->user(),
                    'token' => $req->user()->createToken('auth_in')->plainTextToken,
                ],
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
     * @param Request $req
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function signOut(Request $req)
    {
        try {
            $req->user()->currentAccessToken()->delete();

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
