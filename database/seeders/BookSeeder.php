<?php

namespace Database\Seeders;

use Illuminate\Database\Console\Seeds\WithoutModelEvents;
use Illuminate\Database\Seeder;

class BookSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        // book 1
        \App\Models\Book::create([
            'title' => 'The Alphabet in the Land of Dewi Sri',
            'author' => 'Miss Ani',
            'year' => 2025,
            'genre' => 'Fiction',
            'focus' => 'Alphabet',
            'image' => '/book_cover_1.png',
        ]);

        // book 2
        \App\Models\Book::create([
            'title' => 'The Numbers in the Village of Ten Hills',
            'author' => 'Miss Ani',
            'year' => 2025,
            'genre' => 'Fiction',
            'focus' => 'Number',
            'image' => '/book_cover_2.png',
        ]);

        // book 3
        \App\Models\Book::create([
            'title' => '1984',
            'author' => 'George Orwell',
            'year' => 1949,
            'genre' => 'Dystopian',
            'focus' => 'Alphabet',
            'image' => '/default.jpg',
        ]);

        // book 4
        \App\Models\Book::create([
            'title' => 'Pride and Prejudice',
            'author' => 'Jane Austen',
            'year' => 1813,
            'genre' => 'Romance',
            'focus' => 'Alphabet',
            'status' => false,
            'image' => '/default.jpg',
        ]);

        // book 5
        \App\Models\Book::create([
            'title' => 'Moby Dick',
            'author' => 'Herman Melville',
            'year' => 1851,
            'genre' => 'Adventure',
            'focus' => 'Alphabet',
            'status' => false,
            'image' => '/default.jpg',
        ]);

        // book 6
        \App\Models\Book::create([
            'title' => 'War and Peace',
            'author' => 'Leo Tolstoy',
            'year' => 1869,
            'genre' => 'Historical',
            'focus' => 'Alphabet',
            'status' => false,
            'image' => '/default.jpg',
        ]);

        // book 7
        \App\Models\Book::create([
            'title' => 'Age of Empire',
            'author' => 'Mustoleh',
            'year' => 1922,
            'genre' => 'Comedy',
            'focus' => 'Number',
            'status' => false,
            'image' => '/default.jpg',
        ]);

        // book 8
        \App\Models\Book::create([
            'title' => 'The Great Gatsby',
            'author' => 'Timothy Ronald',
            'year' => 2025,
            'genre' => 'Tragedy',
            'focus' => 'Alphabet',
            'status' => false,
            'image' => '/default.jpg',
        ]);
    }
}
