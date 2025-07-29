<?php

namespace App\Enums;

enum BookProgressEnum: string
{
    case NOT_STARTED = 'not_started';
    case IN_PROGRESS = 'in_progress';
    case COMPLETED = 'completed';
}
