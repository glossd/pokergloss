
### cards:

```js
db.tables.aggregate([{$match:{_id:ObjectId("60817bd8383a2d835945c6e9")}}, {$project: {_id:0, f: "$seats.player.cards.first.str", s: "$seats.player.cards.second.str"}}])
```