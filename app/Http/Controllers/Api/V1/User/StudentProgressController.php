<?php

namespace App\Http\Controllers\Api\V1\User;

use App\Http\Controllers\Controller;
use App\Http\Resources\StudentProgressResource;
use App\Models\User;
use Dedoc\Scramble\Attributes\Group;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;

#[Group('User Progress')]
class StudentProgressController extends Controller
{
    /**
     * Progress
     *
     * Get all student progress data.
     *
     * @operationId getAllStudentProgress
     * @authenticated
     *
     * @return \Illuminate\Http\JsonResponse
     */
    public function index(): JsonResponse
    {
        try {
            $progress = auth()->user()
                ->progress()
                ->with([
                    'group',
                    'group.course',
                    'book',
                ])
                ->get();

            return $this->json(
                message: 'All student progress data retrieved successfully.',
                data: StudentProgressResource::collection($progress),
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Store a newly created resource in storage.
     */
    public function store(Request $request)
    {
        //
    }

    /**
     * Progress of a specific student
     *
     * Get the progress of a specific student by their ID.
     * This method retrieves the progress data for the specified user.
     *
     * @operationId getStudentProgress
     * @authenticated
     * @param User $user
     * @return \Illuminate\Http\JsonResponse
     */
    public function show(User $user): JsonResponse
    {
        try {
            $progress = $user->progress()->get();

            return $this->json(
                message: 'Student progress data retrieved successfully.',
                data: $progress,
            );
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Update the specified resource in storage.
     */
    public function update(Request $request, string $id)
    {
        //
    }

    /**
     * Remove the specified resource from storage.
     */
    public function destroy(string $id)
    {
        //
    }
}
