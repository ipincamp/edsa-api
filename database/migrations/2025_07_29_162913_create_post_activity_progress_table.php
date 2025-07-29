<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     */
    public function up(): void
    {
        Schema::create('post_activity_progress', function (Blueprint $table) {
            $table->primary(['student_progress_id', 'post_activity_id']);
            $table->foreignId('student_progress_id')->constrained('student_progresses')->onDelete('cascade');
            $table->foreignId('post_activity_id')->constrained('post_activities')->onDelete('cascade');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('post_activity_progress');
    }
};
