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
        Schema::table('post_activities', function (Blueprint $table) {
            $table->dropForeign(['book_id']);
            $table->dropUnique(['book_id']);
            $table->dropColumn(['image_url', 'scrambled_word', 'correct_answer']);

            $table->string('type')->after('book_id'); // Cth: 'match_the_picture', 'trace_the_word'
            $table->json('data')->after('type')->nullable(); // Untuk menyimpan data spesifik aktivitas
            $table->integer('order')->after('data')->default(1); // Urutan aktivitas
            $table->integer('points')->after('order')->default(10); // Poin untuk aktivitas

            $table->foreign('book_id')->references('id')->on('books')->onDelete('cascade');
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('post_activities', function (Blueprint $table) {
            $table->dropForeign(['book_id']);

            $table->string('image_url')->nullable();
            $table->string('scrambled_word');
            $table->string('correct_answer');

            $table->dropColumn(['type', 'data', 'order', 'points']);

            $table->unique('book_id');
            $table->foreign('book_id')->references('id')->on('books')->onDelete('cascade');
        });
    }
};
