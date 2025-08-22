<?php

namespace App\Http\Controllers\Api\User;

use App\Enums\BookProgressEnum;
use App\Enums\RolesEnum;
use App\Http\Controllers\Controller;
use App\Http\Requests\Book\Progress\CompleteBookRequest;
use App\Http\Requests\Book\Progress\StartBookRequest;
use App\Http\Requests\Book\Progress\SubmitInteractionRequest;
use App\Http\Requests\Book\Progress\SubmitPostActivityRequest;
use App\Http\Requests\Book\Progress\UpdateLastPageRequest;
use App\Http\Resources\Book\StudentProgressResource;
use App\Models\Interaction;
use App\Models\PostActivity;
use App\Models\StudentProgress;
use App\Models\User;
use Illuminate\Http\Response;
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

            // Cek jika progres ini baru saja dibuat
            if ($progress->wasRecentlyCreated) {
                $progress->update([
                    'status' => BookProgressEnum::IN_PROGRESS->value,
                    'last_page' => 1,
                    'latest_page' => 1,
                ]);
            } elseif ($progress->status === BookProgressEnum::NOT_STARTED->value) {
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

    // Update the last page read by the student
    public function updateLastPage(UpdateLastPageRequest $request)
    {
        try {
            $student = Auth::user();

            $progress = StudentProgress::where('student_id', $student->id)
                ->where('book_id', $request->book_id)
                ->firstOrFail();

            $progress->last_page = $request->page_number;

            if ($request->page_number > $progress->latest_page) {
                $progress->latest_page = $request->page_number;
            }

            $progress->save();

            return $this->sendSuccess(
                message: 'Last page updated successfully.',
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to update last page.',
                statusCode: 500,
            );
        }
    }

    /**
     * Submit interaction answer and save temporary score.
     * Points are only awarded once upon achieving a perfect score.
     *
     * @param  \App\Http\Requests\Book\Progress\SubmitInteractionRequest  $request
     * @return \Illuminate\Http\JsonResponse
     */
    public function submitInteraction(SubmitInteractionRequest $request)
    {
        try {
            $student = Auth::user();
            $interaction = Interaction::with('page.book')->findOrFail($request->interaction_id);

            $progress = StudentProgress::where('student_id', $student->id)
                ->where('book_id', $interaction->page->book->id)
                ->firstOrFail();

            // Cek apakah interaksi ini sudah pernah mendapatkan skor sempurna sebelumnya
            $existingAttempt = $progress->completedInteractions()
                ->where('interaction_id', $interaction->id)
                ->first();

            if ($existingAttempt && $existingAttempt->pivot->correct === $existingAttempt->pivot->total) {
                return $this->sendError(
                    message: 'This interaction has already been completed with a perfect score.',
                    statusCode: 409, // Conflict
                );
            }

            // Simpan atau update skor sementara (terbaru)
            // syncWithoutDetaching akan membuat record baru jika belum ada, atau update jika sudah ada
            $progress->completedInteractions()->syncWithoutDetaching([
                $interaction->id => [
                    'correct' => $request->correct,
                    'total' => $request->total,
                ]
            ]);

            $points_awarded = 0;
            $isPerfectScore = (int) $request->correct === (int) $request->total;

            // Berikan poin HANYA JIKA skor saat ini sempurna DAN ini adalah pertama kalinya sempurna
            if ($isPerfectScore) {
                $points_awarded = $interaction->points;
                $progress->increment('total_points', $points_awarded);
            }

            // Update halaman terakhir yang diakses
            $pageNumber = $interaction->page->page_number;
            $progress->last_page = $pageNumber;
            if ($pageNumber > $progress->latest_page) {
                $progress->latest_page = $pageNumber;
            }
            $progress->save();

            return $this->sendSuccess(
                message: 'Interaction score has been saved.',
                data: [
                    'points_awarded' => $points_awarded,
                    'total_points' => $progress->total_points,
                    'is_perfect_score' => $isPerfectScore,
                ]
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to submit interaction: ' . $e->getMessage(),
                statusCode: 500
            );
        }
    }

    // Submit post-activity answer
    public function submitPostActivity(SubmitPostActivityRequest $request)
    {
        try {
            $student = Auth::user();
            $postActivity = PostActivity::findOrFail($request->post_activity_id);

            $progress = StudentProgress::where('student_id', $student->id)
                ->where('book_id', $postActivity->book_id)
                ->firstOrFail();

            if ($progress->completedPostActivities->contains($postActivity->id)) {
                return $this->sendError(
                    message: 'This post-activity has already been completed.',
                    statusCode: Response::HTTP_CONFLICT,
                );
            }

            if (!$request->is_correct) {
                return $this->sendError(
                    message: 'Answer is incorrect, no points awarded.',
                    statusCode: Response::HTTP_BAD_REQUEST,
                );
            }

            $progress->completedPostActivities()->attach($postActivity->id);

            $progress->increment('total_points', $postActivity->points);

            return $this->sendSuccess(
                message: 'Point awarded for post-activity.',
                data: [
                    'points_awarded' => $postActivity->points,
                    'total_points' => $progress->total_points,
                    'passed' => true,
                ]
            );
        } catch (\Exception $e) {
            return $this->sendError(
                message: 'Failed to submit post-activity: ' . $e->getMessage(),
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
