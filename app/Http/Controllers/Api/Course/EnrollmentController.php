<?php

namespace App\Http\Controllers\Api\Course;

use App\Enums\EnrollmentStatusEnum;
use App\Http\Controllers\Controller;
use App\Http\Requests\Course\EnrollmentRequest;
use App\Models\Course;
use App\Models\Enrollment;
use Illuminate\Http\Response;
use Illuminate\Support\Facades\Auth;

class EnrollmentController extends Controller
{
    public function enroll(EnrollmentRequest $request)
    {
        try {
            $course = Course::where('enrollment_code', $request->enrollment_code)->firstOrFail();
            $student = Auth::user();

            $isAlreadyInGroup = $student->groups()->where('course_id', $course->id)->exists();
            if ($isAlreadyInGroup) {
                return $this->sendError(
                    message: 'You are already enrolled in a group for this course.',
                    statusCode: Response::HTTP_CONFLICT,
                );
            }

            $hasPendingEnrollment = Enrollment::where('user_id', $student->id)
                ->where('course_id', $course->id)
                ->where('status', EnrollmentStatusEnum::PENDING->value)
                ->exists();
            if ($hasPendingEnrollment) {
                return $this->sendError(
                    message: 'You already have a pending enrollment request for this course.',
                    statusCode: Response::HTTP_CONFLICT,
                );
            }

            Enrollment::updateOrCreate(
                ['user_id' => $student->id, 'course_id' => $course->id],
                ['status' => EnrollmentStatusEnum::PENDING->value] // Jika sudah pernah ditolak, bisa enroll lagi
            );

            return $this->sendSuccess(
                statusCode: Response::HTTP_CREATED,
                message: 'Enrollment request for course "' . $course->name . '" has been sent. Waiting for teacher approval.',
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to enroll in course.',
                statusCode: 500
            );
        }
    }
}
