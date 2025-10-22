<?php

return [
    'resources' => [
        'label'                     => 'Activity Log',
        'plural_label'              => 'Activity Logs',
        'hide_restore_action'       => false,
        'restore_action_label'      => 'Restore',
        'hide_resource_action'      => false,
        'hide_restore_model_action' => true,
        'resource_action_label'     => 'View',
        'navigation_item'           => true,
        'navigation_group'          => 'Settings',
        'navigation_icon'           => '',
        'navigation_sort'           => 1,
        'default_sort_column'       => 'id',
        'default_sort_direction'    => 'desc',
        'navigation_count_badge'    => true,
        'resource'                  => \Rmsramos\Activitylog\Resources\ActivitylogResource::class,
    ],
    'date_format'     => 'd/m/Y',
    'datetime_format' => 'd/m/Y H:i:s',
];
