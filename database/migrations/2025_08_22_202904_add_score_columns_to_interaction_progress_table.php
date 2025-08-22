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
        Schema::table('interaction_progress', function (Blueprint $table) {
            $table->unsignedInteger('correct')->default(0)->after('interaction_id');
            $table->unsignedInteger('total')->default(0)->after('correct');
        });
    }

    /**
     * Reverse the migrations.
     */
    public function down(): void
    {
        Schema::table('interaction_progress', function (Blueprint $table) {
            $table->dropColumn(['correct', 'total']);
        });
    }
};
