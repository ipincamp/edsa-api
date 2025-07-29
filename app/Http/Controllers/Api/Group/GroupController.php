<?php

namespace App\Http\Controllers\Api\Group;

use App\Http\Controllers\Controller;
use App\Http\Requests\Group\AssignStudentRequest;
use App\Http\Requests\Group\StoreGroupRequest;
use App\Http\Requests\Group\UpdateGroupRequest;
use App\Http\Resources\Auth\AuthResource;
use App\Http\Resources\Group\GroupResource;
use App\Models\Course;
use App\Models\Group;
use App\Models\User;
use Illuminate\Http\Response;

use Illuminate\Foundation\Auth\Access\AuthorizesRequests;

class GroupController extends Controller
{
    use AuthorizesRequests;

    // Get the list of groups for a specific course
    public function index(Course $course)
    {
        try {
            return $this->sendSuccess(
                message: 'Groups retrieved successfully.',
                data: GroupResource::collection(
                    $course->groups()->with('teachers')->withCount('students')->get(),
                ),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to retrieve groups.',
                statusCode: 500
            );
        }
    }

    // Create a new group for a course
    public function store(StoreGroupRequest $request)
    {
        try {
            $group = Group::create($request->safe()->only(['name', 'course_id']));
            $group->teachers()->attach($request->teacher_ids);

            return $this->sendSuccess(
                statusCode: Response::HTTP_CREATED,
                message: 'Group created successfully.',
                data: new GroupResource($group->load('teachers')),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to create group.',
                statusCode: 500
            );
        }
    }

    // Show a specific group with its related data
    public function show(Group $group)
    {
        try {
            $this->authorize('view', $group);
            $group->load('course', 'teacher', 'students');

            return $this->sendSuccess(
                message: 'Group retrieved successfully.',
                data: new GroupResource($group)
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to retrieve group.',
                statusCode: 500
            );
        }
    }

    // Update an existing group
    public function update(UpdateGroupRequest $request, Group $group)
    {
        try {
            $this->authorize('update', $group);
            $group->update($request->validated());

            return $this->sendSuccess(
                message: 'Group updated successfully.',
                data: new GroupResource($group->load('teacher'))
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to update group.',
                statusCode: 500
            );
        }
    }

    // Delete a group, ensuring it has no students before deletion
    public function destroy(Group $group)
    {
        try {
            $this->authorize('delete', $group);
            if ($group->students()->exists()) {
                return $this->sendError(
                    message: 'Group cannot be deleted because it still has students.',
                    statusCode: Response::HTTP_CONFLICT
                );
            }
            $group->delete();

            return $this->sendSuccess(
                statusCode: Response::HTTP_NO_CONTENT,
                message: 'Group deleted successfully.',
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to delete group.',
                statusCode: 500
            );
        }
    }

    // Assign students to a group
    public function assignStudent(AssignStudentRequest $request, Group $group)
    {
        try {
            $this->authorize('manageStudents', $group);
            $student = User::find($request->student_id);
            // Cek jika siswa sudah ada di grup lain dalam course yang sama
            $group->students()->syncWithoutDetaching($student->id);

            return $this->sendSuccess(
                message: 'Student assigned to group successfully.',
                data: new AuthResource($student)
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to assign student to group.',
                statusCode: 500
            );
        }
    }

    // Remove a student from a group
    public function removeStudent(AssignStudentRequest $request, Group $group)
    {
        try {
            $this->authorize('manageStudents', $group);
            $group->students()->detach($request->student_id);

            return $this->sendSuccess(
                message: 'Student removed from group successfully.',
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to remove student from group.',
                statusCode: 500
            );
        }
    }
}
