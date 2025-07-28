<?php

namespace Database\Seeders;

use Illuminate\Database\Seeder;
use App\Models\Book;
use App\Models\Page;
use App\Models\Interaction;
use App\Models\PostActivity;

class BookSeeder extends Seeder
{
    /**
     * Run the database seeds.
     */
    public function run(): void
    {
        Book::query()->delete(); // Ini akan menghapus semua buku, halaman, interaksi terkait.

        // === BUKU 1: The Alphabet in the Land of Dewi Sri ===
        $book = Book::create([
            'title' => 'The Alphabet in the Land of Dewi Sri',
            'author' => 'Miss Ani',
            'year' => '2025',
            'genre' => 'Fiction',
            'focus' => 'Alphabet',
            'status' => 'published',
            'cover_image' => 'https://placehold.co/400x600/81C784/FFFFFF?text=Dewi+Sri',
            'order_sequence' => 1,
        ]);

        // Halaman 1
        $page1 = Page::create(['book_id' => $book->id, 'page_number' => 1, 'content' => 'Once upon a time, in a beautiful rice-growing village in Indonesia, there was a kind and gentle goddess named Dewi Sri. She was loved by everyone because she taught people how to plant rice and take care of the land.']);
        Interaction::create([
            'page_id' => $page1->id,
            'type' => 'drag_and_drop_vocab',
            'data' => [
                'instruction' => 'Drag the words into the right picture!',
                'words' => ['Dewi Sri', 'mountain', 'rice plant', 'farmer', 'hut', 'buffalo']
            ],
            'points' => 10,
        ]);

        // Halaman 2
        $page2 = Page::create(['book_id' => $book->id, 'page_number' => 2, 'content' => 'One sunny morning, as Dewi Sri walked through the golden fields, she noticed the children were curious about how to read and write.']);
        Interaction::create([
            'page_id' => $page2->id,
            'type' => 'tap_and_count',
            'data' => [
                'instruction' => 'How many children are there around the rice field?',
                'correct_count' => 5 // Asumsi ada 5 gambar anak
            ],
            'points' => 5,
        ]);

        // Halaman 3
        Page::create(['book_id' => $book->id, 'page_number' => 3, 'content' => 'She smiled and said, “Let me teach you the magic of the alphabet, so you can write about our rice fields and tell stories of our village.” Dewi Sri waved her hand, and sparkling letters appeared in the sky. The children were amazed as each letter told its own little story about life in the village.']);

        // Halaman 4
        $page4 = Page::create(['book_id' => $book->id, 'page_number' => 4, 'content' => 'A was for Andong, the cart pulled by buffalo. “Andong helps carry the rice to the market!” Dewi Sri explained.']);
        Interaction::create(['page_id' => $page4->id, 'type' => 'tap_the_letter', 'data' => ['letter' => 'A', 'instruction' => 'Tap the letter that you see, hurry up!'], 'points' => 5]);

        // Halaman 5
        $page5 = Page::create(['book_id' => $book->id, 'page_number' => 5, 'content' => 'B was for Beras, the rice that everyone ate. “Beras gives us strength and happiness,” she said.']);
        Interaction::create(['page_id' => $page5->id, 'type' => 'tap_the_letter', 'data' => ['letter' => 'B', 'instruction' => 'Tap the letter that you see, hurry up!'], 'points' => 5]);

        // Halaman 6
        $page6 = Page::create(['book_id' => $book->id, 'page_number' => 6, 'content' => 'C was for Candi, the temple where people gave thanks. “It reminds us to be grateful for the blessings we have.” The children repeated after her, “A is for Andong, B is for Beras, C is for Candi!”']);
        Interaction::create(['page_id' => $page6->id, 'type' => 'tap_the_letter', 'data' => ['letter' => 'C', 'instruction' => 'Tap the letter that you see, hurry up!'], 'points' => 5]);

        // Halaman 7
        $page7 = Page::create(['book_id' => $book->id, 'page_number' => 7, 'content' => 'As they walked, Dewi Sri showed them more letters: D was for Dewa, the gods who protected the fields. E was for Emas, the golden rice that shone in the sunlight. F was for Fauna, the animals in the fields, like ducks and water buffalo.']);
        Interaction::create(['page_id' => $page7->id, 'type' => 'tap_the_letter', 'data' => ['letters' => ['D', 'E', 'F'], 'instruction' => 'Tap the letter that you see, hurry up!'], 'points' => 10]);

        // Halaman 8
        $page8 = Page::create(['book_id' => $book->id, 'page_number' => 8, 'content' => 'The children laughed when they saw G, which stood for Gong, the big musical instrument that echoed through the village. “We play the gong during festivals!” Dewi Sri said, clapping her hands.']);
        Interaction::create(['page_id' => $page8->id, 'type' => 'tap_the_letter', 'data' => ['letter' => 'G', 'instruction' => 'Tap the letter that you see, hurry up!'], 'points' => 5]);

        // Halaman 9
        $page9 = Page::create(['book_id' => $book->id, 'page_number' => 9, 'content' => 'Soon, the children began spotting letters everywhere: H was for Hujan, the rain that helped the rice grow. I was for Ikan, the fish swimming in the ponds near the fields. J was for Jagung, the corn that grew alongside the rice.']);
        Interaction::create(['page_id' => $page9->id, 'type' => 'tap_the_letter', 'data' => ['letters' => ['H', 'I', 'J'], 'instruction' => 'Tap the letter that you see, hurry up!'], 'points' => 10]);

        // Halaman 10
        Page::create(['book_id' => $book->id, 'page_number' => 10, 'content' => 'Dewi Sri encouraged the children to use their voices. “Can you say, ‘ H is for Hujan, I is for Ikan, J is for Jagung’?” They shouted happily, “ H is for Hujan, I is for Ikan, J is for Jagung!”']);

        // Halaman 11
        $page11 = Page::create(['book_id' => $book->id, 'page_number' => 11, 'content' => 'As the sun set, Dewi Sri introduced the last letters: K for Keris, the traditional dagger. L for Lumbung, the rice barn where the harvest was stored. M for Mangga, the juicy mangoes they loved to eat. Repeat please, “K for Keris, L for Lumbung, M for Mangga!”']);
        Interaction::create(['page_id' => $page11->id, 'type' => 'tap_the_letter', 'data' => ['letters' => ['K', 'L', 'M'], 'instruction' => 'Tap the letter that you see, hurry up!'], 'points' => 10]);

        // Halaman 12
        Page::create(['book_id' => $book->id, 'page_number' => 12, 'content' => 'The children clapped their hands and sang about the letters. “We’ve learned the alphabet with Dewi Sri, and now we can tell stories about our home!” Let’s sing this song together!']);

        // Halaman 13
        $page13 = Page::create(['book_id' => $book->id, 'page_number' => 13, 'content' => 'Dewi Sri smiled warmly and gave them a gift—a magical book where they could write their stories and draw pictures of their favourite letters. “Remember,” she said, “the alphabet is like planting rice. Each letter is a seed that grows into beautiful words and stories.”']);
        Interaction::create(['page_id' => $page13->id, 'type' => 'tap_the_sound', 'data' => ['instruction' => 'Listen carefully and tap the correct letter as you hear it!'], 'points' => 15]);

        // Halaman 14
        Page::create(['book_id' => $book->id, 'page_number' => 14, 'content' => 'From that day on, the children in the village of Dewi Sri became the best storytellers. They wrote about their fields, their animals, and their festivals, sharing their culture with the world. The End.']);

        // Halaman 15 - Post-activity
        // Aktivitas 1 (Memasangkan huruf dengan gambar)
        PostActivity::create([
            'book_id' => $book->id,
            'type' => 'match_the_picture',
            'order' => 1,
            'points' => 20,
            'data' => [
                'instruction' => 'Match the alphabet to the right picture!',
                'pairs' => [
                    ['letter' => 'A', 'image' => 'url/to/andong.png', 'name' => 'Andong'],
                    ['letter' => 'B', 'image' => 'url/to/beras.png', 'name' => 'Beras'],
                    ['letter' => 'C', 'image' => 'url/to/candi.png', 'name' => 'Candi'],
                    // ... tambahkan pasangan lainnya hingga J
                ]
            ],
        ]);

        // Aktivitas 2 (Menebalkan tulisan)
        PostActivity::create([
            'book_id' => $book->id,
            'type' => 'trace_the_word',
            'order' => 2,
            'points' => 15,
            'data' => [
                'instruction' => 'Drag and drop the letter into the right place to write the name of the picture!',
                'image' => 'https://placehold.co/300x200/FFCDB2/4F4F4F?text=Image+of+Andong',
                'word' => 'ANDONG'
            ],
        ]);
    }
}
