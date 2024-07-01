**string**

| in | out |
| --- | --- |
| "abc" | "abc" |
| 123 | "123" |
| 0 | "0" |
| -1 | "-1" |
| 11.11 | "11.11" |
| [] | "" |
| {} | "" |
| null | "" |

**int**

| in | out |
| --- | --- |
| 20 | 20 |
| "20" | 20 |
| 11.11 | 11 |
| 11.6 | 11 |
| 0.1 | 0 |
| -0.1 | 0 |
| [] | 0 |
| {} | 0 |
| null | 0 |
| "" | 0 |

**bool**

| in | out |
| --- | --- |
| true | ture |
| false | false |
| "true" | true |
| "false" | false|
| "1" | true |
| "0" | false |
| 123 | false |
| "123" | false |
| 0.1 | false |
| -0.1 | false |
| [1] | false |
| ["ching"] | false |
| [] | false |
| {} | false |
| null | false |
| "" | false |