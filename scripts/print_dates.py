#!/usr/bin/env python3

import json

from datetime import date, timedelta


def new_slot(name):
    return {
        'Name': name,
        'Limit': 18,
        'Attendees': [],
    }


startDate = date(2021, 4, 7)
endDate   = date(2021, 9, 17)

slots = []
for i in range((endDate - startDate).days + 1):
    day = startDate + timedelta(days=i)
    if (day.weekday() >= 5):
        continue
    slots.append(new_slot(day.strftime("%Y %b %-d 14:00 - 16:00")))
    slots.append(new_slot(day.strftime("%Y %b %-d 17:00 - 19:00")))

print(json.dumps({"Slots": slots}))
