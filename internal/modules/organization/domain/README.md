# Organization Domain

## Overview

The organization module manages the location hierarchy and shift patterns for the plant.
It follows Hexagonal (Ports & Adapters) architecture.

---

## Location

Locations form a tree hierarchy with four kinds:

| Kind      | Description                        |
|-----------|------------------------------------|
| `PLANT`   | Root node — no shift pattern       |
| `AREA`    | Area within a plant                |
| `LINE`    | Production line within an area     |
| `SECTION` | Section within a line              |

Each non-PLANT location holds a `shift_pattern_id` that links it to a `ShiftPattern`.

### Location HTTP Routes

| Method | Path                                                            | Description                                           |
|--------|-----------------------------------------------------------------|-------------------------------------------------------|
| POST   | `/v1/api/organization/locations/`                               | Create a root (PLANT) location                        |
| POST   | `/v1/api/organization/locations/add`                            | Add a child location to a tree                        |
| GET    | `/v1/api/organization/locations/`                               | Get all location trees                                |
| GET    | `/v1/api/organization/locations/:location_code`                 | Get a location tree by root code                      |
| POST   | `/v1/api/organization/locations/:location_code/shift-patterns`  | Create a shift pattern and link it to the location    |

---

## Shift Pattern

A `ShiftPattern` is a rotating sequence of `ShiftEntry` items anchored to a reference date.
The cycle length is derived from the number of entries (e.g. 3 = tri-daily, 7 = weekly).

Shift patterns are always created in the context of a location (`POST .../locations/:location_code/shift-patterns`).
The pattern is persisted and the location's `shift_pattern_id` is updated atomically in the same request.

### ShiftPattern structure

```json
{
  "id": "uuid",
  "name": "Default 3-Shift Rotation",
  "ref_start_date": "2024-01-01",
  "period_days": 3,
  "entries": [
    { "day_index": 0, "name": "Morning",   "start_time": "06:00", "end_time": "14:00", "is_off": false },
    { "day_index": 1, "name": "Afternoon", "start_time": "14:00", "end_time": "22:00", "is_off": false },
    { "day_index": 2, "name": "Night",     "start_time": "22:00", "end_time": "06:00", "is_off": false }
  ],
  "created_at": "2024-01-01T00:00:00Z"
}
```

**Notes:**
- `start_time` / `end_time` use `"HH:MM"` format (24-hour).
- Overnight shifts (where `end_time < start_time`) automatically wrap to the next calendar day.
- Off-day entries set `is_off: true` and omit `start_time` / `end_time`.
- `ref_start_date` uses `"YYYY-MM-DD"` format and anchors the rotation cycle.

### ShiftPattern HTTP Routes

| Method | Path                                                            | Description                                           |
|--------|-----------------------------------------------------------------|-------------------------------------------------------|
| POST   | `/v1/api/organization/locations/:location_code/shift-patterns`  | Create a shift pattern and link it to the location    |
| GET    | `/v1/api/organization/shift-patterns/:id`                       | Get a shift pattern by UUID                           |
| PUT    | `/v1/api/organization/shift-patterns/:id`                       | Update a shift pattern name and entries               |
| DELETE | `/v1/api/organization/shift-patterns/:id`                       | Delete a shift pattern                                |

### Create — Request Body

```json
{
  "name": "3-shift rotation",
  "ref_start_date": "2024-01-01",
  "entries": [
    { "name": "Morning",   "start_time": "06:00", "end_time": "14:00" },
    { "name": "Afternoon", "start_time": "14:00", "end_time": "22:00" },
    { "name": "Night",     "start_time": "22:00", "end_time": "06:00" },
    { "is_off": true }
  ]
}
```

### Update — Request Body

```json
{
  "name": "Updated rotation",
  "entries": [
    { "name": "Morning",   "start_time": "07:00", "end_time": "15:00" },
    { "name": "Afternoon", "start_time": "15:00", "end_time": "23:00" }
  ]
}
```
