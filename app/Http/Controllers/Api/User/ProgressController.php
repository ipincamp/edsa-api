<?php

namespace App\Http\Controllers\Api\User;

use App\Enums\BookProgressEnum;
use App\Enums\RolesEnum;
use App\Http\Controllers\Controller;
use App\Http\Requests\Book\Progress\CompleteBookRequest;
use App\Http\Requests\Book\Progress\StartBookRequest;
use App\Http\Requests\Book\Progress\SubmitInteractionRequest;
use App\Http\Requests\Book\Progress\SubmitPostActivityRequest;
use App\Http\Resources\Book\StudentProgressResource;
use App\Models\Interaction;
use App\Models\PostActivity;
use App\Models\StudentProgress;
use App\Models\User;
use Illuminate\Support\Facades\Auth;

class ProgressController extends Controller
{
    // Start or continue a book's progress
    public function startOrContinueBook(StartBookRequest $request)
    {
        try {
            $student = Auth::user();

            // Cari atau buat progres baru untuk siswa dan buku ini
            $progress = StudentProgress::firstOrCreate(
                ['student_id' => $student->id, 'book_id' => $request->book_id],
                ['status' => BookProgressEnum::IN_PROGRESS->value] // Default value jika baru dibuat
            );

            // Jika statusnya 'not_started', ubah menjadi 'in_progress'
            if ($progress->status === BookProgressEnum::NOT_STARTED->value) {
                $progress->update(['status' => BookProgressEnum::IN_PROGRESS->value]);
            }

            return $this->sendSuccess(
                message: 'Book progress started or continued successfully.',
                data: new StudentProgressResource($progress->load('book')),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to start or continue book progress.',
                statusCode: 500
            );
        }
    }

    // Submit interaction answer
    public function submitInteraction(SubmitInteractionRequest $request)
    {
        try {
            if (!$request->is_correct) {
                return $this->sendError(
                    message: 'Answer is incorrect, no points awarded.',
                    statusCode: 400,
                );
            }

            $student = Auth::user();
            $interaction = Interaction::with('page.book')->findOrFail($request->interaction_id);

            $progress = StudentProgress::where('student_id', $student->id)
                ->where('book_id', $interaction->page->book->id)
                ->firstOrFail();

            $progress->increment('total_points', $interaction->points);
            $progress->update(['last_page' => $interaction->page->page_number]);

            return $this->sendSuccess(
                message: 'Point awarded for interaction.',
                data: [
                    'points_awarded' => $interaction->points,
                    'total_points' => $progress->total_points,
                ]
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to submit interaction.',
                statusCode: 500
            );
        }
    }

    // Submit post-activity answer
    public function submitPostActivity(SubmitPostActivityRequest $request)
    {
        try {
            if (!$request->is_correct) {
                return $this->sendError(
                    message: 'Answer is incorrect, no points awarded.',
                    statusCode: 400,
                );
            }

            $student = Auth::user();
            $postActivity = PostActivity::findOrFail($request->post_activity_id);

            $progress = StudentProgress::where('student_id', $student->id)
                ->where('book_id', $postActivity->book_id)
                ->firstOrFail();

            $progress->increment('total_points', $postActivity->points);

            return $this->sendSuccess(
                message: 'Point awarded for post-activity.',
                data: [
                    'points_awarded' => $postActivity->points,
                    'total_points' => $progress->total_points,
                ]
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to submit post-activity.',
                statusCode: 500
            );
        }
    }

    // Mark a book as completed
    public function completeBook(CompleteBookRequest $request)
    {
        try {
            $student = Auth::user();
            $progress = StudentProgress::where('student_id', $student->id)
                ->where('book_id', $request->book_id)
                ->firstOrFail();

            $progress->update(['status' => BookProgressEnum::COMPLETED->value]);

            // Logika untuk unlock buku selanjutnya bisa ditambahkan di sini

            return $this->sendSuccess(
                message: 'Book progress marked as completed successfully.',
                data: new StudentProgressResource($progress->load('book')),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to complete book progress.',
                statusCode: 500
            );
        }
    }

    // View student progress
    public function viewStudentProgress(User $user)
    {
        try {
            // Pastikan user yang diminta adalah student
            if (!$user->hasRole(RolesEnum::S->value)) {
                return $this->sendError(
                    message: 'User is not a student.',
                    statusCode: 404
                );
            }

            $progress = StudentProgress::where('student_id', $user->id)
                ->with('book')
                ->get();

            return $this->sendSuccess(
                message: 'Student progress retrieved successfully.',
                data: StudentProgressResource::collection($progress),
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to retrieve student progress.',
                statusCode: 500
            );
        }
    }
}
