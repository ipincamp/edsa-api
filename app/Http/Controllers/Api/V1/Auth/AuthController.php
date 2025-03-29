<?php

namespace App\Http\Controllers\Api\V1\Auth;

use App\Http\Controllers\Controller;
use App\Http\Requests\Auth\RequestSignIn;
use App\Http\Requests\Auth\RequestSignUp;
use App\Models\User;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;

class AuthController extends Controller
{
    /**
     * Register
     */
    public function signUp(RequestSignUp $req)
    {
        $credentials = $req->safe()->only(['name', 'email', 'password']);
        $credentials['password'] = bcrypt($credentials['password']);
        $user = User::create($credentials);
        $token = $user->createToken('auth_up')->plainTextToken;

        return response()->json([
            'status' => true,
            'message' => 'Registration successful',
            'data' => [
                'user' => $user,
                'token' => $token,
            ],
        ], 201);
    }

    /**
     * Login
     */
    public function signIn(RequestSignIn $req): \Illuminate\Http\JsonResponse
    {
        $credentials = $req->safe()->only(['email', 'password']);

        if (! Auth::attempt($credentials)) {
            return response()->json([
                'status' => false,
                'message' => 'The provided credentials do not match our records.',
                'data' => null,
            ], 401);
        }

        return response()->json([
            'status' => true,
            'message' => 'Login successful',
            'data' => [
                'user' => $req->user(),
                'token' => $req->user()->createToken('auth_in')->plainTextToken,
            ],
        ], 200);
    }

    /**
     * Logout
     */
    public function signOut(Request $req)
    {
        $req->user()->currentAccessToken()->delete();

        return response()->json([
            'status' => true,
            'message' => 'Logout successful',
            'data' => null,
        ], 200);
    }
}
