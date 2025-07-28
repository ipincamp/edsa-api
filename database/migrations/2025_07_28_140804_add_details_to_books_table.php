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
        Schema::table('books', function (Blueprint $table) {
            $table->string('author')->after('title')->default('EDSA Team');
            $table->year('year')->after('author')->default(date('Y'));
            $table->string('genre')->after('year')->nullable();
            $table->string('focus')->after('genre')->comment('Contoh: alphabet, number, color');
            $table->enum('status', ['published', 'draft', 'locked'])->default('published')->after('focus');
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('books', function (Blueprint $table) {
            $table->dropColumn(['author', 'year', 'genre', 'focus', 'status']);
        });
    }
};
