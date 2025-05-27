<?php

namespace App\Http\Resources;

use App\Models\StudentProgress;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/**
 * @mixin StudentProgress
 */
class StudentProgressResource extends JsonResource
{
    /**
     * Transform the resource into an array.
     *
     * @return array<string, mixed>
     */
    public function toArray(Request $request): array
    {
        return [
            /*
             * The unique identifier for the student progress.
             * @var int
             */
            'id' => $this->id,
            'book' => [
                /*
                 * The unique identifier for the book.
                 * @var int|null
                 */
                'id' => $this->book?->id,
                /*
                 * The title of the book.
                 * @var string|null
                 */
                'title' => $this->book?->title,
                /*
                 * The cover image URL of the book.
                 * @var string|null
                 */
                'image' => $this->book ? config('app.url') . '/assets' . $this->book->image : null,
            ],
            /*
             * The course name associated with the student progress.
             * @var string|null
             */
            'course' => $this->group?->course?->name,
            /*
             * The group name associated with the student progress.
             * @var string|null
             */
            'class' => $this->group?->name,
            /*
             * How many times the student has taken the book.
             * @var int
            */
            'taken' => $this->taken,
            /*
             * The correct points scored in the book.
             * @var int
             */
            'correct' => $this->score_correct,
            /*
             * The incorrect points scored in the book.
             * @var int
             */
            'incorrect' => $this->score_incorrect,
            /*
             * When the student started the book.
             * @var string
             * @format date-time
             * @example 2025-05-19 14:20:00
             */
            'started_at' => $this->time_start,
            /*
             * When the student finished the book.
             * @var string
             * @format date-time
             * @example 2025-05-19 14:20:00
             */
            'finished_at' => $this->time_finish,
            /*
             * When the student progress was first created.
             * @var string
             * @format date-time
             * @example 2025-05-19 14:20:00
             */
            'created_at' => $this->created_at->format('Y-m-d H:i:s'),
            /*
             * When the student progress was last updated.
             * @var string
             * @format date-time
             * @example 2025-05-19 14:20:00
             */
            'updated_at' => $this->updated_at->format('Y-m-d H:i:s'),

        ];
    }
}
