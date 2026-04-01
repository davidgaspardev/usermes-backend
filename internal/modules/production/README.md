# Production Module

Manages production resources, events, and shift tracking within a plant.

---

## Domain Model

```
ShiftPattern (Organization)
    └── EntryForDate(date) → ShiftEntry
            └── TimesForDate(date) → startAt, endAt
                    └── NewShift(patternID, name, startAt, endAt)
                                │
Resource ───────────────────────┤
    └── resCode                 │
                                ▼
                            Event
                      (shiftID = Shift.ID)
```

---

## Event

An `Event` is an immutable record of something that happened to a **Resource** during a **Shift**.
Every state transition of a resource produces a new event — the current state is always the last event.

### Fields

| Field | Type | Description |
|---|---|---|
| `eventType` | `EventType` | What happened (see table below) |
| `status` | `EventStatus` | `PRODUCTION` or `STOP` |
| `statusCode` | `string` | Reason code for a stop (e.g. `6010`). Empty when producing |
| `resCode` | `string` | Resource this event belongs to (e.g. `INJ10`) |
| `prodCode` | `string` | Product/part being produced (e.g. `2101415`). Empty if no item loaded |
| `prodQuantity` | `uint32` | Quantity produced in this period |
| `scrapQuantity` | `uint32` | Quality loss / scrap quantity |
| `retestQuantity` | `uint32` | Quantity sent to retest |
| `toolCode` | `string` | Tool or mould in use. Empty if not applicable |
| `dateStart` | `time.Time` | When this event/period started |
| `dateEnd` | `*time.Time` | When it ended. `nil` means the period is still open |
| `shiftID` | `string` | ID of the `Shift` instance (see Production module) |
| `userID` | `string` | Operator who triggered the event. Empty if unattended |

### EventStatus

| Value | Meaning |
|---|---|
| `PRODUCTION` | Resource is actively producing |
| `STOP` | Resource is stopped (planned, unplanned, or unattended) |

---

## Event Types

### Lifecycle events

| EventType | Status | Description |
|---|---|---|
| `RESOURCE_CREATED` | `STOP` | Resource registered in the system. No operator, no item |
| `CHANGE_SHIFT` | *(carries over)* | Shift boundary crossed. Closes current period, opens new one under the next shift |

### Operator events

| EventType | Status | Description |
|---|---|---|
| `OPERATOR_SIGN_IN` | `PRODUCTION` | Operator logged into the resource (no item loaded) |
| `OPERATOR_SIGN_OUT` | `STOP` | Operator logged out, no item left |
| `OPERATOR_SIGN_IN_WITH_ITEM` | `PRODUCTION` | Operator logged in with a product already loaded |
| `OPERATOR_SIGN_OUT_WITH_ITEM` | `STOP` | Operator logged out but item remains on resource |
| `SWAP_OPERATOR` | `PRODUCTION` | A second operator takes over without interrupting production |

### Production events

| EventType | Status | Description |
|---|---|---|
| `OPERATOR_SET_PRODUCTION` | `PRODUCTION` | Operator manually set resource back to producing |
| `OPERATOR_CHANGE_ITEM` | `PRODUCTION` | Operator swapped the product being produced |
| `OPERATOR_ADD_ITEM` | `PRODUCTION` | Operator loaded a product onto the resource |
| `OPERATOR_REMOVE_ITEM` | `STOP` | Operator removed the product from the resource |

### Stop events

| EventType | Status | Description |
|---|---|---|
| `INSERT_STOP` | `STOP` | Operator recorded a stop reason (e.g. breakdown, setup) |

---

## Event lifecycle example

The table below traces resource `INJ10` across a full working day:

| Time | Status | Status Code | Prod Code | Qty | Scrap | Retest | Shift | User | Event Type | Notes |
|---|---|---|---|---|---|---|---|---|---|---|
| 05:50–06:30 | STOP | 1001 | — | — | — | — | INJ_1_T1 | BR056205 | `RESOURCE_CREATED` | Machine registered, unattended |
| 06:30–07:00 | PRODUCTION | — | 2101415 | 15 | 0 | 0 | T1 | BR056205 | `OPERATOR_SIGN_IN_WITH_ITEM` | Operator logged in with item |
| 07:00–10:00 | PRODUCTION | — | 2101415 | 30 | 0 | 0 | T1 | DA0356879 | `SWAP_OPERATOR` | Second operator took over |
| 10:00–11:00 | STOP | 6010 | 2101415 | 0 | 0 | 0 | T1 | DA0356879 | `INSERT_STOP` | Stop reason recorded |
| 11:00–12:00 | STOP | 6010 | 2101415 | 0 | 0 | 0 | T1 | — | `OPERATOR_SIGN_OUT_WITH_ITEM` | Operator left, item still on machine |
| 11:00–12:00 | STOP | 6010 | 2101415 | 0 | 0 | 0 | T1 | JD056205 | `OPERATOR_SIGN_IN` | New operator signed in |
| 12:00–13:00 | PRODUCTION | — | 2101415 | 0 | 0 | 0 | T1 | JD056205 | `OPERATOR_SET_PRODUCTION` | Operator resumed production |
| 13:00–14:10 | PRODUCTION | — | 500200 | 10 | 0 | 0 | T1 | JD056205 | `OPERATOR_CHANGE_ITEM` | Different product loaded |
| 14:10–14:20 | PRODUCTION | — | 500200 | 10 | 3 | 3 | T2 | JD056205 | `CHANGE_SHIFT` | Shift T1 → T2, period closed |
| 14:20–14:30 | STOP | 99 | — | — | — | — | T2 | JD056205 | `OPERATOR_REMOVE_ITEM` | Item removed from resource |
| 14:30–15:00 | PRODUCTION | — | 800200 | 10 | 3 | 3 | T2 | JD056205 | `OPERATOR_ADD_ITEM` | New item loaded |

### Key observations from the timeline

1. **Current state = last event** — to know what `INJ10` is doing right now, read the most recent event.
2. **`dateEnd = nil` on the open event** — the active period has no end date until a new event arrives.
3. **`CHANGE_SHIFT` seals a period** — when the shift boundary is crossed, the running event is closed with `dateEnd` and a new event opens under the new `shiftID`.
4. **`OPERATOR_SIGN_OUT_WITH_ITEM` + `OPERATOR_SIGN_IN`** can appear at the same time — operator handoff without stopping the item on the machine.
5. **`statusCode` is only meaningful when `status = STOP`** — it identifies the stop reason (e.g. `6010` = planned stop, `99` = item removal, `1001` = machine unattended at creation).
6. **`prodCode` can be `NULL`** — the resource exists but has no item loaded (e.g. right after `RESOURCE_CREATED`).

---

## Shift resolution

When a new event is created the caller must:

1. Fetch the `ShiftPattern` assigned to the resource's location
2. Call `pattern.EntryForDate(now)` → `ShiftEntry`
3. Call `entry.TimesForDate(now)` → `startAt`, `endAt`
4. Create `entity.NewShift(pattern.ID(), entry.Name(), startAt, endAt)`
5. Persist the `Shift` and pass `shift.ID().String()` as `shiftID` to `NewEvent`
