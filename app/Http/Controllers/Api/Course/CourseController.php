<?php

namespace App\Http\Controllers\Api\Course;

use App\Http\Controllers\Controller;
use App\Http\Requests\Course\StoreCourseRequest;
use App\Http\Requests\Course\UpdateCourseRequest;
use App\Http\Resources\Course\CourseResource;
use App\Models\Course;
use Illuminate\Http\Response;

class CourseController extends Controller
{
    // Get a list of all courses
    public function index()
    {
        try {
            $courses = Course::withCount('groups')->paginate(15);

            return $this->sendSuccess(
                message: 'Successfully retrieved courses.',
                data: CourseResource::collection($courses),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to retrieve courses.',
                statusCode: 500
            );
        }
    }

    // Store a new course
    public function store(StoreCourseRequest $request)
    {
        try {
            $course = Course::create($request->validated());

            return $this->sendSuccess(
                message: 'Course created successfully.',
                data: new CourseResource($course),
                statusCode: Response::HTTP_CREATED
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to create course.',
                statusCode: 500
            );
        }
    }

    // Show a specific course
    public function show(Course $course)
    {
        try {
            $course->load('groups');

            return $this->sendSuccess(
                message: 'Successfully retrieved course.',
                data: new CourseResource($course)
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to create course.',
                statusCode: 500
            );
        }
    }

    // Update an existing course
    public function update(UpdateCourseRequest $request, Course $course)
    {
        try {
            $course->update($request->validated());

            return $this->sendSuccess(
                message: 'Course updated successfully.',
                data: new CourseResource($course)
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to update course.',
                statusCode: 500
            );
        }
    }

    // Delete a specific course
    public function destroy(Course $course)
    {
        try {
            if ($course->groups()->exists()) {
                return $this->sendError(
                    message: 'Tidak dapat menghapus course karena masih memiliki grup.',
                    statusCode: Response::HTTP_CONFLICT
                );
            }

            $course->delete();

            return $this->sendSuccess(
                message: 'Course deleted successfully.',
                statusCode: Response::HTTP_NO_CONTENT
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to delete course.',
                statusCode: 500
            );
        }
    }
}
