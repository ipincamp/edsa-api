<?php

namespace App\Http\Resources\User;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;
use App\Models\Book;
use App\Http\Resources\User\BookSummaryResource;

class ProfileResource extends JsonResource
{
    /**
     * Transform the resource into an array.
     *
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        $this->resource->load(['roles', 'groups.course', 'progresses']);

        $firstGroup = $this->groups->first();
        $totalPoints = $this->progresses->sum('total_points');
        $completedBooksCount = $this->progresses->where('status', 'completed')->count();
        // Overall score dihitung dari total poin dibagi jumlah buku yang sudah diselesaikan
        $overallScore = ($completedBooksCount > 0) ? round($totalPoints / $completedBooksCount) : 0;

        return [
            'id' => $this->id,
            'name' => $this->name,
            'email' => $this->email,
            'role' => $this->roles->first()->name,
            'course_name' => $firstGroup ? $firstGroup->course->name : null,
            'group_name' => $firstGroup ? $firstGroup->name : null,
            'overall_score' => $overallScore,
            'joined_at' => $this->created_at->format('Y-m-d H:i:s'),
            'list_buku' => $this->getBookList(),
        ];
    }

    /**
     * Helper function untuk mengambil dan memformat daftar buku.
     *
     * @return \Illuminate\Http\Resources\Json\AnonymousResourceCollection
     */
    private function getBookList()
    {
        if (!$this->hasRole('student')) {
            return BookSummaryResource::collection(Book::orderBy('order_sequence', 'asc')->get());
        }

        $allBooks = Book::orderBy('order_sequence', 'asc')->get();
        $studentProgresses = $this->progresses->keyBy('book_id');

        $lastCompletedBookOrder = Book::whereIn('id', $studentProgresses->where('status', 'completed')->pluck('book_id'))
            ->max('order_sequence') ?? 0;

        $booksWithExtraData = $allBooks->map(function ($book) use ($studentProgresses, $lastCompletedBookOrder) {
            $progress = $studentProgresses->get($book->id);

            $book->is_locked = $book->order_sequence > ($lastCompletedBookOrder + 1);
            $book->total_points = $progress ? $progress->total_points : 0;
            $book->latest_page = $progress ? $progress->latest_page : 0;

            return $book;
        });

        return BookSummaryResource::collection($booksWithExtraData);
    }
}
