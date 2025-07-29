<?php

namespace App\Traits;

use App\Models\Interaction;
use App\Models\PostActivity;

trait AnswerChecker
{
    /**
     * Check answers for interactions within a page.
     *
     * @param Interaction $interaction
     * @param mixed $studentAnswer
     * @return bool
     */
    protected function checkInteractionAnswer(Interaction $interaction, $studentAnswer): bool
    {
        $correctAnswerData = $interaction->data;

        switch ($interaction->type) {
            case 'drag_and_drop_vocab':
            case 'tap_and_count':
            case 'tap_to_count':
            case 'drag_and_drop_count':
                // Untuk tipe ini, asumsikan frontend mengirim jumlah yang benar.
                // Contoh: $studentAnswer = 5
                return (int) $studentAnswer === (int) $correctAnswerData['count'];

            case 'tap_the_letter':
                // Asumsikan frontend mengirim huruf yang ditekan.
                // Contoh: $studentAnswer = "A"
                $letters = $correctAnswerData['letters'] ?? [$correctAnswerData['letter']];
                return in_array(strtoupper($studentAnswer), $letters);

            case 'arrange_letters':
                // Asumsikan frontend mengirim kata yang sudah disusun.
                // Contoh: $studentAnswer = "FATHER"
                return strtolower($studentAnswer) === strtolower($correctAnswerData['word']);

            case 'tap_body_part':
                // Asumsikan frontend mengirim nama bagian tubuh.
                // Contoh: $studentAnswer = "hands"
                return strtolower($studentAnswer) === strtolower($correctAnswerData['body_part']);

                // Tipe lain bisa ditambahkan di sini...

            default:
                return false;
        }
    }

    /**
     * Check answers for post-activity at the end of the book.
     *
     * @param PostActivity $postActivity
     * @param mixed $studentAnswer
     * @return bool
     */
    protected function checkPostActivityAnswer(PostActivity $postActivity, $studentAnswer): bool
    {
        $correctAnswerData = $postActivity->data;

        switch ($postActivity->type) {
            case 'match_picture_to_number':
            case 'label_body_parts':
            case 'family_tree_drag':
                // Untuk tipe mencocokkan/drag-drop, asumsikan frontend mengirim
                // array jawaban yang sudah divalidasi dan kita tinggal return true.
                // Validasi kompleks lebih baik dilakukan di frontend untuk tipe ini.
                return (bool) $studentAnswer;

            case 'order_the_number':
                // Asumsikan frontend mengirim array angka yang sudah diurutkan.
                // Contoh: $studentAnswer = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
                $correctOrder = range($correctAnswerData['range_start'], $correctAnswerData['range_end']);
                return $studentAnswer === $correctOrder;

                // Tipe lain bisa ditambahkan di sini...

            default:
                return false;
        }
    }
}
