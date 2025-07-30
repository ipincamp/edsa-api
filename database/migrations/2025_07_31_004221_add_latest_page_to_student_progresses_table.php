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
        Schema::table('student_progresses', function (Blueprint $table) {
            $table->integer('latest_page')->after('last_page')->default(0)->comment('Halaman terjauh yang pernah dicapai');
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('student_progresses', function (Blueprint $table) {
            $table->dropColumn('latest_page');
        });
    }
};
