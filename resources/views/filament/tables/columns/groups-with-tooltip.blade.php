<div>
    @if ($getRecord()->groups->isNotEmpty())
    <div class="flex flex-wrap gap-1">
        @foreach ($getRecord()->groups as $group)
        <span
            class="fi-badge flex items-center justify-center gap-x-1 rounded-md bg-gray-100 px-2 py-1 text-xs font-medium text-gray-600 ring-1 ring-inset ring-gray-200 dark:bg-gray-500/10 dark:text-gray-400 dark:ring-gray-500/20"
            title="Class: {{ $group->course->name }}">
            {{ $group->name }}
        </span>
        @endforeach
    </div>
    @else
    <span class="text-gray-500 dark:text-gray-400">Unassigned</span>
    @endif
</div>