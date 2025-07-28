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
        // Hapus data lama untuk menghindari duplikasi
        Book::query()->delete();

        // Panggil method untuk membuat setiap buku
        $this->createDewiSriBook();
        $this->createTenHillsBook();
        $this->createBawangPutihBook();
    }

    /**
     * Membuat data untuk buku "The Alphabet in the Land of Dewi Sri".
     */
    private function createDewiSriBook(): void
    {
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

    /**
     * Membuat data untuk buku "The Numbers in the Village of Ten Hills".
     */
    private function createTenHillsBook(): void
    {
        $book2 = Book::create([
            'title' => 'The Numbers in the Village of Ten Hills',
            'author' => 'Miss Ani',
            'year' => '2025',
            'genre' => 'Fiction',
            'focus' => 'Numbers',
            'status' => 'draft',
            'cover_image' => 'https://placehold.co/400x600/89CFF0/FFFFFF?text=Ten+Hills',
            'order_sequence' => 2,
        ]);

        // Halaman 1
        Page::create(['book_id' => $book2->id, 'page_number' => 1, 'content' => 'Once upon a time, in a peaceful village called Ten Hills, there were ten hills where each hill was home to a magical number. The numbers lived together and helped the villagers with their daily tasks. One morning, a little boy named Eza was playing near the hills when he heard a soft voice calling out to him.']);

        // Halaman 2
        $page2 = Page::create(['book_id' => $book2->id, 'page_number' => 2, 'content' => '“Hello, Eza!” said the voice. Eza looked up to see Number 1 standing tall on the first hill.“I am Number 1. I love being first in line! Can you find one rooster on my hill?” Eza looked around and found a bright red rooster. “There it is! One rooster!” he exclaimed.']);
        Interaction::create(['page_id' => $page2->id, 'type' => 'tap_and_say', 'data' => ['instruction_tap' => 'Tap the rooster!', 'sound_on_tap' => 'one', 'instruction_say' => 'Say it! “One rooster!”', 'correct_phrase' => 'one rooster'], 'points' => 10]);

        // Halaman 3
        $page3 = Page::create(['book_id' => $book2->id, 'page_number' => 3, 'content' => 'As Eza walked to the second hill, Number 2 appeared, holding two shiny coconuts.“I am Number 2,” it said with a smile. “I like things in pairs. Can you find two ducks swimming in the pond?” Eza laughed as he spotted the ducks. “One, two! Two ducks!”']);
        Interaction::create(['page_id' => $page3->id, 'type' => 'tap_and_say', 'data' => ['instruction_tap' => 'Can you find the ducks? Tap them!', 'sound_on_tap' => ['one', 'two'], 'instruction_say' => 'Say it! “two ducks!”', 'correct_phrase' => 'two ducks'], 'points' => 10]);

        // Halaman 4
        $page4 = Page::create(['book_id' => $book2->id, 'page_number' => 4, 'content' => 'On the third hill, Number 3 was playing with three colourful kites.“Help me count my kites, Eza!” Number 3 said.Eza pointed to each kite: “One, two, three!”']);
        Interaction::create(['page_id' => $page4->id, 'type' => 'tap_to_count', 'data' => ['instruction' => 'Can you count the kites? Tap them!', 'count' => 3], 'points' => 10]);

        // Halaman 5
        $page5 = Page::create(['book_id' => $book2->id, 'page_number' => 5, 'content' => 'Number 4 was on the fourth hill, building a fence with four strong bamboos. “I need four bamboos; can you take it for me?” it asked. Eza took the bamboos while counting: “One, two, three, four!”']);
        Interaction::create(['page_id' => $page5->id, 'type' => 'drag_and_drop_count', 'data' => ['instruction' => 'Can you move the bamboos to the fence?', 'item_name' => 'bamboo', 'count' => 4], 'points' => 15]);

        // Halaman 6
        $page6 = Page::create(['book_id' => $book2->id, 'page_number' => 6, 'content' => 'Number 5 was grilling five sweet potatoes. “Let’s share them with the villagers!” it said. Eza shared the sweet potatoes while counting, “One, two, three, four, five!”']);
        Interaction::create(['page_id' => $page6->id, 'type' => 'drag_and_drop_count', 'data' => ['instruction' => 'Can you give the grilled sweet potatoes to the villagers?', 'item_name' => 'sweet potato', 'count' => 5], 'points' => 15]);

        // Halaman 7
        $page7 = Page::create(['book_id' => $book2->id, 'page_number' => 7, 'content' => 'Number 6 was teaching six little goats how to jump. “Count them as they leap!” it giggled. Eza counted the goats, “One, two, three, four, five, six!”']);
        Interaction::create(['page_id' => $page7->id, 'type' => 'voice_recognition_leap', 'data' => ['instruction' => 'Call them and they will leap!', 'names' => ['1', '2', '3', '4', '5', '6']], 'points' => 15]);

        // Halaman 8
        $page8 = Page::create(['book_id' => $book2->id, 'page_number' => 8, 'content' => 'Number 7 was painting a picture of the seven colours of the rainbow. “Can you count the colours on the rainbow?” it asked. Eza counted the colours, “One, two, three, four, five, six, seven!”']);
        Interaction::create(['page_id' => $page8->id, 'type' => 'voice_recognition_repeat', 'data' => ['instruction' => 'Say it loudly after me!', 'phrase_to_repeat' => 'one two three four five six seven'], 'points' => 10]);

        // Halaman 9
        $page9 = Page::create(['book_id' => $book2->id, 'page_number' => 9, 'content' => 'When Eza reached the eighth hill, Number 8 was blowing bubbles.“Count my bubbles, Eza!” it said as eight shiny bubbles floated into the sky.Eza counted, “One, two, three, four, five, six, seven, eight!”']);
        Interaction::create(['page_id' => $page9->id, 'type' => 'voice_recognition_pop', 'data' => ['instruction' => 'Say it loudly after me!', 'phrase_to_repeat' => 'one two three four five six seven eight'], 'points' => 10]);

        // Halaman 10
        $page10 = Page::create(['book_id' => $book2->id, 'page_number' => 10, 'content' => 'At the ninth hill, Number 9 was harvesting nine juicy mangoes. “Help me count them so I can share them with the village,” it said.Eza happily counted: “One, two, three, four, five, six, seven, eight, nine mangoes!”']);
        Interaction::create(['page_id' => $page10->id, 'type' => 'drag_and_drop_count', 'data' => ['instruction' => 'Put the mangoes in the basket!', 'item_name' => 'mango', 'count' => 9], 'points' => 15]);

        // Halaman 11
        $page11 = Page::create(['book_id' => $book2->id, 'page_number' => 11, 'content' => 'Finally, Eza climbed the tenth hill, the tallest of them all. Number 10 was standing proudly, holding ten golden coins.“I am Number 10,” it said. “I bring all the numbers together. Can you count all the hills in our village?” Eza turned and counted each hill. “One, two, three, four, five, six, seven, eight, nine, ten hills!”. Once Eza said, “ten hills!” a wooden box filled with colourful numbers cards appear.']);
        Interaction::create(['page_id' => $page11->id, 'type' => 'tap_and_say', 'data' => ['instruction' => 'Tap the mountain one by one and shout loudly, “ten hills!”', 'tap_count' => 10, 'correct_phrase' => 'ten hills'], 'points' => 20]);

        // Halaman 12
        Page::create(['book_id' => $book2->id, 'page_number' => 12, 'content' => 'The numbers cheered and gave Eza the wooden box filled with colourful number cards as a special gift. “Now you can remember us and help others learn to count,” they said. Eza said, “Thank you very much.” From that day on, Eza became the best counter in the village, teaching his friends about the numbers of Ten Hills. The End.']);

        // Post-Activity 1 (dari Slide 13)
        PostActivity::create([
            'book_id' => $book2->id,
            'type' => 'match_picture_to_number',
            'order' => 1,
            'points' => 20,
            'data' => [
                'instruction' => 'Match the picture with the right number!',
                'pairs' => [
                    ['image' => 'url/to/rooster.png', 'count' => 1],
                    ['image' => 'url/to/ducks.png', 'count' => 2],
                    ['image' => 'url/to/kites.png', 'count' => 3],
                ]
            ],
        ]);

        // Post-Activity 2 (dari Slide 14)
        PostActivity::create([
            'book_id' => $book2->id,
            'type' => 'order_the_number',
            'order' => 2,
            'points' => 15,
            'data' => [
                'instruction' => 'Put the number on the mountains orderly!',
                'range_start' => 1,
                'range_end' => 10,
            ],
        ]);
    }

    /**
     * Membuat data untuk buku "Bawang Putih and the Kind Body Parts".
     */
    private function createBawangPutihBook(): void
    {
        $book3 = Book::create([
            'title' => 'Bawang Putih and the Kind Body Parts',
            'author' => 'Miss Ani',
            'year' => '2025',
            'genre' => 'Moral Folklore',
            'focus' => 'Body Parts',
            'status' => 'draft',
            'cover_image' => 'https://placehold.co/400x600/D8BFD8/FFFFFF?text=Bawang+Putih',
            'order_sequence' => 3,
        ]);

        // Halaman 1 - 4 (Narasi)
        Page::create(['book_id' => $book3->id, 'page_number' => 1, 'content' => 'Once upon a time, in a peaceful village, there lived a kind-hearted girl named Bawang Putih. She was hardworking and loved helping others.']);
        Page::create(['book_id' => $book3->id, 'page_number' => 2, 'content' => 'One day, while washing clothes by the river, her favourite scarf slipped into the water and disappeared. As she tried to reach it, an old woman appeared from the woods.']);
        Page::create(['book_id' => $book3->id, 'page_number' => 3, 'content' => '“Bawang Putih,” said the old woman with a kind smile, “I will return your scarf, but only if you visit my house and learn the secrets of your body parts.”']);
        Page::create(['book_id' => $book3->id, 'page_number' => 4, 'content' => 'Curious and respectful, Bawang Putih followed the old woman to her cozy little hut. Inside, the old woman brought out a magical mirror. “Look into this mirror,” she said, “and your body will tell you its secrets.”']);

        // Halaman 5 (Hands)
        $page5 = Page::create(['book_id' => $book3->id, 'page_number' => 5, 'content' => 'First, Bawang Putih’s Hands spoke. “We are your Hands, Bawang Putih. We help you carry baskets, sew clothes, and plant rice. But we also love to help others! Can you clap your Hands to show us you’re ready to help?”...']);
        Interaction::create(['page_id' => $page5->id, 'type' => 'tap_body_part', 'data' => ['instruction' => 'Touch my hands!', 'body_part' => 'hands', 'tap_count' => 2, 'sound' => 'hand'], 'points' => 10]);

        // Halaman 6 (Feet)
        $page6 = Page::create(['book_id' => $book3->id, 'page_number' => 6, 'content' => 'Next, her Feet tapped gently. “We are your Feet! We carry you to the market, the fields, and the river. But remember, we also love to dance when you’re happy! Can you stomp your Feet with joy?”...']);
        Interaction::create(['page_id' => $page6->id, 'type' => 'tap_body_part', 'data' => ['instruction' => 'Touch my feet!', 'body_part' => 'feet', 'tap_count' => 1, 'sound' => 'feet'], 'points' => 10]);

        // Halaman 7 (Eyes)
        $page7 = Page::create(['book_id' => $book3->id, 'page_number' => 7, 'content' => 'Her Eyes sparkled in the mirror. “Bawang Putih, we are your Eyes. We help you see the beauty of the sunrise, the flowing river, and the faces of those you love. Close us sometimes to rest so we can stay strong!”...']);
        Interaction::create(['page_id' => $page7->id, 'type' => 'tap_body_part', 'data' => ['instruction' => 'Touch my eyes!', 'body_part' => 'eyes', 'tap_count' => 2, 'sound' => 'eye'], 'points' => 10]);

        // Halaman 8 (Ears)
        $page8 = Page::create(['book_id' => $book3->id, 'page_number' => 8, 'content' => 'Her Ears perked up. “We are your Ears! We let you hear the laughter of children, the songs of birds, and the advice of elders. Always listen carefully to the sounds of the world. Can you touch your ears?”...']);
        Interaction::create(['page_id' => $page8->id, 'type' => 'tap_body_part', 'data' => ['instruction' => 'Touch my ears!', 'body_part' => 'ears', 'tap_count' => 2, 'sound' => 'ear'], 'points' => 10]);

        // Halaman 9 (Heart)
        $page9 = Page::create(['book_id' => $book3->id, 'page_number' => 9, 'content' => 'The Heart in her chest glowed warmly. “I am your Heart, Bawang Putih. I help you feel kindness and courage. Keep me healthy by sharing love and avoiding anger. Place your hands on your chest to feel your heart inside!”...']);
        Interaction::create(['page_id' => $page9->id, 'type' => 'tap_body_part', 'data' => ['instruction' => 'Touch my heart!', 'body_part' => 'heart', 'tap_count' => 1, 'sound' => 'heart'], 'points' => 10]);

        // Halaman 10 - 12 (Narasi Akhir)
        Page::create(['book_id' => $book3->id, 'page_number' => 10, 'content' => 'After learning about her body parts, the old woman returned Bawang Putih’s scarf and handed her a golden box. “Inside this box is a gift for your kindness and for understanding the wonders of your body,” she said.']);
        Page::create(['book_id' => $book3->id, 'page_number' => 11, 'content' => 'When Bawang Putih opened the box at home, it was filled with fruits and sweets that never ran out, a blessing for her hardworking spirit.']);
        Page::create(['book_id' => $book3->id, 'page_number' => 12, 'content' => 'From that day on, Bawang Putih used her body parts with gratitude, helping her neighbours, dancing at village festivals, and planting more rice fields. She shared her golden gift with everyone, spreading kindness throughout the village. The End.']);

        // Post-Activity 1 (dari Slide 13)
        PostActivity::create([
            'book_id' => $book3->id,
            'type' => 'label_body_parts',
            'order' => 1,
            'points' => 20,
            'data' => [
                'instruction' => 'Name the body parts correctly!',
                'parts' => ['hands', 'feet', 'eyes', 'ears', 'heart']
            ],
        ]);

        // Post-Activity 2 (dari Slide 14)
        PostActivity::create([
            'book_id' => $book3->id,
            'type' => 'group_verbs_to_body_part',
            'order' => 2,
            'points' => 25,
            'data' => [
                'instruction' => 'Group the verbs to the correct body part!',
                'groups' => [
                    ['part' => 'hands', 'verbs' => ['carry', 'sew', 'plant', 'clap']],
                    ['part' => 'feet', 'verbs' => ['walk', 'dance', 'stomp']],
                    ['part' => 'eyes', 'verbs' => ['see', 'close', 'open']],
                    ['part' => 'ears', 'verbs' => ['hear', 'listen']],
                    ['part' => 'heart', 'verbs' => ['feel']],
                ]
            ],
        ]);
    }
}
