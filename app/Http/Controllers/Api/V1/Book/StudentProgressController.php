<?php

namespace App\Http\Controllers\Api\V1\Book;

use App\Http\Controllers\Controller;
use App\Http\Requests\Book\SaveProgressRequest;
use App\Models\StudentProgress;
use Illuminate\Http\Request;

class StudentProgressController extends Controller
{
    /**
     * Display a listing of the resource.
     */
    public function index()
    {
        //
    }

    /**
     * Store a newly created resource in storage.
     */
    public function store(SaveProgressRequest $request)
    {
        try {
            //code...
        } catch (\Exception $e) {
            throw $e;
        }
    }

    /**
     * Display the specified resource.
     */
    public function show(StudentProgress $studentProgress)
    {
        //
    }

    /**
     * Update the specified resource in storage.
     */
    public function update(Request $request, StudentProgress $studentProgress)
    {
        //
    }

    /**
     * Remove the specified resource from storage.
     */
    public function destroy(StudentProgress $studentProgress)
    {
        //
    }
}
