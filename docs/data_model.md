# Data model

> The data model is designed according to a relational database

## Menu entity(`permission`)

| Field       | Description       | Field Type       | Remarks                 |
| ----------- | ----------------- | ---------------- | ----------------------- |
| id          | AutoIncrement ID  | Numerical value  | Primary key             |
| uuid        | UUID              | String           |                         |
| name        | Menu name         | String           |                         |
| sequence    | Sort value        | Numerical value  |                         |
| icon        | Icon              | String           |                         |
| router      | Access routing    | String           |                         |
| hidden      | Hidden menu       | Numerical value  | 0: not hidden 1: hidden |
| parent_id   | Parent ID         | String           |                         |
| parent_path | Parent path       | String           |                         |
| creator     | Creator           | String           |                         |
| created_at  | Creation time     | Time format      |                         |
| updated_at  | Update time       | Time format      |                         |
| deleted_at  | Delete time       | Time format      |                         |

## Menu action associated entity(`menu_action`)

| Field       | Description        | Field Type      | Remarks         |
| ----------- | ------------------ | --------------- | --------------- |
| id          | AutoIncrement ID   | Numerical value | Primary key     |
| menu_id     | Menu ID            | String          |                 |
| code        | Action number      | String          |                 |
| name        | Action name        | String          |                 |
| created_at  | Creation time      | Time format     |                 |
| updated_at  | Update time        | Time format     |                 |
| deleted_at  | Delete time        | Time format     |                 |

## Menu resource associated entity(`menu_resource`)

| Field       | Description       | Field Type       | Remarks         |
| ----------- | ----------------- | ---------------- | --------------- |
| id          | AutoIncrement ID  | Numerical value  | Primary key     |
| menu_id     | Menu ID           | String           |                 |
| code        | Resource number   | String           |                 |
| name        | Resource Name     | String           |                 |
| method      | Request method    | String           |                 |
| path        | Request path      | String           |                 |
| created_at  | Creation time     | Time format      |                 |
| updated_at  | Update time       | Time format      |                 |
| deleted_at  | Delete time       | Time format      |                 |

## Role entity(`role`)

| Field       | Description       | Field Type       | Remarks         |
| ----------- | ----------------- | ---------------- | --------------- |
| id          | AutoIncrement ID  | Numerical value  | Primary key     |
| uuid        | UUID              | String           |                 |
| name        | Role Name         | String           |                 |
| sequence    | Sort value        | Numerical value  |                 |
| memo        | Remarks           | String           |                 |
| creator     | creator           | String           |                 |
| created_at  | Creation time     | Time format      |                 |
| updated_at  | Update time       | Time format      |                 |
| deleted_at  | Delete time       | Time format      |                 |

## Role menu associated entity(`role_menu`)

| Field       | Description          | Field Type      | Remarks                                        |
| ----------- | -------------------- | --------------- | ---------------------------------------------- |
| id          | AutoIncrement ID     | Numerical value | Primary key                                    |
| role_id     | Role ID              | String          |                                                |
| menu_id     | Menu ID              | String          |                                                |
| action      | Action permission    | String          | Action number (multiple separated by comma     |
| resource    | Resource permissions | String          | Resource number (multiple separated by commas) |
| created_at  | Creation time        | Time format     |                                                |
| updated_at  | Update time          | Time format     |                                                |
| deleted_at  | Delete time          | Time format     |                                                |

## User entity(`user`)

| Field      | Description       | Field Type       | Remarks                |
| ---------- | ----------------- | ---------------- | ---------------------- |
| id         | AutoIncrement ID  | Numerical value  | Primary key            |
| uuid       | UUID              | String           |                        |
| user_name  | Login Username    | String           |                        |
| password   | Password          | String           |                        |
| real_name  | Real name         | String           |                        |
| email      | Email             | String           |                        |
| phone      | Phone number      | String           |                        |
| status     | Status            | Numerical value  | 1: enabled 2: disabled |
| creator    | Creator           | String           |                        |
| created_at | Creation time     | Time format      |                        |
| updated_at | Update time       | Time format      |                        |
| deleted_at | Delete time       | Time format      |                        |

## User role association entity(`user_role`)

| Field      | Description       | Field Type       | Remarks     |
| ---------- | ----------------- | ---------------- | ----------- |
| id         | AutoIncrement ID  | Numerical value  | Primary key |
| user_id    | User ID           | String           |             |
| role_id    | Role ID           | String           |             |
| created_at | Creation time     | Time format      |             |
| updated_at | Update time       | Time format      |             |
| deleted_at | Delete time       | Time format      |             |
