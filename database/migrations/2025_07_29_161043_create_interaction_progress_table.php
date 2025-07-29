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
        Schema::create('interaction_progress', function (Blueprint $table) {
            $table->primary(['student_progress_id', 'interaction_id']);
            $table->foreignId('student_progress_id')->constrained('student_progresses')->onDelete('cascade');
            $table->foreignId('interaction_id')->constrained('interactions')->onDelete('cascade');
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('interaction_progress');
    }
};
