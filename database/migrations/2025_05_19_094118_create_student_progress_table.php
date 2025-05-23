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
        Schema::create('student_progress', function (Blueprint $table) {
            $table->id();
            $table->foreignUuid('student_id')->references('id')->on('users')->onDelete('cascade');
            $table->foreignId('book_id')->references('id')->on('books')->onDelete('cascade');
            $table->foreignId('group_id')->nullable()->references('id')->on('groups')->onDelete('cascade');
            $table->integer('score_correct')->default(0);
            $table->integer('score_incorrect')->default(0);
            $table->timestamp('time_start')->nullable();
            $table->timestamp('time_finish')->nullable();
            $table->integer('taken')->default(0);
            $table->timestamps();
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::dropIfExists('student_progress');
    }
};
