<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;

class BookSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        // book 1
        $book1 = \App\Models\Book::create([
            'title' => 'The Alphabet in the Land of Dewi Sri',
            'author' => 'Miss Ani',
            'year' => 2025,
            'genre' => 'Fiction',
            'focus' => 'Alphabet',
            'image' => '/image/book/cover1.png',
        ]);
        $book1->images()->create([
            'path' => '/image/book/word/andong.png',
            'name' => 'ANDONG',
        ]);
        $book1->images()->create([
            'path' => '/image/book/word/candi.png',
            'name' => 'CANDI',
        ]);
        $book1->images()->create([
            'path' => '/image/book/word/dewa.png',
            'name' => 'DEWA',
        ]);
        $book1->images()->create([
            'path' => '/image/book/word/gong.png',
            'name' => 'GONG',
        ]);
        $book1->images()->create([
            'path' => '/image/book/word/ikan.png',
            'name' => 'IKAN',
        ]);
        $book1->images()->create([
            'path' => '/image/book/word/mangga.png',
            'name' => 'MANGGA',
        ]);

        // book 2
        \App\Models\Book::create([
            'title' => 'The Numbers in the Village of Ten Hills',
            'author' => 'Miss Ani',
            'year' => 2025,
            'genre' => 'Fiction',
            'focus' => 'Number',
            'image' => '/image/book/cover2.png',
        ]);

        // book 3
        \App\Models\Book::create([
            'title' => 'Bawang Putih and the Kind Body Parts',
            'author' => 'Miss Ani',
            'year' => 2025,
            'genre' => 'Fiction',
            'focus' => 'Body Parts',
            'image' => '/image/book/cover3.png',
        ]);

        // book 4
        \App\Models\Book::create([
            'title' => 'The Big Mango Tree and the Family of Five',
            'author' => 'Miss Ani',
            'year' => 2025,
            'genre' => 'Fiction',
            'focus' => 'Family',
            'image' => '/image/book/cover4.png',
        ]);

        // book 5
        \App\Models\Book::create([
            'title' => 'Moby Dick',
            'author' => 'Herman Melville',
            'year' => 1851,
            'genre' => 'Adventure',
            'focus' => 'Alphabet',
            'status' => false,
            'image' => '/image/book/default.jpg',
        ]);

        // book 6
        \App\Models\Book::create([
            'title' => 'War and Peace',
            'author' => 'Leo Tolstoy',
            'year' => 1869,
            'genre' => 'Historical',
            'focus' => 'Alphabet',
            'status' => false,
            'image' => '/image/book/default.jpg',
        ]);

        // book 7
        \App\Models\Book::create([
            'title' => 'Age of Empire',
            'author' => 'Mustoleh',
            'year' => 1922,
            'genre' => 'Comedy',
            'focus' => 'Number',
            'status' => false,
            'image' => '/image/book/default.jpg',
        ]);

        // book 8
        \App\Models\Book::create([
            'title' => 'The Great Gatsby',
            'author' => 'Timothy Ronald',
            'year' => 2025,
            'genre' => 'Tragedy',
            'focus' => 'Alphabet',
            'status' => false,
            'image' => '/image/book/default.jpg',
        ]);
    }
}
